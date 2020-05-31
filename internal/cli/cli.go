package cli

import (
	"errors"
	"flag"
)

var (
	//ErrInvalidArgs is returned when we have missing cli arguments
	ErrInvalidArgs = errors.New("cli: missing arguments")
)

//WFWikiParams is used to encapsulate the values CLI args we're expecting
type WFWikiParams struct {
	Module *string
	Query  *string
}

//Parse sets up our expected CLI flags and parses them
func (w *WFWikiParams) Parse(args []string) error {
	cmd := flag.NewFlagSet("wfwiki", flag.ContinueOnError)
	w.Module = cmd.String("module", "", "Specifies the Module to Query")
	w.Query = cmd.String("query", "", "Specifies the item to lookup")

	if err := cmd.Parse(args); err != nil {
		return err
	}

	if *w.Module == "" || *w.Query == "" {
		cmd.PrintDefaults()
		return ErrInvalidArgs
	}

	return nil
}
