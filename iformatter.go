package log

import "bytes"

type Flag int

const (
	FlagTime  = 1
	FlagFile  = 2
	FlagLevel = 4

	FlagStd    = FlagTime | FlagFile | FlagLevel
	FieldSplit = "-"
)

var (
	TimeFormat = "2006/01/02 15:04:05"
)

type LogInstance struct {
	Flag  int    `json:"-"`
	Level string `json:"level"`
	File  string `json:"file"`
	Time  string `json:"time"`
	Msg   string `json:"msg"`
	KV    Fields `json:"-"`
}

type Formatter interface {
	Format(buf *bytes.Buffer, l *LogInstance) (*bytes.Buffer, error)
}
