package main

import (
	"fmt"
	"os"
	"sync"
)

type Logger struct {
	file  *os.File
	mutex sync.Mutex
}

func NewLogger(file *os.File) *Logger {
	return &Logger{
		file:  file,
		mutex: sync.Mutex{},
	}
}

func (l *Logger) Print(str string) (err error) {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	_, err = fmt.Fprint(l.file, str)
	return
}

func (l *Logger) Println(str string) (err error) {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	_, err = fmt.Fprintln(l.file, str)
	return
}
