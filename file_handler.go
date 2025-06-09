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
