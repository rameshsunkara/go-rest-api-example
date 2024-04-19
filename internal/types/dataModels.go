package types

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// OrderStatus represents the status of an order.
type OrderStatus int

const (
	OrderPending OrderStatus = iota
	OrderProcessing
	OrderShipped
	OrderDelivered
	OrderCancelled
)

// OrderUpdate represents an update to an order.
type OrderUpdate struct {
	UpdatedAt string `bson:"updatedAt,omitempty" json:"updatedAt"`
	Notes     string `bson:"notes,omitempty" json:"notes"`
	HandledBy string `bson:"handledBy,omitempty" json:"handledBy"`
}

// Order represents an e-commerce order.
type Order struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"orderId"`
	CreatedAt   string             `bson:"createdAt,omitempty" json:"createdAt"`
	UpdatedAt   string             `bson:"updatedAt,omitempty" json:"updatedAt"`
	Products    []Product          `bson:"products,omitempty" json:"products"`
	User        string             `bson:"user,omitempty" json:"user"`
	TotalAmount float64            `bson:"totalAmount,omitempty" json:"totalAmount"`
	Status      OrderStatus        `bson:"status,omitempty" json:"status"`
	Updates     []OrderUpdate      `bson:"updates,omitempty" json:"updates"`
}

type Product struct {
	Name      string `bson:"name,omitempty" json:"name"`
	UpdatedAt string `bson:"updatedAt,omitempty" json:"updatedAt"`
	Price     uint   `bson:"price,omitempty" json:"price"`
	Status    string `bson:"status,omitempty" json:"status"`
	Remarks   string `bson:"remarks,omitempty" json:"remarks"`
}
