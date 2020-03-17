package main

import (
	"./gatewayserver"
)

func main() {
	var s gatewayserver.Server
	s.Init("config.yml")
	s.RunServer()
}
