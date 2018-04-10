package mqtt

import (
	"encoding/json"
	_ "encoding/json"
	"fmt"
	"github.com/buger/jsonparser"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	packets "github.com/eclipse/paho.mqtt.golang/packets"
	config "github.com/vivekvasvani/mqttmonitor/config"
	server "github.com/vivekvasvani/mqttmonitor/server"
	"io/ioutil"
	"os"
	"strconv"
	"time"
)

type UserData struct {
	Msisdn string `json:"msisdn"`
	UID    string `json:"uid"`
	Token  string `json:"token"`
}

type Latency struct {
	Start int64
	End   int64
}

func (l *Latency) setStart(st int64) {
	l.Start = st
}

func (l *Latency) getStart() int64 {
	return l.Start
}

func (l *Latency) setEnd(end int64) {
	l.End = end
}

func (l *Latency) getEnd() int64 {
	return l.End
}

type MsgLatency struct {
	Start int64
	End   int64
}

func (l *MsgLatency) setStart(st int64) {
	l.Start = st
}

func (l *MsgLatency) getStart() int64 {
	return l.Start
}

func (l *MsgLatency) setEnd(end int64) {
	l.End = end
}

func (l *MsgLatency) getEnd() int64 {
	return l.End
}

var (
	ConnectedClients []mqtt.Client
	latency          Latency
	msgLatency       MsgLatency
	msg              = make(chan packets.ControlPacket)
	user             UserData
	counter          int = 1
	remainingUser    []UserData
)

func GetMqttClient(client *mqtt.ClientOptions) mqtt.Client {
	c := mqtt.NewClient(client)
	return c
}

func getClientOptions(userData UserData) *mqtt.ClientOptions {
	client := mqtt.NewClientOptions()
	client.SetCleanSession(true)
	client.AddBroker("tcp://staging.im.hike.in:1883")
	client.SetDefaultPublishHandler(F)
	client.SetKeepAlive(10 * time.Second)
	client.SetClientID(userData.Msisdn + ":5:true")
	client.SetPingTimeout(1 * time.Minute)
	client.SetOnConnectHandler(ConnHandler)
	client.SetUsername(userData.UID)
	client.SetPassword(userData.Token)
	client.SetAutoReconnect(true)
	client.SetConnectionLostHandler(OnLost)
	client.SetConnectTimeout(1 * time.Minute)
	return client
}

//define a function for the default message handler
var F mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	mType, _ := jsonparser.GetString(msg.Payload(), "t")
	fmt.Println("Message Type In Default Handler : ", mType)
	if mType == "m" {
		msgLatency.setEnd(time.Now().UnixNano())
		fmt.Println("Message Sent latency (ms) :", (msgLatency.getEnd()-msgLatency.getStart())/1000000)
	}

}

var OnLost mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	if err != nil {
		fmt.Println("Connection Lost")
	}
}

var ConnHandler mqtt.OnConnectHandler = func(c mqtt.Client) {
	if c.IsConnected() {
		latency.setEnd(time.Now().UnixNano())
		fmt.Println("Connect latency (ms) :", (latency.getEnd()-latency.getStart())/1000000)

		// message type of st
		msg := "{\"to\": \"" + user.Msisdn + "\", \"t\" : \"st\"}"
		PublishMessage(c, user.UID+"/p", 0, true, msg)

		// message type of bulklastseen
		msgfg := "{ \"d\" : {\"justOpened\" : true, \"bulklastseen\" : false}, \"t\" : \"app\", \"st\":\"fg\"}"
		PublishMessage(c, user.UID+"/p", 0, true, msgfg)

		//message type of m
		messageId := strconv.Itoa(counter)
		message := "random message :" + messageId
		timestamp := fmt.Sprintf("%v", time.Now().UnixNano()/1000000000)
		msgrandom := "{\"t\": \"m\",\"to\": \"u:" + user.UID + "\",\"d\":{\"hm\":\"" + message + "\",\"i\":\"" + messageId + "\", \"ts\":" + timestamp + "}}"
		PublishMessage(c, user.UID+"/p", 1, true, msgrandom)
		counter++
	}

	c.Disconnect(10)
}

func ConnectToMqtt(c mqtt.Client) {
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		fmt.Println("Error in ConnectToMqtt :", token.Error())
	} else {
		fmt.Printf("Connected to server\n")
	}

}

func PublishMessage(c mqtt.Client, topic string, qos byte, retained bool, payload interface{}) { //, msgId uint16) {
	//put current time stamp in channel
	go func() {
		typeM, err := jsonparser.GetString([]byte(payload.(string)), "t")
		if err != nil {
			fmt.Println("Error in parsing publish msg : ", err.Error())
		} else {
			if typeM == "m" {
				currentT := time.Now().UnixNano()
				msgLatency.setStart(currentT)
				config.ChanM <- currentT
			}
		}
	}()

	if token := c.Publish(topic, qos, retained, payload); token.Wait() && token.Error() != nil {
		fmt.Println("Error in ConnectToMqtt :", token.Error())
	} else {
		fmt.Println("Published Successfilly -->", payload.(string))
	}

}

func Init() UserData {
	fmt.Println("Inside Init", len(remainingUser))
	var userS UserData
	if len(remainingUser) > 0 {
		userS = remainingUser[0]
		remainingUser = remainingUser[1:]
	} else {
		fmt.Println("Inside Load")
		LoadUsers()
		userS = remainingUser[0]
		remainingUser = remainingUser[1:]
	}
	return userS
}

func LoadUsers() {
	raw, err := ioutil.ReadFile("./config/users.json")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	errU := json.Unmarshal(raw, &remainingUser)
	if errU != nil {
		fmt.Println(errU.Error())
	}
}

func GetStatsOfMqtt() {
	//mqtt.DEBUG = log.New(os.Stdout, "", 0)
	//mqtt.ERROR = log.New(os.Stdout, "", 0)
	//server.NewServer()
	//wait := make(chan struct{})
	user = Init()
	var opts = getClientOptions(user)
	var c = GetMqttClient(opts)
	latency.setStart(time.Now().UnixNano())
	ConnectToMqtt(c)
	//<-wait

}
