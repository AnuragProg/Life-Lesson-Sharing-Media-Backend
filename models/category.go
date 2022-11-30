package models

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type CategoryRequest struct{
	Title string `json:"title" bson:"title"`
	Description string `json:"description" bson:"description"`
}


type Category struct{
	ID string `json:"_id" bson:"_id"`
	Title string `json:"title" bson:"title"`
	Description string `json:"description" bson:"description"`
}

func (category *CategoryRequest) AddCategory(coll *mongo.Collection)(*mongo.InsertOneResult, error){
	return coll.InsertOne(context.TODO(), category)
}

func (category *Category) UpdateCategory(coll *mongo.Collection) (*mongo.UpdateResult, error){
	id ,err := primitive.ObjectIDFromHex(category.ID)
	if err!=nil{
		return nil, err
	}
	filter := bson.M{"_id":id}
	update := bson.M{
		"$set": bson.M{
			"title":category.Title,
			"description":category.Description,
		},
	}
	return coll.UpdateOne(context.TODO(), filter, update)
}

func GetCategories(coll *mongo.Collection)([]Category, error){
	categories := make([]Category, 0) 

	filter := bson.M{}
	result, err := coll.Find(context.TODO(), filter)	
	if err!=nil{
		return categories, err
	}

	err = result.All(context.TODO(),&categories)
	if err != nil{
		return categories, err
	}
	return categories, nil
}

func GetCategory(categoryId string, coll *mongo.Collection)(*Category, error){
	id, err := primitive.ObjectIDFromHex(categoryId)
	if err!=nil{
		return nil, err 
	}
	filter := bson.M{"_id":id}
	result := coll.FindOne(context.TODO(), filter)
	var category Category
	err = result.Decode(&category)
	if err!=nil{
		return nil, errors.New("category does not exist") 
	}
	return &category, nil
}

func DeleteCategory(categoryId string, coll *mongo.Collection)(*mongo.DeleteResult, error){
	id, err := primitive.ObjectIDFromHex(categoryId)
	if err != nil{
		return nil, err
	}

	filter := bson.M{"_id":id}
	return coll.DeleteOne(context.TODO(), filter)
}