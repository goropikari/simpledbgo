package main

import "github.com/goropikari/simpledbgo/server"

func main() {
	cfg := server.NewConfig()
	s := server.NewServer(cfg)
	s.Run()
}
