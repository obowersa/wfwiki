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

type ModuleContent interface {
	Get() ([]byte, error)
}

type revision struct {
	data string `json:"*"`
}

type pages struct {
	pageid int `json:"pageid"`
	ns int `json:"ns"`
	title string `json:"title"`
	revisions []string `json:"revisions"`
}

type query struct {
	pages map[string]pages `json:"pages"`
}

type container struct {
	query query `json:"query"`
}

func (rp *revision) Unmarshal(data []byte) error {
	var cp *container
	if err := json.Unmarshal(data, &cp); err != nil {
		return fmt.Errorf("%e", err)
	}

	if l := len(cp.query.pages); l != 1 {
		return fmt.Errorf("invalid number of pages. Got: %d, Expected: 1", l)
	}

	//Pages will only be length 1, but the source object has a dynamic name. Look to refactor this
	for _, v := range cp.query.pages {
		if err := json.Unmarshal([]byte(v.revisions[0]), &rp); err != nil {
			return fmt.Errorf("%e", err)
		}
	}

	return nil
}


func JSONToString(m ModuleContent) (string, error) {
	var base revision

	c, err := m.Get()
	if err != nil {
		return "", fmt.Errorf("%e", err)
	}

	if err := json.Unmarshal(c, &base); err != nil {
		return "", fmt.Errorf("%e", err)
	}

	lua.LuaMachine.LoadModule(base.data)
	t := lua.LuaMachine.GetTable()
	data := lua.LuaMachine.ParseTable(&t, "returnJson")

	return data, nil
}
