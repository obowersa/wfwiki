/*
Package modquery provides a client for the Warframe Fandom Wiki's module data.

Usage:

	import "github.com/obowersa/wfwiki/pkg/modquery"

Construct a new WFWiki client, then use the services on the client to access different representations of the data

	client := newWFWiki()

	//Get formatted data for a weapon
	res, _ := client.GetStats("weapon", "Reaper Prime")

Rate Limiting

MediaWiki/Fandom asks for a rate limit of 1 request per second. We achieve this by using the internal/handler library
in this code base. By default a request is processed once per second, with a buffer of up to 10 requests.

TODO: Expose rate limiting errors to the client so they can be handled appropriately

Lua Tables

MediaWiki/Fandom's module/data pages are lua tables. When querying the API a JSON string is returned which holds the
lua code. We parse this table through an embedded lua VM and convert it to JSON before unmarshalling the resulting
[]byte object into a struct
*/
package modquery

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/obowersa/wfwiki/internal/mwmod"
	"github.com/obowersa/wfwiki/internal/ratelimit"
)

type wfmodule interface {
	getURL() string
	getStatsConcat(string) string
}

//WFWiki encapsulates the transport client and the ratelimit handler
type WFWiki struct {
	Client  *http.Client
	Handler ratelimit.Handler
}

type request struct {
	wiki *WFWiki
	url  string
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

//NewWFWiki returns a WFWiki struct. This initialises our http client and ratelimit handler
func NewWFWiki() WFWiki {
	return WFWiki{&http.Client{}, *ratelimit.NewHandler()}
}

//GetStats returns an opinionated set of results for a given module
func (w WFWiki) GetStats(mod string, query string) (string, error) {
	m, err := w.module(mod)
	if err != nil {
		return "", err
	}

	r := request{&w, m.getURL()}

	res, err := w.Handler.Get(&r)
	if err != nil {
		return "", err
	}

	data, err := mwmod.JSONToString(res)
	if err != nil {
		fmt.Printf("%s", err)
	}

	if err := json.Unmarshal([]byte(data), &m); err != nil {
		fmt.Printf("%s", err)
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
	res, err := r.wiki.Client.Get(r.url)
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
