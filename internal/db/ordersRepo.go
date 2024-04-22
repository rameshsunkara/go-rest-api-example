package db

import (
	"context"
	"errors"
	"time"

	"github.com/rameshsunkara/go-rest-api-example/internal/logger"
	"github.com/rameshsunkara/go-rest-api-example/internal/models"
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
	ErrInvalidPOIDDelete     = errors.New("invalid order id")
	ErrUnexpectedDeleteOrder = errors.New("unexpected error occurred while deleting order")
)

// OrdersDataService  - Added for tests/mock.
type OrdersDataService interface {
	Create(ctx context.Context, purchaseOrder *models.Order) (string, error)
	Update(ctx context.Context, purchaseOrder *models.Order) error
	GetAll(ctx context.Context) (*[]models.Order, error)
	GetByID(ctx context.Context, id string) (*models.Order, error)
	DeleteByID(ctx context.Context, id string) (int64, error)
}

// OrdersRepo - Implements OrdersDataService.
type OrdersRepo struct {
	collection *mongo.Collection
	logger     logger.AppLogger
}

func NewOrdersRepo(db MongoDatabase) *OrdersRepo {
	iDBSvc := &OrdersRepo{
		collection: db.Collection(OrdersCollection),
	}
	return iDBSvc
}

func (o *OrdersRepo) Create(ctx context.Context, po *models.Order) (string, error) {
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
	o.logger.Info().Str("order_id", result.InsertedID.(primitive.ObjectID).Hex()).Msg("created new order")
	return result.InsertedID.(primitive.ObjectID).Hex(), nil
}

func (o *OrdersRepo) Update(ctx context.Context, po *models.Order) error {
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

func (o *OrdersRepo) GetAll(ctx context.Context) (*[]models.Order, error) {
	if vErr := validate(o.collection); vErr != nil {
		return nil, vErr
	}

	filter := bson.M{}
	findOptions := options.Find()
	findOptions.SetLimit(DefaultPageSize)

	cursor, err := o.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, err
	}
	var results []models.Order
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	return &results, nil
}

func (o *OrdersRepo) GetByID(ctx context.Context, id string) (*models.Order, error) {
	if err := validate(o.collection); err != nil {
		return nil, err
	}

	oID, err := primitive.ObjectIDFromHex(id)
	if oID.IsZero() || err != nil {
		return nil, ErrInvalidPOIDUpdate
	}

	filter := bson.D{primitive.E{Key: "_id", Value: oID}}

	var result models.Order
	err = o.collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrPOIDNotFound
		}
		return nil, err
	}

	return &result, nil
}

func (o *OrdersRepo) DeleteByID(ctx context.Context, id string) (int64, error) {
	if err := validate(o.collection); err != nil {
		return 0, err
	}

	oID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return 0, ErrInvalidPOIDDelete
	}
	filter := bson.D{primitive.E{Key: "_id", Value: oID}}

	res, err := o.collection.DeleteOne(ctx, filter)
	if err != nil {
		return 0, ErrUnexpectedDeleteOrder
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
