package httpapi

import (
	"context"
	"testing"
)

func TestDatabaseChecker_Check(t *testing.T) {
	service := NewService(1234)
	service.SetDatabase(&mockDB{})

	underTest := &databaseChecker{service}
	if err := underTest.Check(context.Background()); err != nil {
		t.Fatalf("Expected no error for a good database, got %v", err)
	}

	service.SetDatabase(&mockErrorDB{})
	if err := underTest.Check(context.Background()); err == nil {
		t.Fatalf("Expected an error for an erroring database, got none")
	}
}

func TestDatabaseChecker_WithCancelledContext(t *testing.T) {
	service := NewService(1234)
	service.SetDatabase(&mockDB{})

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	underTest := &databaseChecker{service}
	if err := underTest.Check(ctx); err == nil {
		t.Fatalf("Expected an error for a cancelled context, got none")
	}
}
