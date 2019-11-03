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

var _ repository.Segment = SegmentRepository{}

type SegmentRepository struct {
	db  *mongo.Database
	col *mongo.Collection
}

func (r SegmentRepository) FindAll(ctx context.Context, offset, limit *int64) ([]*flaggio.Segment, error) {
	cursor, err := r.col.Find(ctx, bson.M{}, &options.FindOptions{
		Skip:  offset,
		Limit: limit,
		Sort:  bson.M{"_id": 1},
	})
	if err != nil {
		return nil, err
	}

	var segments []*flaggio.Segment
	for cursor.Next(ctx) {
		var f segmentModel
		// decode the document
		if err := cursor.Decode(&f); err != nil {
			return nil, err
		}
		segments = append(segments, f.asSegment())
	}

	// check if the cursor encountered any errors while iterating
	if err := cursor.Err(); err != nil {
		return nil, err
	}
	return segments, nil
}

func (r SegmentRepository) FindByID(ctx context.Context, idHex string) (*flaggio.Segment, error) {
	id, err := primitive.ObjectIDFromHex(idHex)
	if err != nil {
		return nil, err
	}
	filter := bson.M{"_id": id}

	var f segmentModel
	if err := r.col.FindOne(ctx, filter).Decode(&f); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.NotFound("segment")
		}
		return nil, err
	}
	return f.asSegment(), nil
}

func (r SegmentRepository) FindByKey(ctx context.Context, key string) (*flaggio.Segment, error) {
	// filter for the segment key
	filter := bson.M{"key": key}

	var f segmentModel
	if err := r.col.FindOne(ctx, filter).Decode(&f); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.NotFound("segment")
		}
		return nil, err
	}
	return f.asSegment(), nil
}

func (r SegmentRepository) Create(ctx context.Context, f flaggio.NewSegment) (*flaggio.Segment, error) {
	id := primitive.NewObjectID()
	_, err := r.col.InsertOne(ctx, &segmentModel{
		ID:          id,
		CreatedAt:   time.Now(),
		Name:        f.Name,
		Description: f.Description,
		Rules:       []segmentRuleModel{},
	})
	if err != nil {
		return nil, err
	}
	return r.FindByID(ctx, id.Hex())
}

func (r SegmentRepository) Update(ctx context.Context, idHex string, f flaggio.UpdateSegment) error {
	id, err := primitive.ObjectIDFromHex(idHex)
	if err != nil {
		return err
	}
	mods := bson.M{
		"updatedAt": time.Now(),
	}
	if f.Name != nil {
		mods["name"] = *f.Name
	}
	if f.Description != nil {
		mods["description"] = *f.Description
	}
	if len(mods) == 0 {
		return errors.ErrNothingToUpdate
	}
	filter := bson.M{"_id": id}
	update := bson.M{"$set": mods}
	res, err := r.col.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if res.ModifiedCount == 0 {
		return errors.NotFound("segment")
	}
	return nil
}

func (r SegmentRepository) Delete(ctx context.Context, idHex string) error {
	id, err := primitive.ObjectIDFromHex(idHex)
	if err != nil {
		return err
	}
	res, err := r.col.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return errors.NotFound("segment")
	}
	return nil
}

func NewMongoSegmentRepository(ctx context.Context, db *mongo.Database) (*SegmentRepository, error) {
	col := db.Collection("segments")
	_, err := col.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys:    bson.M{"rules._id": 1},
			Options: options.Index().SetUnique(true).SetSparse(true).SetBackground(false),
		},
		{
			Keys:    bson.M{"rules.constraints._id": 1},
			Options: options.Index().SetUnique(true).SetSparse(true).SetBackground(false),
		},
	})
	if err != nil {
		return nil, err
	}
	return &SegmentRepository{
		db:  db,
		col: col,
	}, nil
}
