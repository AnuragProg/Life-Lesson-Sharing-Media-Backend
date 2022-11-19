package models

import(
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
) 

type PersonalLifeLessonRequest struct {
	UserId       string   `json:"userId" bson:"userId"`
	Username     string   `json:"username" bson:"username"`
	Title        string   `json:"title" bson:"title"`
	Learning     string   `json:"learning" bson:"learning"`
	RelatedStory string   `json:"relatedStory" bson:"relatedStory"`
	CreatedOn    uint32   `json:"createdOn" bson:"createdOn"`
	CategoryId   string   `json:"categoryId" bson:"categoryId"`
	Likes        []string `json:"likes" bson:"likes"`
	Comments     []string `json:"comments" bson:"comments"`
}

type PersonalLifeLesson struct {
	ID           string   `json:"_id" bson:"_id"`
	UserId       string   `json:"userId" bson:"userId"`
	Username     string   `json:"username" bson:"username"`
	Title        string   `json:"title" bson:"title"`
	Learning     string   `json:"learning" bson:"learning"`
	RelatedStory string   `json:"relatedStory" bson:"relatedStory"`
	CreatedOn    uint32   `json:"createdOn" bson:"createdOn"`
	CategoryId   string   `json:"categoryId" bson:"categoryId"`
	Likes        []string `json:"likes" bson:"likes"`
	Comments     []string `json:"comments" bson:"comments"`
}

func (pll *PersonalLifeLessonRequest) AddPll(coll *mongo.Collection)(*mongo.InsertOneResult, error){
	return coll.InsertOne(context.TODO(), pll)
}


func (pll *PersonalLifeLesson) UpdatePll(coll *mongo.Collection)(*mongo.UpdateResult, error){
	pllId, err := primitive.ObjectIDFromHex(pll.ID)
	if err!=nil{
		return nil, err
	}
	filter := bson.M{
		"_id": pllId,
	}
	update := bson.M{
		"$set": bson.M{
			"username" : pll.Username,
			"title" : pll.Title,
			"learning" : pll.Learning,
			"relatedStory" : pll.RelatedStory,
			"categoryId" : pll.CategoryId,
		},
	}
	return coll.UpdateOne(context.TODO(), filter, update)
}


func DeletePll(pllId string, coll *mongo.Collection) (*mongo.DeleteResult, error){
	id, err := primitive.ObjectIDFromHex(pllId)
	if err !=nil{
		return nil, err
	}
	filter := bson.M{"_id": id,}
	return coll.DeleteOne(context.TODO(), filter)
}

func GetPlls(coll *mongo.Collection)([]PersonalLifeLesson, error){
	var plls[]PersonalLifeLesson
	filter := bson.M{}
	cursor, err := coll.Find(context.TODO(), filter)
	if err != nil{
		return make([]PersonalLifeLesson, 0), err
	}
	err = cursor.All(context.TODO(), &plls)

	if err != nil{
		return make([]PersonalLifeLesson, 0), err
	}

	if plls == nil{
		return make([]PersonalLifeLesson, 0), nil 
	}
	return plls, nil

}

func GetPll(pplId string, coll *mongo.Collection) (*PersonalLifeLesson, error){
	id, err := primitive.ObjectIDFromHex(pplId)
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