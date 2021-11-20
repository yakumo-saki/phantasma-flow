package objects

type LogData struct {
	Logs []LogDatum
}

type LogDatum struct {
	DateTime string `comment:"ISO8601"`
	Message  string
}
