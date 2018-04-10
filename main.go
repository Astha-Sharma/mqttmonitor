package main
import
(
	server "github.com/vivekvasvani/mqttmonitor/server"
	)

)



func main() {
	//mqtt.DEBUG = log.New(os.Stdout, "", 0)
	//mqtt.ERROR = log.New(os.Stdout, "", 0)
	wait := make(chan struct{})
	server.NewServer()
	<-wait

}
