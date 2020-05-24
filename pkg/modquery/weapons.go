package modquery

import (
	"fmt"
)

const (
	weaponURL = "https://warframe.fandom.com/api.php?action=query&prop=revisions&rvprop=content&format=json&formatversion=2&titles=Module%3AWeapons%2Fdata"
)

type damage struct {
	Impact   float64 `json:"Impact"`
	Puncture float64 `json:"Puncture"`
	Slash    float64 `json:"Slash"`
}

type normalAttack struct {
	Damage         damage  `json:"Damage"`
	CritChance     float64 `json:"CritChance"`
	CritMultiplier float64 `json:"CritMultiplier"`
	StatusChance   float64 `json:"StatusChance"`
	FireRate       float64 `json:"FireRate"`
}

type weapon struct {
	BlockAngle      int          `json:"BlockAngle"`
	ComboDur        int          `json:"ComboDur"`
	FollowThrough   float64      `json:"FollowThrough"`
	MeleeRange      float64      `json:"MeleeRange"`
	SlamAttack      float64      `json:"SlamAttack"`
	SlamRadialDmg   float64      `json:"SlamRadialDmg"`
	SlamRadius      float64      `json:"SlamRadius"`
	HeavyAttack     float64      `json:"HeavyAttack"`
	WindUp          float64      `json:"WindUp"`
	HeavySlamAttack float64      `json:"HeavySlamAttack"`
	HeavyRadialDmg  float64      `json:"HeavyRadialDmg"`
	HeavySlamRadius float64      `json:"HeavySlamRadius"`
	Name            string       `json:"Name"`
	Cost            cost         `json:"Cost, omitempty"`
	Class           string       `json:"Class"`
	Conclave        bool         `json:"Conclave"`
	Disposition     float64      `json:"Disposition"`
	Image           string       `json:"Image"`
	Introduced      string       `json:"Introduced"`
	Mastery         int          `json:"Mastery"`
	NormalAttack    normalAttack `json:"NormalAttack"`
	SlideAttack     int          `json:"SlideAttack"`
	StancePolarity  string       `json:"StancePolarity"`
	Traits          []string     `json:"Traits"`
	Type            string       `json:"Type"`
	Users           []string     `json:"Users"`
}

type stances struct {
	Name     string `json:"Name"`
	Class    string `json:"Class"`
	Polarity string `json:"Polarity,omitempty"`
	Image    string `json:"Image"`
	PvP      bool   `json:"PvP,omitempty"`
	Weapon   string `json:"weapon,omitempty"`
	Link     string `json:"Link,omitempty"`
}

type augments struct {
	Name     string   `json:"Name"`
	Category string   `json:"Category"`
	Source   string   `json:"Source"`
	Weapons  []string `json:"Weapons"`
}
type WeaponData struct {
	IgnoreInCount []string          `json:"IgnoreInCount"`
	Weapons       map[string]weapon `json:"Weapons"`
	Stances       []stances         `json:"Stances"`
	Augments      []augments        `json:"Augments"`
}

func (w WeaponData) getURL() string {
	return weaponURL
}

func (w WeaponData) getStats(name string) weapon {
	return w.Weapons[name]
}

func (w WeaponData) getStatsConcat(name string) string {
	if _, ok := w.Weapons[name]; ok {
		return fmt.Sprintf("%s: %+v\n", name, w.Weapons[name])
	} else {
		return fmt.Sprintf("No warframe named: %s found", name)
	}
}
