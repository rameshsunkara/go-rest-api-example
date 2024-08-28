package data

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// OrderStatus represents the status of an order.
type OrderStatus string

const (
	OrderPending    OrderStatus = "OrderPending"
	OrderProcessing OrderStatus = "OrderProcessing"
	OrderShipped    OrderStatus = "OrderShipped"
	OrderDelivered  OrderStatus = "OrderDelivered"
	OrderCancelled  OrderStatus = "OrderCancelled"
)

// Order represents the structure of an order.
type Order struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"orderId"`
	Version     int64              `json:"version" bson:"version"`
	CreatedAt   time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt   time.Time          `json:"updatedAt" bson:"updatedAt"`
	Products    []Product          `json:"products" bson:"products"`
	User        string             `json:"user" bson:"user"`
	TotalAmount float64            `json:"totalAmount" bson:"totalAmount"`
	Status      OrderStatus        `json:"status" bson:"status"`
	Updates     []OrderUpdate      `json:"updates" bson:"updates"`
}

// OrderUpdate represents the structure of an order update.
type OrderUpdate struct {
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
	Notes     string    `json:"notes" bson:"notes"`
	HandledBy string    `json:"handledBy" bson:"handledBy"`
}

// Product represents the structure of a product.
type Product struct {
	Name      string    `json:"name" bson:"name"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
	Price     float64   `json:"price" bson:"price"`
	Status    string    `json:"status" bson:"status"`
	Remarks   string    `json:"remarks" bson:"remarks"`
	Quantity  uint64    `json:"quantity"`
}
