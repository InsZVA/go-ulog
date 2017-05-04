package main

import "log"

type consoleModule struct {}

func (con *consoleModule) produce(data map[string][]byte) {
	log.Println(parseIPPacket(data["raw"]))
}
