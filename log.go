package log

import (
	"fmt"
	"runtime"
	"strings"
	"sync"
	"time"
)

var (
	MaxBytesPerLog         = 1024 * 3
	globalLogInstallBuffer *sync.Pool

	_std = NewLogger(nil, FlagStd)
)

type Fields map[string]interface{}

type Logger struct {
	level                Level
	flag                 Flag
	handlers             []Handler
	kv                   Fields
	enableEscapeNewLines bool // 是否允许转义的换行符，默认不允许，会将转义的换行符和回车符替换成||
}

func init() {
	globalLogInstallBuffer = &sync.Pool{
		New: func() interface{} {
			return &LogInstance{}
		},
	}
}

// NewLogger 新建logger实例
func NewLogger(h Handler, flag Flag) *Logger {
	l := new(Logger)
	l.level = LevelInfo
	l.handlers = []Handler{h}
	l.flag = flag
	l.kv = make(Fields, 5)

	return l
}

func (l *Logger) SetLevel(level Level) {
	l.level = level
}

func (l *Logger) SetLevelS(levelStr string) {
	l.SetLevel(LevelString[strings.ToLower(levelStr)])
}

func (l *Logger) Output(callDepth int, level Level, format string, v ...interface{}) {
	if l.level > level {
		return
	}

	msg := fmt.Sprintf(format, v...)
	l.OutputMsg(callDepth, level, msg)
}

func (l *Logger) OutputMsg(depth int, level Level, msg string) {
	var fileLine, currentTime, slevel string

	if l.flag&FlagTime > 0 {
		currentTime = time.Now().Format(TimeFormat)
	}
	if l.flag&FlagLevel > 0 {
		slevel = LevelNames[int(level)]
	}
	if l.flag&FlagFile > 0 {
		_, file, line, ok := runtime.Caller(depth)
		if !ok {
			file = "???"
			line = 0
		} else {
			v := strings.Split(file, "/")
			idx := len(v) - 3
			if idx < 0 {
				idx = 0
			}
			file = strings.Join(v[idx:], "/")
		}
		fileLine = fmt.Sprintf("%s:[%d]", file, line)
	}

	if len(msg) > MaxBytesPerLog {
		msg = fmt.Sprintf("%s... data too long, soucre-length=%d",
			msg[0:MaxBytesPerLog], len(msg))
	}

	if !l.enableEscapeNewLines {
		msg = strings.Replace(msg, "\r", "||", -1)
		msg = strings.Replace(msg, "\n", "||", -1)
	}

	log := globalLogInstallBuffer.Get().(*LogInstance)
	log.Flag = int(l.flag)
	log.File = fileLine
	log.Level = slevel
	log.KV = l.kv
	log.Time = currentTime
	log.Msg = msg

	for _, h := range l.handlers {
		if h != nil {
			h.AsyncWrite(log)
		}
	}
}

func (l *Logger) Trace(format string, v ...interface{}) {
	l.Output(3, LevelTrace, format, v...)
}

func (l *Logger) Debug(format string, v ...interface{}) {
	l.Output(3, LevelDebug, format, v...)
}

func (l *Logger) Info(format string, v ...interface{}) {
	l.Output(3, LevelInfo, format, v...)
}

func (l *Logger) Warn(format string, v ...interface{}) {
	l.Output(3, LevelWarn, format, v...)
}

func (l *Logger) Error(format string, v ...interface{}) {
	l.Output(3, LevelError, format, v...)
}

func (l *Logger) Fatal(format string, v ...interface{}) {
	l.Output(3, LevelFatal, format, v...)
}

func (l *Logger) Buss(format string, v ...interface{}) {
	l.Output(3, LevelBuss, format, v...)
}

func Trace(format string, v ...interface{}) {
	_std.Output(3, LevelTrace, format, v...)
}

func Debug(format string, v ...interface{}) {
	_std.Output(3, LevelDebug, format, v...)
}

func Info(format string, v ...interface{}) {
	_std.Output(3, LevelInfo, format, v...)
}

func Warn(format string, v ...interface{}) {
	_std.Output(3, LevelWarn, format, v...)
}

func Error(format string, v ...interface{}) {
	_std.Output(3, LevelError, format, v...)
}

func Fatal(format string, v ...interface{}) {
	_std.Output(3, LevelFatal, format, v...)
}

func Buss(format string, v ...interface{}) {
	_std.Output(3, LevelBuss, format, v...)
}
