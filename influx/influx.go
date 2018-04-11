package influx

import (
	"fmt"
	"github.com/influxdata/influxdb/client/v2"
	"sync"
	"time"
)

const (
	MyDB       = "mqttmonitoring"
	INFLUX_URL = "http://got.hike.in:8086"
)

var (
	MqttServer string = "mqtt.im.hike.in"
	instance   client.Client
	once       sync.Once
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
}
