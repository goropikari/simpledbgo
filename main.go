package main

import "github.com/goropikari/simpledbgo/server"

func main() {
	s := server.NewServer()
	s.Run()
}
