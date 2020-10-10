package httpapi

import (
	"net/http"

	"github.com/muffix/relayr-challenge/internal/database"
)

// offerSearchRequest is the struct representing the POST request body to the endpoint
type offerSearchRequest struct {
	ProductName string `json:"product"`
	Category string `json:"category"`
}

// offerSearchResponse is the struct representing responses to searches
type offerSearchResponse struct {
	Name string        `json:"name"`
	Category string    `json:"category"`
	Offers []offerData `json:"offers"`
}

type offerData struct{
	Supplier string `json:"supplier"`
	Price float32 `json:"price"`
}

// offerErrorResponse is the response in case an error occurs
type offerErrorResponse struct {
	Error string `json:"error"`
}

// handleOfferSearch returns an http.HandlerFunc for the home page endpoint
func (s *Service) handleOfferSearch() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		request := offerSearchRequest{}
		err := s.decode(w, r, &request)
		if err != nil {
			s.respond(w, r, offerErrorResponse{err.Error()}, http.StatusBadRequest)
			return
		}

		offers, err := s.offers.Get(request.ProductName, request.Category)
		if err != nil {
			s.respond(w, r, offerErrorResponse{err.Error()}, http.StatusInternalServerError)
			return
		}

		response := offerSearchResponse{
			Name: request.ProductName,
			Category: request.Category,
		}

		for _, o := range offers {
			response.Offers = append(response.Offers, offerData{o.Supplier, o.Price})
		}
		s.respond(w, r, response, http.StatusOK)
	}
}

type offer struct {
	Product string `json:"product"`
	Category string `json:"category"`
	offerData
}

type offerRequest offer
type offerResponse struct{
	ImportedOffers int `json:"importedOffersCount"`
}

func (s *Service) handleOffer() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		request := offerRequest{}
		err := s.decode(w, r, &request)
		if err != nil {
			s.respond(w, r, offerErrorResponse{err.Error()}, http.StatusBadRequest)
			return
		}

		err = s.offers.Insert(request.Product, request.Category, request.Supplier, request.Price)
		if err != nil {
			s.respond(w, r, offerErrorResponse{err.Error()}, http.StatusInternalServerError)
			return
		}

		s.respond(w, r, offerResponse{1}, http.StatusOK)
	}
}

type offerBatchRequest []offerRequest
type offerBatchResponse offerResponse

func (s *Service) handleOfferBatch() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		request := offerBatchRequest{}
		err := s.decode(w, r, &request)
		if err != nil {
			s.respond(w, r, offerErrorResponse{err.Error()}, http.StatusBadRequest)
			return
		}

		// Map the request data to the database model
		offerModels := make([]database.Offer, len(request))
		for i, offer := range request {
			offerModels[i] = database.Offer{
				Product:  offer.Product,
				Category: offer.Category,
				Supplier: offer.Supplier,
				Price:    offer.Price,
			}
		}


		err = s.offers.InsertMultiple(offerModels)
		if err != nil {
			s.respond(w, r, offerErrorResponse{err.Error()}, http.StatusInternalServerError)
			return
		}


		s.respond(w, r, offerBatchResponse{len(request)}, http.StatusOK)
	}
}
