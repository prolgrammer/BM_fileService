package repositories

import (
	"app/internal/entities"
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
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
	pipeline := []bson.M{
		{
			"$match": bson.M{
				"_id":             userId,
				"categories.name": categoryName,
			},
		},
		{
			"$unwind": "$categories",
		},
		{
			"$match": bson.M{
				"categories.name": categoryName,
			},
		},
		{
			"$replaceRoot": bson.M{
				"newRoot": "$categories",
			},
		},
	}

	cursor, err := m.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return entities.Category{}, err
	}
	defer cursor.Close(ctx)

	if !cursor.Next(ctx) {
		return entities.Category{}, ErrCategoryNotFound
	}

	var category entities.Category
	if err := cursor.Decode(&category); err != nil {
		return entities.Category{}, err
	}

	fmt.Println(category)

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
	update := bson.M{"$pull": bson.M{"categories": bson.M{"name": categoryName}}}

	result, err := m.collection.UpdateOne(ctx, filer, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return ErrCategoryNotFound
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
