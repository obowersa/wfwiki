package modquery

import (
	"encoding/json"
	"fmt"
	"math"
	"strconv"
)

const (
	weaponURL = "https://warframe.fandom.com/api.php?action=query&prop=revisions&rvprop=content&format=json&formatversion=2&titles=Module%3AWeapons%2Fdata"
)

type heavyAttack struct {
	Damage string
}

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

func (h heavyAttack) String() string {
	return fmt.Sprintf("%s", h.Damage)
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

func (w WeaponData) getURL() string {
	return weaponURL
}

func (w WeaponData) getStats(name string) weapon {
	return w.Weapons[name]
}

func (w weapon) getDamage() string {
	//TODO: Testing, need to create weapon dummy, and move ifs to switch?
	d := w.NormalAttack.Damage
	damage := math.Round(d.Slash + d.Impact + d.Puncture)
	slashPer := math.Round(d.Slash * (100 / damage))
	impactPer := math.Round(d.Impact * (100 / damage))
	punctPer := math.Round(d.Puncture * (100 / damage))

	damageFmt := fmt.Sprintf("Damage: %.2f", damage)
	if slashPer == 50.00 || impactPer == 50.00 || punctPer == 50.00 {
		if slashPer == 50.00 && impactPer == 50.00 {
			return fmt.Sprintf("%s: [Slash: %.2f%%/Impact: %.2f%%]", damageFmt, slashPer, impactPer)
		} else if slashPer == 50.00 && punctPer == 50.00 {
			return fmt.Sprintf("%s: [Slash: %.2f%%/Puncture: %.2f%%]", damageFmt, slashPer, punctPer)
		} else {
			return fmt.Sprintf("%s: [Impact: %.2f%%/Puncture: %.2f%%]", damageFmt, impactPer, punctPer)
		}
	} else if slashPer > impactPer && slashPer > punctPer {
		return fmt.Sprintf("%s: [Slash: %.2f%%]", damageFmt, slashPer)
	} else if impactPer > punctPer {
		return fmt.Sprintf("%s: [Impact: %.2f%%]", damageFmt, impactPer)
	}

	return fmt.Sprintf("%s: [Puncture: %.2f%%]", damageFmt, punctPer)

}
func (w WeaponData) getStatsConcat(name string) string {
	if _, ok := w.Weapons[name]; ok {
		wWeapon := w.Weapons[name]
		return fmt.Sprintf("%s: [Mastery: %d, Type: %s, Class: %s, NormalAttack: [CritChance: %d%%, CritMultiplier: %.2f, StatusChance: %d%%, %s, HeavyAttack: %s FireRate: %.2f]]",
			name,
			wWeapon.Mastery,
			wWeapon.Type,
			wWeapon.Class,
			int(wWeapon.NormalAttack.CritChance*100),
			wWeapon.NormalAttack.CritMultiplier,
			int(wWeapon.NormalAttack.StatusChance*100),
			wWeapon.getDamage(),
			wWeapon.HeavyAttack,
			wWeapon.NormalAttack.FireRate)
	}
	return fmt.Sprintf("No weapon named: %s found", name)
}
