package db

import (
	"context"
	"errors"
	"time"

	"github.com/bogdanutanu/go-rest-api-example/internal/models/data"
	"github.com/bogdanutanu/go-rest-api-example/pkg/logger"
	"github.com/bogdanutanu/go-rest-api-example/pkg/mongodb"
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
	ErrUnexpectedGetOrder    = errors.New("unexpected error occurred while fetching order")
	ErrInvalidID             = errors.New("failed to assert inserted ID as ObjectID")
)

// OrdersDataService defines the interface for order data operations.
type OrdersDataService interface {
	Create(ctx context.Context, purchaseOrder *data.Order) (string, error)
	Update(ctx context.Context, purchaseOrder *data.Order) error
	GetAll(ctx context.Context, limit int64) (*[]data.Order, error)
	GetByID(ctx context.Context, id primitive.ObjectID) (*data.Order, error)
	DeleteByID(ctx context.Context, id primitive.ObjectID) error
}

// OrdersRepo implements OrdersDataService using MongoDB.
type OrdersRepo struct {
	collection *mongo.Collection
	logger     logger.Logger
}

// NewOrdersRepo creates a new OrdersRepo.
func NewOrdersRepo(lgr logger.Logger, db mongodb.MongoDatabase) (*OrdersRepo, error) {
	if lgr == nil || db == nil {
		return nil, errors.New("missing required inputs to create OrdersRepo")
	}
	return &OrdersRepo{
		collection: db.Collection(OrdersCollection),
		logger:     lgr,
	}, nil
}

// Create inserts a new order into the collection.
func (o *OrdersRepo) Create(ctx context.Context, po *data.Order) (string, error) {
	if err := validateCollection(o.collection); err != nil {
		return "", err
	}
	if !po.ID.IsZero() {
		return "", ErrInvalidPOIDCreate
	}

	result, err := o.collection.InsertOne(ctx, po)
	if err != nil {
		o.logger.Error().Err(err).Msg("failed to create order")
		return "", ErrFailedToCreateOrder
	}
	insertedID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return "", ErrInvalidID
	}
	o.logger.Info().Str("orderId", insertedID.Hex()).Msg("created new order")
	return insertedID.Hex(), nil
}

// Update modifies an existing order.
func (o *OrdersRepo) Update(ctx context.Context, po *data.Order) error {
	if err := validateCollection(o.collection); err != nil {
		return err
	}
	oID, err := primitive.ObjectIDFromHex(po.ID.Hex())
	if err != nil || oID.IsZero() {
		return ErrInvalidPOIDUpdate
	}
	po.UpdatedAt = time.Now()

	filter := bson.D{{Key: "_id", Value: po.ID}}
	update := bson.D{{Key: "$set", Value: po}}
	result, err := o.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		o.logger.Error().Err(err).Msg("failed to update order")
		return ErrUnexpectedUpdateOrder
	}
	if result.MatchedCount == 0 {
		o.logger.Info().Msg("order id for update not found")
		return ErrPOIDNotFound
	}
	return nil
}

// GetAll retrieves all orders up to the specified limit.
func (o *OrdersRepo) GetAll(ctx context.Context, limit int64) (*[]data.Order, error) {
	if err := validateCollection(o.collection); err != nil {
		return nil, err
	}
	findOptions := options.Find().SetLimit(limit)
	cursor, err := o.collection.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		o.logger.Error().Err(err).Msg("failed to find orders")
		return nil, ErrUnexpectedGetOrder
	}
	var results []data.Order
	if err = cursor.All(ctx, &results); err != nil {
		o.logger.Error().Err(err).Msg("failed to decode orders")
		return nil, ErrUnexpectedGetOrder
	}
	return &results, nil
}

// GetByID retrieves an order by its ObjectID.
func (o *OrdersRepo) GetByID(ctx context.Context, oID primitive.ObjectID) (*data.Order, error) {
	if err := validateCollection(o.collection); err != nil {
		return nil, err
	}
	filter := bson.D{{Key: "_id", Value: oID}}
	var result data.Order
	err := o.collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrPOIDNotFound
		}
		o.logger.Error().Err(err).Msg("failed to get order by id")
		return nil, ErrUnexpectedGetOrder
	}
	return &result, nil
}

// DeleteByID removes an order by its ObjectID.
func (o *OrdersRepo) DeleteByID(ctx context.Context, id primitive.ObjectID) error {
	if err := validateCollection(o.collection); err != nil {
		return err
	}
	filter := bson.D{{Key: "_id", Value: id}}
	res, err := o.collection.DeleteOne(ctx, filter)
	if err != nil {
		o.logger.Error().Err(err).Msg("failed to delete order")
		return ErrUnexpectedDeleteOrder
	}
	if res.DeletedCount == 0 {
		return ErrPOIDNotFound
	}
	return nil
}

// validateCollection checks if the collection is initialized.
func validateCollection(collection *mongo.Collection) error {
	if collection == nil {
		return ErrInvalidInitialization
	}
	return nil
}
