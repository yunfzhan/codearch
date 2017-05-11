package main

import (
	"log"
	"fmt"
    //"io"
	"os"
    "sync"
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
    fmt.Println("")
    fmt.Println("cwscan -m <VC project or makefile> -p")
    fmt.Println("Recursively scan all project files.")
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
    recursive bool //扫描所有工程中的文件调用关系
	project Option
	file Option
    output Option
}

var gargs GArguments
var ig bool //ignore case search
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
        case "p":
            gargs.recursive=true
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
		if err:=buildVCProject(make); err!=nil {
			log.Fatal(err)
		}
        ig=true
	} else if strings.HasSuffix(make, ".pro") {
        // QT project

    } else {
        // Makefile
        if err:=buildMakefile(make); err!=nil {
            log.Fatal(err)
        }
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
func  createGraphNode(fname string, parent string) IContent {
    r:=new(DotGraphContent)
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
    //nodes=append(nodes, buff.String())
    r.AddString(buff.String())
    if parent=="" {
        return r
    }
    buff.Reset()
    _, pnameNoPath:=filepath.Split(parent)
    pnameNoExt:=getPureFileName(pnameNoPath)
    buff.WriteString("node_"+pnameNoExt+" [label=\""+pnameNoPath+"\"]")
    buff.Reset()
    buff.WriteString("node_"+fnameNoExt)
    buff.WriteString(" -> ")
    buff.WriteString("node_"+pnameNoExt)
    //nodes=append(nodes, buff.String())
    r.AddString(buff.String())
    return r
}

func createGraph(o IOutputWriter, c IContent) {
    err:=o.Open()
    fmt.Printf("Generated %s\n",o.Name())
    if err!=nil {
        log.Fatal(err)
    }
    defer o.Close()
    o.WriteString("digraph cpp_graph {\n")
    o.WriteString("\tnode [color=\"blue\"]\n")
    o.Write(c.Read())
    //for i:=0; i<len(nodes); i++ {
    //    o.WriteString(nodes[i]+"\n")
    //}
    o.WriteString("}\n")
}

func scanSingleFile(file string, o IOutputWriter, wg *sync.WaitGroup) {
    if wg!=nil {
        defer wg.Done()
    }
    //递归扫描指定文件中的头文件包括关系
    scanner:=new(CodeReference)
    scanner.Init(file)
    c:=new(DotGraphContent)
    gLookupTable.Walk(scanner, c, createGraphNode)
    createGraph(o, c)
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

    gLookupTable.cwd, _=os.Getwd()

    if !gargs.project.flag && gargs.project.param=="" {
		log.Fatal("Need a project name.")
	}

    if gargs.output.flag && gargs.output.param!="" {
        //指定输出文件，目前可以是输出dot文件或者sqlite数据库
        log.Fatal("Not specified any output file name.")
    }
	parseMakefile(gargs.project.param)
    //转到应用所在的目录，防止在解析工程时更改了路径对后面的搜索的影响
    os.Chdir(gLookupTable.cwd)

	if gargs.clean {
		//列出不在工程文件中的文件
        root, _ := filepath.Split(gargs.project.param)
        filepath.Walk(root, cleanFn)
	} else if gargs.file.flag && gargs.file.param!=""{
        o:=new(FileWriter)
        if !gargs.output.flag {
            o.name=gargs.file.param+".dot"
        } else {
            o.name=gargs.output.param
        }
        scanSingleFile(gargs.file.param, o, nil)
	} else if gargs.recursive {
        root, _:=filepath.Split(gargs.project.param)
        var wg sync.WaitGroup
        wg.Add(len(gLookupTable.Files))
        os.Chdir(root)
        for k, _:=range gLookupTable.Files {
            abspath, _:=filepath.Abs(k)
            o:=new(FileWriter)
            o.name=abspath+".dot"
            go scanSingleFile(abspath, o, &wg)
        }
        wg.Wait()
    } else {
		log.Fatal("Missing file name to scan.")
	}
}
