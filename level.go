package log

type Level int

const (
	LevelTrace = 0
	LevelDebug = 1
	LevelInfo  = 2
	LevelWarn  = 3
	LevelError = 4
	LevelFatal = 5
	LevelBuss  = 6
)

var LevelString = map[string]Level{
	"trace": LevelTrace,
	"debug": LevelDebug,
	"info":  LevelInfo,
	"warn":  LevelWarn,
	"error": LevelError,
	"fatal": LevelFatal,
	"buss":  LevelBuss,
}

var LevelNames = [7]string{
	"TRACE", "DEBUG", "INFO", "WARN", "ERROR", "FATAL", "BUSS",
}
