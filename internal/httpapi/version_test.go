package httpapi

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestVersion(t *testing.T) {
	s := NewService(1234)
	req := httptest.NewRequest("GET", "http://testsite.local/", nil)

	// create writer to record the response we get
	w := httptest.NewRecorder()

	// call the home page handler
	s.handleVersion()(w, req)
	resp := w.Result()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Got bad status code %d, want %d", resp.StatusCode, http.StatusOK)
	}

	// decode JSON and check contents
	got := versionResponse{}
	err := json.NewDecoder(resp.Body).Decode(&got)
	if err != nil {
		t.Fatal(err)
	}

	want := versionResponse{
		LaunchDate: got.LaunchDate,
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Got incorrect response %s, want %s", got, want)
	}
}
