package main

import (
    "os"
    "runtime"
    //"fmt"
    "strings"
    "bufio"
    "regexp"
    "sync"
    "path/filepath"
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
    head *Node
    tail *Node
}

func (q *Queue) push(t string, parent string) {
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
    visits map[string]string
}

func (cr *CodeReference) Init(fname string) {
    cr.visits=make(map[string]string)
    cr.scanningQueue.head=nil
    cr.scanningQueue.tail=nil
    if !filepath.IsAbs(fname) {
        fname, _=filepath.Abs(fname)
    }
    cr.scanningQueue.push(fname, "")
}

/**************************************************
*         检索文件中的头文件                       
**************************************************/
func readIncludes(fname string) ([]string, error) {
    if strings.HasPrefix(fname, "$$") {
        return nil, nil
    }
    f, err:=os.Open(fname)
    if err!=nil {
        return nil, err
    }

    defer f.Close()

    regInclude:=regexp.MustCompile(`#include\s+[<\"].+[>\"]`)
    regHeader:=regexp.MustCompile(`[<\"].+?[>\"]`)
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


/**************************************************
*         存储工程文件中包含路径和文件的结构体                       
**************************************************/
type LookupTable struct {
    ignoreCase bool
	Paths []string
	Files map[string]string

    Scanner CodeReference
}


var gLookupTable LookupTable

func (g *LookupTable) iContains(name string) bool {
    for k, _:=range gLookupTable.Files {
        if strings.ToUpper(name)==strings.ToUpper(k) {
            return true
        }
    }
    return false
}

func (g *LookupTable) Contains(name string, ignorecase bool) bool {
    var ok bool
    if ignorecase {
        ok=g.iContains(name)
    } else {
        _, ok=g.Files[name]
    }
    return ok
}

/**************************************************
*         在头文件目录中查找文件                       
**************************************************/
func (g *LookupTable) searchInDirectories(fname string) string {
    var result string=""
    var wg sync.WaitGroup
    wg.Add(len(g.Paths))

    for i:=0; i<len(g.Paths); i++ {
        go func(idx int) {
            defer wg.Done()
            filepath.Walk(g.Paths[idx], func(path string, info os.FileInfo, err error) error {
                _, f:=filepath.Split(path)
                if !g.ignoreCase && f==fname {
                    result=path
                } else if g.ignoreCase && strings.EqualFold(f, fname) {
                    result=path
                }
                return nil
            })
        }(i)
    }
    wg.Wait()
    return result
}

func (g *LookupTable) Walk(_callback func(fname string, parent string)) {
    g.ignoreCase=runtime.GOOS=="windows"
    for !g.Scanner.scanningQueue.empty() {
        fname, parent:=g.Scanner.scanningQueue.pop()
        _, ok:=g.Scanner.visits[fname]
        if ok {
            continue
        }
        _callback(fname, parent)
        lines, err:=readIncludes(fname)
        if err!=nil {
            break
        }
        g.Scanner.visits[fname]=""
        for i:=0; lines!=nil && i<len(lines); i++ {
            line:=g.searchInDirectories(lines[i])
            if line=="" {
                //没有找到文件，需要特殊标记
                g.Scanner.scanningQueue.push("$$"+lines[i], fname)
            } else {
                g.Scanner.scanningQueue.push(line,fname)
            }
        }
    }
}
