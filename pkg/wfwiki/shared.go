package wfwiki

import (
	"encoding/json"
	"strconv"
)

type marketCost struct {
	MCost string
}

//Parts struct shared by multiple modules
type parts struct {
	Name  string `json:"Name"`
	Type  string `json:"Type"`
	Count int    `json:"Count"`
}

//Cost struct shared by multiple modules
type cost struct {
	Credits    int        `json:"Credits,omitempty"`
	BPCost     int        `json:"BPCost,omitempty"`
	MarketCost marketCost `json:"MarketCost,omitempty"`
	Rush       int        `json:"Rush,omitempty"`
	Time       int        `json:"Time,omitempty"`
	Parts      []parts    `json:"Parts,omitempty"`
}

func (h *marketCost) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		var i int
		if err := json.Unmarshal(data, &i); err != nil {
			return err
		}

		s = strconv.Itoa(i)
	}

	h.MCost = s

	return nil
}
