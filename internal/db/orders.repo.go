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
	ErrInvalidInitialization = errors.New("invalid initialization")
	ErrInvalidPOIDCreate = errors.New("order id should be empty")
	ErrInvalidPOIDUpdate = errors.New("invalid order id")
	ErrUnexpectedUpdateOrder = errors.New("unexpected error occurred while updating order")
	ErrPOIDNotFound = errors.New("purchase order doesn't exist with given id")
	ErrFailedToCreateOrder  = errors.New("failed to create order")
	ErrInvalidPOIDDelete = errors.New("invalid order id")
	ErrorUnexpectedDeleteOrder = errors.New("unexpected error occurred while deleting order")
)

// OrdersDataService  - Added for tests/mock
type OrdersDataService interface {
	Create(ctx context.Context, purchaseOrder *types.Order) (string, error)
	Update(ctx context.Context, purchaseOrder *types.Order) error
	GetAll(ctx context.Context) (*[]types.Order, error)
	GetById(ctx context.Context, id string) (types.Order, error)
	DeleteById(ctx context.Context, id string) (int64, error)
}

// OrdersRepo - Implements OrdersDataService
type OrdersRepo struct {
	collection *mongo.Collection
}

func NewOrdersRepo(db MongoDatabase) *OrdersRepo {
	iDBSvc := &OrdersRepo{
		collection: db.Collection(OrdersCollection),
	}
	return iDBSvc
}

func (ordDataSvc *OrdersRepo) Create(ctx context.Context, po *types.Order) (string, error) {
	if err := validate(ordDataSvc.collection); err != nil {
		return "", err
	}
	if !po.ID.IsZero() {
		return "", ErrInvalidPOIDCreate
	}
	po.UpdatedAt = util.CurrentISOTime()

	result, err := ordDataSvc.collection.InsertOne(ctx, po)
	if err != nil {
		log.Err(err).Msg("error occurred while creating order")
		return "", ErrFailedToCreateOrder
	}
	log.Info().Msgf("created new order: %s", result.InsertedID.(primitive.ObjectID).Hex())
	return result.InsertedID.(primitive.ObjectID).Hex(), nil
}

func (ordDataSvc *OrdersRepo) Update(ctx context.Context, po *types.Order) error {
	if err := validate(ordDataSvc.collection); err != nil {
		return err
	}

	oID, err := primitive.ObjectIDFromHex(po.ID.Hex())
	if oID.IsZero() || err != nil {
		return ErrInvalidPOIDUpdate
	}

	po.UpdatedAt = util.CurrentISOTime()

	filter := bson.D{primitive.E{Key: "_id", Value: po.ID}}
	update := bson.D{primitive.E{Key: "$set", Value: po}}
	result, err := ordDataSvc.collection.UpdateOne(ctx, filter, update, nil)

	if err != nil {
		log.Err(err).Msg("error occurred while updating order")
		return ErrUnexpectedUpdateOrder
	}

	if result.MatchedCount == 0 {
		log.Info().Msg("order id given for updating the order is not found")
		return ErrPOIDNotFound
	}

	return nil
}

func (ordDataSvc *OrdersRepo) GetAll(ctx context.Context) (*[]types.Order, error) {
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

func (ordDataSvc *OrdersRepo) GetById(ctx context.Context, id string) (*types.Order, error) {
	if err := validate(ordDataSvc.collection); err != nil {
		return nil, err
	}

	oID, err := primitive.ObjectIDFromHex(id)
	if oID.IsZero() || err != nil {
		return nil, ErrInvalidPOIDUpdate
	}

	filter := bson.D{primitive.E{Key: "_id", Value: oID}}

	var result types.Order
	err = ordDataSvc.collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrPOIDNotFound
		}
		return nil, err
	}

	return &result, nil
}

func (ordDataSvc *OrdersRepo) DeleteById(ctx context.Context, id string) (int64, error) {
	if err := validate(ordDataSvc.collection); err != nil {
		return 0, err
	}

	oID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return 0, ErrInvalidPOIDDelete
	}
	filter := bson.D{primitive.E{Key: "_id", Value: oID}}

	res, err := ordDataSvc.collection.DeleteOne(ctx, filter)
	if err != nil {
		return 0, ErrorUnexpectedDeleteOrder
	}

	if res.DeletedCount == 0 {
		return 0, ErrPOIDNotFound
	}
	return res.DeletedCount, nil
}

func validate(collection *mongo.Collection) error {
	if collection == nil {
		return ErrInvalidInitialization
	}
	return nil
}
