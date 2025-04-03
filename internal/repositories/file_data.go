package repositories

import (
	"app/internal/entities"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type fileMongoRepository struct {
	collection *mongo.Collection
}

func NewFileMongoRepository(collection *mongo.Collection) FileRepository {
	return &fileMongoRepository{collection: collection}
}

func (f *fileMongoRepository) CreateFile(ctx context.Context, data entities.File) error {
	file := bson.M{
		"name":        data.Name,
		"description": data.Description,
		"size":        data.Size,
		"type":        data.Type,
		"version":     data.Version,
		"created_at":  data.CreatedAt,
		"categories":  data.Categories,
	}

	_, err := f.collection.InsertOne(ctx, file)

	return err
}

func (f *fileMongoRepository) SelectFile(ctx context.Context, categoryId, folderName, fileName string) (entities.File, error) {
	var file entities.File
	filter := bson.M{
		"name":                   fileName,
		"categories.category_id": categoryId,
		"categories.folders":     entities.CreateFolder(folderName),
	}

	err := f.collection.FindOne(ctx, filter).Decode(&file)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return entities.File{}, ErrFileNotFound
		}
		return entities.File{}, err
	}

	return file, nil
}

func (f *fileMongoRepository) SelectFiles(ctx context.Context, categoryId, folderName string) ([]entities.File, error) {
	filter := bson.M{
		"categories.category_id": categoryId,
		"categories.folders":     entities.CreateFolder(folderName),
	}

	cursor, err := f.collection.Find(ctx, filter)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrFileNotFound
		}
		return nil, err
	}
	defer cursor.Close(ctx)

	var files []entities.File
	err = cursor.All(ctx, &files)
	if err != nil {
		return nil, err
	}

	return files, nil

}

func (f *fileMongoRepository) DeleteFile(ctx context.Context, categoryId, folderName, fileName string) error {
	filter := bson.M{
		"name":                   fileName,
		"categories.category_id": categoryId,
		"categories.folders":     entities.CreateFolder(folderName),
	}

	res, err := f.collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return ErrFileNotFound
	}

	return nil
}

func (f *fileMongoRepository) CheckFileExists(ctx context.Context, categoryId, folderName, fileName string) (bool, error) {
	filter := bson.M{
		"name":                   fileName,
		"categories.category_id": categoryId,
		"categories.folders":     entities.CreateFolder(folderName),
	}

	err := f.collection.FindOne(ctx, filter).Err()
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}
