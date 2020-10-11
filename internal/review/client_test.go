package review

import (
	"fmt"
	"testing"
)

func TestClient_Suppliers(t *testing.T) {
	suppliers := make([]string, 100)
	for i := 0; i < 100; i++ {
		suppliers[i] = fmt.Sprintf("client_%d", i)
	}

	reviewerUnderTest := &Client{}

	reviews, err := reviewerUnderTest.Suppliers(suppliers)

	if err != nil {
		t.Fatalf("Expected no error retrieving reviews, got %v", err)
	}

	if len(reviews) != 100 {
		t.Fatalf("Expected 100 reviews, got %d", len(reviews))
	}

	for _, sup := range suppliers {
		if _, ok := reviews[sup]; !ok {
			t.Fatalf("Expected score for supplier %s, but didn't find one", sup)
		}
	}
}
