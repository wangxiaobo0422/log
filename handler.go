package log

/*
Handler表示将日志写入到某种io/设备
*/

type Handler interface {
	Write(p []byte) (n int, err error)
	Close() error
	AsyncWrite(instance *LogInstance)
}
