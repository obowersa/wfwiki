package modquery

import (
	"fmt"
)

const (
	modURL = "https://warframe.fandom.com/api.php?action=query&prop=revisions&rvprop=content&format=json&formatversion=2&titles=Module%3AMods%2Fdata"
)

type mods struct {
	Image        string `json:"Image"`
	Name         string `json:"Name"`
	Polarity     string `json:"Polarity"`
	Rarity       string `json:"Rarity"`
	Transmutable bool   `json:"Transmutable"`
}

type modData struct {
	IgnoreInCount []string        `json:"IgnoreInCount"`
	Mods          map[string]mods `json:"Mods"`
}

func (w modData) getURL() string {
	return modURL
}

func (w modData) getStatsConcat(name string) string {
	if _, ok := w.Mods[name]; ok {
		wMod := w.Mods[name]
		return fmt.Sprintf("%s: [Polarity: %s, Rarity: %s]", name, wMod.Polarity, wMod.Rarity)
	}

	return fmt.Sprintf("No mod named: %s found", name)
}
