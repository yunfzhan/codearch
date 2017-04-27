package main

import (
    "os"
    //"fmt"
    "bytes"
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
    Nodes []string
}

func (cr *CodeReference) Init(fname string) {
    cr.visits=make(map[string]string)
    cr.scanningQueue.head=nil
    cr.scanningQueue.tail=nil
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

func getPureFileName(fname string) string {
    fname=strings.Replace(fname, ".", "_", -1)
    fname=strings.Replace(fname, "\\", "_", -1)
    fname=strings.Replace(fname, "/", "_", -1)
    return fname
}
/**************************************************
*         生成图的结点                       
* 输出结点前加上node_，因为dot文件的结点名字首字母
* 不能是非字母
**************************************************/
func (cr *CodeReference) createGraphNode(fname string, parent string) {
    var buff bytes.Buffer
    var fnameNoExt string
    if strings.HasPrefix(fname, "$$") {
        fnameNoExt=getPureFileName(fname)
        fnameNoExt=fnameNoExt[2:]
        buff.WriteString("node_"+fnameNoExt+" [color=\"red\", label=\""+fname[2:]+"\"]")
    } else {
        _, fnameNoPath:=filepath.Split(fname)
        fnameNoExt=getPureFileName(fnameNoPath)
        buff.WriteString("node_"+fnameNoExt+" [label=\""+fnameNoPath+"\"]")
    }
    cr.Nodes=append(cr.Nodes, buff.String())
    if parent=="" {
        return
    }
    buff.Reset()
    _, pnameNoPath:=filepath.Split(parent)
    pnameNoExt:=getPureFileName(pnameNoPath)
    buff.WriteString("node_"+pnameNoExt+" [label=\""+pnameNoPath+"\"]")
    buff.Reset()
    buff.WriteString("node_"+fnameNoExt)
    buff.WriteString(" -> ")
    buff.WriteString("node_"+pnameNoExt)
    cr.Nodes=append(cr.Nodes, buff.String())
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
                if f==fname {
                    result=path
                }
                return nil
            })
        }(i)
    }
    wg.Wait()
    return result
}

func (g *LookupTable) Walk() {
    for !g.Scanner.scanningQueue.empty() {
        fname, parent:=g.Scanner.scanningQueue.pop()
        _, ok:=g.Scanner.visits[fname]
        if ok {
            continue
        }
        g.Scanner.createGraphNode(fname, parent)
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
