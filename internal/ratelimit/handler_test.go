package ratelimit

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type testEndpoint struct{}

func (e testEndpoint) Call() ([]byte, error) {
	return []byte("TestingEndpoint"), nil
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

	c := NewHandler()
	req := testEndpoint{}
	testStart := time.Now()

	for n := 1; n < 5; n++ {
		_, err := c.Get(&req)
		if err != nil {
			t.Error(err)
		}
	}

	testDuration := time.Since(testStart)

	if testDuration < (baseDuration * 5) {
		t.Errorf("rate limiting failed, testDuration: %s, baseDuration: %s", testDuration, baseDuration)
	}
}
