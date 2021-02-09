package graph

import (
	"context"

	"github.com/cmsgov/easi-app/pkg/graph/model"
	"github.com/cmsgov/easi-app/pkg/storage"
)

//go:generate go run github.com/99designs/gqlgen

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

// Resolver is a resolver.
type Resolver struct {
	store          *storage.Store
	createTestDate func(context.Context, *model.TestDate) (*model.TestDate, error)
}

// NewResolver constructs a resolver
func NewResolver(
	store *storage.Store,
	createTestDate func(context.Context, *model.TestDate) (*model.TestDate, error),
) *Resolver {
	return &Resolver{store: store, createTestDate: createTestDate}
}
