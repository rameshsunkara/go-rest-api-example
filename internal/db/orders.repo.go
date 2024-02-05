package db

import (
	"context"
	"errors"

	"github.com/rameshsunkara/go-rest-api-example/internal/util"
	"github.com/rs/zerolog/log"

	"github.com/rameshsunkara/go-rest-api-example/internal/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	OrdersCollection = "purchaseorders"
	PageSize         = 100
)

// OrdersDataService  TODO: Strong type method definitions
type OrdersDataService interface {
	Create(ctx context.Context, purchaseOrder interface{}) (*mongo.InsertOneResult, error)
	Update(ctx context.Context, purchaseOrder interface{}) (int64, error)
	GetAll(ctx context.Context) (interface{}, error)
	GetById(ctx context.Context, id string) (interface{}, error)
	DeleteById(ctx context.Context, id string) (int64, error)
}

func NewOrderDataService(db MongoDatabase) OrdersDataService {
	iDBSvc := &ordersRepo{
		collection: db.Collection(OrdersCollection),
	}
	return iDBSvc
}

// ordersRepo - Implements OrdersDataService
type ordersRepo struct {
	collection *mongo.Collection
}

func (ordDataSvc *ordersRepo) Create(ctx context.Context, po interface{}) (*mongo.InsertOneResult, error) {
	if vErr := validate(ordDataSvc.collection); vErr != nil {
		return nil, vErr
	}
	purchaseOrder := po.(*types.Order)
	if !purchaseOrder.ID.IsZero() {
		return nil, errors.New("invalid request")
	}
	purchaseOrder.LastUpdatedAt = util.CurrentISOTime()

	result, err := ordDataSvc.collection.InsertOne(ctx, purchaseOrder)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Update - Create and Update can be merged using upsert, but this is to demonstrate CRUD rest API so ...
func (ordDataSvc *ordersRepo) Update(ctx context.Context, po interface{}) (int64, error) {
	if vErr := validate(ordDataSvc.collection); vErr != nil {
		return 0, vErr
	}

	purchaseOrder := po.(*types.Order)
	if primitive.ObjectID.IsZero(purchaseOrder.ID) || !primitive.IsValidObjectID(purchaseOrder.ID.Hex()) {
		return 0, errors.New("invalid request")
	}

	purchaseOrder.LastUpdatedAt = util.CurrentISOTime()

	opts := options.Update().SetUpsert(true)
	filter := bson.D{primitive.E{Key: "_id", Value: purchaseOrder.ID}}
	update := bson.D{primitive.E{Key: "$set", Value: purchaseOrder}}
	result, err := ordDataSvc.collection.UpdateOne(ctx, filter, update, opts)

	if err != nil {
		log.Err(err).Msg("Error occurred while updating order")
	}

	if result.MatchedCount != 0 {
		log.Info().Msg("matched and replaced an existing document")
		return result.MatchedCount, nil
	}

	if result.UpsertedCount != 0 {
		log.Info().Msg("inserted a new order with ID")
		return result.MatchedCount, nil
	}

	return 0, nil
}

func (ordDataSvc *ordersRepo) GetAll(ctx context.Context) (interface{}, error) {
	if vErr := validate(ordDataSvc.collection); vErr != nil {
		return nil, vErr
	}

	filter := bson.M{}
	options := options.Find()
	options.SetLimit(PageSize)

	cursor, err := ordDataSvc.collection.Find(ctx, filter, options)
	if err != nil {
		return nil, err
	}
	var results []types.Order
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	return &results, nil
}

func (ordDataSvc *ordersRepo) GetById(ctx context.Context, id string) (interface{}, error) {
	if vErr := validate(ordDataSvc.collection); vErr != nil {
		return nil, vErr
	}

	docID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("bad request")
	}
	filter := bson.D{primitive.E{Key: "_id", Value: docID}}

	var result types.Order
	error := ordDataSvc.collection.FindOne(ctx, filter).Decode(&result)
	if error != nil {
		if error == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, error
	}

	return &result, nil
}

func (ordDataSvc *ordersRepo) DeleteById(ctx context.Context, id string) (int64, error) {
	if vErr := validate(ordDataSvc.collection); vErr != nil {
		return 0, vErr
	}

	docID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return 0, errors.New("bad request")
	}
	filter := bson.D{primitive.E{Key: "_id", Value: docID}}

	res, error := ordDataSvc.collection.DeleteOne(ctx, filter)
	if error != nil {
		return 0, error
	}

	return res.DeletedCount, nil
}

func validate(collection *mongo.Collection) error {
	if collection == nil {
		return errors.New("collection is not defined")
	}
	return nil
}
