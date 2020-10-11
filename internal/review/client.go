package review

import (
	"math"
	"math/rand"
)

// Reviewer is an interface for fetching reviews for suppliers
type Reviewer interface {
	Suppliers(supplierNames []string) (map[string]float32, error)
}

// Client is a reviews client
//
// The scores are the average values out of the ratings given by customers between 1 and 5.
// This is obviously a dummy client. A real one would call out to the actual service and return
// real scores. We're happy enough with some sweet randomness here.
type Client struct{}

// randomScore returns a score between (1,5]
func (c *Client) randomScore() float32 {
	return float32(math.Round((6-(float64(rand.Intn(100))/20.0))*100) / 100)
}

// Suppliers returns the review score for the given suppliers
func (c *Client) Suppliers(supplierNames []string) (map[string]float32, error) {
	reviews := make(map[string]float32)
	for _, supplier := range supplierNames {
		reviews[supplier] = c.randomScore()
	}
	return reviews, nil
}
