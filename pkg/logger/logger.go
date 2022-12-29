package logger

import (
	"io"
	"log"
	"os"
)

func New(filePath string) (*log.Logger, error) {
	f, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}
	wrt := io.MultiWriter(os.Stdout, f)
	l := log.New(wrt, "govisa: ", log.LstdFlags)
	l.SetOutput(wrt)
	return l, nil
}
