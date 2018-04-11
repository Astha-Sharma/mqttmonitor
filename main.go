package main

import (
	server "github.com/vivekvasvani/mqttmonitor/server"
)

func main() {
	wait := make(chan struct{})
	server.NewServer()
	<-wait
}
