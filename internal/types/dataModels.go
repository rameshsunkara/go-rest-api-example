package types

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Order struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"order_id"`
	UpdatedAt string             `bson:"updatedAt,omitempty" json:"updatedAt"`
	Products  []Product          `bson:"products,omitempty" json:"products"`
}

type Product struct {
	Name      string `bson:"name,omitempty" json:"name"`
	UpdatedAt string `bson:"updatedAt,omitempty" json:"updatedAt"`
	Price     uint   `bson:"price,omitempty" json:"price"`
	Status    string `bson:"status,omitempty" json:"status"`
	Remarks   string `bson:"remarks,omitempty" json:"remarks"`
}
