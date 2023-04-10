//+build ignore

package main

import (
	"log"
	"path/filepath"

	"github.com/spf13/cobra/doc"
	"github.com/wso2/product-apim-tooling/import-export-cli/cmd/mi"
)

func main() {
	log.Println("Generating docs...")
	mi.MICmd.DisableAutoGenTag = true

	err := doc.GenMarkdownTree(mi.MICmd, "docs")
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Generating MI bash completions...")
	err = mi.MICmd.GenBashCompletionFile(filepath.FromSlash("./shell-completions/mi_bash_completions.sh"))
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Generating MI zsh completions...")
	err = mi.MICmd.GenZshCompletionFile(filepath.FromSlash("./shell-completions/mi_zsh_completions.sh"))
	if err != nil {
		log.Fatal(err)
	}
}
