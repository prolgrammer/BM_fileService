package repositories

import (
	"app/internal/entities"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type folderMongoRepository struct {
	collection *mongo.Collection
}

func NewFolderMongoRepository(collection *mongo.Collection) FolderRepository {
	return &folderMongoRepository{
		collection: collection,
	}
}

func (f folderMongoRepository) CreateFolder(ctx context.Context, userId, categoryName, folderName string) error {
	filter := bson.M{"user_id": userId, "name": categoryName}

	newFolder := entities.Folder{
		Name:  folderName,
		Files: []entities.File{},
	}

	update := bson.M{"$push": bson.M{"folders": newFolder}}

	_, err := f.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}

func (f folderMongoRepository) SelectFolder(ctx context.Context, userId, categoryName, folderName string) (entities.Folder, error) {
	var category entities.Category
	filter := bson.M{"user_id": userId, "name": categoryName, "folders.name": folderName}
	projection := bson.M{"folders.$": 1}

	err := f.collection.FindOne(ctx, filter, options.FindOne().SetProjection(projection)).Decode(&category)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return entities.Folder{}, ErrFolderNotFound
		}
		return entities.Folder{}, err
	}

	if len(category.Folders) == 0 {
		return entities.Folder{}, ErrFolderNotFound
	}

	return category.Folders[0], nil
}

func (f folderMongoRepository) SelectFolders(ctx context.Context, userId, categoryName string) ([]entities.Folder, error) {
	var category entities.Category
	filter := bson.M{"user_id": userId, "name": categoryName}

	err := f.collection.FindOne(ctx, filter).Decode(&category)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrCategoryNotFound
		}
	}

	return category.Folders, nil
}

func (f folderMongoRepository) DeleteFolder(ctx context.Context, userId, categoryName, folderName string) error {
	filter := bson.M{"user_id": userId, "name": categoryName, "folders.name": folderName}
	update := bson.M{"$pull": bson.M{"folders": bson.M{"name": folderName}}}

	res, err := f.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if res.MatchedCount == 0 {
		return ErrFolderNotFound
	}

	return nil
}

func (f folderMongoRepository) CheckFolderExists(ctx context.Context, userId, categoryName, folderName string) (bool, error) {
	filter := bson.M{"user_id": userId, "name": categoryName, "folders.name": folderName}
	err := f.collection.FindOne(ctx, filter).Err()
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}
