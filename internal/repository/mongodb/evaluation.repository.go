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

var _ repository.Evaluation = (*EvaluationRepository)(nil)

// EvaluationRepository implements repository.Flag interface using mongodb.
type EvaluationRepository struct {
	db  *mongo.Database
	col *mongo.Collection
}

// FindAllByUserID returns all previous flag evaluations for a given user ID.
func (r *EvaluationRepository) FindAllByUserID(ctx context.Context, userID string, search *string, offset, limit *int64) (*flaggio.EvaluationResults, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "MongoEvaluationRepository.FindAllByUserID")
	defer span.Finish()

	filter := bson.M{"userId": userID}
	if search != nil {
		filter["flagKey"] = primitive.Regex{Pattern: regexp.QuoteMeta(*search), Options: "i"}
	}

	cursor, err := r.col.Find(ctx, filter, &options.FindOptions{
		Skip:      offset,
		Limit:     limit,
		Sort:      bson.M{"flagKey": 1},
		Collation: &options.Collation{Locale: "en"},
	})
	if err != nil {
		return nil, err
	}

	var evals flaggio.EvaluationList
	for cursor.Next(ctx) {
		var e evaluationModel
		// decode the document
		if err := cursor.Decode(&e); err != nil {
			return nil, err
		}
		evals = append(evals, e.asEvaluation())
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

	return &flaggio.EvaluationResults{
		Evaluations: evals,
		Total:       int(total),
	}, nil
}

// FindAllByReqHash returns all previous flag evaluations for a given request hash.
func (r *EvaluationRepository) FindAllByReqHash(ctx context.Context, reqHash string) (flaggio.EvaluationList, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "MongoEvaluationRepository.FindAllByReqHash")
	defer span.Finish()

	filter := bson.M{"requestHash": reqHash}
	cursor, err := r.col.Find(ctx, filter, &options.FindOptions{
		Sort:      bson.M{"flagKey": 1},
		Collation: &options.Collation{Locale: "en"},
	})
	if err != nil {
		return nil, err
	}

	var evals flaggio.EvaluationList
	for cursor.Next(ctx) {
		var e evaluationModel
		// decode the document
		if err := cursor.Decode(&e); err != nil {
			return nil, err
		}
		evals = append(evals, e.asEvaluation())
	}

	// check if the cursor encountered any errors while iterating
	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return evals, nil
}

// FindByReqHashAndFlagKey returns a previous flag evaluation for a given request hash and flag key.
func (r *EvaluationRepository) FindByReqHashAndFlagKey(ctx context.Context, reqHash, flagKey string) (*flaggio.Evaluation, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "MongoEvaluationRepository.FindByReqHashAndFlagKey")
	defer span.Finish()

	filter := bson.M{"requestHash": reqHash, "flagKey": flagKey}

	var e evaluationModel
	if err := r.col.FindOne(ctx, filter).Decode(&e); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.NotFound("evaluation")
		}
		return nil, err
	}
	return e.asEvaluation(), nil
}

// FindByID returns a previous flag evaluation by its ID.
func (r *EvaluationRepository) FindByID(ctx context.Context, idHex string) (*flaggio.Evaluation, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "MongoEvaluationRepository.FindByID")
	defer span.Finish()

	id, err := primitive.ObjectIDFromHex(idHex)
	if err != nil {
		return nil, err
	}

	var e evaluationModel
	if err := r.col.FindOne(ctx, bson.M{"_id": id}).Decode(&e); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.NotFound("evaluation")
		}
		return nil, err
	}
	return e.asEvaluation(), nil
}

// ReplaceOne creates or replaces one evaluation for a user ID.
func (r *EvaluationRepository) ReplaceOne(ctx context.Context, userID string, eval *flaggio.Evaluation) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "MongoEvaluationRepository.ReplaceOne")
	defer span.Finish()

	// delete current evaluation
	flgID, err := primitive.ObjectIDFromHex(eval.FlagID)
	if err != nil {
		return err
	}
	_, err = r.col.DeleteOne(ctx, bson.M{"userId": userID, "flagId": flgID})
	if err != nil {
		return err
	}

	// prepare list of evaluations to insert
	_, err = r.col.InsertOne(ctx, &evaluationModel{
		ID:          primitive.NewObjectID(),
		FlagID:      flgID,
		FlagKey:     eval.FlagKey,
		FlagVersion: eval.FlagVersion,
		RequestHash: eval.RequestHash,
		UserID:      userID,
		Value:       eval.Value,
		CreatedAt:   time.Now(),
	})
	return err
}

// ReplaceAll creates or replaces evaluations for a combination of user and request hash.
func (r *EvaluationRepository) ReplaceAll(ctx context.Context, userID, reqHash string, evals flaggio.EvaluationList) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "MongoEvaluationRepository.ReplaceAll")
	defer span.Finish()

	// get list of evaluated flag ids
	evalsToDelete := make([]primitive.ObjectID, len(evals))
	for idx, eval := range evals {
		flgID, err := primitive.ObjectIDFromHex(eval.FlagID)
		if err != nil {
			return err
		}
		evalsToDelete[idx] = flgID
	}
	// delete current
	_, err := r.col.DeleteMany(ctx, bson.M{"userId": userID, "flagId": bson.M{"$in": evalsToDelete}})
	if err != nil {
		return err
	}

	// prepare list of evaluations to insert
	evalsToInsert := make([]interface{}, len(evals))
	for idx, eval := range evals {
		flgID, err := primitive.ObjectIDFromHex(eval.FlagID)
		if err != nil {
			return err
		}
		evalsToInsert[idx] = &evaluationModel{
			ID:          primitive.NewObjectID(),
			FlagID:      flgID,
			FlagKey:     eval.FlagKey,
			FlagVersion: eval.FlagVersion,
			RequestHash: eval.RequestHash,
			UserID:      userID,
			Value:       eval.Value,
			CreatedAt:   time.Now(),
		}
	}
	_, err = r.col.InsertMany(ctx, evalsToInsert)
	return err
}

// Delete deletes evaluations for a user.
func (r *EvaluationRepository) DeleteAllByUserID(ctx context.Context, userID string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "MongoEvaluationRepository.DeleteAllByUserID")
	defer span.Finish()

	_, err := r.col.DeleteMany(ctx, bson.M{"userId": userID})
	return err
}

// DeleteByID deletes an evaluation by its ID.
func (r *EvaluationRepository) DeleteByID(ctx context.Context, idHex string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "MongoEvaluationRepository.DeleteByID")
	defer span.Finish()

	id, err := primitive.ObjectIDFromHex(idHex)
	if err != nil {
		return err
	}

	_, err = r.col.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

// NewEvaluationRepository returns a new evaluation repository that uses mongodb as underlying storage.
// It also creates all needed indexes, if they don't yet exist.
func NewEvaluationRepository(ctx context.Context, db *mongo.Database) (repository.Evaluation, error) {
	col := db.Collection("evaluations")
	_, err := col.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "userId", Value: 1}, {Key: "flagId", Value: 1}},
			Options: options.Index().SetUnique(true).SetBackground(false),
		},
	})
	if err != nil {
		return nil, err
	}
	return &EvaluationRepository{
		db:  db,
		col: col,
	}, nil
}
