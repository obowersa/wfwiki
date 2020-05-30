package wfwiki

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

const (
	weaponURL = "https://warframe.fandom.com/api.php?action=query&prop=revisions&rvprop=content&format=json&formatversion=2&titles=Module%3AWeapons%2Fdata"
)

type heavyAttack struct {
	Damage string
}

type normalDamage struct {
	damageType map[string]float64
}

type normalAttack struct {
	AttackName     string       `json:"AttackName,omitempty"`
	Damage         normalDamage `json:"Damage"`
	CritChance     float64      `json:"CritChance"`
	CritMultiplier float64      `json:"CritMultiplier"`
	StatusChance   float64      `json:"StatusChance"`
	FireRate       float64      `json:"FireRate"`
}

type weapon struct {
	BlockAngle      int          `json:"BlockAngle"`
	ComboDur        int          `json:"ComboDur"`
	FollowThrough   float64      `json:"FollowThrough"`
	MeleeRange      float64      `json:"MeleeRange"`
	SlamAttack      float64      `json:"SlamAttack"`
	SlamRadialDmg   float64      `json:"SlamRadialDmg"`
	SlamRadius      float64      `json:"SlamRadius"`
	HeavyAttack     heavyAttack  `json:"HeavyAttack"`
	WindUp          float64      `json:"WindUp"`
	HeavySlamAttack float64      `json:"HeavySlamAttack"`
	HeavyRadialDmg  float64      `json:"HeavyRadialDmg"`
	HeavySlamRadius float64      `json:"HeavySlamRadius"`
	Name            string       `json:"Name"`
	Cost            cost         `json:"Cost,omitempty"`
	Class           string       `json:"Class"`
	Conclave        bool         `json:"Conclave"`
	Disposition     float64      `json:"Disposition"`
	Image           string       `json:"Image"`
	Introduced      string       `json:"Introduced"`
	Mastery         int          `json:"Mastery"`
	NormalAttack    normalAttack `json:"NormalAttack"`
	SecondaryAttack normalAttack `json:"SecondaryAttack"`
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
type weaponData struct {
	IgnoreInCount []string          `json:"IgnoreInCount"`
	Weapons       map[string]weapon `json:"Weapons"`
	Stances       []stances         `json:"Stances"`
	Augments      []augments        `json:"Augments"`
}

func (h heavyAttack) String() string {
	if h.Damage == "" {
		return "None"
	}

	return h.Damage
}

func (h *heavyAttack) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		var i int
		if err := json.Unmarshal(data, &i); err != nil {
			return err
		}

		s = strconv.Itoa(i)
	}

	h.Damage = s

	return nil
}

func (n *normalDamage) UnmarshalJSON(data []byte) error {
	n.damageType = make(map[string]float64)
	if err := json.Unmarshal(data, &n.damageType); err != nil {
		return err
	}

	return nil
}

func (n normalDamage) totalDamage() string {
	var f float64
	for _, v := range n.damageType {
		f += v
	}

	return fmt.Sprintf("%.0f", f)
}

func (n normalDamage) damagePercent() (string, error) {
	var fm []string

	f := 0.0

	x, err := strconv.ParseFloat(n.totalDamage(), 64)
	if err != nil {
		return "", err
	}

	for k, v := range n.damageType {
		if v > f {
			f = v
			fm = nil
			fm = append(fm, k)
		} else if v == f {
			fm = append(fm, k)
		}
	}

	sort.Strings(fm) //We want a consistent order for viewing

	d := f * (100 / x)

	return fmt.Sprintf("%s: %.0f%%", strings.Join(fm, "/"), d), nil
}

func (w weaponData) getURL() string {
	return weaponURL
}

func (w weapon) getDamage() string {
	d := w.NormalAttack.Damage.totalDamage()

	v, err := w.NormalAttack.Damage.damagePercent()
	if err != nil {
		return "Unknown %"
	}

	return fmt.Sprintf("Damage: %s (%s)", d, v)
}

func (w weaponData) getStatsConcat(name string) string {
	if _, ok := w.Weapons[name]; ok {
		wWeapon := w.Weapons[name]
		fmt.Println("TEST")
		fmt.Println(wWeapon.SecondaryAttack)

		return fmt.Sprintf("%s: [Mastery: %d, Class: %s, NormalAttack: [%s, CritChance: %d%%, CritMultiplier: %.2f, StatusChance: %d%%, FireRate: %.2f], HeavyAttack: %s, SecondaryAttack: %v]",
			name,
			wWeapon.Mastery,
			wWeapon.Class,
			wWeapon.getDamage(),
			int(wWeapon.NormalAttack.CritChance*100),
			wWeapon.NormalAttack.CritMultiplier,
			int(wWeapon.NormalAttack.StatusChance*100),
			wWeapon.NormalAttack.FireRate,
			wWeapon.HeavyAttack,
			wWeapon.SecondaryAttack,
		)
	}

	return fmt.Sprintf("No weapon named: %s found", name)
}
