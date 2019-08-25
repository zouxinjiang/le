package clog

type (
	LogLevel   uint64
	DataFormat string
)

const (
	Lvl_Info LogLevel = 1 << iota
	Lvl_Warning
	Lvl_Error
	Lvl_Debug
	FMT_Json  DataFormat = "json"
	FMT_Plain DataFormat = "plain"
)

var str2lvl = map[LogLevel]string{
	Lvl_Debug:   "DEBUG",
	Lvl_Error:   "ERROR",
	Lvl_Info:    "INFO",
	Lvl_Warning: "WARNING",
}

func (ll LogLevel) String() string {
	return str2lvl[ll]
}

type FormatFunc func(level LogLevel, skip int) string
