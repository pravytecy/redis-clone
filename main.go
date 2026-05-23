package main

import (
	"flag"
	"fmt"

	"github.com/pravytecy/redis-clone/config"
	"github.com/pravytecy/redis-clone/server"
)

func main() {
	fmt.Println("Redis clone starting...")
	setUpFlags()
	fmt.Println("intializing jerry ...")
	server.RunTcpSyncServer()
}

func setUpFlags() {
	flag.StringVar(&config.Host, "host", "0.0.0.0", "host for the jerry server")
	flag.IntVar(&config.Port, "port", 7379, "port for the server")
	flag.Parse()
}
