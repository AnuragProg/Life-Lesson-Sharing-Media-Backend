package models

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)


type UserRequest struct{
	Username string `json:"username" bson:"username"`
	Email string `json:"email" bson:"email"`
	Photo string `json:"photo,omitempty" bson:"photo"`
	JoinedOn uint32 `json:"joinedOn" bson:"joinedOn"`
}

type User struct{
	ID string `json:"_id" bson:"_id"`
	Username string `json:"username" bson:"username"`
	Email string `json:"email" bson:"email"`
	Photo string `json:"photo,omitempty" bson:"photo,omitempty"`
	JoinedOn uint32 `json:"joinedOn" bson:"joinedOn"`
}

func (user *UserRequest) AddUser(coll *mongo.Collection) (*mongo.InsertOneResult, error){
	return coll.InsertOne(context.TODO(), user)
}

func (user *User) UpdateUser(coll *mongo.Collection) (*mongo.UpdateResult, error){
	id, err:= primitive.ObjectIDFromHex(user.ID)
	if err!=nil{
		return nil, err
	}
	filter := bson.D{{Key: "_id",Value: id}}
	update := bson.M{
		"$set":bson.M{
			"username": user.Username,
			"email" : user.Email,
			"photo" : user.Photo,
		},
	}
	return coll.UpdateOne(context.TODO(), filter, update)
} 

func DeleteUser(userId string, coll *mongo.Collection) (*mongo.DeleteResult, error){
	id , err:= primitive.ObjectIDFromHex(userId)
	if err!=nil{
		return nil, err
	}
	return coll.DeleteOne(context.TODO(), bson.M{"_id": id})	
}

func GetUsers(coll *mongo.Collection) ([]User, error) {
	var users []User
	cursor, err := coll.Find(context.TODO(), bson.D{})
	if err !=nil{
		return users, err 
	}
	result := cursor.All(context.TODO(), &users)
	return users, result
}