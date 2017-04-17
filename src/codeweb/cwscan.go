package main

import (
	"log"
	"fmt"
	"os"
	"strings"
	"errors"
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
	flag string
	param string
}

type GArguments struct{
	clean bool
	project Option
	file Option
}

type FilesArch struct{
	paths []string
	files [][]string
}

var gargs GArguments
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
			gargs.file.flag="f"
		case "m":
			gargs.project.flag="m"
		default:
			return errors.New("Invalid option: "+arg)
		}
	} else {
		if gargs.file.flag=="f" && gargs.file.param=="" {
			gargs.file.param=arg
		} else if gargs.project.flag=="m" && gargs.project.param=="" {
			gargs.project.param=arg
		} else {
			return errors.New("Unknown arguments: "+arg)
		}
	}
	return nil
}

func parseMakefile(make string) {
	if strings.HasSuffix(make, ".vcxproj") || strings.HasSuffix(make, ".vcproj") {
		// VC project
		if 	err:=buildVCProject(make); err!=nil {
			log.Fatal(err)
		}
	}
}

func main(){
	argnum := len(os.Args)
	if argnum==1 || argnum>4 {
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

	if gargs.project.flag!="m" && gargs.project.param=="" {
		log.Fatal("Need a project name.")
	}

	parseMakefile(gargs.project.param)

	if gargs.clean {
		//Clean up abundant files of the project
		
	} else if gargs.file.flag=="f" && gargs.file.param!=""{

	} else {
		log.Fatal("Missing file name to scan.")
	}
}
