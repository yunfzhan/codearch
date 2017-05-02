package main

import (
    //"io"
    "os"
    //"errors"
)

type IOutputWriter interface {
    Name() string
    Open() error
    Write(p []byte) (n int, err error)
    WriteString(s string) (n int, err error)
    Close()
}

type FileWriter struct {
    name string
    fp *os.File
}

func (f *FileWriter) Name() string {
    return f.name
}

func (f *FileWriter) Open() error {
    var err error
    f.fp, err = os.Create(f.name)
    return err
}

func (f *FileWriter) Write(p []byte) (n int, err error) {
    return f.fp.Write(p)
}

func (f *FileWriter) WriteString(s string) (n int, err error) {
    return f.fp.WriteString(s)
}

func (f *FileWriter) Close(){
    f.fp.Close()
}
