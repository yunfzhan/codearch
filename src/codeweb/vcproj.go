package main

import (
	"fmt"
	//"os"
	//"bufio"
	"encoding/xml"
	"io/ioutil"
	//"bytes"
)

type VCProject struct {
	ItemDefinitionGroup []ItemDefinitionGroup
	ItemGroup ItemGroup
}

type ItemDefinitionGroup struct {
	Condition string `xml:",attr"`
	ClCompile ClCompile
}

type ClCompile struct {
	//Definitions string `xml:"PreprocessorDefinitions,chardata"`
	Directories string `xml:"AdditionalIncludeDirectories, chardata"`
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

	fmt.Printf("item definition: %v\n", result.ItemDefinitionGroup)
	//fmt.Printf("files: %v\n", result.ItemGroup)
	return nil
}
