package repository

import (
	"context"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"github.com/nileshshrs/infinite-storage/model"
)

type UserRepo struct {
	collection *mongo.Collection
}

func NewUserRepository(col *mongo.Collection) *UserRepo {
	return &UserRepo{collection: col}
}

func (r *UserRepo) Create(user *model.User) error {
	_, err := r.collection.InsertOne(context.Background(), user)
	return err
}

func (r *UserRepo) FindByEmail(email string) (*model.User, error) {
	var user model.User
	err := r.collection.FindOne(context.Background(), bson.M{"email": strings.ToLower(email)}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepo) FindByUsername(username string) (*model.User, error) {
	var user model.User
	err := r.collection.FindOne(context.Background(), bson.M{"username": strings.ToLower(username)}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
