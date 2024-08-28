package external

import (
	"github.com/rameshsunkara/go-rest-api-example/internal/models/data"
)

// APIError represents the structure of an API error response.
type APIError struct {
	HTTPStatusCode int    `json:"httpStatusCode"`
	Message        string `json:"message"`
	DebugID        string `json:"debugId"`
	ErrorCode      string `json:"errorCode"`
}

// OrderInput represents the structure of input for creating or updating an order.
type OrderInput struct {
	Products []ProductInput `json:"products" binding:"required"`
}

// ProductInput represents the structure of input for creating or updating a product.
type ProductInput struct {
	Name     string  `json:"name" binding:"required"`
	Price    float64 `json:"price" binding:"required"`
	Quantity uint64  `json:"quantity" binding:"required"`
}

// Order represents the structure of an order.
type Order struct {
	ID          string             `json:"orderId"`
	Version     int64              `json:"version"`
	CreatedAt   string             `json:"createdAt"`
	UpdatedAt   string             `json:"updatedAt"`
	Products    []data.Product     `json:"products"`
	User        string             `json:"user"`
	TotalAmount float64            `json:"totalAmount"`
	Status      data.OrderStatus   `json:"status"`
	Updates     []data.OrderUpdate `json:"updates"`
}
