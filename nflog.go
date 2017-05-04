// nflog module
// This module is an input module that read the netlink nflog events, then parse them and transfer to outputs.

package main

import (
	"log"
	"github.com/Alkorin/nflog"
	"errors"
	"encoding/binary"
)

var (
	ERR_EXIT_SENT = errors.New("nflog: an error signal has sent.")
)

type nflogModule struct {
	group uint16
	copyRange uint16
	running bool
	exiting chan func ()
	outputs []outputModule
}

func newNflogInput(group uint16) *nflogModule {
	return &nflogModule{
		group: group,
		copyRange: 64,
		running: false,
		exiting: make(chan func(), 1),
	}
}

func (nfl *nflogModule) addOutput(output outputModule) {
	if nfl.outputs == nil {
		nfl.outputs = []outputModule{}
	}
	nfl.outputs = append(nfl.outputs, output)
}

func (nfl *nflogModule) _parse(m nflog.Msg) map[string][]byte {
	info := make(map[string][]byte)
	info["family"] = []byte{m.Family}
	info["prefix"] = []byte(m.Prefix)
	info["raw"] = m.Payload
	if m.InDev != nil {
		indev := make([]byte, 4)
		binary.BigEndian.PutUint32(indev, *m.InDev)
		info["indev"] = indev
	}
	if m.OutDev != nil {
		outdev := make([]byte, 4)
		binary.BigEndian.PutUint32(outdev, *m.OutDev)
		info["outdev"] = outdev
	}

	return info
}

func (nfl *nflogModule) _work(n *nflog.NFLog) {
	defer func () {
		nfl.running = false
	} ()

	errorLogged := false
	for {
		select {
		case m := <-n.Messages():
			result := nfl._parse(m)
			for _, out := range nfl.outputs {
				out.produce(result)
			}
		case e := <-n.Errors():
			if !errorLogged {
				errorLogged = true
				log.Printf("we encountered a error: %s, we may lose some events.\n", e.Error())
			}
		case f := <-nfl.exiting:
			f()
			return
		}
	}
}

func (nfl *nflogModule) start() error {
	if nfl.running {
		return nil
	}

	conf := nflog.NewConfig()
	conf.Groups = []uint16{nfl.group}
	if nfl.copyRange > 0 {
		conf.CopyRange = nfl.copyRange
	} else {
		conf.CopyRange = 64
	}
	conf.Return.Errors = true

	n, err := nflog.New(conf)
	if err != nil {
		return err
	}

	nfl.running = true
	go nfl._work(n)
	return nil
}

func (nfl *nflogModule) stop(callback func()) error {
	select {
	case nfl.exiting <- callback:
		return nil
	default:
		return ERR_EXIT_SENT
	}
	// unreachable
	return nil
}