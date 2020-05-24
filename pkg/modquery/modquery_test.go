package modquery

import (
	"fmt"
	"testing"
)
//TODO: Expand testing, include mock of API
func Test(t *testing.T) {
	//fmt.Println(processTest())
}

var moduleTests = []struct {
	value string
	expected string
	err error
}{
	{"weapon", "https://warframe.fandom.com/api.php?action=query&prop=revisions&rvprop=content&format=json&formatversion=2&titles=Module%3AWeapons%2Fdata", nil },
	{"fail", "", nil },
}
func TestGetModule(t *testing.T){
	for _, tt := range moduleTests {
		t.Run(tt.value, func(t *testing.T) {
			res, err := getModule(tt.value)
			if err != nil {
				if tt.value != "fail" {
					t.Errorf("Unexpected results: Got: %s, Expected: %s, Error: %s", res.getURL(), tt.expected, err)
				}
			} else if res.getURL() != tt.expected && err != tt.err {
				t.Errorf("Unexpected results: Got: %s, Expected: %s, Error: %s", res.getURL(), tt.expected, err)
			}
		})
	}
}

func TestGetStats(t *testing.T){
	x := GetStats("warframe", "Ash")
	fmt.Println(x)
}