package graph

import (
	"context"

	"github.com/cmsgov/easi-app/pkg/models"
	"github.com/cmsgov/easi-app/pkg/storage"
	"github.com/cmsgov/easi-app/pkg/upload"

	"github.com/google/uuid"
)

//go:generate go run github.com/99designs/gqlgen

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

// Resolver is a resolver.
type Resolver struct {
	store    *storage.Store
	service  ResolverService
	s3Client *upload.S3Client
}

// ResolverService holds service methods for use in resolvers
type ResolverService struct {
	CreateTestDate                             func(context.Context, *models.TestDate) (*models.TestDate, error)
	AddGRTFeedback                             func(context.Context, *models.GRTFeedback, *models.Action, models.SystemIntakeStatus) (*models.GRTFeedback, error)
	AuthorizeUserIsReviewTeamOrIntakeRequester func(ctx context.Context, existingIntake *models.SystemIntake) (bool, error)
	CreateActionUpdateStatus                   func(context.Context, *models.Action, uuid.UUID, models.SystemIntakeStatus, bool) (*models.SystemIntake, error)
	IssueLifecycleID                           func(context.Context, *models.SystemIntake, *models.Action) (*models.SystemIntake, error)
	RejectIntake                               func(context.Context, *models.SystemIntake, *models.Action) (*models.SystemIntake, error)
}

// NewResolver constructs a resolver
func NewResolver(
	store *storage.Store,
	service ResolverService,
	s3Client *upload.S3Client,
) *Resolver {
	return &Resolver{store: store, service: service, s3Client: s3Client}
}
