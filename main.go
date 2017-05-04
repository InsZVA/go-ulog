package main

func main() {
	nfl := newNflogInput(2)
	con := &consoleModule{}
	nfl.addOutput(con)
	nfl.start()
	for{}
}
