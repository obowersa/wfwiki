package main

import (
	"fmt"
	"os"

	"github.com/obowersa/wfwiki/pkg/wfwiki"

	"github.com/obowersa/wfwiki/internal/cli"
)

func main() {
	args := cli.WFWikiParams{}
	if err := args.Parse(os.Args[1:]); err != nil {
		os.Exit(1)
	}

	w := wfwiki.NewWFWiki()

	res, err := w.GetStats(*args.Module, *args.Query)
	if err != nil {
		fmt.Println("%w", err)
		os.Exit(1)
	}

	fmt.Println(res)
}
