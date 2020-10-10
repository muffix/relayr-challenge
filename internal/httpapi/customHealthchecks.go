package httpapi

import (
	"context"
	"fmt"
)

type customChecker struct{}

func (customChecker) Check(ctx context.Context) error {
	// Fill in meaningful checks here.
	return nil
}

type databaseChecker struct {
	service *Service
}

// Check tests if we can read from the database without error
//
// Respects a potential timeout from the healthcheck config
func (c *databaseChecker) Check(ctx context.Context) error {
	errorChannel := make(chan error)

	go func() {
		_, err := c.service.offers.Get("test", "test")
		errorChannel <- err
	}()

	select {
	case <-ctx.Done():
		return fmt.Errorf("cancelled while waiting for database")
	case err := <-errorChannel:
		return err
	}
}
