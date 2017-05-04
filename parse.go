package main

import (
	"encoding/binary"
	"bytes"
	"fmt"
	"os"
	"log"
	"bufio"
)

var IP_PROTOCOL [256]string

func init() {
	// Load ip protocol string name
	f, err := os.OpenFile("ip_protocol.txt", os.O_RDONLY, 0666)
	if err != nil {
		log.Println(err)
	} else {
		i := 0
		reader := bufio.NewReader(f)
		for l, _, err := reader.ReadLine(); i < 256 && err == nil; l, _, err = reader.ReadLine() {
			IP_PROTOCOL[i] = string(l)
			i++
		}
		f.Close()
	}
}

type iphdr struct {
	Version_inl uint8 // version:4 inl:4
	Tos         uint8
	Tot_len     uint16
	Id          uint16
	Frag_off    uint16
	Ttl         uint8
	Protocol    uint8
	Check       uint16
	Sa          [4]uint8
	Da          [4]uint8
}

func parseIPPacket(raw []byte) map[string]string {
	// TODO: use raw array to replace map
	ret := make(map[string]string)
	reader := bytes.NewReader(raw)
	iph := iphdr{}
	err := binary.Read(reader, binary.BigEndian, &iph)
	if err != nil {
		ret["ERR"] = err.Error()
		return ret
	}
	ret["IP_SADDR"] = fmt.Sprintf("%d.%d.%d.%d",
		iph.Sa[0], iph.Sa[1], iph.Sa[2], iph.Sa[3])
	ret["IP_DADDR"] = fmt.Sprintf("%d.%d.%d.%d",
		iph.Da[0], iph.Da[1], iph.Da[2], iph.Da[3])
	ret["PROTOCOL"] = IP_PROTOCOL[iph.Protocol]
	ret["CHECKSUM"] = fmt.Sprintf("%d", iph.Check)
	ret["TTL"] = fmt.Sprintf("%d", iph.Ttl)
	ret["TOS"] = fmt.Sprintf("%d", iph.Tos)
	ret["TOT_LEN"] = fmt.Sprintf("%d", iph.Tot_len)
	ret["ID"] = fmt.Sprintf("%d", iph.Id)
	ret["FRAG_OFF"] = fmt.Sprintf("%d", iph.Frag_off)
	return ret
}
