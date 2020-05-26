package modquery

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/obowersa/wfwiki/internal/lua"
)

type wfmodule interface {
	getURL() string
	getStatsConcat(string) string
}

type wfwiki struct {
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

//Parts struct shared by multiple modules
type parts struct {
	Name  string `json:"Name"`
	Type  string `json:"Type"`
	Count int    `json:"Count"`
}

//Cost struct shared by multiple modules
type cost struct {
	Credits    int     `json:"Credits,omitempty"`
	BPCost     int     `json:"BPCost,omitempty"`
	MarketCost float64 `json:"MarketCost,omitempty"`
	Rush       int     `json:"Rush,omitempty"`
	Time       int     `json:"Time,omitempty"`
	Parts      []parts `json:"Parts,omitempty"`
}

func getWikiContent(w wfmodule) wfwiki {
	var module wfwiki

	req, err := http.NewRequest("GET", w.getURL(), nil)
	if err != nil {
		panic(err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	if err := json.Unmarshal(body, &module); err != nil {
		panic(err)
	}

	return module
}

func getModule(n string) (wfmodule, error) {
	n = strings.ToLower(n)
	switch n {
	case "weapon":
		return new(weaponData), nil
	case "warframe":
		return new(warframeData), nil
	case "mod":
		return new(modData), nil
	default:
		return nil, fmt.Errorf("%s does not match supported modules", n)
	}
}

func getWikiLua(m wfmodule) (string, error) {
	var l string

	wikiJSON := getWikiContent(m)
	if len(wikiJSON.Query.Pages) != 1 {
		return "", fmt.Errorf("too many pages for single request")
	}

	for _, v := range wikiJSON.Query.Pages {
		l += v.Revisions[0].Data
	}

	return l, nil
}

//GetStats queries the wiki module specified by mod, and returns the stats about the object specified
//by query
func GetStats(mod string, query string) string {
	module, err := getModule(mod)
	if err != nil {
		panic(err)
	}

	moduleLua, err := getWikiLua(module)
	if err != nil {
		panic(err)
	}

	lua.LuaMachine.LoadModule(moduleLua)
	t := lua.LuaMachine.GetTable()
	data := lua.LuaMachine.ParseTable(&t, "returnJson")

	if err := json.Unmarshal([]byte(data), &module); err != nil {
		fmt.Println(err)
	}

	return module.getStatsConcat(query)
}
