package log

import "os"

/*
Handler表示将日志写入到某种io/设备
*/

type Handler interface {
	Write(p []byte) (n int, err error)
	Close() error
	AsyncWrite(instance *LogInstance)
}

func newStdHandler() *StreamHandler {
	h, _ := newStreamHandler(os.Stdout)
	return h
}
