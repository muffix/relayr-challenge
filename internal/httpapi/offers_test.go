package httpapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/muffix/relayr-challenge/internal/database"
)

const (
	offerBody       = `{"product":"Towel","category":"Must Haves","supplier":"Hitchhiker Essentials","price":42}`
	offerSearchBody = `{"product":"Towel", "category":"Must Haves"}`
)

// mockDB is a mock of the database which never errors, but returns some mock data
type mockDB struct{}

func (mock *mockDB) Insert(_, _, _ string, _ float32) error  { return nil }
func (mock *mockDB) InsertMultiple(_ []database.Offer) error { return nil }
func (mock *mockDB) Close() error                            { return nil }
func (mock *mockDB) Get(_, _ string) ([]database.Offer, error) {
	return []database.Offer{
		{Product: "Towel", Category: "Must Haves", Supplier: "Hitchhiker Essentials, just more expensive", Price: 44},
		{Product: "Towel", Category: "Must Haves", Supplier: "Hitchhiker Essentials", Price: 42},
		{Product: "Towel", Category: "Must Haves", Supplier: "Hitchhiker Knockoffs", Price: 42},
	}, nil
}

// mockErrorDB is a mock of the database which always errors
type mockErrorDB struct{}

func (mock *mockErrorDB) Insert(_, _, _ string, _ float32) error    { return fmt.Errorf("error") }
func (mock *mockErrorDB) InsertMultiple(_ []database.Offer) error   { return fmt.Errorf("error") }
func (mock *mockErrorDB) Close() error                              { return fmt.Errorf("error") }
func (mock *mockErrorDB) Get(_, _ string) ([]database.Offer, error) { return nil, fmt.Errorf("error") }

type mockReviewer struct{}

func (m *mockReviewer) Suppliers(supplierNames []string) (map[string]float32, error) {
	scores := make(map[string]float32)
	for _, sup := range supplierNames {
		if sup == "Hitchhiker Knockoffs" {
			scores[sup] = 1
		} else {
			scores[sup] = 3
		}

	}
	return scores, nil
}

type mockErrorReviewer struct{}

func (m *mockErrorReviewer) Suppliers(supplierNames []string) (map[string]float32, error) {
	return nil, fmt.Errorf("error")
}

func prepareTestRequest(requestBody string) (*httptest.ResponseRecorder, *http.Request) {
	req := httptest.NewRequest(
		"POST",
		"http://testsite.local/",
		bytes.NewReader([]byte(requestBody)),
	)

	// create writer to record the response we get
	w := httptest.NewRecorder()
	return w, req
}

// offerErrorScenario is a helper function for error scenarios
//
// Posts the requestBody to the handler and makes sure we get the expected HTTP response code
func offerErrorScenario(t *testing.T, h http.HandlerFunc, requestBody string, expectedStatusCode int) {
	w, req := prepareTestRequest(requestBody)

	// call the offer handler
	h(w, req)
	resp := w.Result()

	if resp.StatusCode != expectedStatusCode {
		t.Fatalf("Got bad status code %d, want %d", resp.StatusCode, expectedStatusCode)
	}

	// decode JSON and check contents
	got := offerErrorResponse{}
	err := json.NewDecoder(resp.Body).Decode(&got)
	if err != nil {
		t.Fatal(err)
	}

	if got.Error == "" {
		t.Fatal("Expected error message in response, got nothing")
	}
}

// offerSuccessScenario is a helper function for successful requests
//
// Posts the requestBody to the handler and compares the wanted struct to the one that it got.
// Makes sure we're responding with an HTTP 200 response code.
// Fails the test in case of an error or unexpected status code.
func offerSuccessScenario(t *testing.T, h http.HandlerFunc, requestBody string, got interface{}, want interface{}) {
	w, req := prepareTestRequest(requestBody)

	// call the offer handler
	h(w, req)
	resp := w.Result()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Got bad status code %d, want %d", resp.StatusCode, http.StatusOK)
	}

	// decode JSON and check contents
	err := json.NewDecoder(resp.Body).Decode(got)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Got incorrect response %v, want %v", got, want)
	}
}

func TestOfferHandler(t *testing.T) {
	service := NewService(1234)
	service.SetDatabase(&mockDB{})
	offerSuccessScenario(t, service.handleOffer(), offerBody, &offerResponse{}, &offerResponse{ImportedOffers: 1})
}

func TestOfferHandler_withInvalidBody(t *testing.T) {
	service := NewService(1234)
	service.SetDatabase(&mockDB{})
	offerErrorScenario(t, service.handleOffer(), "I'm not JSON", http.StatusBadRequest)
}

func TestOfferHandler_withDBError(t *testing.T) {
	service := NewService(1234)
	service.SetDatabase(&mockErrorDB{})
	offerErrorScenario(t, service.handleOffer(), offerBody, http.StatusInternalServerError)
}

func TestBatchOfferHandler(t *testing.T) {
	service := NewService(1234)
	service.SetDatabase(&mockDB{})
	batchRequestBody := fmt.Sprintf("[%s, %s]", offerBody, offerBody)
	offerSuccessScenario(
		t,
		service.handleOfferBatch(),
		batchRequestBody,
		&offerResponse{},
		&offerResponse{ImportedOffers: 2},
	)
}

func TestBatchOfferHandler_withInvalidBody(t *testing.T) {
	service := NewService(1234)
	service.SetDatabase(&mockDB{})
	offerErrorScenario(t, service.handleOfferBatch(), "I'm not JSON", http.StatusBadRequest)
}

func TestBatchOfferHandler_withDBError(t *testing.T) {
	service := NewService(1234)
	service.SetDatabase(&mockErrorDB{})
	batchRequestBody := fmt.Sprintf("[%s, %s]", offerBody, offerBody)
	offerErrorScenario(
		t,
		service.handleOfferBatch(),
		batchRequestBody,
		http.StatusInternalServerError,
	)
}

func TestOfferSearch(t *testing.T) {
	service := NewService(1234)
	service.SetDatabase(&mockDB{})
	service.SetReviewer(&mockReviewer{})

	offerSuccessScenario(
		t,
		service.handleOfferSearch(),
		offerSearchBody,
		&offerSearchResponse{},
		&offerSearchResponse{
			Name:     "Towel",
			Category: "Must Haves",
			Offers: []offerData{
				{
					Supplier:    "Hitchhiker Essentials",
					ReviewScore: 3,
					Price:       42,
				},
				{
					Supplier:    "Hitchhiker Knockoffs",
					ReviewScore: 1,
					Price:       42,
				},
				{
					Supplier:    "Hitchhiker Essentials, just more expensive",
					ReviewScore: 3,
					Price:       44,
				},
			},
		},
	)
}

func TestOfferSearch_withInvalidBody(t *testing.T) {
	service := NewService(1234)
	service.SetDatabase(&mockDB{})
	offerErrorScenario(
		t,
		service.handleOfferSearch(),
		"I'm not JSON",
		http.StatusBadRequest,
	)
}

func TestOfferSearch_withDBError(t *testing.T) {
	service := NewService(1234)
	service.SetDatabase(&mockErrorDB{})
	offerErrorScenario(
		t,
		service.handleOfferSearch(),
		offerSearchBody,
		http.StatusInternalServerError,
	)
}

func TestOfferSearch_withReviewerError(t *testing.T) {
	service := NewService(1234)
	service.SetDatabase(&mockDB{})
	service.SetReviewer(&mockErrorReviewer{})
	offerErrorScenario(
		t,
		service.handleOfferSearch(),
		offerSearchBody,
		http.StatusBadGateway,
	)
}
