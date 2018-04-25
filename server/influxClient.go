package server

import (
	"encoding/json"
	"fmt"
	"github.com/influxdata/influxdb/client/v2"
	"github.com/valyala/fasthttp"
	_ "github.com/vjeantet/jodaTime"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	CONNLATENCY = "ConnectLatency"
	MyDB        = "mqttmonitoring"
	username    = ""
	password    = ""
	ONE         = "0-50"
	TWO         = "51-100"
	THREE       = "101-500"
	FOUR        = "501-1000"
	FIVE        = "1001-5000"
	SIX         = "5000+"
	TOTAL       = "TOTAL"
)

var (
	instance       client.Client
	once           sync.Once
	timeMapConnect = make(map[string]int64)
	timeMapPubAck  = make(map[string]int64)
	timeMapMsgSent = make(map[string]int64)
)

func GetInflxInstance() client.Client {
	once.Do(func() {
		c, err := client.NewHTTPClient(client.HTTPConfig{
			Addr: "http://35.198.242.78:8086",
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
		ArrayC   []ResponseO
		ArrayP   []ResponseO
		ArrayM   []ResponseO
		MinMaxC  Stats
		MinMaxP  Stats
		MinMaxM  Stats
		duration string
	)

	duration = string(ctx.QueryArgs().Peek("duration")[:])
	connectL := client.NewQuery("SELECT * FROM ConnectLatency WHERE time > now() - "+duration+" ORDER BY time ASC", MyDB, "ns")
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

	connectS := client.NewQuery("SELECT MIN(latency), MAX(latency) FROM ConnectLatency WHERE time > now() - "+duration, MyDB, "ns")
	if response, err := GetInflxInstance().Query(connectS); err == nil && response.Error() == nil {
		for i := range response.Results[0].Series[0].Values {
			min, _ := response.Results[0].Series[0].Values[i][1].(json.Number).Int64()
			max, _ := response.Results[0].Series[0].Values[i][2].(json.Number).Int64()
			MinMaxC = Stats{min, max}
		}
	} else {
		fmt.Println(err, response)
	}

	connectP := client.NewQuery("SELECT * FROM PubAckLatency WHERE time > now() - "+duration+" ORDER BY time ASC", MyDB, "ns")
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

	pubackS := client.NewQuery("SELECT MIN(latency), MAX(latency) FROM PubAckLatency WHERE time > now() - "+duration, MyDB, "ns")
	if response, err := GetInflxInstance().Query(pubackS); err == nil && response.Error() == nil {
		for i := range response.Results[0].Series[0].Values {
			min, _ := response.Results[0].Series[0].Values[i][1].(json.Number).Int64()
			max, _ := response.Results[0].Series[0].Values[i][2].(json.Number).Int64()
			MinMaxP = Stats{min, max}
		}
	} else {
		fmt.Println(err, response)
	}

	msgSent := client.NewQuery("SELECT * FROM MessageSentLatency WHERE time > now() - "+duration+" ORDER BY time ASC", MyDB, "ns")
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

	msgSentS := client.NewQuery("SELECT MIN(latency), MAX(latency) FROM MessageSentLatency WHERE time > now() - "+duration, MyDB, "ns")
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
	var (
		CardsConnect24        Cards
		CardsConnectLastWeek  Cards
		CardsPubAck24         Cards
		CardsPubAckLastWeek   Cards
		CardsMsg24            Cards
		CardsMsgLastWeek      Cards
		OutageConnect24       = []string{}
		OutageConnectLastWeek = []string{}
		OutagePubAck24        = []string{}
		OutagePubAckLastWeek  = []string{}
		OutageMsg24           = []string{}
		OutageMsgLastWeek     = []string{}
	)

	connectC := client.NewQuery("SELECT PERCENTILE(latency, 99), COUNT(*) FROM ConnectLatency WHERE time > now() - 24h", MyDB, "ns")
	if response, err := GetInflxInstance().Query(connectC); err == nil && response.Error() == nil {
		percentile, _ := response.Results[0].Series[0].Values[0][1].(json.Number).Int64()
		count, _ := response.Results[0].Series[0].Values[0][2].(json.Number).Int64()
		CardsConnect24.setPer(percentile)
		CardsConnect24.setTotal(count)
	} else {
		fmt.Println(err, response)
	}

	connectP := client.NewQuery("SELECT PERCENTILE(latency, 99), COUNT(*) FROM PubAckLatency WHERE time > now() - 24h", MyDB, "ns")
	if response, err := GetInflxInstance().Query(connectP); err == nil && response.Error() == nil {
		percentile, _ := response.Results[0].Series[0].Values[0][1].(json.Number).Int64()
		count, _ := response.Results[0].Series[0].Values[0][2].(json.Number).Int64()
		CardsPubAck24.setPer(percentile)
		CardsPubAck24.setTotal(count)
	} else {
		fmt.Println(err, response)
	}

	connectM := client.NewQuery("SELECT PERCENTILE(latency, 99), COUNT(*) FROM MessageSentLatency WHERE time > now() - 24h", MyDB, "ns")
	if response, err := GetInflxInstance().Query(connectM); err == nil && response.Error() == nil {
		percentile, _ := response.Results[0].Series[0].Values[0][1].(json.Number).Int64()
		count, _ := response.Results[0].Series[0].Values[0][2].(json.Number).Int64()
		CardsMsg24.setPer(percentile)
		CardsMsg24.setTotal(count)
	} else {
		fmt.Println(err, response)
	}

	/*
		FconnectC := client.NewQuery("SELECT COUNT(*) FROM ConnectLatency WHERE (latency > 5000 OR latency < 0) AND time > now() - 24h", MyDB, "ns")
		if response, err := GetInflxInstance().Query(FconnectC); err == nil && response.Error() == nil {
			fmt.Println(response)
			if len(response.Results[0].Series) == 0 {
				CardsConnect24.setFailure(0)
			} else {
				failures, _ := response.Results[0].Series[0].Values[0][1].(json.Number).Int64()
				CardsConnect24.setFailure(failures)
			}
		} else {
			fmt.Println(err, response)
		}
	*/

	OutageFconnectC := client.NewQuery("SELECT * FROM ConnectLatency WHERE (latency > 5000 OR latency < 0) AND time > now() - 24h", MyDB, "ns")
	if response, err := GetInflxInstance().Query(OutageFconnectC); err == nil && response.Error() == nil {
		if len(response.Results[0].Series) != 0 {
			for i := range response.Results[0].Series[0].Values {
				t, _ := response.Results[0].Series[0].Values[i][0].(json.Number).Int64()
				latency, _ := response.Results[0].Series[0].Values[i][1].(json.Number).Int64()
				OutageConnect24 = append(OutageConnect24, strings.Split(time.Unix(0, t).String(), ".")[0]+" : "+strconv.FormatInt(latency, 10)+" ms")
			}
		}
		if len(response.Results[0].Series) == 0 {
			CardsConnect24.setFailure(0)
		} else {
			CardsConnect24.setFailure(int64(len(response.Results[0].Series[0].Values)))
		}
	} else {
		fmt.Println(err, response)
	}

	/*
		FconnectP := client.NewQuery("SELECT COUNT(*) FROM PubAckLatency WHERE (latency > 5000 OR latency < 0) AND time > now() - 24h", MyDB, "ns")
		if response, err := GetInflxInstance().Query(FconnectP); err == nil && response.Error() == nil {
			if len(response.Results[0].Series) == 0 {
				CardsPubAck24.setFailure(0)
			} else {
				failures, _ := response.Results[0].Series[0].Values[0][1].(json.Number).Int64()
				CardsPubAck24.setFailure(failures)
			}
		} else {
			fmt.Println(err, response)
		}
	*/

	OutagePuback := client.NewQuery("SELECT * FROM PubAckLatency WHERE (latency > 5000 OR latency < 0) AND time > now() - 24h", MyDB, "ns")
	if response, err := GetInflxInstance().Query(OutagePuback); err == nil && response.Error() == nil {
		if len(response.Results[0].Series) != 0 {
			for i := range response.Results[0].Series[0].Values {
				t, _ := response.Results[0].Series[0].Values[i][0].(json.Number).Int64()
				latency, _ := response.Results[0].Series[0].Values[i][1].(json.Number).Int64()
				OutagePubAck24 = append(OutagePubAck24, strings.Split(time.Unix(0, t).String(), ".")[0]+" : "+strconv.FormatInt(latency, 10)+" ms")
			}
		}
		if len(response.Results[0].Series) == 0 {
			CardsPubAck24.setFailure(0)
		} else {
			CardsPubAck24.setFailure(int64(len(response.Results[0].Series[0].Values)))
		}
	} else {
		fmt.Println(err, response)
	}

	/*
		FconnectM := client.NewQuery("SELECT COUNT(*) FROM MessageSentLatency WHERE (latency > 5000 OR latency < 0) AND time > now() - 24h", MyDB, "ns")
		if response, err := GetInflxInstance().Query(FconnectM); err == nil && response.Error() == nil {
			if len(response.Results[0].Series) == 0 {
				CardsMsg24.setFailure(0)
			} else {
				failures, _ := response.Results[0].Series[0].Values[0][1].(json.Number).Int64()
				CardsMsg24.setFailure(failures)
			}
		} else {
			fmt.Println(err, response)
		}
	*/

	OutageMessageSent := client.NewQuery("SELECT * FROM MessageSentLatency WHERE (latency > 5000 OR latency < 0) AND time > now() - 24h", MyDB, "ns")
	if response, err := GetInflxInstance().Query(OutageMessageSent); err == nil && response.Error() == nil {
		for i := range response.Results[0].Series[0].Values {
			if len(response.Results[0].Series) != 0 {
				t, _ := response.Results[0].Series[0].Values[i][0].(json.Number).Int64()
				latency, _ := response.Results[0].Series[0].Values[i][1].(json.Number).Int64()
				OutageMsg24 = append(OutageMsg24, strings.Split(time.Unix(0, t).String(), ".")[0]+" : "+strconv.FormatInt(latency, 10)+" ms")
			}
		}
		if len(response.Results[0].Series) == 0 {
			CardsMsg24.setFailure(0)
		} else {
			CardsMsg24.setFailure(int64(len(response.Results[0].Series[0].Values)))
		}

	} else {
		fmt.Println(err, response)
	}

	LWconnectC := client.NewQuery("SELECT PERCENTILE(latency, 99), COUNT(*) FROM ConnectLatency WHERE time > now() - 7d", MyDB, "ns")
	if response, err := GetInflxInstance().Query(LWconnectC); err == nil && response.Error() == nil {
		percentile, _ := response.Results[0].Series[0].Values[0][1].(json.Number).Int64()
		count, _ := response.Results[0].Series[0].Values[0][2].(json.Number).Int64()
		CardsConnectLastWeek.setPer(percentile)
		CardsConnectLastWeek.setTotal(count)
	} else {
		fmt.Println(err, response)
	}

	LWconnectP := client.NewQuery("SELECT PERCENTILE(latency, 99), COUNT(*) FROM PubAckLatency WHERE time > now() - 7d", MyDB, "ns")
	if response, err := GetInflxInstance().Query(LWconnectP); err == nil && response.Error() == nil {
		percentile, _ := response.Results[0].Series[0].Values[0][1].(json.Number).Int64()
		count, _ := response.Results[0].Series[0].Values[0][2].(json.Number).Int64()
		CardsPubAckLastWeek.setPer(percentile)
		CardsPubAckLastWeek.setTotal(count)
	} else {
		fmt.Println(err, response)
	}

	LWconnectM := client.NewQuery("SELECT PERCENTILE(latency, 99), COUNT(*) FROM MessageSentLatency WHERE time > now() - 7d", MyDB, "ns")
	if response, err := GetInflxInstance().Query(LWconnectM); err == nil && response.Error() == nil {
		percentile, _ := response.Results[0].Series[0].Values[0][1].(json.Number).Int64()
		count, _ := response.Results[0].Series[0].Values[0][2].(json.Number).Int64()
		CardsMsgLastWeek.setPer(percentile)
		CardsMsgLastWeek.setTotal(count)
	} else {
		fmt.Println(err, response)
	}

	/*
		LWFconnectC := client.NewQuery("SELECT COUNT(*) FROM ConnectLatency WHERE (latency > 5000 OR latency < 0) AND time > now() - 7d", MyDB, "ns")
		if response, err := GetInflxInstance().Query(LWFconnectC); err == nil && response.Error() == nil {
			fmt.Println(response)
			if len(response.Results[0].Series) == 0 {
				CardsConnectLastWeek.setFailure(0)
			} else {
				failures, _ := response.Results[0].Series[0].Values[0][1].(json.Number).Int64()
				CardsConnectLastWeek.setFailure(failures)
			}
		} else {
			fmt.Println(err, response)
		}
	*/

	OutageConnectLasWweek := client.NewQuery("SELECT * FROM ConnectLatency WHERE (latency > 5000 OR latency < 0) AND time > now() - 7d", MyDB, "ns")
	if response, err := GetInflxInstance().Query(OutageConnectLasWweek); err == nil && response.Error() == nil {
		if len(response.Results[0].Series) != 0 {
			for i := range response.Results[0].Series[0].Values {
				t, _ := response.Results[0].Series[0].Values[i][0].(json.Number).Int64()
				latency, _ := response.Results[0].Series[0].Values[i][1].(json.Number).Int64()
				OutageConnectLastWeek = append(OutageConnectLastWeek, strings.Split(time.Unix(0, t).String(), ".")[0]+" : "+strconv.FormatInt(latency, 10)+" ms")
			}
		}
		if len(response.Results[0].Series) == 0 {
			CardsConnectLastWeek.setFailure(0)
		} else {
			CardsConnectLastWeek.setFailure(int64(len(response.Results[0].Series[0].Values)))
		}

	} else {
		fmt.Println(err, response)
	}

	/*
		LWFconnectP := client.NewQuery("SELECT COUNT(*) FROM PubAckLatency WHERE (latency > 5000 OR latency < 0) AND time > now() - 7d", MyDB, "ns")
		if response, err := GetInflxInstance().Query(LWFconnectP); err == nil && response.Error() == nil {
			if len(response.Results[0].Series) == 0 {
				CardsPubAckLastWeek.setFailure(0)
			} else {
				failures, _ := response.Results[0].Series[0].Values[0][1].(json.Number).Int64()
				CardsPubAckLastWeek.setFailure(failures)
			}
		} else {
			fmt.Println(err, response)
		}
	*/

	OutagePubAckLasWweek := client.NewQuery("SELECT * FROM PubAckLatency WHERE (latency > 5000 OR latency < 0) AND time > now() - 7d", MyDB, "ns")
	if response, err := GetInflxInstance().Query(OutagePubAckLasWweek); err == nil && response.Error() == nil {
		if len(response.Results[0].Series) != 0 {
			for i := range response.Results[0].Series[0].Values {
				t, _ := response.Results[0].Series[0].Values[i][0].(json.Number).Int64()
				latency, _ := response.Results[0].Series[0].Values[i][1].(json.Number).Int64()
				OutagePubAckLastWeek = append(OutagePubAckLastWeek, strings.Split(time.Unix(0, t).String(), ".")[0]+" : "+strconv.FormatInt(latency, 10)+" ms")
			}
		}
		if len(response.Results[0].Series) == 0 {
			CardsPubAckLastWeek.setFailure(0)
		} else {
			CardsPubAckLastWeek.setFailure(int64(len(response.Results[0].Series[0].Values)))
		}
	} else {
		fmt.Println(err, response)
	}

	/*
		LWFconnectM := client.NewQuery("SELECT COUNT(*) FROM MessageSentLatency WHERE (latency > 5000 OR latency < 0) AND time > now() - 7d", MyDB, "ns")
		if response, err := GetInflxInstance().Query(LWFconnectM); err == nil && response.Error() == nil {
			if len(response.Results[0].Series) == 0 {
				CardsMsgLastWeek.setFailure(0)
			} else {
				failures, _ := response.Results[0].Series[0].Values[0][1].(json.Number).Int64()
				CardsMsgLastWeek.setFailure(failures)
			}
		} else {
			fmt.Println(err, response)
		}
	*/

	Outage24Hours := client.NewQuery("SELECT * FROM MessageSentLatency WHERE (latency > 5000 OR latency < 0) AND time > now() - 7d", MyDB, "ns")
	if response, err := GetInflxInstance().Query(Outage24Hours); err == nil && response.Error() == nil {
		if len(response.Results[0].Series) != 0 {
			for i := range response.Results[0].Series[0].Values {
				t, _ := response.Results[0].Series[0].Values[i][0].(json.Number).Int64()
				latency, _ := response.Results[0].Series[0].Values[i][1].(json.Number).Int64()
				OutageMsgLastWeek = append(OutageMsgLastWeek, strings.Split(time.Unix(0, t).String(), ".")[0]+" : "+strconv.FormatInt(latency, 10)+" ms")
			}
		}
		if len(response.Results[0].Series) == 0 {
			CardsMsgLastWeek.setFailure(0)
		} else {
			CardsMsgLastWeek.setFailure(int64(len(response.Results[0].Series[0].Values)))
		}

	} else {
		fmt.Println(err, response)
	}

	CardsConnect24.setUptime(float32(100) - float32(CardsConnect24.getFailure())/float32(CardsConnect24.getTotal()))
	CardsConnect24.setOutage(OutageConnect24)

	CardsPubAck24.setUptime(float32(100) - float32(CardsPubAck24.getFailure())/float32(CardsPubAck24.getTotal()))
	CardsPubAck24.setOutage(OutagePubAck24)

	CardsMsg24.setUptime(float32(100) - float32(CardsMsg24.getFailure())/float32(CardsMsg24.getTotal()))
	CardsMsg24.setOutage(OutageMsg24)

	CardsConnectLastWeek.setUptime(float32(100) - float32(CardsConnectLastWeek.getFailure())/float32(CardsConnectLastWeek.getTotal()))
	CardsConnectLastWeek.setOutage(OutageConnectLastWeek)

	CardsPubAckLastWeek.setUptime(float32(100) - float32(CardsPubAckLastWeek.getFailure())/float32(CardsPubAckLastWeek.getTotal()))
	CardsPubAckLastWeek.setOutage(OutagePubAckLastWeek)

	CardsMsgLastWeek.setUptime(float32(100) - float32(CardsMsgLastWeek.getFailure())/float32(CardsMsgLastWeek.getTotal()))
	CardsMsgLastWeek.setOutage(OutageMsgLastWeek)

	responseToSend := UpTimeStruct{CardsConnect24, CardsPubAck24, CardsMsg24, CardsConnectLastWeek, CardsPubAckLastWeek, CardsMsgLastWeek}
	fmt.Println(ToJsonString(responseToSend))
	ctx.Response.Header.Set("Content-Type", "application/json")
	ctx.Response.SetBodyString(ToJsonString(responseToSend))
}

func TimeDistribution(ctx *fasthttp.RequestCtx) {
	var (
		CONNS   []RangeS
		PUBACK  []RangeS
		MSGSENT []RangeS
	)

	var (
		connect0to50      int64
		connect51to100    int64
		connect101to500   int64
		connect501to1000  int64
		connect1001to5000 int64
		connect5001plus   int64
	)

	var (
		puback0to50      int64
		puback51to100    int64
		puback101to500   int64
		puback501to1000  int64
		puback1001to5000 int64
		puback5001plus   int64
	)

	var (
		msg0to50      int64
		msg51to100    int64
		msg101to500   int64
		msg501to1000  int64
		msg1001to5000 int64
		msg5001plus   int64
	)

	Common := client.NewQuery("SELECT * FROM ConnectLatency WHERE time > now() - 7d", MyDB, "ns")
	CommonPubAck := client.NewQuery("SELECT * FROM PubAckLatency WHERE time > now() - 7d", MyDB, "ns")
	CommonMsgSent := client.NewQuery("SELECT * FROM MessageSentLatency WHERE time > now() - 7d", MyDB, "ns")

	/*
		Crange0To50 := client.NewQuery("SELECT count(*) FROM ConnectLatency WHERE time > now() - 7d AND latency > 0 AND latency <= 50", MyDB, "ns")
		Crange51To100 := client.NewQuery("SELECT count(*) FROM ConnectLatency WHERE time > now() - 7d AND latency > 51 AND latency <= 100", MyDB, "ns")
		Crange101To500 := client.NewQuery("SELECT count(*) FROM ConnectLatency WHERE time > now() - 7d AND latency > 101 AND latency <= 500", MyDB, "ns")
		Crange501To1000 := client.NewQuery("SELECT count(*) FROM ConnectLatency WHERE time > now() - 7d AND latency > 501 AND latency <= 1000", MyDB, "ns")
		Crange1001To5000 := client.NewQuery("SELECT count(*) FROM ConnectLatency WHERE time > now() - 7d AND latency > 1001 AND latency <= 5000", MyDB, "ns")
		Crange5001Plus := client.NewQuery("SELECT count(*) FROM ConnectLatency WHERE time > now() - 7d AND latency latency >= 5001", MyDB, "ns")

		Prange0To50 := client.NewQuery("SELECT count(*) FROM PubAckLatency WHERE time > now() - 7d AND latency > 0 AND latency <= 50", MyDB, "ns")
		Prange51To100 := client.NewQuery("SELECT count(*) FROM PubAckLatency WHERE time > now() - 7d AND latency > 51 AND latency <= 100", MyDB, "ns")
		Prange101To500 := client.NewQuery("SELECT count(*) FROM PubAckLatency WHERE time > now() - 7d AND latency > 101 AND latency <= 500", MyDB, "ns")
		Prange501To1000 := client.NewQuery("SELECT count(*) FROM PubAckLatency WHERE time > now() - 7d AND latency > 501 AND latency <= 1000", MyDB, "ns")
		Prange1001To5000 := client.NewQuery("SELECT count(*) FROM PubAckLatency WHERE time > now() - 7d AND latency > 1001 AND latency <= 5000", MyDB, "ns")
		Prange5001Plus := client.NewQuery("SELECT count(*) FROM PubAckLatency WHERE time > now() - 7d AND latency latency >= 5001", MyDB, "ns")

		Mrange0To50 := client.NewQuery("SELECT count(*) FROM MessageSentLatency WHERE time > now() - 7d AND latency > 0 AND latency <= 50", MyDB, "ns")
		Mrange51To100 := client.NewQuery("SELECT count(*) FROM MessageSentLatency WHERE time > now() - 7d AND latency > 51 AND latency <= 100", MyDB, "ns")
		Mrange101To500 := client.NewQuery("SELECT count(*) FROM MessageSentLatency WHERE time > now() - 7d AND latency > 101 AND latency <= 500", MyDB, "ns")
		Mrange501To1000 := client.NewQuery("SELECT count(*) FROM MessageSentLatency WHERE time > now() - 7d AND latency > 501 AND latency <= 1000", MyDB, "ns")
		Mrange1001To5000 := client.NewQuery("SELECT count(*) FROM MessageSentLatency WHERE time > now() - 7d AND latency > 1001 AND latency <= 5000", MyDB, "ns")
		Mrange5001Plus := client.NewQuery("SELECT count(*) FROM MessageSentLatency WHERE time > now() - 7d AND latency latency >= 5001", MyDB, "ns")
	*/

	/* Connect */
	if response, err := GetInflxInstance().Query(Common); err == nil && response.Error() == nil {
		for i := range response.Results[0].Series[0].Values {
			min, _ := response.Results[0].Series[0].Values[i][1].(json.Number).Int64()
			switch {
			case min > 0 && min <= 50:
				connect0to50 = connect0to50 + 1
			case min > 51 && min <= 100:
				connect51to100 = connect51to100 + 1
			case min > 101 && min <= 500:
				connect101to500 = connect101to500 + 1
			case min > 501 && min <= 1000:
				connect501to1000 = connect501to1000 + 1
			case min > 1001 && min <= 5000:
				connect1001to5000 = connect1001to5000 + 1
			case min > 5000:
				connect5001plus = connect5001plus + 1
			}
		}
		timeMapConnect["0-50"] = connect0to50
		timeMapConnect["51-100"] = connect51to100
		timeMapConnect["101-500"] = connect101to500
		timeMapConnect["501-1000"] = connect501to1000
		timeMapConnect["1001-5000"] = connect1001to5000
		timeMapConnect["5000+"] = connect5001plus

	} else {
		fmt.Println(err, response)
	}

	/* puback */
	if response, err := GetInflxInstance().Query(CommonPubAck); err == nil && response.Error() == nil {
		for i := range response.Results[0].Series[0].Values {
			min, _ := response.Results[0].Series[0].Values[i][1].(json.Number).Int64()
			switch {
			case min > 0 && min <= 50:
				puback0to50 = puback0to50 + 1
			case min > 51 && min <= 100:
				puback51to100 = puback51to100 + 1
			case min > 101 && min <= 500:
				puback101to500 = puback101to500 + 1
			case min > 501 && min <= 1000:
				puback501to1000 = puback501to1000 + 1
			case min > 1001 && min <= 5000:
				puback1001to5000 = puback1001to5000 + 1
			case min > 5000:
				puback5001plus = puback5001plus + 1
			}
		}
		timeMapPubAck["0-50"] = puback0to50
		timeMapPubAck["51-100"] = puback51to100
		timeMapPubAck["101-500"] = puback101to500
		timeMapPubAck["501-1000"] = puback501to1000
		timeMapPubAck["1001-5000"] = puback1001to5000
		timeMapPubAck["5000+"] = puback5001plus

	} else {
		fmt.Println(err, response)
	}

	/* Message Sent */
	if response, err := GetInflxInstance().Query(CommonMsgSent); err == nil && response.Error() == nil {
		for i := range response.Results[0].Series[0].Values {
			min, _ := response.Results[0].Series[0].Values[i][1].(json.Number).Int64()
			switch {
			case min > 0 && min <= 50:
				msg0to50 = msg0to50 + 1
			case min > 51 && min <= 100:
				msg51to100 = msg51to100 + 1
			case min > 101 && min <= 500:
				msg101to500 = msg101to500 + 1
			case min > 501 && min <= 1000:
				msg501to1000 = msg501to1000 + 1
			case min > 1001 && min <= 5000:
				msg1001to5000 = msg1001to5000 + 1
			case min > 5000:
				msg5001plus = msg5001plus + 1
			}
		}
		timeMapMsgSent["0-50"] = msg0to50
		timeMapMsgSent["51-100"] = msg51to100
		timeMapMsgSent["101-500"] = msg101to500
		timeMapMsgSent["501-1000"] = msg501to1000
		timeMapMsgSent["1001-5000"] = msg1001to5000
		timeMapMsgSent["5000+"] = msg5001plus

	} else {
		fmt.Println(err, response)
	}

	/* Connect
	if response, err := GetInflxInstance().Query(Crange0To50); err == nil && response.Error() == nil {
		if len(response.Results[0].Series) == 0 {
			timeMapConnect["0-50"] = 0
		} else {
			count, _ := response.Results[0].Series[0].Values[0][1].(json.Number).Int64()
			timeMapConnect["0-50"] = count
		}
	}
	if response, err := GetInflxInstance().Query(Crange51To100); err == nil && response.Error() == nil {
		if len(response.Results[0].Series) == 0 {
			timeMapConnect["51-100"] = 0
		} else {
			count, _ := response.Results[0].Series[0].Values[0][1].(json.Number).Int64()
			timeMapConnect["51-100"] = count
		}
	}
	if response, err := GetInflxInstance().Query(Crange101To500); err == nil && response.Error() == nil {
		if len(response.Results[0].Series) == 0 {
			timeMapConnect["101-500"] = 0
		} else {
			count, _ := response.Results[0].Series[0].Values[0][1].(json.Number).Int64()
			timeMapConnect["101-500"] = count
		}
	}
	if response, err := GetInflxInstance().Query(Crange501To1000); err == nil && response.Error() == nil {
		if len(response.Results[0].Series) == 0 {
			timeMapConnect["501-1000"] = 0
		} else {
			count, _ := response.Results[0].Series[0].Values[0][1].(json.Number).Int64()
			timeMapConnect["501-1000"] = count
		}
	}
	if response, err := GetInflxInstance().Query(Crange1001To5000); err == nil && response.Error() == nil {
		if len(response.Results[0].Series) == 0 {
			timeMapConnect["1001-5000"] = 0
		} else {
			count, _ := response.Results[0].Series[0].Values[0][1].(json.Number).Int64()
			timeMapConnect["1001-5000"] = count
		}
	}
	if response, err := GetInflxInstance().Query(Crange5001Plus); err == nil && response.Error() == nil {
		if len(response.Results[0].Series) == 0 {
			timeMapConnect["5000+"] = 0
		} else {
			count, _ := response.Results[0].Series[0].Values[0][1].(json.Number).Int64()
			timeMapConnect["5000+"] = count
		}
	}

	*/

	/* Pub Ack
	if response, err := GetInflxInstance().Query(Prange0To50); err == nil && response.Error() == nil {
		if len(response.Results[0].Series) == 0 {
			timeMapPubAck["0-50"] = 0
		} else {
			count, _ := response.Results[0].Series[0].Values[0][1].(json.Number).Int64()
			timeMapPubAck["0-50"] = count
		}
	}
	if response, err := GetInflxInstance().Query(Prange51To100); err == nil && response.Error() == nil {
		if len(response.Results[0].Series) == 0 {
			timeMapPubAck["51-100"] = 0
		} else {
			count, _ := response.Results[0].Series[0].Values[0][1].(json.Number).Int64()
			timeMapPubAck["51-100"] = count
		}
	}
	if response, err := GetInflxInstance().Query(Prange101To500); err == nil && response.Error() == nil {
		if len(response.Results[0].Series) == 0 {
			timeMapPubAck["101-500"] = 0
		} else {
			count, _ := response.Results[0].Series[0].Values[0][1].(json.Number).Int64()
			timeMapPubAck["101-500"] = count
		}
	}
	if response, err := GetInflxInstance().Query(Prange501To1000); err == nil && response.Error() == nil {
		if len(response.Results[0].Series) == 0 {
			timeMapPubAck["501-1000"] = 0
		} else {
			count, _ := response.Results[0].Series[0].Values[0][1].(json.Number).Int64()
			timeMapPubAck["501-1000"] = count
		}
	}
	if response, err := GetInflxInstance().Query(Prange1001To5000); err == nil && response.Error() == nil {
		if len(response.Results[0].Series) == 0 {
			timeMapPubAck["1001-5000"] = 0
		} else {
			count, _ := response.Results[0].Series[0].Values[0][1].(json.Number).Int64()
			timeMapPubAck["1001-5000"] = count
		}
	}
	if response, err := GetInflxInstance().Query(Prange5001Plus); err == nil && response.Error() == nil {
		if len(response.Results[0].Series) == 0 {
			timeMapPubAck["5000+"] = 0
		} else {
			count, _ := response.Results[0].Series[0].Values[0][1].(json.Number).Int64()
			timeMapPubAck["5000+"] = count
		}
	}

	*/

	/* MessageSent Ack
	if response, err := GetInflxInstance().Query(Mrange0To50); err == nil && response.Error() == nil {
		if len(response.Results[0].Series) == 0 {
			timeMapMsgSent["0-50"] = 0
		} else {
			count, _ := response.Results[0].Series[0].Values[0][1].(json.Number).Int64()
			timeMapMsgSent["0-50"] = count
		}
	}
	if response, err := GetInflxInstance().Query(Mrange51To100); err == nil && response.Error() == nil {
		if len(response.Results[0].Series) == 0 {
			timeMapMsgSent["51-100"] = 0
		} else {
			count, _ := response.Results[0].Series[0].Values[0][1].(json.Number).Int64()
			timeMapMsgSent["51-100"] = count
		}
	}
	if response, err := GetInflxInstance().Query(Mrange101To500); err == nil && response.Error() == nil {
		if len(response.Results[0].Series) == 0 {
			timeMapMsgSent["101-500"] = 0
		} else {
			count, _ := response.Results[0].Series[0].Values[0][1].(json.Number).Int64()
			timeMapMsgSent["101-500"] = count
		}
	}
	if response, err := GetInflxInstance().Query(Mrange501To1000); err == nil && response.Error() == nil {
		if len(response.Results[0].Series) == 0 {
			timeMapMsgSent["501-1000"] = 0
		} else {
			count, _ := response.Results[0].Series[0].Values[0][1].(json.Number).Int64()
			timeMapMsgSent["501-1000"] = count
		}
	}
	if response, err := GetInflxInstance().Query(Mrange1001To5000); err == nil && response.Error() == nil {
		if len(response.Results[0].Series) == 0 {
			timeMapMsgSent["1001-5000"] = 0
		} else {
			count, _ := response.Results[0].Series[0].Values[0][1].(json.Number).Int64()
			timeMapMsgSent["1001-5000"] = count
		}
	}
	if response, err := GetInflxInstance().Query(Mrange5001Plus); err == nil && response.Error() == nil {
		if len(response.Results[0].Series) == 0 {
			timeMapMsgSent["5000+"] = 0
		} else {
			count, _ := response.Results[0].Series[0].Values[0][1].(json.Number).Int64()
			timeMapMsgSent["5000+"] = count
		}
	}
	*/

	CONNS = append(CONNS, RangeS{ONE, timeMapConnect[ONE]})
	CONNS = append(CONNS, RangeS{TWO, timeMapConnect[TWO]})
	CONNS = append(CONNS, RangeS{THREE, timeMapConnect[THREE]})
	CONNS = append(CONNS, RangeS{FOUR, timeMapConnect[FOUR]})
	CONNS = append(CONNS, RangeS{FIVE, timeMapConnect[FIVE]})
	CONNS = append(CONNS, RangeS{SIX, timeMapConnect[SIX]})
	CONNS = append(CONNS, RangeS{TOTAL, timeMapConnect[SIX] + timeMapConnect[FIVE] + timeMapConnect[FOUR] + timeMapConnect[THREE] + timeMapConnect[TWO] + timeMapConnect[ONE]})

	PUBACK = append(PUBACK, RangeS{ONE, timeMapPubAck[ONE]})
	PUBACK = append(PUBACK, RangeS{TWO, timeMapPubAck[TWO]})
	PUBACK = append(PUBACK, RangeS{THREE, timeMapPubAck[THREE]})
	PUBACK = append(PUBACK, RangeS{FOUR, timeMapPubAck[FOUR]})
	PUBACK = append(PUBACK, RangeS{FIVE, timeMapPubAck[FIVE]})
	PUBACK = append(PUBACK, RangeS{SIX, timeMapPubAck[SIX]})
	PUBACK = append(PUBACK, RangeS{TOTAL, timeMapPubAck[SIX] + timeMapPubAck[FIVE] + timeMapPubAck[FOUR] + timeMapPubAck[THREE] + timeMapPubAck[TWO] + timeMapPubAck[ONE]})

	MSGSENT = append(MSGSENT, RangeS{ONE, timeMapMsgSent[ONE]})
	MSGSENT = append(MSGSENT, RangeS{TWO, timeMapMsgSent[TWO]})
	MSGSENT = append(MSGSENT, RangeS{THREE, timeMapMsgSent[THREE]})
	MSGSENT = append(MSGSENT, RangeS{FOUR, timeMapMsgSent[FOUR]})
	MSGSENT = append(MSGSENT, RangeS{FIVE, timeMapMsgSent[FIVE]})
	MSGSENT = append(MSGSENT, RangeS{SIX, timeMapMsgSent[SIX]})
	MSGSENT = append(MSGSENT, RangeS{TOTAL, timeMapMsgSent[SIX] + timeMapMsgSent[FIVE] + timeMapMsgSent[FOUR] + timeMapMsgSent[THREE] + timeMapMsgSent[TWO] + timeMapMsgSent[ONE]})

	responseToSend := ResponseDist{CONNS, PUBACK, MSGSENT}
	fmt.Println(ToJsonString(responseToSend))
	ctx.Response.Header.Set("Content-Type", "application/json")
	ctx.Response.SetBodyString(ToJsonString(responseToSend))

}

func ToJsonString(p interface{}) string {
	bytes, err := json.Marshal(p)
	if err != nil {
		fmt.Println("Error", err.Error())
	}
	return string(bytes)
}
