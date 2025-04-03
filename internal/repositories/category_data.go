package repositories

import (
	"app/internal/entities"
	"context"
	"errors"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type categoryMongoRepository struct {
	collection *mongo.Collection
}

func NewCategoryDataRepository(collection *mongo.Collection) CategoryRepository {
	return &categoryMongoRepository{collection: collection}
}

func (m *categoryMongoRepository) CreateCategory(ctx context.Context, userId, categoryName string) error {
	category := entities.Category{
		Id:      uuid.New().String(),
		Name:    categoryName,
		Folders: []entities.Folder{},
	}

	filter := bson.M{"_id": userId}
	update := bson.M{"$push": bson.M{"categories": category}}
	_, err := m.collection.UpdateOne(ctx, filter, update)

	return err
}

func (m *categoryMongoRepository) SelectCategory(ctx context.Context, userId, categoryName string) (entities.Category, error) {
	var category entities.Category
	filer := bson.M{"_id": userId, "categories.name": categoryName}
	opts := options.FindOne().SetProjection(bson.M{"categories.$": 1})

	err := m.collection.FindOne(ctx, filer, opts).Decode(&category)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return entities.Category{}, ErrCategoryNotFound
		}
		return entities.Category{}, err
	}

	return category, nil
}

func (m *categoryMongoRepository) SelectAllCategories(ctx context.Context, userId string) ([]entities.Category, error) {
	var account entities.Account
	filer := bson.M{"_id": userId}
	err := m.collection.FindOne(ctx, filer).Decode(&account)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrCategoryNotFound
		}
		return nil, err
	}

	return account.Categories, nil
}

func (m *categoryMongoRepository) DeleteCategory(ctx context.Context, userId, categoryName string) error {
	filer := bson.M{"_id": userId, "categories.name": categoryName}
	_, err := m.collection.DeleteOne(ctx, filer)

	if err != nil {
		return err
	}

	return nil
}

func (m *categoryMongoRepository) CheckCategoryExists(ctx context.Context, userId, categoryName string) (bool, error) {
	filer := bson.M{"_id": userId, "categories.name": categoryName}
	err := m.collection.FindOne(ctx, filer).Err()
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}
