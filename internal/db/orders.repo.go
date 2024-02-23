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
	OrdersCollection = "purchaseOrders"
	DefaultPageSize  = 100
)

var (
	ErrInvalidPurchaseOrderId = errors.New("invalid purchase order id")
	ErrFailedToCreateOrder    = errors.New("failed to create order")
)

// OrdersDataService  - Added for mocking purpose
type OrdersDataService interface {
	Create(ctx context.Context, purchaseOrder *types.Order) (string, error)
	Update(ctx context.Context, purchaseOrder *types.Order) (int64, error)
	GetAll(ctx context.Context) (interface{}, error)
	GetById(ctx context.Context, id string) (interface{}, error)
	DeleteById(ctx context.Context, id string) (int64, error)
}

func NewOrdersRepo(db MongoDatabase) *OrdersRepo {
	iDBSvc := &OrdersRepo{
		collection: db.Collection(OrdersCollection),
	}
	return iDBSvc
}

// OrdersRepo - Implements OrdersDataService
type OrdersRepo struct {
	collection *mongo.Collection
}

func (ordDataSvc *OrdersRepo) Create(ctx context.Context, po *types.Order) (string, error) {
	if vErr := validate(ordDataSvc.collection); vErr != nil {
		return "", vErr
	}
	if !po.ID.IsZero() {
		return "", ErrInvalidPurchaseOrderId
	}
	po.UpdatedAt = util.CurrentISOTime()

	result, err := ordDataSvc.collection.InsertOne(ctx, po)
	if err != nil {
		log.Err(err).Msg("error occurred while creating order")
		return "", ErrFailedToCreateOrder
	}
	log.Info().Msgf("Inserted a single document: %s", result.InsertedID.(primitive.ObjectID).Hex())
	return result.InsertedID.(primitive.ObjectID).Hex(), nil
}

// Update - Create and Update can be merged using upsert, but this is to demonstrate CRUD rest API so ...
func (ordDataSvc *OrdersRepo) Update(ctx context.Context, po *types.Order) (int64, error) {
	if vErr := validate(ordDataSvc.collection); vErr != nil {
		return 0, vErr
	}

	if primitive.ObjectID.IsZero(po.ID) || !primitive.IsValidObjectID(po.ID.Hex()) {
		return 0, errors.New("invalid request")
	}

	po.UpdatedAt = util.CurrentISOTime()

	opts := options.Update().SetUpsert(true)
	filter := bson.D{primitive.E{Key: "_id", Value: po.ID}}
	update := bson.D{primitive.E{Key: "$set", Value: po}}
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

func (ordDataSvc *OrdersRepo) GetAll(ctx context.Context) (interface{}, error) {
	if vErr := validate(ordDataSvc.collection); vErr != nil {
		return nil, vErr
	}

	filter := bson.M{}
	findOptions := options.Find()
	findOptions.SetLimit(DefaultPageSize)

	cursor, err := ordDataSvc.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, err
	}
	var results []types.Order
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	return &results, nil
}

func (ordDataSvc *OrdersRepo) GetById(ctx context.Context, id string) (interface{}, error) {
	if vErr := validate(ordDataSvc.collection); vErr != nil {
		return nil, vErr
	}

	docID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("bad request")
	}
	filter := bson.D{primitive.E{Key: "_id", Value: docID}}

	var result types.Order
	err2 := ordDataSvc.collection.FindOne(ctx, filter).Decode(&result)
	if err2 != nil {
		if errors.Is(err2, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err2
	}

	return &result, nil
}

func (ordDataSvc *OrdersRepo) DeleteById(ctx context.Context, id string) (int64, error) {
	if vErr := validate(ordDataSvc.collection); vErr != nil {
		return 0, vErr
	}

	docID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return 0, errors.New("bad request")
	}
	filter := bson.D{primitive.E{Key: "_id", Value: docID}}

	res, err2 := ordDataSvc.collection.DeleteOne(ctx, filter)
	if err2 != nil {
		return 0, err2
	}

	return res.DeletedCount, nil
}

func validate(collection *mongo.Collection) error {
	if collection == nil {
		return errors.New("collection is not defined")
	}
	return nil
}
