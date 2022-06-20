package users

import (
	"context"
	"errors"
	"github.com/adrianbanasiak/go-users-service/internal/shared"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	ErrInsertFailed = errors.New("insert query failed")
	ErrQueryFailed  = errors.New("query failed")
	ErrNotFound     = errors.New("not found")
)

func NewMongoRepository(log shared.Logger, db *mongo.Database) *MongoRepository {
	collection := db.Collection("users")
	return &MongoRepository{log: log, collection: collection}
}

type MongoRepository struct {
	log        shared.Logger
	collection *mongo.Collection
}

func (r *MongoRepository) ChangeEmail(ctx context.Context, user User) error {
	_, err := r.collection.UpdateByID(ctx, user.ID, bson.D{
		{"$set", bson.M{"email": user.Email, "updated_at": user.UpdatedAt}},
	})

	if err != nil {
		r.log.Errorw("failed to change user email in collection", "error", err)
		return ErrQueryFailed
	}

	return nil
}

func (r *MongoRepository) FindByID(ctx context.Context, userID uuid.UUID) (User, error) {
	var u User
	err := r.collection.FindOne(ctx, bson.D{{"_id", userID}}).Decode(&u)
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			return User{}, ErrNotFound
		}

		r.log.Errorw("failed to fetch user from collection",
			"error", err)
		return User{}, ErrQueryFailed
	}

	return u, nil
}

func (r *MongoRepository) FindByEmail(ctx context.Context, email string) (User, error) {
	var u User
	sr := r.collection.FindOne(ctx, bson.D{{"email", email}})

	err := sr.Decode(&u)
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			return User{}, ErrNotFound
		}

		r.log.Errorw("failed to fetch user from collection",
			"error", err)
		return User{}, ErrQueryFailed
	}

	return u, nil
}

func (r *MongoRepository) Create(ctx context.Context, user User) (User, error) {
	_, err := r.collection.InsertOne(ctx, user)
	if err != nil {
		r.log.Errorw("failed to persist user in collection",
			"userID", user.ID,
			"error", err)

		return User{}, ErrInsertFailed
	}

	return user, nil
}

func (r *MongoRepository) Delete(ctx context.Context, ID uuid.UUID) error {
	res, err := r.collection.DeleteOne(ctx, bson.M{"_id": ID})
	if err != nil {
		r.log.Errorw("failed to delete user from collection",
			"userID", ID,
			"error", err)
		return ErrQueryFailed
	}

	if res.DeletedCount != 1 {
		return ErrNotFound
	}

	return nil
}

func (r *MongoRepository) FindPaginated(ctx context.Context, query FindUserQuery, page, items int) ([]User, error) {
	i := int64(items)
	skip := int64(page*items - items)

	var filter bson.M
	if query.Country != "" {
		filter = bson.M{
			"country_code": query.Country,
		}
	}
	cur, err := r.collection.Find(ctx, filter, &options.FindOptions{Limit: &i, Skip: &skip})
	if err != nil {
		r.log.Errorw("failed to list users in collection", "error", err)
		return nil, ErrQueryFailed
	}

	res := make([]User, 0)
	err = cur.All(ctx, &res)
	if err != nil {
		r.log.Errorw("failed to list users in collection", "error", err)
		return nil, ErrQueryFailed
	}

	return res, nil
}
