package modquery

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/obowersa/wfwiki/internal/lua"
)

type wfmodule interface {
	getURL() string
	getStatsConcat(string) string
}

type wikiJSON struct {
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

type wikiResponse struct {
	response *http.Response
	err      error
}

type wikiRequest struct {
	url    string
	respCh chan *wikiResponse
	ctx    context.Context
}

type requestHandler struct {
	client     *http.Client
	rate       time.Duration
	reqCh      chan *wikiRequest
	reqTimeout time.Duration
}

func newRequestHandler() *requestHandler {
	client := &http.Client{}
	rate := time.Second
	requests := make(chan *wikiRequest, 10)
	reqTimeout := 5 * time.Second
	r := requestHandler{client, rate, requests, reqTimeout}

	go r.processRequests()

	return &r
}

func (r *requestHandler) processRequests() {
	throttle := time.NewTicker(r.rate)

	for req := range r.reqCh {
		<-throttle.C
		//NOTE: Move getWikiContent function to be a method of requestHandler?
		//NOTE: Reconsider this when testing turns out to be a pain
		go r.fulfillRequest(req)
	}
}

func (r *requestHandler) fulfillRequest(wReq *wikiRequest) {
	//TODO: Add in hook to caching algorithm. Want to avoid parsing the lua code into the VM for every response
	//TODO: Potentially impact cache check within process reqCh. Return cache before waiting for ticker so repeated reqCh can be fulfilled outside of the 1/s tick loop
	res, err := r.client.Get(wReq.url)
	//TODO: Refactor select statement based on: https://blog.golang.org/pipelines
	select {
	case <-wReq.ctx.Done():
		return
	default:
		wReq.respCh <- &wikiResponse{res, err}
	}
}

func (r *requestHandler) handleRequest(req *wikiRequest) (resp *wikiResponse) {
	//TODO: Look into moving getModule func into requestHandler method sig
	c := make(chan *wikiResponse)
	ctx, cancel := context.WithCancel(context.Background())

	defer close(c)

	r.reqCh <- &wikiRequest{req.url, c, ctx}
	select {
	case n := <-c:
		cancel()
		resp = n
	case <-time.After(r.reqTimeout):
		cancel()
		//TODO: Refactor error handling
		resp = &wikiResponse{err: fmt.Errorf("timeout waiting for response from: %s", req.url)}
	}

	return
}

type WFWiki struct {
	httpReq *requestHandler
}

func NewWFWiki() WFWiki {
	return WFWiki{newRequestHandler()}
}

//Refactor the below into seperate functions
func (w *WFWiki) request(m string) (wfmodule, string, error){
	var wikiBase wikiJSON

	module, err := getModule(m)
	if err != nil {
		//TODO: Handle this gracefully
		panic(err)
	}

	//Wiki request creator ?
	req := wikiRequest{module.getURL(), nil, nil}
	res := w.httpReq.handleRequest(&req).response

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	if err := json.Unmarshal(body, &wikiBase); err != nil {
		panic(err)
	}

	var l string

	if len(wikiBase.Query.Pages) != 1 {
		return module, "", fmt.Errorf("too many pages for single request")
	}

	for _, v := range wikiBase.Query.Pages {
		l += v.Revisions[0].Data
	}

	return module, l, nil
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


//GetStats queries the wiki module specified by mod, and returns the stats about the object specified
//by query
func (w WFWiki) GetStats(mod string, query string) string {
	module, moduleLua, err := w.request(mod)
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
