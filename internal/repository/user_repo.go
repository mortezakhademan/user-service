package repository

import (
	"context"
	componentsList "git.ramooz.org/ramooz/golang-components/paginated-list"
	"github.com/mortezakhademan/user-service-sample/internal/domain"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type UserRepository interface {
	// Insert inserts a new user and returns the created record's ID (hex string).
	Insert(ctx context.Context, u *domain.User) (bson.ObjectID, error)

	// Update updates existing user by its ID (u.ID must be set).
	Update(ctx context.Context, u *domain.User) error

	// Get returns a user by ID (hex string).
	Get(ctx context.Context, id bson.ObjectID) (*domain.User, error)

	// Delete removes a user by ID (hex string).
	Delete(ctx context.Context, id bson.ObjectID) error

	// List returns all users. You can extend with pagination/filter params later.
	List(ctx context.Context, list *componentsList.List) ([]*domain.User, error)
}
