package repositories

import (
	"app/internal/entities"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type accountMongoRepository struct {
	collection *mongo.Collection
}

func NewAccountMongoRepository(collection *mongo.Collection) AccountRepository {
	return &accountMongoRepository{collection: collection}
}

func (a *accountMongoRepository) CreateAccount(ctx context.Context, userId string) error {
	account := entities.Account{
		Id:         userId,
		Categories: []entities.Category{},
	}

	_, err := a.collection.InsertOne(ctx, account)
	return err
}

func (a *accountMongoRepository) SelectAccount(ctx context.Context, userId string) (entities.Account, error) {
	var account entities.Account
	filter := bson.D{{"_id", userId}}
	err := a.collection.FindOne(ctx, filter).Decode(&account)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return entities.Account{}, ErrAccountNotFound
		}
		return entities.Account{}, err
	}
	return account, nil
}
