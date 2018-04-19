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

type UpTimeStruct struct {
	Connect24Hours  Cards `json:"connect24hours"`
	Puback24Hours   Cards `json:"puback24hours"`
	Msgsent24Hours  Cards `json:"msgsent24hours"`
	Connectlastweek Cards `json:"connectlastweek"`
	Pubacklastweek  Cards `json:"pubacklastweek"`
	Msgsentlastweek Cards `json:"msgsentlastweek"`
}

type Cards struct {
	P99th      int64    `json:"P99th"`
	Uptime     float32  `json:"uptime"`
	Failure    int64    `json:"failure"`
	Count      int64    `json:"total"`
	OutageTime []string `json:"outagetime"`
}

func (c *Cards) setPer(value int64) {
	c.P99th = value
}

func (c *Cards) setUptime(value float32) {
	c.Uptime = value
}

func (c *Cards) setFailure(value int64) {
	c.Failure = value
}

func (c *Cards) setTotal(value int64) {
	c.Count = value
}

func (c *Cards) getPer() int64 {
	return c.P99th
}

func (c *Cards) getUptime() float32 {
	return c.Uptime
}

func (c *Cards) getFailure() int64 {
	return c.Failure
}

func (c *Cards) getTotal() int64 {
	return c.Count
}

func (c *Cards) setOutage(value []string) {
	c.OutageTime = value
}

func (c *Cards) getOutage() []string {
	return c.OutageTime
}

type ResponseDist struct {
	Connectlastweek []RangeS `json:"connectlastweek"`
	Pubacklastweek  []RangeS `json:"pubacklastweek"`
	Messagelastweek []RangeS `json:"messagelastweek"`
}

type RangeS struct {
	Range string `json:"range"`
	Count int64  `json:"count"`
}
