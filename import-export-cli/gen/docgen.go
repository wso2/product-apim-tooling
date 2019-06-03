//+build ignore

package main

import (
	"github.com/spf13/cobra/doc"
	"github.com/wso2/product-apim-tooling/import-export-cli/cmd"
	"log"
)

func main() {
	err := doc.GenMarkdownTree(cmd.RootCmd, "./docs")
	if err != nil {
		log.Fatal(err)
	}
}
