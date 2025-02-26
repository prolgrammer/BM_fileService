package repositories

import (
	"app/internal/entities"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongoRepository struct {
	collection *mongo.Collection
}

func NewCategoryDataRepository(collection *mongo.Collection) CategoryRepository {
	return &mongoRepository{collection: collection}
}

func (m mongoRepository) CreateCategory(ctx context.Context, userId, categoryName string) error {
	category := entities.Category{
		Name:    categoryName,
		UserId:  userId,
		Folders: []entities.Folder{},
	}

	_, err := m.collection.InsertOne(ctx, category)

	return err
}

func (m mongoRepository) SelectCategory(ctx context.Context, userId, categoryName string) (entities.Category, error) {
	var category entities.Category
	filer := bson.M{"user_id": userId, "name": categoryName}

	err := m.collection.FindOne(ctx, filer).Decode(&category)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return entities.Category{}, ErrCategoryNotFound
		}
		return category, err
	}

	return category, nil
}

func (m mongoRepository) SelectAllCategories(ctx context.Context, userId string) ([]entities.Category, error) {
	filer := bson.M{"user_id": userId}
	cursor, err := m.collection.Find(ctx, filer)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var categories []entities.Category
	err = cursor.All(ctx, &categories)
	if err != nil {
		return nil, err
	}

	return categories, nil
}

func (m mongoRepository) DeleteCategory(ctx context.Context, userId, category string) error {
	filer := bson.M{"user_id": userId, "name": category}
	_, err := m.collection.DeleteOne(ctx, filer)

	if err != nil {
		return err
	}

	return nil
}
