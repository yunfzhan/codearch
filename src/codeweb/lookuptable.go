package main

import (
    "strings"
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

func (cr *CodeReference) Walk() {

     for !cr.scanningQueue.empty() {

     }
}

/**************************************************
*                                                  
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
