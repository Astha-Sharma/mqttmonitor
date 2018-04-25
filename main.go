package main

import (
	"github.com/vivekvasvani/mqttmonitor/config"
	server "github.com/vivekvasvani/mqttmonitor/server"
	"github.com/vivekvasvani/mqttmonitor/sql"
)

func main() {
	wait := make(chan struct{})
	config.InitializeConfig("config/config.yml")
	mysql.InitMysql()
	server.NewServer()
	<-wait
}
