package models

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)


// Initial request from user
type UserRequest struct{
	Username string `json:"username" bson:"username"`
	Email string `json:"email" bson:"email"`
	Password string `json:"password" bson:"password"`
	Photo string `json:"photo,omitempty" bson:"photo"`
	JoinedOn uint32 `json:"joinedOn" bson:"joinedOn"`
}

// For Deciding whether requested user is admin or not
// Will be saved on the db
type UserIntermediate struct{
	Username string `json:"username" bson:"username"`
	Email string `json:"email" bson:"email"`
	Password string `json:"password" bson:"password"`
	Photo string `json:"photo,omitempty" bson:"photo"`
	JoinedOn uint32 `json:"joinedOn" bson:"joinedOn"`
	IsAdmin bool `json:"isAdmin" bson:"isAdmin"`
}

type User struct{
	ID string `json:"_id" bson:"_id"`
	Username string `json:"username" bson:"username"`
	Email string `json:"email" bson:"email"`
	Password string `json:"password" bson:"password"`
	Photo string `json:"photo,omitempty" bson:"photo,omitempty"`
	JoinedOn uint32 `json:"joinedOn" bson:"joinedOn"`
	LastToken string `json:"token,omitempty" bson:"token,omitempty"`
	IsAdmin bool `json:"isAdmin" bson:"isAdmin"`
}


func (user *UserRequest)ToUserIntermediate(isAdmin bool)(*UserIntermediate){
	return &UserIntermediate{
		Username: user.Username,
		Email: user.Email,
		Password: user.Password,
		Photo: user.Photo,
		JoinedOn: user.JoinedOn,
		IsAdmin: isAdmin,
	}
}


func (user *UserIntermediate) AddUser(coll *mongo.Collection) (*mongo.InsertOneResult, error){

	// checking if user already exists or not 
	var u User
	filter := bson.M{"email":user.Email}
	result := coll.FindOne(context.TODO(), filter)
	err := result.Decode(&u)
	if err == nil{
		return nil, errors.New("User already exists")
	}	

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
			"password" :user.Password, 
			"photo" : user.Photo,
		},
	}
	return coll.UpdateOne(context.TODO(), filter, update)
} 

func DeleteUser(userId string, pllColl *mongo.Collection, commentColl *mongo.Collection) (*mongo.DeleteResult, error){
	id , err:= primitive.ObjectIDFromHex(userId)
	if err!=nil{
		return nil, err
	}
	filter := bson.M{"_id":id}
	return pllColl.DeleteOne(context.TODO(), filter)	
}

func GetUsers(coll *mongo.Collection) ([]User, error) {
	var users []User = make([]User, 0)
	cursor, err := coll.Find(context.TODO(), bson.D{})
	if err !=nil{
		return users, err 
	}
	result := cursor.All(context.TODO(), &users)
	return users, result
}

func GetUserById(userId string, coll *mongo.Collection)(*User, error){
	id, err := primitive.ObjectIDFromHex(userId)
	if err != nil{
		return nil, err 
	}
	filter := bson.M{"_id":id}
	result := coll.FindOne(context.TODO(), filter)

	var user User
	err = result.Decode(&user)
	if err != nil{
		return nil, err
	}
	return &user, nil
}

func GetUsersById(userIds []string, coll *mongo.Collection) ([]User, error) {
	users := make([]User, 0)
	for _, userId := range userIds{
		user, err:= GetUserById(userId, coll)
		if err != nil{
			continue
		}
		users = append(users, *user)
	}
	return users, nil
}