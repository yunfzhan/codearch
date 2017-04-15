package main

import (
	"fmt"
	"os"
	"strings"
	"errors"
)

func printHelp(){
	fmt.Println("**************************************************************")
	fmt.Println("*                   CWScan v1.0 Help                         *")
	fmt.Println("**************************************************************")
	fmt.Println("")
	fmt.Println("cwscan -[options: d|f|c] [arg] dotfile")
	fmt.Println("\n")
	fmt.Println("d\t\t\targ is a folder name.")
	fmt.Println("f\t\t\targ is a file name.")
	fmt.Println("c\t\t\targ is a constant or function name.")
	fmt.Println("output is generated .dot file name.")
}

type GArguments struct{
	flag string		// d|f|c - 区分输入的名称类型
	infile string	//根据上面的类型，这个可以是目录，文件名或函数与常量
	outfile string	//要生成的.dot文件名
}

var gargs GArguments
//var argMap map[string][]string

func processCmdLine(arg string, flag bool) error {
	if flag {
		if arg!="d" && arg!="f" && arg!="c" {
			return errors.New("Invalid option: "+arg)
		}
		gargs.flag=arg
	} else {
		if gargs.infile=="" {
			gargs.infile=arg
		} else if gargs.outfile==""{ 
			gargs.outfile=arg
		} else {
			return errors.New("Unknown argument!")
		}
	}
	return nil
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
			fmt.Println(err)
			return
		}
	}

	if gargs.flag=="" {
		fmt.Println("Missing option name!")
		return
	} else if gargs.infile=="" {
		fmt.Println("Missing input name!")
		return
	} else if gargs.outfile=="" {
		fmt.Println("Missing .dot name!")
		return
	}
	fmt.Printf("option: %v, input: %v, output: %v\n", gargs.flag, gargs.infile, gargs.outfile);
	
}
