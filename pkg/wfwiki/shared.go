package wfwiki

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
