package main

import (
	//"fmt"
	//"reflect"
    "sort"
    //"strconv"
	"errors"
	"strings"
	"encoding/xml"
	"io/ioutil"
	//"bytes"
    "sync"
    "os"
    "path/filepath"
)

/*******************************************************************************
*
*           XML数据结构 
*******************************************************************************/
type VCProject struct {
	ItemDefinitionGroup []ItemDefinitionGroup
	ItemGroup ItemGroup
}

type ItemDefinitionGroup struct {
	Condition string `xml:",attr"`
	ClCompile string `xml:",innerxml"`
}

type ItemGroup struct {
	Cpp []Compile `xml:"ClCompile"`
	Header []Include `xml:"ClInclude"`
}

type Compile struct {
	Include string `xml:"Include,attr"`
}

type Include struct {
	Include string `xml:"Include,attr"`
}

/*
	Unmarshal方法无法解析出路径，所以使用原始的方法来解析出路径并存储到
	结果当中。
*/
func parseIncludePaths(project *VCProject) bool {
    result:=false
	for i:=0; i<len(project.ItemDefinitionGroup); i++ {
		inputReader:=strings.NewReader(project.ItemDefinitionGroup[i].ClCompile)
		decoder:=xml.NewDecoder(inputReader)
		bMetIncludes:=false
		for t,e:=decoder.Token(); e==nil; t, e=decoder.Token(){
			switch token := t.(type) {
				// 处理元素开始（标签）
			case xml.StartElement:
				name := token.Name.Local
				if name=="AdditionalIncludeDirectories" {
					bMetIncludes=true
				}
				// 处理元素结束（标签）
			case xml.EndElement:
				//fmt.Printf("Token of '%s' end\n", token.Name.Local)
				// 处理字符数据（这里就是元素的文本）
			case xml.CharData:
				if bMetIncludes {
					content := string([]byte(token))
					//fmt.Printf("This is the content: %v\n", content)
					e=errors.New("Tag found")
					project.ItemDefinitionGroup[i].ClCompile=content
					bMetIncludes=false
                    result=true
				}
			}
		}
	}
    return result
}

func removeDuplicateStrings(stringArray *[]string) []string {
    sort.Strings(*stringArray)
    duplicate:=""
    var result []string
    for i:=0; i<len(*stringArray); i++ {
        if duplicate!=(*stringArray)[i] {
            duplicate=(*stringArray)[i]
            result=append(result, duplicate)
        }
    }

    return result
}

func buildSearchPaths(wg *sync.WaitGroup, paths []ItemDefinitionGroup, workdir string) {
    defer wg.Done()
    // 改变当前的工作路径以便后面获取绝对路径的函数有效
    os.Chdir(workdir)
    for i:=0; i<len(paths); i++ {
        arr:=strings.Split(paths[i].ClCompile, ";")
        for j:=0; j<len(arr); j++ {
            abspath, _:=filepath.Abs(arr[j])
            gLookupTable.Paths=append(gLookupTable.Paths, abspath)
        }
    }
    gLookupTable.Paths=removeDuplicateStrings(&gLookupTable.Paths)
}

func buildSearchFiles(wg *sync.WaitGroup, files ItemGroup) {
    defer wg.Done()
    gLookupTable.Files = make(map[string]string)
    //分析CPP文件
    for i:=0; i<len(files.Cpp); i++ {
        //因为在Windows系统中，如果文件名包括路径，那么分隔符一定是\，所以在内部处理时统一换成/。
        files.Cpp[i].Include=strings.Replace(files.Cpp[i].Include, "\\",string(os.PathSeparator), -1)
        dir, file:=filepath.Split(files.Cpp[i].Include)
        gLookupTable.Files[file]=dir
    }
    //分析头文件
    for i:=0; i<len(files.Header); i++ {
        files.Header[i].Include=strings.Replace(files.Header[i].Include, "\\",string(os.PathSeparator), -1)
        dir, file:=filepath.Split(files.Header[i].Include)
        gLookupTable.Files[file]=dir
    }
}

func buildVCProject(fname string) error {
    // 区分当前文件的路径和文件名
    dir, _:=filepath.Split(fname)
    var err error;
    if dir=="" {
        //没有指定目录时使用当前路径
        if dir, err=os.Getwd(); err!=nil {
            return err
        }
    }

	// 从文件读取，如可以如下：
	content, err := ioutil.ReadFile(fname)
	if err!=nil {
		return err
	}
	var result VCProject
	err=xml.Unmarshal(content, &result)
	if err!=nil {
		return err
	}

    hasInclude:=parseIncludePaths(&result)
	//fmt.Printf("item definition: %v\n", result.ItemDefinitionGroup[0].ClCompile)
	//fmt.Printf("files: %v\n", result)

    var wg sync.WaitGroup
    if hasInclude {
        wg.Add(2)
        go buildSearchPaths(&wg, result.ItemDefinitionGroup, dir)
    } else {
        wg.Add(1)
    }
    go buildSearchFiles(&wg, result.ItemGroup)
    wg.Wait()

    gLookupTable.Paths=append(gLookupTable.Paths, dir)
    //fmt.Println(gLookupTable)
	return nil
}