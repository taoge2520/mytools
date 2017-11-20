package main

import (
	"fmt"
	"mytools/file"
)

//test
func main() {
	err := file.GetLine("./0.txt")
	if err != nil {
		fmt.Println(err)
	}
	err = file.AppendToFile("./0.txt", "test for file operator")
	if err != nil {
		fmt.Println(err)
	}
}
