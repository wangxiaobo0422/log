package log

/* 整个log包，默认只配置一个全局IHandleIOWriteThread
   如果log对象设置了多个handler，则每个handler最好都
   设置一个对应的IHandleIOWriteThread。
*/

type DropLogCallbackFunc func(l *LogInstance, drop int)

// IHandleIOWriteThread
type IHandleIOWriteThread interface {
	AsyncWrite(h Handler, fmt Formatter, log *LogInstance)
	SetDropCallback(f DropLogCallbackFunc)
	SetLimitCallback(f DropLogCallbackFunc)
	SetLimiter(token float64, burst int)
	Close()
}
