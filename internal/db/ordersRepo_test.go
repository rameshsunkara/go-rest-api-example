package db_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/bogdanutanu/go-rest-api-example/internal/db"
	"github.com/bogdanutanu/go-rest-api-example/internal/db/mocks"
	"github.com/bogdanutanu/go-rest-api-example/internal/models/data"
	"github.com/bogdanutanu/go-rest-api-example/internal/utilities"
	"github.com/bogdanutanu/go-rest-api-example/pkg/logger"
	"github.com/bogdanutanu/go-rest-api-example/pkg/mongodb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

var testLgr = logger.New("debug", os.Stdout)

func TestNewOrdersRepo(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		lgr     logger.Logger
		db      mongodb.MongoDatabase
		wantErr bool
	}{
		{
			name:    "success",
			lgr:     testLgr,
			db:      &mocks.MockMongoDataBase{},
			wantErr: false,
		},
		{
			name:    "nil logger",
			lgr:     nil,
			db:      &mocks.MockMongoDataBase{},
			wantErr: true,
		},
		{
			name:    "nil db",
			lgr:     testLgr,
			db:      nil,
			wantErr: true,
		},
		{
			name:    "nil logger and db",
			lgr:     nil,
			db:      nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			repo, err := db.NewOrdersRepo(tt.lgr, tt.db)
			if tt.wantErr {
				require.Error(t, err)
				assert.Nil(t, repo)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, repo)
			}
		})
	}
}

func TestValidate(t *testing.T) {
	t.Parallel()
	ds, err := db.NewOrdersRepo(testLgr, &mocks.MockMongoDataBase{})
	require.NoError(t, err)

	testCases := []struct {
		name     string
		testFunc func() error
		wantErr  error
	}{
		{
			name: "Create with invalid initialization",
			testFunc: func() error {
				_, cErr := ds.Create(context.Background(), &data.Order{})
				return cErr
			},
			wantErr: db.ErrInvalidInitialization,
		},
		{
			name: "GetAll with invalid initialization",
			testFunc: func() error {
				_, gErr := ds.GetAll(context.Background(), 10)
				return gErr
			},
			wantErr: db.ErrInvalidInitialization,
		},
		{
			name: "GetByID with invalid initialization",
			testFunc: func() error {
				_, gErr := ds.GetByID(context.Background(), primitive.NewObjectID())
				return gErr
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
			t.Parallel()
			tFuncErr := tc.testFunc()
			require.Error(t, tFuncErr)
			assert.Equal(t, tc.wantErr, tFuncErr)
		})
	}
}

func TestOrdersRepoCreate(t *testing.T) {
	t.Parallel()
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

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
				TotalAmount: utilities.CalculateTotalAmount([]data.Product{{Name: "Product 1", Price: 10.0, Quantity: 2}}),
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
				TotalAmount: utilities.CalculateTotalAmount([]data.Product{{Name: "Product 1", Price: 10.0, Quantity: 2}}),
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
			repo, err := db.NewOrdersRepo(testLgr, mt.DB)
			if err != nil {
				t.Errorf("failed to create repo")
				return
			}
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

func TestOrdersRepoUpdate(t *testing.T) {
	t.Parallel()
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

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
				TotalAmount: utilities.CalculateTotalAmount([]data.Product{{Name: "Product 1", Price: 10.0, Quantity: 2}}),
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
				TotalAmount: utilities.CalculateTotalAmount([]data.Product{{Name: "Product 1", Price: 10.0, Quantity: 2}}),
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
				TotalAmount: utilities.CalculateTotalAmount([]data.Product{{Name: "Product 1", Price: 10.0, Quantity: 2}}),
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
				TotalAmount: utilities.CalculateTotalAmount([]data.Product{{Name: "Product 1", Price: 10.0, Quantity: 2}}),
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
			repo, repoErr := db.NewOrdersRepo(testLgr, mt.DB)
			if repoErr != nil {
				t.Errorf("failed to create repo")
				return
			}
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

func TestOrdersRepoGetByID(t *testing.T) {
	t.Parallel()
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

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
			repo, repoErr := db.NewOrdersRepo(testLgr, mt.DB)
			if repoErr != nil {
				t.Errorf("failed to create repo")
				return
			}
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

func TestOrdersRepoDeleteByID(t *testing.T) {
	t.Parallel()
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

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
			repo, repoErr := db.NewOrdersRepo(testLgr, mt.DB)
			if repoErr != nil {
				t.Errorf("failed to create repo")
				return
			}
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

func TestOrdersRepoGetAll(t *testing.T) {
	t.Parallel()
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

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
			repo, repoErr := db.NewOrdersRepo(testLgr, mt.DB)
			if repoErr != nil {
				t.Errorf("failed to create repo")
				return
			}
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
