package main

import (
	"fmt"

	"github.com/elanq/banker/parser"
)

func main() {
	path := "/Users/eq/Downloads/sample.csv"
	reader := parser.NewBankReader()
	err := reader.ReadFile(path)
	if err != nil {
		fmt.Println(err)
		fmt.Println("Can't read file")
	}
}
