package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	filePath := flag.String("file", "N/A", "file path to be processed")

	flag.Parse()

	if *filePath == "N/A" {
		println("Please specify the file path to be processed using the '-file' flag.")
		println()
		println(" example: file-to-byte -f /somewhere/over/the/rainbow/file.out")
		os.Exit(0)
	}

	filedata, err := ioutil.ReadFile(*filePath)
	if err != nil {
		panic(err)
	}

	fmt.Print("var filedata = []byte{")
	for i, v := range filedata {
		if i > 0 {
			fmt.Print(", ")
		}
		fmt.Print(v)
	}
	fmt.Println("}")
}
