package mongodb

import (
	"context"
	"regexp"
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

var _ repository.Flag = (*FlagRepository)(nil)

// FlagRepository implements repository.Flag interface using mongodb.
type FlagRepository struct {
	db  *mongo.Database
	col *mongo.Collection
}

// FindAll returns a list of flags, based on an optional offset and limit.
func (r *FlagRepository) FindAll(ctx context.Context, search *string, offset, limit *int64) (*flaggio.FlagResults, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "MongoFlagRepository.FindAll")
	defer span.Finish()

	filter := bson.M{}
	if search != nil {
		filter["$or"] = []bson.M{
			{"key": primitive.Regex{Pattern: regexp.QuoteMeta(*search), Options: "i"}},
			{"$text": bson.M{"$search": *search}},
		}
	}
	cursor, err := r.col.Find(ctx, filter, &options.FindOptions{
		Skip:      offset,
		Limit:     limit,
		Sort:      bson.M{"key": 1},
		Collation: &options.Collation{Locale: "en"},
	})
	if err != nil {
		return nil, err
	}

	var flags []*flaggio.Flag
	for cursor.Next(ctx) {
		var f flagModel
		// decode the document
		if err := cursor.Decode(&f); err != nil {
			return nil, err
		}
		flags = append(flags, f.asFlag())
	}

	// check if the cursor encountered any errors while iterating
	if err := cursor.Err(); err != nil {
		return nil, err
	}

	// get the total results
	total, err := r.col.CountDocuments(ctx, filter)
	if err != nil {
		return nil, err
	}

	return &flaggio.FlagResults{
		Flags: flags,
		Total: int(total),
	}, nil
}

// FindByID returns a flag that has a given ID.
func (r *FlagRepository) FindByID(ctx context.Context, idHex string) (*flaggio.Flag, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "MongoFlagRepository.FindByID")
	defer span.Finish()

	id, err := primitive.ObjectIDFromHex(idHex)
	if err != nil {
		return nil, err
	}
	filter := bson.M{"_id": id}

	var f flagModel
	if err := r.col.FindOne(ctx, filter).Decode(&f); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.NotFound("flag")
		}
		return nil, err
	}
	return f.asFlag(), nil
}

// FindByKey returns a flag that has a given key.
func (r *FlagRepository) FindByKey(ctx context.Context, key string) (*flaggio.Flag, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "MongoFlagRepository.FindByKey")
	defer span.Finish()

	// filter for the flag key
	filter := bson.M{"key": key}

	var f flagModel
	if err := r.col.FindOne(ctx, filter).Decode(&f); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.NotFound("flag")
		}
		return nil, err
	}
	return f.asFlag(), nil
}

// Create creates a new flag.
func (r *FlagRepository) Create(ctx context.Context, f flaggio.NewFlag) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "MongoFlagRepository.Create")
	defer span.Finish()

	id := primitive.NewObjectID()
	_, err := r.col.InsertOne(ctx, &flagModel{
		ID:          id,
		CreatedAt:   time.Now(),
		Key:         f.Key,
		Name:        f.Name,
		Description: f.Description,
		Enabled:     false,
		Version:     1,
		Variants:    []variantModel{},
		Rules:       []flagRuleModel{},
	})
	if err != nil {
		return "", err
	}
	return id.Hex(), nil
}

// Update updates a flag.
func (r *FlagRepository) Update(ctx context.Context, idHex string, f flaggio.UpdateFlag) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "MongoFlagRepository.Update")
	defer span.Finish()

	id, err := primitive.ObjectIDFromHex(idHex)
	if err != nil {
		return err
	}
	mods := bson.M{
		"updatedAt": time.Now(),
	}
	if f.Key != nil {
		mods["key"] = *f.Key
	}
	if f.Name != nil {
		mods["name"] = *f.Name
	}
	if f.Description != nil {
		mods["description"] = *f.Description
	}
	if f.Enabled != nil {
		mods["enabled"] = *f.Enabled
	}
	if f.DefaultVariantWhenOn != nil {
		oid, err := primitive.ObjectIDFromHex(*f.DefaultVariantWhenOn)
		if err != nil {
			return err
		}
		mods["defaultVariantWhenOn"] = oid
	}
	if f.DefaultVariantWhenOff != nil {
		oid, err := primitive.ObjectIDFromHex(*f.DefaultVariantWhenOff)
		if err != nil {
			return err
		}
		mods["defaultVariantWhenOff"] = oid
	}
	if len(mods) == 0 {
		return errors.BadRequest("nothing to update")
	}
	filter := bson.M{"_id": id}
	update := bson.M{
		"$set": mods,
		"$inc": bson.M{"version": 1},
	}
	res, err := r.col.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if res.ModifiedCount == 0 {
		return errors.NotFound("flag")
	}
	return nil
}

// Delete deletes a flag.
func (r *FlagRepository) Delete(ctx context.Context, idHex string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "MongoFlagRepository.Delete")
	defer span.Finish()

	id, err := primitive.ObjectIDFromHex(idHex)
	if err != nil {
		return err
	}
	res, err := r.col.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return errors.NotFound("flag")
	}
	return nil
}

// NewFlagRepository returns a new flag repository that uses mongodb as underlying storage.
// It also creates all needed indexes, if they don't yet exist.
func NewFlagRepository(ctx context.Context, db *mongo.Database) (repository.Flag, error) {
	col := db.Collection("flags")
	_, err := col.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "key", Value: 1}},
			Options: options.Index().SetUnique(true).SetBackground(false),
		},
		{
			Keys:    bson.D{{Key: "variants._id", Value: 1}},
			Options: options.Index().SetUnique(true).SetSparse(true).SetBackground(false),
		},
		{
			Keys:    bson.D{{Key: "variants.key", Value: 1}},
			Options: options.Index().SetUnique(true).SetSparse(true).SetBackground(false),
		},
		{
			Keys:    bson.D{{Key: "rules._id", Value: 1}},
			Options: options.Index().SetUnique(true).SetSparse(true).SetBackground(false),
		},
		{
			Keys:    bson.D{{Key: "rules.distributions._id", Value: 1}},
			Options: options.Index().SetUnique(true).SetSparse(true).SetBackground(false),
		},
		{
			Keys:    bson.D{{Key: "rules.constraints._id", Value: 1}},
			Options: options.Index().SetUnique(true).SetSparse(true).SetBackground(false),
		},
		{
			Keys:    bson.D{{Key: "name", Value: "text"}},
			Options: options.Index().SetBackground(false),
		},
	})
	if err != nil {
		return nil, err
	}
	return &FlagRepository{
		db:  db,
		col: col,
	}, nil
}
