package main

import (
    "os"
    "strings"
    "bufio"
    "regexp"
)
/**************************************************
*           待扫描文件队列                                
**************************************************/
type Node struct {
    value string
    next *Node
}

type Queue struct {
    visits map[string]int
    head *Node
    tail *Node
}

func (q *Queue) push(t string) {
    _, ok:=q.visits[t]
    if ok {
        return
    }

    p:=&Node{value: t, next: nil}
    if q.head==nil {
        q.head=p
        q.tail=p
    } else {
        q.tail.next=p
        q.tail=p
    }
}

func (q *Queue) pop() string {
    result:=q.head.value
    q.head=q.head.next
    return result
}

func (q *Queue) empty() bool {
    return q.head==q.tail
}

/**************************************************
*          扫描实现                                    
**************************************************/
type CodeReference struct {
    scanningQueue Queue
}

func (cr *CodeReference) Init(fname string) {
    cr.scanningQueue.visits=make(map[string]int)
    cr.scanningQueue.head=nil
    cr.scanningQueue.tail=nil
    cr.scanningQueue.push(fname)
}

func readIncludes(fname string) ([]string, error) {
    f, err:=os.Open(fname)
    if err!=nil {
        return nil, err
    }

    defer f.Close()

    reg:=regexp.MustCompile(`#include\s+[<\"].+[>\"]`)
    var lines []string
    scanner:=bufio.NewScanner(f)
    for scanner.Scan() {
        line:=scanner.Text()
        line=reg.FindString(line)
        if line!="" {
            lines=append(lines, line)
        }
    }

    if scanner.Err()!=nil {
        return nil, scanner.Err()
    }

    return lines, nil
}

func (cr *CodeReference) Walk() {

     for !cr.scanningQueue.empty() {
        fname:=cr.scanningQueue.pop()
        readIncludes(fname)
     }
}

/**************************************************
*         存储工程文件中包含路径和文件的结构体                       
**************************************************/
type LookupTable struct {
	Paths []string
	Files map[string]string

    Scanner CodeReference
}


var gLookupTable LookupTable

func (g LookupTable) iContains(name string) bool {
    for k, _:=range gLookupTable.Files {
        if strings.ToUpper(name)==strings.ToUpper(k) {
            return true
        }
    }
    return false
}

func (g LookupTable) Contains(name string, ignorecase bool) bool {
    var ok bool
    if ignorecase {
        ok=g.iContains(name)
    } else {
        _, ok=g.Files[name]
    }
    return ok
}
