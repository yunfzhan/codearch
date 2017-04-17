package main

import (
	"fmt"
	"os"
	"bufio"
)

func buildVCProject(fname string) error {
	fmt.Println(fname)
	f, err:=os.Open(fname)
	defer f.Close()
	if err!=nil {
		return err	
	}

	var lines []string
	scanner:=bufio.NewScanner(f)
	for scanner.Scan() {
		lines=append(lines, scanner.Text())
	}

	if scanner.Err()!=nil {
		return scanner.Err()
	}

	return nil
}
