package repositories

import (
	"app/internal/entities"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
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

	update := bson.M{"$push": bson.M{"categories.$.folders": newFolder}}

	_, err := f.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}

func (f folderMongoRepository) SelectFolder(ctx context.Context, userId, categoryName, folderName string) (entities.Folder, error) {
	pipeline := []bson.M{
		{
			"$match": bson.M{
				"_id":                     userId,
				"categories.name":         categoryName,
				"categories.folders.name": folderName,
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
			"$unwind": "$categories.folders",
		},
		{
			"$match": bson.M{
				"categories.folders.name": folderName,
			},
		},
		{
			"$replaceRoot": bson.M{
				"newRoot": "$categories.folders",
			},
		},
	}

	cursor, err := f.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return entities.Folder{}, err
	}
	defer cursor.Close(ctx)

	if !cursor.Next(ctx) {
		return entities.Folder{}, ErrFolderNotFound
	}

	var folder entities.Folder
	if err := cursor.Decode(&folder); err != nil {
		return entities.Folder{}, err
	}

	return folder, nil
}

func (f folderMongoRepository) SelectFolders(ctx context.Context, userId, categoryName string) ([]entities.Folder, error) {
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
			"$project": bson.M{
				"_id":     0,
				"folders": "$categories.folders",
			},
		},
	}
	cursor, err := f.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if !cursor.Next(ctx) {
		return nil, ErrCategoryNotFound
	}

	type result struct {
		Folders []entities.Folder `bson:"folders"`
	}

	var res result
	if err := cursor.Decode(&res); err != nil {
		return nil, err
	}

	return res.Folders, nil
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
