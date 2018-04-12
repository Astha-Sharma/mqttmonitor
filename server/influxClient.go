package server

import (
	"encoding/json"
	"fmt"
	"github.com/influxdata/influxdb/client/v2"
	"github.com/valyala/fasthttp"
	_ "github.com/vjeantet/jodaTime"
	"log"
	"strings"
	"sync"
	"time"
)

const (
	CONNLATENCY = "ConnectLatency"
	MyDB        = "mqttmonitoring"
	username    = ""
	password    = ""
)

var (
	instance client.Client
	once     sync.Once
)

func GetInflxInstance() client.Client {
	once.Do(func() {
		c, err := client.NewHTTPClient(client.HTTPConfig{
			Addr: "http://got.hike.in:8086",
		})
		if err != nil {
			log.Fatalf("Error creating InfluxDB Client: ", err.Error())
		}
		instance = c
	})
	return instance
}

func GetConnectLatency(ctx *fasthttp.RequestCtx) {
	var (
		ArrayC  []ResponseO
		ArrayP  []ResponseO
		ArrayM  []ResponseO
		MinMaxC Stats
		MinMaxP Stats
		MinMaxM Stats
	)

	connectL := client.NewQuery("SELECT * FROM ConnectLatency WHERE time > now() - 1h ORDER BY time ASC", MyDB, "ns")
	if response, err := GetInflxInstance().Query(connectL); err == nil && response.Error() == nil {
		for i := range response.Results[0].Series[0].Values {
			t, _ := response.Results[0].Series[0].Values[i][0].(json.Number).Int64()
			latency, _ := response.Results[0].Series[0].Values[i][1].(json.Number).Int64()
			server := response.Results[0].Series[0].Values[i][2].(string)
			//temp := time.Unix(0, t).Add((time.Hour*time.Duration(5) + time.Minute*time.Duration(30)))
			//fmt.Println("Temp-->", temp)
			Res := ResponseO{strings.Split(time.Unix(0, t).String(), ".")[0], latency, server}
			ArrayC = append(ArrayC, Res)
		}
	} else {
		fmt.Println(err, response)
	}

	connectS := client.NewQuery("SELECT MIN(latency), MAX(latency) FROM ConnectLatency WHERE time > now() - 1h", MyDB, "ns")
	if response, err := GetInflxInstance().Query(connectS); err == nil && response.Error() == nil {
		for i := range response.Results[0].Series[0].Values {
			min, _ := response.Results[0].Series[0].Values[i][1].(json.Number).Int64()
			max, _ := response.Results[0].Series[0].Values[i][2].(json.Number).Int64()
			MinMaxC = Stats{min, max}
		}
	} else {
		fmt.Println(err, response)
	}

	connectP := client.NewQuery("SELECT * FROM PubAckLatency WHERE time > now() - 1h ORDER BY time ASC", MyDB, "ns")
	if response, err := GetInflxInstance().Query(connectP); err == nil && response.Error() == nil {
		for i := range response.Results[0].Series[0].Values {
			t, _ := response.Results[0].Series[0].Values[i][0].(json.Number).Int64()
			latency, _ := response.Results[0].Series[0].Values[i][1].(json.Number).Int64()
			server := response.Results[0].Series[0].Values[i][2].(string)
			Res := ResponseO{strings.Split(time.Unix(0, t).String(), ".")[0], latency, server}
			ArrayP = append(ArrayP, Res)
		}
	} else {
		fmt.Println(err, response)
	}

	pubackS := client.NewQuery("SELECT MIN(latency), MAX(latency) FROM PubAckLatency WHERE time > now() - 1h", MyDB, "ns")
	if response, err := GetInflxInstance().Query(pubackS); err == nil && response.Error() == nil {
		for i := range response.Results[0].Series[0].Values {
			min, _ := response.Results[0].Series[0].Values[i][1].(json.Number).Int64()
			max, _ := response.Results[0].Series[0].Values[i][2].(json.Number).Int64()
			MinMaxP = Stats{min, max}
		}
	} else {
		fmt.Println(err, response)
	}

	msgSent := client.NewQuery("SELECT * FROM MessageSentLatency WHERE time > now() - 1h ORDER BY time ASC", MyDB, "ns")
	if response, err := GetInflxInstance().Query(msgSent); err == nil && response.Error() == nil {
		for i := range response.Results[0].Series[0].Values {
			t, _ := response.Results[0].Series[0].Values[i][0].(json.Number).Int64()
			latency, _ := response.Results[0].Series[0].Values[i][1].(json.Number).Int64()
			server := response.Results[0].Series[0].Values[i][2].(string)
			Res := ResponseO{strings.Split(time.Unix(0, t).String(), ".")[0], latency, server}
			ArrayM = append(ArrayM, Res)
		}
	} else {
		fmt.Println(err, response)
	}

	msgSentS := client.NewQuery("SELECT MIN(latency), MAX(latency) FROM MessageSentLatency WHERE time > now() - 1h", MyDB, "ns")
	if response, err := GetInflxInstance().Query(msgSentS); err == nil && response.Error() == nil {
		for i := range response.Results[0].Series[0].Values {
			min, _ := response.Results[0].Series[0].Values[i][1].(json.Number).Int64()
			max, _ := response.Results[0].Series[0].Values[i][2].(json.Number).Int64()
			MinMaxM = Stats{min, max}
		}
	} else {
		fmt.Println(err, response)
	}

	ResponseAB := ResponseA{ArrayC, MinMaxC, ArrayP, MinMaxP, ArrayM, MinMaxM}
	ctx.Response.Header.Set("Content-Type", "application/json")
	ctx.Response.SetBodyString(ToJsonString(ResponseAB))
}

func GetUpDownTime(ctx *fasthttp.RequestCtx) {

}

func ToJsonString(p interface{}) string {
	bytes, err := json.Marshal(p)
	if err != nil {
		fmt.Println("Error", err.Error())
	}
	return string(bytes)
}
