package server

import (
	"encoding/json"
	_ "encoding/json"
	"fmt"
	"github.com/buger/jsonparser"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	packets "github.com/eclipse/paho.mqtt.golang/packets"
	config "github.com/vivekvasvani/mqttmonitor/config"
	influx "github.com/vivekvasvani/mqttmonitor/influx"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"
)

const (
	charset = "abcdefghijklmnopqrstuvwxyz" +
		"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

var (
	latency       Latency
	msgLatency    MsgLatency
	msg           = make(chan packets.ControlPacket)
	connChan      = make(chan packets.ControlPacket)
	user          UserData
	counter       int64 = 1
	remainingUser []UserData
	c             mqtt.Client
	Environment   string
	mqttBroker    string
	seededRand    *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))
)

func GetMqttClient(client *mqtt.ClientOptions) mqtt.Client {
	return mqtt.NewClient(client)
}

func getClientOptions(userData UserData) *mqtt.ClientOptions {
	client := mqtt.NewClientOptions()
	client.SetCleanSession(true)
	client.AddBroker(mqttBroker)
	client.SetDefaultPublishHandler(F)
	client.SetKeepAlive(10 * time.Second)
	client.SetClientID(userData.Msisdn + ":5:true")
	client.SetPingTimeout(1 * time.Minute)
	client.SetOnConnectHandler(ConnHandler)
	client.SetUsername(userData.UID)
	client.SetPassword(userData.Token)
	client.SetAutoReconnect(false)
	client.SetConnectionLostHandler(OnLost)
	client.SetConnectTimeout(1 * time.Minute)
	return client
}

//define a function for the default message handler
var F mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	mType, _ := jsonparser.GetString(msg.Payload(), "t")
	fmt.Println("Message Type In Default Handler : ", mType, string(msg.Payload()))
	if mType == "m" {
		msgLatency.setEnd(time.Now().UnixNano())
		fmt.Println("Message Sent latency (ms) :", (msgLatency.getEnd()-msgLatency.getStart())/1000000)
		influx.PushData("MessageSentLatency", (msgLatency.getEnd()-msgLatency.getStart())/1000000)
	}

}

var OnLost mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	if err != nil {
		fmt.Println("Connection Lost :", err.Error())
	}
	//client.Disconnect(0)
}

var ConnHandler mqtt.OnConnectHandler = func(c mqtt.Client) {
	if c.IsConnected() {
		fmt.Println("Inside OnConnectHandler...")
		//latency.setEnd(time.Now().UnixNano())
		/*
			// message type of st
			msg := "{\"to\": \"" + user.Msisdn + "\", \"t\" : \"st\"}"
			PublishMessage(c, user.UID+"/p", 0, true, msg)

			// message type of bulklastseen
			msgfg := "{ \"d\" : {\"justOpened\" : true, \"bulklastseen\" : false}, \"t\" : \"app\", \"st\":\"fg\"}"
			PublishMessage(c, user.UID+"/p", 0, true, msgfg)
		*/
	}
}

func ConnectToMqtt(c mqtt.Client) {
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		fmt.Println("Error in ConnectToMqtt :", token.Error())
	} else {
		fmt.Printf("Connected to server\n")
	}

}

func SendTypeM() {
	messageId := strconv.FormatInt(counter, 10)
	message := String(20)
	timestamp := fmt.Sprintf("%v", time.Now().UnixNano()/1000000000)
	msgrandom := "{\"t\": \"m\",\"to\": \"u:" + user.UID + "\",\"d\":{\"hm\":\"" + message + "\",\"i\":\"" + messageId + "\", \"ts\":" + timestamp + "}}"
	PublishMessage(c, user.UID+"/p", 1, false, msgrandom)
	counter++
}

func PublishMessage(c mqtt.Client, topic string, qos byte, retained bool, payload interface{}) { //, msgId uint16) {
	//put current time stamp in channel
	//message type of m
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
		fmt.Println("Error in PublishMessage :", token.Error())
	} else {
		fmt.Println("Published Successfilly -->", payload.(string))
	}
}

func Init() UserData {
	var userS UserData
	if len(remainingUser) > 0 {
		userS = remainingUser[0]
		remainingUser = remainingUser[1:]
	} else {
		LoadUsers()
		userS = remainingUser[0]
		remainingUser = remainingUser[1:]
	}
	return userS
}

func LoadUsers() {
	var fileName string
	if Environment == "stag" {
		fileName = "./config/stag-users.json"
	} else {
		fileName = "./config/prod-users.json"
	}
	fmt.Println("File name :", fileName)
	raw, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	errU := json.Unmarshal(raw, &remainingUser)
	if errU != nil {
		fmt.Println(errU.Error())
	}
}

func StringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func String(length int) string {
	return StringWithCharset(length, charset)
}

func GetStatsOfMqtt(env string) {
	//mqtt.DEBUG = log.New(os.Stdout, "", 0)
	mqtt.ERROR = log.New(os.Stdout, "", 0)

	//If client is already connected disconnect the same
	if c != nil && c.IsConnected() {
		c.Disconnect(0)
	}

	//fmt.Println(c.IsConnected())
	//set the environment
	Environment = env

	//Set mqtt broker
	if Environment == "stag" {
		mqttBroker = "tcp://staging.im.hike.in:1883"
	} else {
		mqttBroker = "tcp://mqtt.im.hike.in:5222"
	}

	//Load the user and connect to mqtt broker
	func() {
		user = Init()
		opts := getClientOptions(user)
		c = GetMqttClient(opts)
		latency.setStart(time.Now().UnixNano())
		go func() { config.ChanConAck <- time.Now().UnixNano() }()
		ConnectToMqtt(c)
		SendTypeM()
	}()
}
