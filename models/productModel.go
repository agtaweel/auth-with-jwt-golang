package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Product struct {
	ID          primitive.ObjectID `bson:"_id"`
	Name        *string            `json:"name" validate: "required,min=2,max=100"`
	Description *string            `json:"description" validate: "required,min=2,max=500"`
	Price       *float32           `json:"price" validate: "required,min=0.1"`
	Buyer       *string            `json:"buyer" vaildate:"required,min=5"`
	Phone       *string            `json:"phone" validate:"required,min=6"`
	Created_at  *time.Time         `json:"created_at"`
	Updated_at  *time.Time         `json:"updated_at"`
}
