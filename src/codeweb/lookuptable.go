package main

import (
    "os"
    "fmt"
    "bytes"
    "strings"
    "bufio"
    "regexp"
)
/**************************************************
*           待扫描文件队列                                
**************************************************/
type Node struct {
    value string
    parent string
    next *Node
}

type Queue struct {
    visits map[string]int
    head *Node
    tail *Node
}

func (q *Queue) push(t string, parent string) {
    _, ok:=q.visits[t]
    if ok {
        return
    }

    p:=&Node{value: t, parent: parent, next: nil}
    if q.head==nil {
        q.head=p
        q.tail=p
    } else {
        q.tail.next=p
        q.tail=p
    }
}

func (q *Queue) pop() (string,string) {
    v:=q.head.value
    p:=q.head.parent
    q.head=q.head.next
    return v, p
}

func (q *Queue) empty() bool {
    return q.head==nil
}

/**************************************************
*          扫描实现                                    
**************************************************/
type CodeReference struct {
    scanningQueue Queue
    Nodes []string
}

func (cr *CodeReference) Init(fname string) {
    cr.scanningQueue.visits=make(map[string]int)
    cr.scanningQueue.head=nil
    cr.scanningQueue.tail=nil
    cr.scanningQueue.push(fname, "")
}

func readIncludes(fname string) ([]string, error) {
    f, err:=os.Open(fname)
    if err!=nil {
        return nil, err
    }

    defer f.Close()

    regInclude:=regexp.MustCompile(`#include\s+[<\"].+[>\"]`)
    regHeader:=regexp.MustCompile(`[<\"].+[>\"]`)
    var lines []string
    scanner:=bufio.NewScanner(f)
    for scanner.Scan() {
        line:=scanner.Text()
        line=regInclude.FindString(line)
        if line!="" {
            header:=regHeader.FindString(line)
            if header!="" {
                lines=append(lines, header[1:len(header)-1])
            }
        }
    }

    if scanner.Err()!=nil {
        return nil, scanner.Err()
    }

    return lines, nil
}

func (cr *CodeReference) createGraphNode(fname string, parent string) {
    var buff bytes.Buffer
    // Node attribute defined here
    // Format: fname [attr=...]
    buff.WriteString(fname)
    buff.WriteString(" -> ")
    buff.WriteString(parent)
    cr.Nodes=append(cr.Nodes, buff.String())
}

func (cr *CodeReference) Walk() {
    for !cr.scanningQueue.empty() {
        fname, parent:=cr.scanningQueue.pop()
        cr.createGraphNode(fname, parent)
        lines, err:=readIncludes(fname)
        if err!=nil {
            break
        }
        for i:=0; i<len(lines); i++ {
            fmt.Printf("Push %s %s\n", lines[i], fname)
            cr.scanningQueue.push(lines[i],fname)
        }
    }

    fmt.Println(cr.Nodes)
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
