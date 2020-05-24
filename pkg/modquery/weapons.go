package modquery

import (
	"encoding/json"
	"fmt"
)

const(
	weaponURL  = "https://warframe.fandom.com/api.php?action=query&prop=revisions&rvprop=content&format=json&formatversion=2&titles=Module%3AWeapons%2Fdata"
)
type Damage struct {
	Impact   float64 `json:"Impact"`
	Puncture float64 `json:"Puncture"`
	Slash    float64 `json:"Slash"`
}

type NormalAttack struct {
	Damage         Damage  `json:"Damage"`
	CritChance     float64 `json:"CritChance"`
	CritMultiplier float64     `json:"CritMultiplier"`
	StatusChance   float64 `json:"StatusChance"`
	FireRate       float64 `json:"FireRate"`
}

type Parts struct {
	Name  string `json:"Name"`
	Type  string `json:"Type"`
	Count int    `json:"Count"`
}

type Cost struct {
	Credits    int     `json:"Credits"`
	BPCost     int     `json:"BPCost"`
	MarketCost json.Number
	Rush       int     `json:"Rush"`
	Time       int     `json:"Time"`
	Parts      []Parts `json:"Parts"`
}

type Weapon struct {
	BlockAngle      int          `json:"BlockAngle"`
	ComboDur        int          `json:"ComboDur"`
	FollowThrough   float64      `json:"FollowThrough"`
	MeleeRange      float64      `json:"MeleeRange"`
	SlamAttack      float64          `json:"SlamAttack"`
	SlamRadialDmg   float64          `json:"SlamRadialDmg"`
	SlamRadius      float64          `json:"SlamRadius"`
	HeavyAttack     float64          `json:"HeavyAttack"`
	WindUp          float64      `json:"WindUp"`
	HeavySlamAttack float64          `json:"HeavySlamAttack"`
	HeavyRadialDmg  float64          `json:"HeavyRadialDmg"`
	HeavySlamRadius float64          `json:"HeavySlamRadius"`
	Name            string       `json:"Name"`
	Cost            Cost         `json:"Cost"`
	Class           string       `json:"Class"`
	Conclave        bool         `json:"Conclave"`
	Disposition     float64      `json:"Disposition"`
	Image           string       `json:"Image"`
	Introduced      string       `json:"Introduced"`
	Mastery         int          `json:"Mastery"`
	NormalAttack    NormalAttack `json:"NormalAttack"`
	SlideAttack     int          `json:"SlideAttack"`
	StancePolarity  string       `json:"StancePolarity"`
	Traits          []string     `json:"Traits"`
	Type            string       `json:"Type"`
	Users           []string     `json:"Users"`
}
type WeaponData struct {
	IgnoreInCount []string          `json:"IgnoreInCount"`
	Weapons       map[string]Weapon `json:"Weapons"`
	Stances       []struct {
		Name     string `json:"Name"`
		Class    string `json:"Class"`
		Polarity string `json:"Polarity,omitempty"`
		Image    string `json:"Image"`
		PvP      bool   `json:"PvP,omitempty"`
		Weapon   string `json:"Weapon,omitempty"`
		Link     string `json:"Link,omitempty"`
	} `json:"Stances"`
	Augments []struct {
		Name     string   `json:"Name"`
		Category string   `json:"Category"`
		Source   string   `json:"Source"`
		Weapons  []string `json:"Weapons"`
	} `json:"Augments"`
}

func (w WeaponData) getURL() string {
	return weaponURL
}

func (w WeaponData) getStats(name string) Weapon {
	return w.Weapons[name]
}

func (w WeaponData) getStatsConcat(name string) string {
	if _, ok := w.Weapons[name]; ok {
		return fmt.Sprintf("%s: %+v\n", name, w.Weapons[name])
	} else {
		return fmt.Sprintf("No weapon named: %s found", name)
	}
}