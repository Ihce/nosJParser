package main

import (
	"fmt"
	deserializer "nosj_parser/internal/deserialize"
	"os"
)

func main() {
	os.Exit(run())
}

func run() int {
	var args = os.Args
	if len(args) != 2 {
		fmt.Fprintln(os.Stderr, "ERROR -- Could not resolve file path")
		return 66
	}

	filePath := os.Args[1]
	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "ERROR -- Could not read file")
		return 66
	}
	result, err := deserializer.Deserialize(string(fileContent))

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 66
	}

	fmt.Print(result)
	return 0
}
