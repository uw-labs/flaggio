package mongodb

import (
	"context"
	"time"

	"github.com/victorkohl/flaggio/internal/errors"
	"github.com/victorkohl/flaggio/internal/flaggio"
	"github.com/victorkohl/flaggio/internal/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var _ repository.Flag = FlagRepository{}

type FlagRepository struct {
	db  *mongo.Database
	col *mongo.Collection
}

func (r FlagRepository) FindAll(ctx context.Context, offset, limit *int64) ([]*flaggio.Flag, error) {
	cursor, err := r.col.Find(ctx, bson.M{}, &options.FindOptions{
		Skip:  offset,
		Limit: limit,
		Sort:  bson.M{"_id": 1},
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
	return flags, nil
}

func (r FlagRepository) FindByID(ctx context.Context, idHex string) (*flaggio.Flag, error) {
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

func (r FlagRepository) FindByKey(ctx context.Context, key string) (*flaggio.Flag, error) {
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

func (r FlagRepository) Create(ctx context.Context, f flaggio.NewFlag) (*flaggio.Flag, error) {
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
		return nil, err
	}
	return r.FindByID(ctx, id.Hex())
}

func (r FlagRepository) Update(ctx context.Context, idHex string, f flaggio.UpdateFlag) error {
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
	if len(mods) == 0 {
		return errors.ErrNothingToUpdate
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

func (r FlagRepository) Delete(ctx context.Context, idHex string) error {
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

func NewMongoFlagRepository(ctx context.Context, db *mongo.Database) (*FlagRepository, error) {
	col := db.Collection("flags")
	_, err := col.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys:    bson.M{"variants._id": 1},
			Options: options.Index().SetUnique(true).SetBackground(false),
		},
		{
			Keys:    bson.M{"rules._id": 1},
			Options: options.Index().SetUnique(true).SetBackground(false),
		},
		{
			Keys:    bson.M{"rules.distributions._id": 1},
			Options: options.Index().SetUnique(true).SetBackground(false),
		},
		{
			Keys:    bson.M{"rules.constraints._id": 1},
			Options: options.Index().SetUnique(true).SetBackground(false),
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
