package main

type outputModule interface {
	produce(map[string][]byte)
}

type inputModule interface {
	start() error
	running() int
	stop(callback func())
	addOutput(outputModule)
}
