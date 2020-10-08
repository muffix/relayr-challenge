package httpapi

import (
	"context"
)

type customChecker struct{}

func (customChecker) Check(ctx context.Context) error {
	// Fill in meaningful checks here.
	return nil
}
