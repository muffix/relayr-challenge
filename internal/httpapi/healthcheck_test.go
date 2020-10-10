package httpapi

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/muffix/relayr-challenge/internal/database"
)

func TestReadiness(t *testing.T) {
	s := NewService(1234)
	s.SetDatabase(&mockDB{})
	req := httptest.NewRequest("GET", "http://testsite.local/", nil)

	// create writer to record the response we get
	w := httptest.NewRecorder()
	s.handleReadiness()(w, req)

	resp := w.Result()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Got bad status code %d, want %d", resp.StatusCode, http.StatusOK)
	}

	// decode JSON and check contents
	got := healthcheckResponse{}
	err := json.NewDecoder(resp.Body).Decode(&got)
	if err != nil {
		t.Fatal(err)
	}

	want := healthcheckResponse{Status: "OK"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Expected %s, got %s", want, got)
	}
}

func TestLiveness(t *testing.T) {
	testCases := []struct {
		db             database.Offers
		want           healthcheckResponse
		expectedStatus int
	}{
		{
			db:             &mockDB{},
			want:           healthcheckResponse{Status: "OK"},
			expectedStatus: http.StatusOK,
		},
		{
			db: &mockErrorDB{},
			want: healthcheckResponse{
				Status: "Service Unavailable",
				Errors: map[string]string{"database": "error"},
			},
			expectedStatus: http.StatusServiceUnavailable,
		},
	}
	s := NewService(1234)

	for _, testCase := range testCases {
		s.SetDatabase(testCase.db)
		req := httptest.NewRequest("GET", "http://testsite.local/", nil)

		// create writer to record the response we get
		w := httptest.NewRecorder()
		s.handleLiveness()(w, req)

		resp := w.Result()

		if resp.StatusCode != testCase.expectedStatus {
			t.Fatalf("Got bad status code %d, want %d", resp.StatusCode, testCase.expectedStatus)
		}

		// decode JSON and check contents
		got := healthcheckResponse{}
		err := json.NewDecoder(resp.Body).Decode(&got)
		if err != nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual(got, testCase.want) {
			t.Fatalf("Expected %s, got %s", testCase.want, got)
		}
	}
}
