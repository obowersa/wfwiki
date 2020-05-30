/*
Package mwmod converts the response from a mediawiki API request for a Module/data page and returns a json string
with the pages content. Behind the scenes this uses an embedded lua VM to execute the table and convert it to json
*/
package mwmod

import (
	"encoding/json"
	"fmt"

	"github.com/obowersa/wfwiki/internal/lua"
)

/*
TODO: Revisit ModuleInterface
ModuleInterface is intended to be used with the JSOnToString implementation. Need to revisit this
type ModuleContent interface {
	Get() ([]byte, error)
}*/

type container struct {
	Query struct {
		Pages map[string]struct {
			Pageid    int    `json:"pageid"`
			Ns        int    `json:"ns"`
			Title     string `json:"title"`
			Revisions []struct {
				Data string `json:"*"`
			} `json:"revisions"`
		} `json:"pages"`
	} `json:"query"`
}

/*
JSONToString unmarshalls a MediaWiki/Fandom API response for a module/data page, then calls out to a Lua library
to process the remaining data table into a JSON string.
*/
func JSONToString(m []byte) (string, error) {
	var cp container

	var base string

	if err := json.Unmarshal(m, &cp); err != nil {
		return "", fmt.Errorf("%e", err)
	}

	//Pages will only be length 1, but the source object has a dynamic name. Look to refactor this
	for _, v := range cp.Query.Pages {
		base = v.Revisions[0].Data
	}

	lua.LuaMachine.LoadModule(base)
	t := lua.LuaMachine.GetTable()
	data := lua.LuaMachine.ParseTable(&t, "returnJson")

	return data, nil
}
