package main

import (
    //"io"
    "os"
    //"errors"
)

type FileWriter struct {
    Name string
    Err error
    fp *os.File
}

func (f *FileWriter) Open() {
    f.fp, f.Err = os.Create(f.Name)
}

func (f *FileWriter) Write(p []byte) (n int, err error) {
    return f.fp.Write(p)
}

func (f *FileWriter) Close(){
    f.fp.Close()
}
