package db

import (
	"context"
	"errors"
	"time"

	"github.com/rameshsunkara/go-rest-api-example/internal/logger"
	"github.com/rameshsunkara/go-rest-api-example/internal/models/data"
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
	ErrInvalidPOIDCreate     = errors.New("order id should be empty")
	ErrInvalidPOIDUpdate     = errors.New("invalid order id")
	ErrUnexpectedUpdateOrder = errors.New("unexpected error occurred while updating order")
	ErrPOIDNotFound          = errors.New("purchase order doesn't exist with given id")
	ErrFailedToCreateOrder   = errors.New("failed to create order")
	ErrUnexpectedDeleteOrder = errors.New("unexpected error occurred while deleting order")
)

// OrdersDataService  - Added for tests/mock.
type OrdersDataService interface {
	Create(ctx context.Context, purchaseOrder *data.Order) (string, error)
	Update(ctx context.Context, purchaseOrder *data.Order) error
	GetAll(ctx context.Context, limit int64) (*[]data.Order, error)
	GetByID(ctx context.Context, id primitive.ObjectID) (*data.Order, error)
	DeleteByID(ctx context.Context, id primitive.ObjectID) error
}

// OrdersRepo - Implements OrdersDataService.
type OrdersRepo struct {
	collection *mongo.Collection
	logger     *logger.AppLogger
}

func NewOrdersRepo(db MongoDatabase, lgr *logger.AppLogger) *OrdersRepo {
	iDBSvc := &OrdersRepo{
		collection: db.Collection(OrdersCollection),
		logger:     lgr,
	}
	return iDBSvc
}

func (o *OrdersRepo) Create(ctx context.Context, po *data.Order) (string, error) {
	if err := validate(o.collection); err != nil {
		return "", err
	}
	if !po.ID.IsZero() {
		return "", ErrInvalidPOIDCreate
	}

	result, err := o.collection.InsertOne(ctx, po)
	if err != nil {
		o.logger.Error().Err(err).Msg("error occurred while creating order")
		return "", ErrFailedToCreateOrder
	}
	o.logger.Info().Str("orderId", result.InsertedID.(primitive.ObjectID).Hex()).Msg("created new order")
	return result.InsertedID.(primitive.ObjectID).Hex(), nil
}

func (o *OrdersRepo) Update(ctx context.Context, po *data.Order) error {
	if err := validate(o.collection); err != nil {
		return err
	}

	oID, err := primitive.ObjectIDFromHex(po.ID.Hex())
	if oID.IsZero() || err != nil {
		return ErrInvalidPOIDUpdate
	}

	po.UpdatedAt = time.Now()

	filter := bson.D{primitive.E{Key: "_id", Value: po.ID}}
	update := bson.D{primitive.E{Key: "$set", Value: po}}
	result, err := o.collection.UpdateOne(ctx, filter, update, nil)

	if err != nil {
		o.logger.Error().Err(err).Msg("error occurred while updating order")
		return ErrUnexpectedUpdateOrder
	}

	if result.MatchedCount == 0 {
		o.logger.Info().Msg("order id given for updating the order is not found")
		return ErrPOIDNotFound
	}

	return nil
}

func (o *OrdersRepo) GetAll(ctx context.Context, limit int64) (*[]data.Order, error) {
	if vErr := validate(o.collection); vErr != nil {
		return nil, vErr
	}

	filter := bson.M{}
	findOptions := options.Find()
	findOptions.SetLimit(limit)

	cursor, err := o.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, err
	}
	var results []data.Order
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	return &results, nil
}

func (o *OrdersRepo) GetByID(ctx context.Context, oID primitive.ObjectID) (*data.Order, error) {
	if err := validate(o.collection); err != nil {
		return nil, err
	}

	filter := bson.D{primitive.E{Key: "_id", Value: oID}}
	var result data.Order
	err := o.collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrPOIDNotFound
		}
		return nil, err
	}

	return &result, nil
}

func (o *OrdersRepo) DeleteByID(ctx context.Context, id primitive.ObjectID) error {
	if err := validate(o.collection); err != nil {
		return err
	}

	filter := bson.D{primitive.E{Key: "_id", Value: id}}

	res, err := o.collection.DeleteOne(ctx, filter)
	if err != nil {
		return ErrUnexpectedDeleteOrder
	}

	if res.DeletedCount == 0 {
		return ErrPOIDNotFound
	}
	return nil
}

func validate(collection *mongo.Collection) error {
	if collection == nil {
		return ErrInvalidInitialization
	}
	return nil
}
