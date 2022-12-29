package logger

import (
	"fmt"
	"io"
	"log"
	"os"
)

func New(dir string) (*log.Logger, error) {
	var l *log.Logger
	f, err := os.OpenFile(fmt.Sprintf("%s/error.log", dir), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	wrt := io.MultiWriter(os.Stdout, f)
	l.SetOutput(wrt)
	return l, nil
}
