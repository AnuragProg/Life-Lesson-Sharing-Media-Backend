package models

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

/*
Redundant Fields(Will be taking on controllers) :-
1. userId (Will be replaced after authentication in controller)
2. likes
3. comments
4. categoryId
*/
type PersonalLifeLessonRequest struct {
	Title        string   `json:"title" bson:"title"`
	Learning     string   `json:"learning" bson:"learning"`
	RelatedStory string   `json:"relatedStory" bson:"relatedStory"`
	CategoryId   string   `json:"categoryId" bson:"categoryId"`
}

type PersonalLifeLessonUpdateRequest struct {
	ID 			 string    `json:"_id" bson:"_id"`
	Title        string   `json:"title" bson:"title"`
	Learning     string   `json:"learning" bson:"learning"`
	RelatedStory string   `json:"relatedStory" bson:"relatedStory"`
	CategoryId   string   `json:"categoryId" bson:"categoryId"`
}

type PersonalLifeLessonRequestIntermediate struct {
	UserId       string   `json:"userId" bson:"userId"`
	Username     string   `json:"username" bson:"username"`
	Title        string   `json:"title" bson:"title"`
	Learning     string   `json:"learning" bson:"learning"`
	RelatedStory string   `json:"relatedStory" bson:"relatedStory"`
	CreatedOn    int64   `json:"createdOn" bson:"createdOn"`
	CategoryId   string   `json:"categoryId" bson:"categoryId"`
}

type PersonalLifeLesson struct {
	ID           string   `json:"_id" bson:"_id"`
	UserId       string   `json:"userId" bson:"userId"`
	Username     string   `json:"username" bson:"username"`
	Title        string   `json:"title" bson:"title"`
	Learning     string   `json:"learning" bson:"learning"`
	RelatedStory string   `json:"relatedStory" bson:"relatedStory"`
	CreatedOn    int64    `json:"createdOn" bson:"createdOn"`
	CategoryId   string   `json:"categoryId" bson:"categoryId"`
	Likes        []string `json:"likes,omitempty" bson:"likes"`
	Comments     []string `json:"comments,omitempty" bson:"comments"`
}

func DeletePll(pllId string, coll *mongo.Collection) error{
	id, err := primitive.ObjectIDFromHex(pllId)
	if err != nil{
		return errors.New("not a personal life lesson id")
	}
	filter := bson.M{"_id":id}
	result,err := coll.DeleteOne(context.TODO(), filter)	
	if err != nil{
		return err
	}
	if result.DeletedCount == 0 {
		return errors.New("no personal life lesson deleted")
	}
	return nil
}

func (pll *PersonalLifeLessonRequest)ToPersonalLifeLessonRequestIntermediate(userId, username string)*PersonalLifeLessonRequestIntermediate{
	return &PersonalLifeLessonRequestIntermediate{
		UserId: userId,
		Username: username,
		Title: pll.Title,
		Learning: pll.Learning,
		RelatedStory: pll.RelatedStory,
		CreatedOn: time.Now().Unix(),
	}
}

/*
Adds Single Personal Life Lesson post
*/
func (pll *PersonalLifeLessonRequestIntermediate) AddPll(coll *mongo.Collection)(*mongo.InsertOneResult, error){
	return coll.InsertOne(context.TODO(), pll)
}

/*
Updates Single Personal Life Lesson post
*/
func (pll *PersonalLifeLessonUpdateRequest) UpdatePll(userId, username string, coll *mongo.Collection)(*mongo.UpdateResult, error){
	pllId, err := primitive.ObjectIDFromHex(pll.ID)
	if err!=nil{
		return nil, err
	}
	filter := bson.M{
		"_id": pllId,
	}
	update := bson.M{
		"$set": bson.M{
			"username" : username,
			"title" : pll.Title,
			"learning" : pll.Learning,
			"relatedStory" : pll.RelatedStory,
			"categoryId" : pll.CategoryId,
		},
	}
	return coll.UpdateOne(context.TODO(), filter, update)
}

/*
Returns all Personal Life Lesson posts
*/
func GetPlls(coll *mongo.Collection)([]PersonalLifeLesson, error){
	plls := make([]PersonalLifeLesson, 0)

	filter := bson.M{}
	cursor, err := coll.Find(context.TODO(), filter)
	if err != nil{
		return plls , nil
	}
	err = cursor.All(context.TODO(), &plls)
	
	return plls,err 
}

/*
Returns @pllId corresponding Personal Life Lesson post
*/
func GetPll(pllId string, coll *mongo.Collection) (*PersonalLifeLesson, error){
	id, err := primitive.ObjectIDFromHex(pllId)
	if err!=nil{
		return nil, err
	}
	var pll PersonalLifeLesson
	filter := bson.M{ "_id": id }
	findResult := coll.FindOne(context.TODO(), filter)
	err =findResult.Decode(&pll)
	if err!=nil{
		return nil, err
	}
	return &pll, nil
}

func LikePlls(pllIds []string, userId string, pllColl *mongo.Collection){
	pllObjectIds := make(map[string]primitive.ObjectID, len(pllIds))
	for _, pllId := range pllIds{
		id, err := primitive.ObjectIDFromHex(pllId)
		if err != nil{
			continue
		}
		pllObjectIds[pllId] = id
	}

	for _, id := range pllObjectIds{
		go pllColl.UpdateOne(context.TODO(), bson.M{"_id":id}, bson.M{"$addToSet":bson.M{"likes":userId}})
	}
}
func DislikePlls(pllIds []string, userId string, pllColl *mongo.Collection){
	pllObjectIds := make(map[string]primitive.ObjectID, len(pllIds))
	for _, pllId := range pllIds{
		id, err := primitive.ObjectIDFromHex(pllId)
		if err != nil{
			continue
		}
		pllObjectIds[pllId] = id
	}
	for _, id := range pllObjectIds{
		go pllColl.UpdateOne(context.TODO(), bson.M{"_id":id}, bson.M{"$pull":bson.M{"likes":userId}})
	}
}