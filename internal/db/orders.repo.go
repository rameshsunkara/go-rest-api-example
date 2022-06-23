package db

import (
	"context"
	"errors"
	"github.com/rs/zerolog/log"

	"github.com/rameshsunkara/go-rest-api-example/internal/models"
	"github.com/rameshsunkara/go-rest-api-example/pkg/util"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	OrdersCollection = "purchaseorders"
	PageSize         = 100
)

type OrdersDataService struct {
	collection *mongo.Collection
}

func NewOrderDataService(db MongoDatabase) *OrdersDataService {
	iDBSvc := &OrdersDataService{
		collection: db.Collection(OrdersCollection),
	}
	return iDBSvc
}

func (ordDataSvc *OrdersDataService) Create(po interface{}) (*mongo.InsertOneResult, error) {
	if vErr := validate(ordDataSvc.collection); vErr != nil {
		return nil, vErr
	}
	purchaseOrder := po.(*models.Order)
	if !purchaseOrder.ID.IsZero() {
		return nil, errors.New("invalid request")
	}
	purchaseOrder.LastUpdatedAt = util.CurrentISOTime()

	result, err := ordDataSvc.collection.InsertOne(context.TODO(), purchaseOrder)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Update - Create and Update can be merged using upsert, but this is to demonstrate CRUD rest API so ...
func (ordDataSvc *OrdersDataService) Update(po interface{}) (int64, error) {
	if vErr := validate(ordDataSvc.collection); vErr != nil {
		return 0, vErr
	}

	purchaseOrder := po.(*models.Order)
	if primitive.ObjectID.IsZero(purchaseOrder.ID) {
		return 0, errors.New("invalid request")
	}

	purchaseOrder.LastUpdatedAt = util.CurrentISOTime()

	opts := options.Update().SetUpsert(true)
	filter := bson.D{primitive.E{Key: "_id", Value: purchaseOrder.ID}}
	update := bson.D{primitive.E{Key: "$set", Value: purchaseOrder}}
	result, err := ordDataSvc.collection.UpdateOne(context.TODO(), filter, update, opts)

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

func (ordDataSvc *OrdersDataService) GetAll() (interface{}, error) {
	if vErr := validate(ordDataSvc.collection); vErr != nil {
		return nil, vErr
	}

	filter := bson.M{}
	options := options.Find()
	options.SetLimit(PageSize)

	cursor, err := ordDataSvc.collection.Find(context.TODO(), filter, options)
	if err != nil {
		return nil, err
	}
	var results []models.Order
	if err = cursor.All(context.TODO(), &results); err != nil {
		return nil, err
	}

	return &results, nil
}

func (ordDataSvc *OrdersDataService) GetById(id string) (interface{}, error) {
	if vErr := validate(ordDataSvc.collection); vErr != nil {
		return nil, vErr
	}

	isValidId := primitive.IsValidObjectID(id)
	docID, err := primitive.ObjectIDFromHex(id)
	if err != nil || !isValidId {
		return nil, errors.New("bad request")
	}
	filter := bson.D{primitive.E{Key: "_id", Value: docID}}

	var result models.Order
	error := ordDataSvc.collection.FindOne(context.TODO(), filter).Decode(&result)
	if error != nil {
		if error == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, error
	}

	return &result, nil
}

func (ordDataSvc *OrdersDataService) DeleteById(id string) (int64, error) {
	if vErr := validate(ordDataSvc.collection); vErr != nil {
		return 0, vErr
	}

	isValidId := primitive.IsValidObjectID(id)
	docID, err := primitive.ObjectIDFromHex(id)
	if err != nil || !isValidId {
		return 0, errors.New("bad request")
	}
	filter := bson.D{primitive.E{Key: "_id", Value: docID}}

	res, error := ordDataSvc.collection.DeleteOne(context.TODO(), filter)
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
