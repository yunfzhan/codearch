package main

import (
	"log"
	"fmt"
    //"io"
	"os"
    "bytes"
	"strings"
	"errors"
    "path/filepath"
)
/*******************************************************************************
*
*            帮助函数
*******************************************************************************/
func printHelp(){
	fmt.Println("**************************************************************")
	fmt.Println("*                   CWScan v1.0 Help                         *")
	fmt.Println("**************************************************************")
	fmt.Println("")
	fmt.Println("cwscan -c -m <VC project or makefile.>")
	fmt.Println("List files not in project and clean up.")
	fmt.Println("")
	fmt.Println("cwscan -m <VC project or makefile> -f <file name to scan>")
	fmt.Println("Create topology of included files of a file.")
}

/*******************************************************************************
*
*            定义数据结构和全局变量
*******************************************************************************/
type Option struct{
	flag bool
	param string
}

type GArguments struct{
	clean bool
	project Option
	file Option
    output Option
}

var gargs GArguments
var ig bool //ignore case search
var nodes []string
//var argMap map[string][]string

/*******************************************************************************
*
*            实现函数定义
*******************************************************************************/
func processCmdLine(arg string, flag bool) error {
	if flag {
		switch arg {
		case "c":
			gargs.clean=true
		case "f":
			gargs.file.flag=true
		case "m":
			gargs.project.flag=true
        case "o":
            gargs.output.flag=true
		default:
			return errors.New("Invalid option: "+arg)
		}
	} else {
		if gargs.file.flag && gargs.file.param=="" {
			gargs.file.param=arg
		} else if gargs.project.flag && gargs.project.param=="" {
			gargs.project.param=arg
		} else if gargs.output.flag && gargs.output.param=="" {
            gargs.output.param=arg
        } else {
			return errors.New("Unknown arguments: "+arg)
		}
	}
	return nil
}

func parseMakefile(make string) {
	if strings.HasSuffix(make, ".vcxproj") || strings.HasSuffix(make, ".vcproj") {
		// VC project
		if	err:=buildVCProject(make); err!=nil {
			log.Fatal(err)
		}
        ig=true
	} else {
        ig=false
    }
}
/*
    搜索工程文件中不存在的源代码和头文件
*/
func cleanFn(path string, info os.FileInfo, err error) error {
    filename:=strings.ToLower(info.Name())
    if info.Mode().IsRegular() && (strings.HasSuffix(filename, ".cpp") || strings.HasSuffix(filename, ".h") || strings.HasSuffix(filename, ".cxx") || strings.HasSuffix(filename, ".hpp")) && !gLookupTable.Contains(info.Name(), ig/*在Windows系统中文件名是不区分大小写，而Linux却不是。所以需要这个参数*/) {
        fmt.Printf("Not found %s\n", info.Name())
    }
    return nil
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
func  createGraphNode(fname string, parent string) {
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
    nodes=append(nodes, buff.String())
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
    nodes=append(nodes, buff.String())
}

func createGraph(o *FileWriter) {
    o.Open()
    fmt.Printf("Generated %s\n",o.Name)
    if o.Err!=nil {
        log.Fatal(o.Err)
    }
    defer o.Close()
    o.Write([]byte("digraph cpp_graph {\n"))
    o.Write([]byte("\tnode [color=\"blue\"]\n"))
    for i:=0; i<len(nodes); i++ {
        o.Write([]byte(nodes[i]+"\n"))
    }
    o.Write([]byte("}\n"))
}

func main(){
	argnum := len(os.Args)
	if argnum==1 || argnum>5 {
		//No arguments
		printHelp()
		return
	}

	//Process command line
	for i:=1; i<argnum; i++ {
		arg:=os.Args[i]
		flag:=strings.HasPrefix(arg, "-")
		if flag {
			arg=arg[1:]
		}
		err:=processCmdLine(arg,flag)
		if err!=nil {
			// meet errors
			log.Fatal(err)
		}
	}

    if !gargs.project.flag && gargs.project.param=="" {
		log.Fatal("Need a project name.")
	}

    var o *FileWriter
    if gargs.output.flag && gargs.output.param!="" {
        //指定输出文件，目前可以是输出dot文件或者sqlite数据库
        o=&FileWriter{gargs.output.param, nil, nil}
    }
	parseMakefile(gargs.project.param)

	if gargs.clean {
		//列出不在工程文件中的文件
        root, _ := filepath.Split(gargs.project.param)
        filepath.Walk(root, cleanFn)
	} else if gargs.file.flag && gargs.file.param!=""{
        //递归扫描指定文件中的头文件包括关系
        gLookupTable.Scanner.Init(gargs.file.param)
        gLookupTable.Walk(createGraphNode)
        fmt.Println(gargs.file.param)
        if !gargs.output.flag {
            o=&FileWriter{gargs.file.param+".dot", nil, nil}
        }
        createGraph(o)
	} else {
		log.Fatal("Missing file name to scan.")
	}
}
