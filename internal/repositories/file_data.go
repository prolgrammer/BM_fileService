package repositories

import (
	"app/internal/entities"
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type fileMongoRepository struct {
	collection *mongo.Collection
}

func NewFileMongoRepository(collection *mongo.Collection) FileRepository {
	return &fileMongoRepository{collection: collection}
}

func (f *fileMongoRepository) CreateFile(ctx context.Context, userId, categoryName, folderName string, data entities.File) error {
	filter := bson.M{
		"user_id":            userId,
		"name":               categoryName,
		"folders.name":       folderName,
		"folders.files.name": bson.M{"$ne": folderName},
	}

	update := bson.M{
		"$push": bson.M{
			"folders.$.files": data,
		},
	}

	_, err := f.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}

func (f *fileMongoRepository) SelectFile(ctx context.Context, userId, categoryName, folderName, nameFile string) (entities.File, error) {
	pipeline := []bson.M{
		{
			"$match": bson.M{
				"user_id":      userId,
				"name":         categoryName,
				"folders.name": folderName,
			},
		},
		{
			"$unwind": "$folders",
		},
		{
			"$match": bson.M{
				"folders.name": folderName,
			},
		},
		{
			"$unwind": "$folders.files",
		},
		{
			"$match": bson.M{
				"folders.files.name": nameFile,
			},
		},
		{
			"$replaceRoot": bson.M{
				"newRoot": "$folders.files",
			},
		},
	}

	cursor, err := f.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return entities.File{}, err
	}
	defer cursor.Close(ctx)

	var file entities.File
	for cursor.Next(ctx) {
		if err := cursor.Decode(&file); err != nil {
			return entities.File{}, err
		}
		return file, nil
	}

	return entities.File{}, ErrFileNotFound

}

func (f *fileMongoRepository) SelectFiles(ctx context.Context, userId, categoryName, folderName string) ([]entities.File, error) {
	fmt.Println("Repository categoryName", categoryName)
	fmt.Println("Repository folderName", folderName)
	pipeline := []bson.M{
		{
			"$match": bson.M{
				"user_id":      userId,
				"name":         categoryName,
				"folders.name": folderName,
			},
		},
		{
			"$unwind": "$folders",
		},
		{
			"$match": bson.M{
				"folders.name": folderName,
			},
		},
		{
			"$project": bson.M{
				"files": "$folders.files",
				"_id":   0,
			},
		},
	}

	fmt.Println(pipeline)

	cursor, err := f.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []entities.File

	for cursor.Next(ctx) {
		if err := cursor.Decode(&results); err != nil {
			return nil, err
		}
		return results, nil
	}

	return nil, ErrFolderNotFound
}

func (f *fileMongoRepository) DeleteFile(ctx context.Context, userId, categoryName, folderName, fileName string) error {
	filter := bson.M{
		"user_id":      userId,
		"name":         categoryName,
		"folders.name": folderName,
	}
	update := bson.M{
		"$pull": bson.M{
			"folders.$.files": bson.M{"name": fileName},
		},
	}

	res, err := f.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return ErrFileNotFound
	}

	return nil
}

func (f *fileMongoRepository) CheckFileExists(ctx context.Context, userId, categoryName, folderName, fileName string) (bool, error) {
	filter := bson.M{
		"user_id":            userId,
		"name":               categoryName,
		"folders.name":       folderName,
		"folders.files.name": fileName,
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
