package mwmod

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/obowersa/wfwiki/internal/ratelimit"
)

type datasource interface {
	Get(string) ([]byte, error)
}

type defaultds struct {
	Client *http.Client
}

func (d defaultds) Get(s string) ([]byte, error) {
	res, err := d.Client.Get(s)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return body, err
}

//Wiki encapsulates the transport client and the ratelimit handler
type Wiki struct {
	Client  datasource
	Handler ratelimit.Handler
}

//NewWiki returns a Wiki struct. This initialises our http client and ratelimit handler
func NewWiki(d datasource) *Wiki {
	if d == nil {
		d = defaultds{&http.Client{}}
	}

	return &Wiki{d, *ratelimit.NewHandler()}
}

//Request uses an input which represents an endpoint for the underlying datasource in order to request data
func (w Wiki) Request(s string) ([]byte, error) {
	r := request{&w, s}

	res, err := w.Handler.Get(&r)
	if err != nil {
		return nil, err
	}

	data, err := JSONToString(res)
	if err != nil {
		fmt.Printf("%s", err)
	}

	return []byte(data), nil
}

type request struct {
	wiki *Wiki
	url  string
}

func (r request) Call() ([]byte, error) {
	res, err := r.wiki.Client.Get(r.url)
	if err != nil {
		return nil, err
	}

	return res, nil
}
