package main

import (
	"fmt"
	//"reflect"
	"errors"
	"strings"
	"encoding/xml"
	"io/ioutil"
	//"bytes"
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
func parseIncludePaths(project *VCProject) {
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
				}
			}
		}
	}
}

func buildVCProject(fname string) error {
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

	parseIncludePaths(&result)
	//fmt.Printf("item definition: %v\n", result.ItemDefinitionGroup[0].ClCompile)
	fmt.Printf("files: %v\n", result)
	return nil
}
