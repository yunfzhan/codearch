package main

import (
    //"io"
    "os"
    //"fmt"
)

type IContent interface {
    AddString(line string)
    Add(c IContent)
    Read() []byte
}

type DotGraphContent struct {
    buff []byte
}

func (d *DotGraphContent) AddString(line string) {
    d.buff=append(d.buff, []byte(line+"\n")...)
}

func (d *DotGraphContent) Add(c IContent) {
    d.buff=append(d.buff, c.Read()...)
}

func (d *DotGraphContent) Read() []byte {
    return d.buff
}

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
