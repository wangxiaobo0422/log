package log

import "io"

type StreamHandler struct {
	w           io.Writer
	writeThread IHandleIOWriteThread
	fmt         Formatter
}

func newStreamHandler(w io.Writer) (*StreamHandler, error) {
	h := &StreamHandler{
		w:           w,
		writeThread: nil,
		fmt:         _defaultFormatter,
	}

	return h, nil
}

func (s *StreamHandler) Clone() Handler {
	c := new(StreamHandler)
	c.w = s.w
	c.writeThread = s.writeThread
	c.fmt = s.fmt

	return c
}

func (s *StreamHandler) AsyncWrite(l *LogInstance) {
	if s.writeThread != nil {
		s.writeThread.AsyncWrite(s, s.fmt, l)
	} else {
		_globalWriteThread.AsyncWrite(s, s.fmt, l)
	}
}

func (s *StreamHandler) Write(b []byte) (int, error) {
	return s.w.Write(b)
}

func (s *StreamHandler) Close() error {
	if s.writeThread != nil {
		s.writeThread.Close()
	}
	if wc, ok := s.w.(io.WriteCloser); ok {
		return wc.Close()
	}
	return nil
}
