package db_test

import (
	"context"
	"testing"
	"time"

	"github.com/rameshsunkara/go-rest-api-example/internal/db"
	"github.com/rameshsunkara/go-rest-api-example/internal/db/mocks"
	"github.com/rameshsunkara/go-rest-api-example/internal/logger"
	"github.com/rameshsunkara/go-rest-api-example/internal/models"
	"github.com/rameshsunkara/go-rest-api-example/internal/models/data"
	"github.com/rameshsunkara/go-rest-api-example/internal/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

func TestNewOrderDataService(t *testing.T) {
	lgr := logger.Setup(models.ServiceEnv{Name: "test"})
	ds := db.NewOrdersRepo(&mocks.MockMongoDataBase{}, lgr)
	assert.Implements(t, (*db.OrdersDataService)(nil), ds)
}

func TestValidate(t *testing.T) {
	lgr := logger.Setup(models.ServiceEnv{Name: "test"})
	ds := db.NewOrdersRepo(&mocks.MockMongoDataBase{}, lgr)

	testCases := []struct {
		name     string
		testFunc func() error
		wantErr  error
	}{
		{
			name: "Create with invalid initialization",
			testFunc: func() error {
				_, err := ds.Create(context.Background(), &data.Order{})
				return err
			},
			wantErr: db.ErrInvalidInitialization,
		},
		{
			name: "GetAll with invalid initialization",
			testFunc: func() error {
				_, err := ds.GetAll(context.Background(), 10)
				return err
			},
			wantErr: db.ErrInvalidInitialization,
		},
		{
			name: "GetByID with invalid initialization",
			testFunc: func() error {
				_, err := ds.GetByID(context.Background(), primitive.NewObjectID())
				return err
			},
			wantErr: db.ErrInvalidInitialization,
		},
		{
			name: "Update with invalid initialization",
			testFunc: func() error {
				return ds.Update(context.Background(), &data.Order{})
			},
			wantErr: db.ErrInvalidInitialization,
		},
		{
			name: "DeleteByID with invalid initialization",
			testFunc: func() error {
				return ds.DeleteByID(context.Background(), primitive.NewObjectID())
			},
			wantErr: db.ErrInvalidInitialization,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.testFunc()
			require.Error(t, err)
			assert.Equal(t, tc.wantErr, err)
		})
	}
}

func TestOrdersRepo_Create(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	lgr := logger.Setup(models.ServiceEnv{Name: "test"})

	tests := []struct {
		name    string
		order   *data.Order
		mock    func(mt *mtest.T)
		wantErr error
	}{
		{
			name: "Success",
			order: &data.Order{
				Version:     1,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
				Products:    []data.Product{{Name: "Product 1", Price: 10.0, Quantity: 2}},
				User:        "test@example.com",
				Status:      data.OrderPending,
				TotalAmount: util.CalculateTotalAmount([]data.Product{{Name: "Product 1", Price: 10.0, Quantity: 2}}),
			},
			mock: func(mt *mtest.T) {
				mt.AddMockResponses(mtest.CreateSuccessResponse())
			},
			wantErr: nil,
		},
		{
			name: "InvalidData",
			order: &data.Order{
				ID: primitive.NewObjectID(),
			},
			mock:    func(_ *mtest.T) {},
			wantErr: db.ErrInvalidPOIDCreate,
		},
		{
			name: "InsertError",
			order: &data.Order{
				Version:     1,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
				Products:    []data.Product{{Name: "Product 1", Price: 10.0, Quantity: 2}},
				User:        "test@example.com",
				Status:      data.OrderPending,
				TotalAmount: util.CalculateTotalAmount([]data.Product{{Name: "Product 1", Price: 10.0, Quantity: 2}}),
			},
			mock: func(mt *mtest.T) {
				mt.AddMockResponses(mtest.CreateWriteErrorsResponse(mtest.WriteError{Code: 11000}))
			},
			wantErr: db.ErrFailedToCreateOrder,
		},
	}

	for _, tt := range tests {
		mt.Run(tt.name, func(mt *mtest.T) {
			tt.mock(mt)
			repo := db.NewOrdersRepo(mt.DB, lgr)
			resultID, err := repo.Create(context.TODO(), tt.order)
			if tt.wantErr != nil {
				require.Error(t, err)
				assert.Equal(t, tt.wantErr, err)
			} else {
				require.NoError(t, err)
				assert.NotEmpty(t, resultID)
			}
		})
	}
}

func TestOrdersRepo_Update(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	lgr := logger.Setup(models.ServiceEnv{Name: "test"})

	tests := []struct {
		name    string
		order   *data.Order
		mock    func(mt *mtest.T)
		wantErr error
	}{
		{
			name: "Success",
			order: &data.Order{
				ID:          primitive.NewObjectID(),
				Version:     1,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
				Products:    []data.Product{{Name: "Product 1", Price: 10.0, Quantity: 2}},
				User:        "test@example.com",
				Status:      data.OrderPending,
				TotalAmount: util.CalculateTotalAmount([]data.Product{{Name: "Product 1", Price: 10.0, Quantity: 2}}),
			},
			mock: func(mt *mtest.T) {
				mt.AddMockResponses(mtest.CreateSuccessResponse(bson.E{Key: "n", Value: 1}))
			},
			wantErr: nil,
		},
		{
			name: "NonExistingID",
			order: &data.Order{
				ID:          primitive.NewObjectID(),
				Version:     1,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
				Products:    []data.Product{{Name: "Product 1", Price: 10.0, Quantity: 2}},
				User:        "test@example.com",
				Status:      data.OrderPending,
				TotalAmount: util.CalculateTotalAmount([]data.Product{{Name: "Product 1", Price: 10.0, Quantity: 2}}),
			},
			mock: func(mt *mtest.T) {
				mt.AddMockResponses(mtest.CreateSuccessResponse())
			},
			wantErr: db.ErrPOIDNotFound,
		},
		{
			name: "ZeroID",
			order: &data.Order{
				Version:     1,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
				Products:    []data.Product{{Name: "Product 1", Price: 10.0, Quantity: 2}},
				User:        "test@example.com",
				Status:      data.OrderPending,
				TotalAmount: util.CalculateTotalAmount([]data.Product{{Name: "Product 1", Price: 10.0, Quantity: 2}}),
			},
			mock: func(mt *mtest.T) {
				mt.AddMockResponses(mtest.CreateSuccessResponse())
			},
			wantErr: db.ErrInvalidPOIDUpdate,
		},
		{
			name: "UpdateError",
			order: &data.Order{
				ID:          primitive.NewObjectID(),
				Version:     1,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
				Products:    []data.Product{{Name: "Product 1", Price: 10.0, Quantity: 2}},
				User:        "test@example.com",
				Status:      data.OrderPending,
				TotalAmount: util.CalculateTotalAmount([]data.Product{{Name: "Product 1", Price: 10.0, Quantity: 2}}),
			},
			mock: func(mt *mtest.T) {
				mt.AddMockResponses(mtest.CreateWriteErrorsResponse(mtest.WriteError{Code: 11000}))
			},
			wantErr: db.ErrUnexpectedUpdateOrder,
		},
	}

	for _, tt := range tests {
		mt.Run(tt.name, func(mt *mtest.T) {
			tt.mock(mt)
			repo := db.NewOrdersRepo(mt.DB, lgr)
			err := repo.Update(context.TODO(), tt.order)
			if tt.wantErr != nil {
				require.Error(t, err)
				assert.Equal(t, tt.wantErr, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestOrdersRepo_GetByID(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	lgr := logger.Setup(models.ServiceEnv{Name: "test"})
	testOID := primitive.NewObjectID()
	tests := []struct {
		name    string
		orderID primitive.ObjectID
		mock    func(mt *mtest.T)
		wantErr error
	}{
		{
			name:    "Success",
			orderID: testOID,
			mock: func(mt *mtest.T) {
				mt.AddMockResponses(mtest.CreateCursorResponse(1, "ordersdb.orders", mtest.FirstBatch, bson.D{
					{Key: "_id", Value: testOID},
					{Key: "user", Value: "test@example.com"},
				}))
			},
			wantErr: nil,
		},
		{
			name:    "InvalidID",
			orderID: primitive.NilObjectID,
			mock: func(mt *mtest.T) {
				mt.AddMockResponses(mtest.CreateCursorResponse(0, "ordersdb.orders", mtest.FirstBatch))
			},
			wantErr: db.ErrPOIDNotFound,
		},
		{
			name:    "FetchError",
			orderID: primitive.NilObjectID,
			mock: func(mt *mtest.T) {
				mt.AddMockResponses(mtest.CreateWriteErrorsResponse(mtest.WriteError{Code: 11000}))
			},
			wantErr: db.ErrUnexpectedGetOrder,
		},
	}

	for _, tt := range tests {
		mt.Run(tt.name, func(mt *mtest.T) {
			tt.mock(mt)
			repo := db.NewOrdersRepo(mt.DB, lgr)
			result, err := repo.GetByID(context.TODO(), tt.orderID)
			if tt.wantErr != nil {
				require.Error(t, err)
				assert.Equal(t, tt.wantErr, err)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.orderID, result.ID)
			}
		})
	}
}

func TestOrdersRepo_DeleteByID(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	lgr := logger.Setup(models.ServiceEnv{Name: "test"})

	tests := []struct {
		name    string
		orderID primitive.ObjectID
		mock    func(mt *mtest.T)
		wantErr error
	}{
		{
			name:    "Success",
			orderID: primitive.NewObjectID(),
			mock: func(mt *mtest.T) {
				mt.AddMockResponses(mtest.CreateSuccessResponse(bson.E{Key: "n", Value: 1}))
			},
			wantErr: nil,
		},
		{
			name:    "InvalidID",
			orderID: primitive.NilObjectID,
			mock: func(mt *mtest.T) {
				mt.AddMockResponses(mtest.CreateSuccessResponse())
			},
			wantErr: db.ErrPOIDNotFound,
		},
		{
			name:    "DeleteError",
			orderID: primitive.NilObjectID,
			mock: func(mt *mtest.T) {
				mt.AddMockResponses(mtest.CreateWriteErrorsResponse(mtest.WriteError{Code: 11000}))
			},
			wantErr: db.ErrUnexpectedDeleteOrder,
		},
	}

	for _, tt := range tests {
		mt.Run(tt.name, func(mt *mtest.T) {
			tt.mock(mt)
			repo := db.NewOrdersRepo(mt.DB, lgr)
			err := repo.DeleteByID(context.TODO(), tt.orderID)
			if tt.wantErr != nil {
				require.Error(t, err)
				assert.Equal(t, tt.wantErr, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestOrdersRepo_GetAll(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	lgr := logger.Setup(models.ServiceEnv{Name: "test"})
	oCollName := "ordersdb.orders"
	tests := []struct {
		name    string
		limit   int64
		mock    func(mt *mtest.T)
		wantErr error
		wantLen int
	}{
		{
			name:  "Success",
			limit: 10,
			mock: func(mt *mtest.T) {
				find := mtest.CreateCursorResponse(1, oCollName, mtest.FirstBatch, bson.D{
					{Key: "_id", Value: primitive.NewObjectID()},
					{Key: "user", Value: "test@example.com"},
				})
				killCursors := mtest.CreateCursorResponse(0, oCollName, mtest.NextBatch)
				mt.AddMockResponses(find, killCursors)
			},
			wantErr: nil,
			wantLen: 1,
		},
		{
			name:  "NoData",
			limit: 10,
			mock: func(mt *mtest.T) {
				find := mtest.CreateCursorResponse(1, oCollName, mtest.FirstBatch)
				mt.AddMockResponses(find)
			},
			wantErr: nil,
			wantLen: 0,
		},
	}

	for _, tt := range tests {
		mt.Run(tt.name, func(mt *mtest.T) {
			tt.mock(mt)
			repo := db.NewOrdersRepo(mt.DB, lgr)
			results, err := repo.GetAll(context.TODO(), tt.limit)
			if tt.wantErr != nil {
				require.Error(t, err)
				assert.Equal(t, tt.wantErr, err)
			} else {
				require.NoError(t, err)
				assert.Len(t, *results, tt.wantLen)
			}
		})
	}
}
