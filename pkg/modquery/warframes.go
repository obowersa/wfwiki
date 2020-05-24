package modquery

import "fmt"

const(
	warframeURL  = "https://warframe.fandom.com/api.php?action=query&prop=revisions&rvprop=content&format=json&formatversion=2&titles=Module%3AWarframes%2Fdata"
)

type Warframe     struct {
			Armor        int      `json:"Armor"`
			AuraPolarity string   `json:"AuraPolarity"`
			Conclave     bool     `json:"Conclave"`
			Energy       int      `json:"Energy"`
			Health       int      `json:"Health"`
			Image        string   `json:"Image"`
			Portrait     string   `json:"Portrait"`
			Name         string   `json:"Name"`
			Polarities   []string `json:"Polarities"`
			Shield       int      `json:"Shield"`
			Sprint       float64  `json:"Sprint"`
			Introduced   string   `json:"Introduced"`
			Sex          string   `json:"Sex"`
			Vaulted      bool     `json:"Vaulted"`
			MainCost     struct {
				Credits    int `json:"Credits"`
				BPCost     int `json:"BPCost"`
				MarketCost int `json:"MarketCost"`
				Rush       int `json:"Rush"`
				Time       int `json:"Time"`
				Parts      []struct {
					Name  string `json:"Name"`
					Type  string `json:"Type"`
					Count int    `json:"Count"`
				} `json:"Parts"`
			} `json:"MainCost"`
			NeuroCost struct {
				Credits int `json:"Credits"`
				Rush    int `json:"Rush"`
				Time    int `json:"Time"`
				Parts   []struct {
					Name  string `json:"Name"`
					Type  string `json:"Type"`
					Count int    `json:"Count"`
				} `json:"Parts"`
			} `json:"NeuroCost"`
			ChassisCost struct {
				Credits int `json:"Credits"`
				Rush    int `json:"Rush"`
				Time    int `json:"Time"`
				Parts   []struct {
					Name  string `json:"Name"`
					Type  string `json:"Type"`
					Count int    `json:"Count"`
				} `json:"Parts"`
			} `json:"ChassisCost"`
			SystemCost struct {
				Credits int `json:"Credits"`
				Rush    int `json:"Rush"`
				Time    int `json:"Time"`
				Parts   []struct {
					Name  string `json:"Name"`
					Type  string `json:"Type"`
					Count int    `json:"Count"`
				} `json:"Parts"`
			} `json:"SystemCost"`
}

type WarframeData struct {
	IgnoreInCount []string          `json:"IgnoreInCount"`
	Warframes 	map[string]Warframe
}

func (w WarframeData) getURL() string {
	return warframeURL
}


func (w WarframeData) getStatsConcat(name string) string {
	if _,ok := w.Warframes[name]; ok {
		return fmt.Sprintf("%s: %+v\n", name, w.Warframes[name])
	} else {
		return fmt.Sprintf("No warframe named: %s found", name)
	}

}
