//+build ignore

package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra/doc"
	"github.com/wso2/product-apim-tooling/import-export-cli/cmd"
)

func main() {
	log.Println("Generating docs...")
	err := doc.GenMarkdownTree(cmd.RootCmd, "docs")
	if err != nil {
		log.Fatal(err)
	}

	err = os.MkdirAll("shell-completions", os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Generating bash completions...")
	err = cmd.RootCmd.GenBashCompletionFile(filepath.FromSlash("./shell-completions/apimcli_bash_completions.sh"))
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Generating zsh completions...")
	err = cmd.RootCmd.GenZshCompletionFile(filepath.FromSlash("./shell-completions/apimcli_zsh_completions.sh"))
	if err != nil {
		log.Fatal(err)
	}
}
