package log

import (
	"bytes"
	"log"
	"os"
	"sync"
	"sync/atomic"
)

/* 整个log包，默认只配置一个全局IHandleIOWriteThread
   如果log对象设置了多个handler，则每个handler最好都
   设置一个对应的IHandleIOWriteThread。
*/

var stdErrLog = log.New(os.Stderr, "[go/log]", log.Ldate|log.Ltime|log.Lshortfile)

const (
	_4k = 4096
	_8k = 8192
)

type DropLogCallbackFunc func(l *LogInstance, drop int)

// IHandleIOWriteThread
type IHandleIOWriteThread interface {
	AsyncWrite(h Handler, fmt Formatter, log *LogInstance)
	Close()
}

type hdlWriter struct {
	handler   Handler
	formatter Formatter
	log       *LogInstance
}

type HandlerIOWriteThread struct {
	name                string
	close               bool
	quit                chan bool
	handlerWriterChan   chan *hdlWriter
	handlerWriterBuffer *sync.Pool

	writeBuffer *bytes.Buffer

	dropCnt     int64
	writeCnt    int64
	asyncSumCnt int64
	errPrintCnt int64
	limitCnt    int64

	wg                   sync.WaitGroup
	dropLogCallbackFunc  DropLogCallbackFunc
	limitLogCallbackFunc DropLogCallbackFunc
}

func NewHandlerIOWriteThread(name string, cap int) *HandlerIOWriteThread {
	h := new(HandlerIOWriteThread)

	h.name = name
	h.quit = make(chan bool, 1)
	h.handlerWriterChan = make(chan *hdlWriter, cap)

	h.handlerWriterBuffer = &sync.Pool{
		New: func() interface{} {
			return &hdlWriter{}
		},
	}

	h.writeBuffer = bytes.NewBuffer(make([]byte, 0, _8k))
	h.wg.Add(1)

	go h.run()

	return h
}

//go:norace
func (h *HandlerIOWriteThread) AsyncWrite(hl Handler, fmt Formatter, log *LogInstance) {
	if hl == nil || fmt == nil {
		return
	}

	atomic.AddInt64(&h.asyncSumCnt, 1)

	hw := h.handlerWriterBuffer.Get().(*hdlWriter)
	hw.handler = hl
	hw.formatter = fmt
	hw.log = log

	select {
	case h.handlerWriterChan <- hw:
		return
	default:
		atomic.AddInt64(&h.dropCnt, 1)
	}
}

func (h *HandlerIOWriteThread) doFormat(hdw *hdlWriter, buf *bytes.Buffer) {
	defer func() {
		globalLogInstallBuffer.Put(hdw.log)
		h.handlerWriterBuffer.Put(hdw)
	}()

	if hdw.formatter == nil {
		hdw.formatter = _defaultFormatter
	}

	if _, err := hdw.formatter.Format(buf, hdw.log); err != nil {
		if hdw.formatter != _defaultFormatter {
			hdw.formatter.Format(buf, hdw.log)
		} else {
			return
		}
	}
	h.writeCnt += 1
}

func (h *HandlerIOWriteThread) doWrite(hdw *hdlWriter) {
	if hdw == nil {
		return
	}

	pBuf := h.writeBuffer
	var next *hdlWriter

	for hdw != nil {
		// 写buffer
		h.doFormat(hdw, pBuf)
		if pBuf.Len() >= _4k {
			// 大于阈值直接写
			hdw.handler.Write(pBuf.Bytes())
			pBuf.Reset()
		}
		// 没有关闭并且还有日志
		if !h.close && len(h.handlerWriterChan) > 0 {
			next = <-h.handlerWriterChan
			// 下个handler不是同一个并且还没有写
			if next.handler != hdw.handler && pBuf.Len() > 0 {
				// 现在这个handler直接写日志
				hdw.handler.Write(pBuf.Bytes())
				// 重置buffer
				pBuf.Reset()
			}
			// 同一个handler,继续写buffer
			hdw = next
		} else {
			break
		}
	}

	if pBuf.Len() > 0 {
		_, err := hdw.handler.Write(pBuf.Bytes())
		if err != nil {
			stdErrLog.Printf("Log[%s] write[%s] error[%v]\n",
				h.name, pBuf.String(), err)
		}
		pBuf.Reset()
	}
}

func (h *HandlerIOWriteThread) run() {
	defer h.wg.Done()

	var stop = false
	var hdw *hdlWriter

	for {
		select {
		case <-h.quit:
			// 设置状态位
			stop = true
		case hdw = <-h.handlerWriterChan:
			// 写日志
			h.doWrite(hdw)
		}

		if !stop {
			continue
		}

		// 已停止并且已全部处理完日志，直接返回
		if stop && len(h.handlerWriterChan) == 0 {
			return
		}
	}
}

func (h *HandlerIOWriteThread) Close() {
	if h.close {
		return
	}
	h.close = true
	select {
	case h.quit <- true:
		// 等待线程处理完剩余的日志
		h.wg.Wait()
	default:
		return
	}
}
