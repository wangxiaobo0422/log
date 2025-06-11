package log

import (
	"os"
	"path"
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

const (
	WhenSecond = 0
	WhenMinute = 1
	WhenHour   = 2
	WhenDay    = 3
)
