package log

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
)

type when int8

const (
	WhenSecond when = iota
	WhenMinute      = 1
	WhenHour        = 2
	WhenDay         = 3
)

type FileHandler struct {
	*StreamHandler
	fd       *os.File
	fileName string
}

func NewFileHandler(fileName string) (*FileHandler, error) {
	dir := path.Dir(fileName)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err = os.Mkdir(dir, 0777); err != nil {
			if !os.IsExist(err) {
				return nil, err
			}
		}
	}

	f, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	sh, _ := newStreamHandler(f)

	return &FileHandler{
		StreamHandler: sh,
		fd:            f,
		fileName:      fileName,
	}, nil
}

func (f *FileHandler) Clone() Handler {
	c := new(FileHandler)
	c.StreamHandler = f.StreamHandler.Clone().(*StreamHandler)
	c.fd = f.fd
	c.fileName = f.fileName

	return c
}

func (f *FileHandler) Close() error {
	if f.fd != nil {
		f.fd.Close()
	}

	return nil
}

type TimeRotatingFileHandler struct {
	*FileHandler

	logName    string
	logDir     string
	interval   int64
	suffix     string
	rolloverAt int64
	keepLogNum int
}

func NewTimeRotatingFileHandler(fileName string, w when, interval int, keepLogNum int) (*TimeRotatingFileHandler, error) {
	fh, err := NewFileHandler(fileName)
	if err != nil {
		return nil, err
	}

	if keepLogNum == 0 {
		keepLogNum = 5
	}

	h := new(TimeRotatingFileHandler)
	h.FileHandler = fh
	h.logName = filepath.Base(fileName)
	h.logDir = filepath.Dir(fileName)
	h.keepLogNum = keepLogNum

	switch w {
	case WhenSecond:
		h.interval = 1
		h.suffix = "2016-01-02_15-04-05"
	case WhenMinute:
		h.interval = 60
		h.suffix = "2016-01-02_15-04"
	case WhenHour:
		h.interval = 60 * 60
		h.suffix = "2016-01-02_15"
	case WhenDay:
		h.interval = 24 * 3600
		h.suffix = "2016-01-02"
	default:
		return nil, fmt.Errorf("invalid when_rotate: %d", w)
	}

	h.interval = h.interval * int64(interval)
	finfo, _ := h.fd.Stat()
	// 每隔一段时间就写到新文件
	h.rolloverAt = finfo.ModTime().Unix() + h.interval

	return h, nil
}
