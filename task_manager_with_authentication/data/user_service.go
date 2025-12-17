package data

import (
	"context"
	"errors"
	"time"

	"task_manager/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var UserCollection *mongo.Collection

// Initialize the user collection after DB connects
func InitUserCollection() {
	UserCollection = Client.Database("task_manager_db").Collection("users")
}

// -------------------------------
// USER UTILITIES
// -------------------------------

// HashPassword encrypts the plain password
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// CheckPassword compares raw vs hashed password
func CheckPassword(hashedPwd, plainPwd string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPwd), []byte(plainPwd))
	return err == nil
}

// -------------------------------
// USER CRUD
// -------------------------------

// CreateUser adds a new user to MongoDB
func CreateUser(input models.RegisterInput) (models.User, error) {

	// Check if username already exists
	_, err := GetUserByUsername(input.Username)
	if err == nil {
		return models.User{}, errors.New("username already exists")
	}

	// Hash password
	hashedPwd, err := HashPassword(input.Password)
	if err != nil {
		return models.User{}, err
	}

	// Determine if this is the first user → make admin
	role := "user"
	count, _ := CountUsers()
	if count == 0 {
		role = "admin"
	}

	user := models.User{
		ID:       primitive.NewObjectID(),
		Username: input.Username,
		Password: hashedPwd,
		Role:     role,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = UserCollection.InsertOne(ctx, user)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

// GetUserByUsername finds a user by username
func GetUserByUsername(username string) (models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user models.User
	err := UserCollection.FindOne(ctx, bson.M{"username": username}).Decode(&user)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

// GetUserByID finds a user given its ObjectID
func GetUserByID(id primitive.ObjectID) (models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user models.User
	err := UserCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

// CountUsers returns how many users exist
func CountUsers() (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	count, err := UserCollection.CountDocuments(ctx, bson.M{})
	return count, err
}

// PromoteUser upgrades a user to admin role
func PromoteUser(id primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	update := bson.M{
		"$set": bson.M{"role": "admin"},
	}

	result, err := UserCollection.UpdateOne(ctx, bson.M{"_id": id}, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return errors.New("user not found")
	}

	return nil
}
