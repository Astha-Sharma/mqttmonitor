package influx

import (
	"fmt"
	"github.com/influxdata/influxdb/client/v2"
	_ "github.com/nlopes/slack"
	fast "github.com/vivekvasvani/mqttmonitor/client"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	MyDB          = "mqttmonitoring"
	INFLUX_URL    = "http://got.hike.in:8086"
	SLACK_WEBHOOK = ""
)

var (
	MqttServer string = "mqtt.im.hike.in"
	instance   client.Client
	once       sync.Once
	header     = make(map[string]string)
)

func GetInflxInstance() client.Client {
	once.Do(func() {
		c, err := client.NewHTTPClient(client.HTTPConfig{
			Addr: INFLUX_URL,
		})
		if err != nil {
			fmt.Println("Error creating InfluxDB Client: ", err.Error())
		}
		instance = c
	})
	return instance
}

func PushData(db string, latency int64) {
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database: MyDB,
	})
	if err != nil {
		fmt.Println(err)
	}

	// Create a point and add to batch
	fields := map[string]interface{}{
		"latency": latency,
		"server":  MqttServer,
	}
	loc, _ := time.LoadLocation("Asia/Kolkata")
	now := time.Now().In(loc)
	pt, err := client.NewPoint(db, nil, fields, now)
	if err != nil {
		fmt.Println(err)
	}
	bp.AddPoint(pt)

	// Write the batch
	if err := GetInflxInstance().Write(bp); err != nil {
		fmt.Println(err)
	}

	if latency >= 5000 {
		PostInSlack(db, latency)
	}

}

func PostInSlack(db string, latency int64) {
	var alertType string
	switch db {
	case "ConnectLatency":
		alertType = "Connect Latency is HIGH"
	case "PubAckLatency":
		alertType = "PubAck Latency is HIGH"
	case "MessageSentLatency":
		alertType = "MessageSent Latency is HIGH"
	}

	options := make([]string, 2)
	options[0] = alertType
	options[1] = strconv.FormatInt(latency, 10) + " ms"
	payload := SubstParams(options, GetPayload("alert.json"))
	fast.HitRequest(SLACK_WEBHOOK, "POST", header, payload)
}

func GetPayload(payloadPath string) string {
	if payloadPath != "" {
		dir, _ := os.Getwd()
		templateData, _ := ioutil.ReadFile(dir + "/payloads/" + payloadPath)
		return string(templateData)
	} else {
		return ""
	}
}

func SubstParams(sessionMap []string, textData string) string {
	for i, value := range sessionMap {
		if strings.ContainsAny(textData, "${"+strconv.Itoa(i)) {
			textData = strings.Replace(textData, "${"+strconv.Itoa(i)+"}", value, -1)
		}
	}
	return textData
}
