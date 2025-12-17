package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// User represents an application user stored in MongoDB
type User struct {
	ID       primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Username string             `json:"username" bson:"username"`
	Password string             `json:"-" bson:"password"` // Hashed password
	Role     string             `json:"role" bson:"role"` // "admin" or "user"
}

// Used when the user registers
type RegisterInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Used when the user logs in
type LoginInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}
