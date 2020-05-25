package modquery

import "fmt"

const (
	warframeURL = "https://warframe.fandom.com/api.php?action=query&prop=revisions&rvprop=content&format=json&formatversion=2&titles=Module%3AWarframes%2Fdata"
)

type Warframe struct {
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
	MainCost     cost     `json:"MainCost"`
	NeuroCost    cost     `json:"NeuroCost"`
	ChassisCost  cost     `json:"ChassisCost"`
	SystemCost   cost     `json:"SystemCost"`
}

type WarframeData struct {
	IgnoreInCount []string `json:"IgnoreInCount"`
	Warframes     map[string]Warframe
}

func (w WarframeData) getURL() string {
	return warframeURL
}

func (w WarframeData) getStatsConcat(name string) string {
	if _, ok := w.Warframes[name]; ok {
		wframe := w.Warframes[name]
		return fmt.Sprintf("%s: [Armor: %d, Shield: %d, Health: %d, Energy: %d]", name, wframe.Armor, wframe.Shield, wframe.Health, wframe.Energy)
	}
	return fmt.Sprintf("No warframe named: %s found", name)
}
