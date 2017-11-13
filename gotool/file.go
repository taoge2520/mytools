package gotool

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
)

func FileLine() (string, int) {
	_, file, line, ok := runtime.Caller(1) //这里的1 为一次func封装，多一次就要加1
	if !ok {
		file = "???"
		line = 0
	}
	short := file
	for i := len(file) - 1; i > 0; i-- {
		if file[i] == '/' {
			short = file[i+1:]
			break
		}
	}
	file = short
	return file, line
}

func CreatePidFile(path string) error {
	dir, filename := filepath.Split(path)
	if dir != "" {
		if err := os.MkdirAll(dir, 0666); err != nil {
			return err
		}
	}
	if filename == "" {
		return errors.New("missing filename")
	}
	data := []byte(strconv.Itoa(os.Getpid()))
	return ioutil.WriteFile(path, data, 666)
}

func OpenFile(name string, flag int, perm os.FileMode) (*os.File, error) {
	dir, _ := filepath.Split(name)
	if dir != "" {
		if err := os.MkdirAll(dir, perm); err != nil {
			return nil, err
		}
	}
	return os.OpenFile(name, flag, perm)
}
