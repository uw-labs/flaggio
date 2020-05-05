package mongodb

import (
	"context"
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

var _ repository.Segment = (*SegmentRepository)(nil)

// SegmentRepository implements repository.Segment interface using mongodb.
type SegmentRepository struct {
	db  *mongo.Database
	col *mongo.Collection
}

// FindAll returns a list of segments, based on an optional offset and limit.
func (r *SegmentRepository) FindAll(ctx context.Context, offset, limit *int64) ([]*flaggio.Segment, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "MongoSegmentRepository.FindAll")
	defer span.Finish()

	cursor, err := r.col.Find(ctx, bson.M{}, &options.FindOptions{
		Skip:      offset,
		Limit:     limit,
		Sort:      bson.M{"name": 1},
		Collation: &options.Collation{Locale: "en"},
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

// FindByID returns a segment that has a given ID.
func (r *SegmentRepository) FindByID(ctx context.Context, idHex string) (*flaggio.Segment, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "MongoSegmentRepository.FindByID")
	defer span.Finish()

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

// Create creates a new segment.
func (r *SegmentRepository) Create(ctx context.Context, f flaggio.NewSegment) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "MongoSegmentRepository.Create")
	defer span.Finish()

	id := primitive.NewObjectID()
	_, err := r.col.InsertOne(ctx, &segmentModel{
		ID:          id,
		CreatedAt:   time.Now(),
		Name:        f.Name,
		Description: f.Description,
		Rules:       []segmentRuleModel{},
	})
	if err != nil {
		return "", err
	}
	return id.Hex(), nil
}

// Update updates a segment.
func (r *SegmentRepository) Update(ctx context.Context, idHex string, f flaggio.UpdateSegment) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "MongoSegmentRepository.Update")
	defer span.Finish()

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
		return errors.BadRequest("nothing to update")
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

// Delete deletes a segment.
func (r *SegmentRepository) Delete(ctx context.Context, idHex string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "MongoSegmentRepository.Delete")
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
		return errors.NotFound("segment")
	}
	return nil
}

// NewSegmentRepository returns a new segment repository that uses mongodb as underlying storage.
// It also creates all needed indexes, if they don't yet exist.
func NewSegmentRepository(ctx context.Context, db *mongo.Database) (repository.Segment, error) {
	col := db.Collection("segments")
	_, err := col.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "rules._id", Value: 1}},
			Options: options.Index().SetUnique(true).SetSparse(true).SetBackground(false),
		},
		{
			Keys:    bson.D{{Key: "rules.constraints._id", Value: 1}},
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
