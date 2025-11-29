package repository

import (
	"context"
	"errors"
	componentsList "git.ramooz.org/ramooz/golang-components/paginated-list"
	"git.ramooz.org/ramooz/golang-components/paginated-list/mongodb/models"
	"github.com/mortezakhademan/user-service-sample/internal/domain"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"time"
)

// MongoUserRepository is a MongoDB implementation of UserRepository.
type MongoUserRepository struct {
	collection *mongo.Collection
	timeout    time.Duration
}

// NewMongoUserRepository creates a new repository bound to a MongoDB collection.
func NewMongoUserRepository(db *mongo.Database, collectionName string) *MongoUserRepository {
	return &MongoUserRepository{
		collection: db.Collection(collectionName),
		timeout:    5 * time.Second,
	}
}

// Insert implements repository.Insert.
// If u.ID is empty, a new ObjectID is generated. The stored Mongo _id is ObjectID and returned as hex string.
func (r *MongoUserRepository) Insert(ctx context.Context, u *domain.User) (bson.ObjectID, error) {
	// create a timeout context
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	var err error
	if u.ID.IsZero() {
		u.ID = bson.NewObjectID()
	}

	_, err = r.collection.InsertOne(ctx, u)
	if err != nil {
		return bson.NilObjectID, err
	}

	return u.ID, nil
}

// Update updates an existing user document by _id (hex string).
func (r *MongoUserRepository) Update(ctx context.Context, u *domain.User) error {
	if u.ID.IsZero() {
		return errors.New("user id is required for update")
	}
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()
	setUpdate := bson.M{
		"name": u.Name,
	}
	update := bson.M{
		"$set": setUpdate,
	}
	if u.Phone != "" {
		setUpdate["phone"] = u.Phone
	} else {
		update["$unset"] = bson.M{"phone": ""}
	}

	res, err := r.collection.UpdateByID(ctx, u.ID, update)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}

// Get finds a user by hex id and returns domain.User.
func (r *MongoUserRepository) Get(ctx context.Context, id bson.ObjectID) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	user := &domain.User{}
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// Delete removes a user by hex id.
func (r *MongoUserRepository) Delete(ctx context.Context, id bson.ObjectID) error {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	res, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}

// List returns paginated users.
func (r *MongoUserRepository) List(ctx context.Context, list *componentsList.List) ([]*domain.User, error) {
	users := []*domain.User{}
	list.RunQuery(ctx, r.collection, map[string]*models.ColumnInfo{
		"id":    models.NewObjectIDColumnInfo("_id"),
		"name":  models.NewTextColumnInfo("name"),
		"phone": models.NewTextColumnInfo("phone"),
	}, "-id", &users)
	return users, nil
}
