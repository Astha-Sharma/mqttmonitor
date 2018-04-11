package server

type StatusMessage struct {
	StatusCode string `json:"statusCode"`
	StatusType string `json:"statusType"`
	Message    string `json:"message"`
}

//Response structure
type Response struct {
	Status StatusMessage `json:"status"`
	Data   interface{}   `json:"data,omitempty"`
}

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

type ResponseO struct {
	Time    string `json:"time,omitempty"`
	Latency int64  `json:"latency,omitempty"`
	Server  string `json:"server,omitempty"`
}

type Stats struct {
	Min int64 `json:"min,omitempty"`
	Max int64 `json:"max,omitempty"`
}

type ResponseA struct {
	ArrayC []ResponseO `json:"connect,omitempty"`
	StatsC Stats       `json:"connectstats,omitempty"`
	ArrayP []ResponseO `json:"puback,omitempty"`
	StatsP Stats       `json:"pubackstats,omitempty"`
	ArrayM []ResponseO `json:"messagesent,omitempty"`
	StatsM Stats       `json:"messagesentstats,omitempty"`
}
