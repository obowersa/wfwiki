package modquery

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

//TODO: Expand testing, include mock of API
var getModuleTests = []struct {
	value    string
	expected string
	err      error
}{
	{"weapon", weaponURL, nil},
	{"warframe", warframeURL, nil},
	{"mod", modURL, nil},
	{"fail", "", nil},
}

func TestGetModule(t *testing.T) {
	for _, tt := range getModuleTests {
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

func TestGetStats(t *testing.T) {
	var getStatsTests = []struct {
		module   string
		query    string
		expected string
		err      error
	}{
		{"weapon", "Sigma & Octantis", "Sigma & Octantis", nil},
		{"warframe", "Ash", "Ash", nil},
		{"mod", "Abating Link", "Abating Link", nil},
	}

	n := NewWFWiki()

	for _, tt := range getStatsTests {
		t.Run(tt.module, func(t *testing.T) {
			res := n.GetStats(tt.module, tt.query)
			if !strings.HasPrefix(res, tt.expected) {
				t.Errorf("Base result does not match query: Got: %s, Expected: %s", res, tt.expected)
			}

		})
	}
}

//Expand testing to cover nil cases. Look into table testing and how to make use of for this
//TODO: Add time to first byte support for testing duration of call
func TestRequestHandler(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	defer srv.Close()

	baseStart := time.Now()

	for n := 1; n < 5; n++ {
		_, err := http.Get(srv.URL)
		if err != nil {
			t.Error(err)
		}
	}

	baseDuration := time.Since(baseStart)

	c := newRequestHandler()
	req := wikiRequest{srv.URL, nil, nil}
	testStart := time.Now()

	for n := 1; n < 5; n++ {
		res := c.handleRequest(&req)
		if res.err != nil {
			t.Error(res.err)
		}
	}

	testDuration := time.Since(testStart)

	if testDuration < (baseDuration * 5) {
		t.Errorf("rate limiting failed, testDuration: %s, baseDuration: %s", testDuration, baseDuration)
	}
}
