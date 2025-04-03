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
	filter := bson.M{
		"_id":             userId,
		"categories.name": categoryName,
	}

	newFolder := entities.CreateFolder(folderName)

	update := bson.M{"$push": bson.M{"categories.folders.$": newFolder}}

	_, err := f.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}

func (f folderMongoRepository) SelectFolder(ctx context.Context, userId, categoryName, folderName string) (entities.Folder, error) {
	var folder entities.Folder
	filter := bson.M{"_id": userId, "categories.name": categoryName, "categories.folders.name": folderName}
	projection := bson.M{"categories.folders.$": 1}

	err := f.collection.FindOne(ctx, filter, options.FindOne().SetProjection(projection)).Decode(&folder)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return entities.Folder{}, ErrFolderNotFound
		}
		return entities.Folder{}, err
	}

	return folder, nil
}

func (f folderMongoRepository) SelectFolders(ctx context.Context, userId, categoryName string) ([]entities.Folder, error) {
	var folders []entities.Folder
	filter := bson.M{"_id": userId, "categories.name": categoryName}

	cursor, err := f.collection.Find(ctx, filter)
	if err != nil {
		return nil, ErrCategoryNotFound
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var folder entities.Folder
		err := cursor.Decode(&folder)
		if err != nil {
			return nil, err
		}
		folders = append(folders, folder)
	}

	return folders, nil
}

func (f folderMongoRepository) DeleteFolder(ctx context.Context, userId, categoryName, folderName string) error {
	filter := bson.M{"_id": userId, "categories.name": categoryName, "categories.folders.name": folderName}
	update := bson.M{"$pull": bson.M{
		"categories.$.folders": bson.M{"name": folderName},
	}}

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
	filter := bson.M{"_id": userId, "categories.name": categoryName, "categories.folders.name": folderName}
	err := f.collection.FindOne(ctx, filter).Err()
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}
