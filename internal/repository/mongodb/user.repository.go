package mongodb

import (
	"context"
	"regexp"
	"strings"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/victorkt/flaggio/internal/errors"
	"github.com/victorkt/flaggio/internal/flaggio"
	"github.com/victorkt/flaggio/internal/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var _ repository.User = (*UserRepository)(nil)

// UserRepository implements repository.Flag interface using mongodb.
type UserRepository struct {
	db  *mongo.Database
	col *mongo.Collection
}

// FindAll returns a list of users, based on an optional offset and limit.
func (r *UserRepository) FindAll(ctx context.Context, search *string, offset, limit *int64) (*flaggio.UserResults, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "MongoUserRepository.FindAll")
	defer span.Finish()

	filter := bson.M{}
	if search != nil {
		filter["_id"] = primitive.Regex{Pattern: regexp.QuoteMeta(*search), Options: "i"}
	}
	cursor, err := r.col.Find(ctx, filter, &options.FindOptions{
		Skip:  offset,
		Limit: limit,
		Sort:  bson.M{"_id": 1},
	})
	if err != nil {
		return nil, err
	}

	var users []*flaggio.User
	for cursor.Next(ctx) {
		var u userModel
		// decode the document
		if err := cursor.Decode(&u); err != nil {
			return nil, err
		}
		u.Context = sanitizeUserContextPrefixKey(u.Context, "%", "$")
		users = append(users, u.asUser())
	}

	// check if the cursor encountered any errors while iterating
	if err := cursor.Err(); err != nil {
		return nil, err
	}

	total, err := r.col.CountDocuments(ctx, filter)
	if err != nil {
		return nil, err
	}

	return &flaggio.UserResults{
		Users: users,
		Total: int(total),
	}, nil
}

// FindByID returns a user by its id.
func (r *UserRepository) FindByID(ctx context.Context, id string) (*flaggio.User, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "MongoUserRepository.FindByID")
	defer span.Finish()

	filter := bson.M{"_id": id}

	var u userModel
	if err := r.col.FindOne(ctx, filter).Decode(&u); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.NotFound("user")
		}
		return nil, err
	}
	u.Context = sanitizeUserContextPrefixKey(u.Context, "%", "$")
	return u.asUser(), nil
}

// Replace creates or updates a user.
func (r *UserRepository) Replace(ctx context.Context, userID string, userCtx flaggio.UserContext) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "MongoUserRepository.Replace")
	defer span.Finish()

	_, err := r.col.UpdateOne(ctx, bson.M{"_id": userID}, bson.M{"$set": &userModel{
		ID:        userID,
		Context:   sanitizeUserContextPrefixKey(userCtx, "$", "%"),
		UpdatedAt: time.Now(),
	}}, options.Update().SetUpsert(true))
	return err
}

// Delete deletes a user.
func (r *UserRepository) Delete(ctx context.Context, userID string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "MongoUserRepository.Delete")
	defer span.Finish()

	_, err := r.col.DeleteOne(ctx, bson.M{"_id": userID})
	return err
}

// sanitizeUserContextPrefixKey replaces all user context keys that starts with `old`
// and replaces them with `new`. this is used because mongodb doesn't allow fields
// starting with $ character, so we replace it with another character for storage
func sanitizeUserContextPrefixKey(userCtx flaggio.UserContext, oldStr, newStr string) flaggio.UserContext {
	usrCtx := make(flaggio.UserContext, len(userCtx))
	for key, value := range userCtx {
		if strings.HasPrefix(key, oldStr) {
			key = strings.Replace(key, oldStr, newStr, 1)
		}
		usrCtx[key] = value
	}
	return usrCtx
}

// NewUserRepository returns a new user repository that uses mongodb as underlying storage.
// It also creates all needed indexes, if they don't yet exist.
func NewUserRepository(ctx context.Context, db *mongo.Database) (repository.User, error) {
	col := db.Collection("users")
	return &UserRepository{
		db:  db,
		col: col,
	}, nil
}
