package modquery

import (
	"encoding/json"
	"fmt"
	"github.com/obowersa/wfwiki/internal/mwmod"
	"github.com/obowersa/wfwiki/internal/ratelimit"
	"io/ioutil"
	"net/http"
	"strings"
)

//TODO: Refactor this when I figure out client side API
var WikiBase WFWiki

type wfmodule interface {
	getURL() string
	getStatsConcat(string) string
}

type WFWiki struct {
	client *http.Client
	handler ratelimit.Handler
}

type request struct {
	wiki *WFWiki
	url string
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


func init() {
	WikiBase = newWFWiki()
}

func newWFWiki() WFWiki {
	return WFWiki{&http.Client{},*ratelimit.NewHandler()}
}

//Refactor the below into seperate functions
func (w WFWiki) GetStats(mod string, query string) (string, error) {

	m, err := w.module(mod)
	if err != nil {
		return "", err
	}

	r := request{&w,m.getURL()}

	res, err := w.handler.Get(&r)
	if err != nil {
		return "", err
	}

	data, err := mwmod.JSONToString(res)
	if err != nil {
		fmt.Errorf("%s", err)
	}


	if err := json.Unmarshal([]byte(data), &m); err != nil {
		fmt.Println("TEST")
	}
	return m.getStatsConcat(query), nil
}

func (w WFWiki) module(n string) (wfmodule, error) {
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

func (r request) Call() ([]byte, error) {
	res, err := r.wiki.client.Get(r.url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}