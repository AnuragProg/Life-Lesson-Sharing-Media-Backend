package models

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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

/*
Deletes Personal Life Lesson posts and corresponding comments associated with UserId
*/
func deletePllsAndCorrespondingCommentsOfUserId(userId string, pllColl *mongo.Collection, commentColl *mongo.Collection) error {
	var wg sync.WaitGroup

	errMessage := ""
	wg.Add(1)
	go func(){
		filter:= bson.M{ "userId": userId }
		result, err := pllColl.DeleteMany(context.TODO(), filter)  
		if err != nil{
			errMessage += err.Error()
		}
		if result.DeletedCount == 0{
			errMessage += "unable to delete personal life lessons"
		}
		wg.Done()
	}()

	wg.Add(1)
	go func(){
		filter := bson.M{"userId": userId}
		result, err := commentColl.DeleteMany(context.TODO(), filter)
		if err != nil{
			errMessage += err.Error()
		}
		if result.DeletedCount == 0{
			errMessage += "unable to delete comments"
		}
		wg.Done()
	}()

	wg.Wait()
	if len(errMessage) == 0{
		return nil
	}else{
		return errors.New(errMessage)
	}
}


/*
Adds Single Personal Life Lesson post
*/
func (pll *PersonalLifeLessonRequest) AddPll(coll *mongo.Collection)(*mongo.InsertOneResult, error){
	return coll.InsertOne(context.TODO(), pll)
}


/*
Updates Single Personal Life Lesson post
*/
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


/*
Deletes Single Personal Life Lesson post as well as All Comment's from Comment's Table 
*/
func DeletePllAndCommentsFromTable(
	pllId string, pllColl *mongo.Collection,
	commentColl *mongo.Collection,
) (*mongo.DeleteResult, error){
	id, err := primitive.ObjectIDFromHex(pllId)
	if err !=nil{
		return nil, err
	}
	filter := bson.M{"_id": id,}

	var pll PersonalLifeLesson
	result, err := pllColl.Find(context.TODO(), filter)

	if err != nil{
		return nil, err
	}
	err = result.Decode(&pll)
	if err != nil{
		return nil, err
	}

	// Deleting associated comments
	var wg sync.WaitGroup
	for _, commentId := range pll.Comments{
		wg.Add(1)
		go func(commentId string){
			DeleteComment(commentId, commentColl)
			wg.Done()
		}(commentId)
	}
	wg.Wait()	
	return pllColl.DeleteOne(context.TODO(), filter)
}


/*
Returns all Personal Life Lesson posts
*/
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

/*
Returns @pllId corresponding comment's 
*/
func GetPllComments(pllId string, coll *mongo.Collection) ([]Comment, error){
	comments := make([]Comment, 0)
	id, err := primitive.ObjectIDFromHex(pllId)
	if err != nil{
		return comments, err
	}
	filter := bson.M{"_id":id}
	result := coll.FindOne(context.TODO(), filter)
	var pll PersonalLifeLesson
	err = result.Decode(&pll)
	if err != nil{
		return comments, nil
	}
	return getComments(pll.Comments, coll)
}


/*
Adds Single Comment to Personal Life Lesson post's comment's list as well as to the Comments table
*/
func (comment *CommentRequest)AddCommentToTableAndList(
	pllId string, pllColl *mongo.Collection, commentColl *mongo.Collection,
)error{
	result,err := comment.addComment(commentColl)
	if err != nil{
		return err
	}
	allocatedCommentId := result.InsertedID.(primitive.ObjectID).Hex()
	fmt.Println("Allocated Comment Id is", allocatedCommentId)
	_, err = addCommentToPllCommentList(pllId, allocatedCommentId, pllColl)
	if err !=nil{
		return err
	}
	return nil
}

/*
Deletes Single Comment from Personal Life Lesson post's comment's list as well as from the Comment's table
*/
func DeleteCommentFromTableAndList(
	pllId string, commentId string, 
	pllColl *mongo.Collection, commentColl *mongo.Collection,
)error{
	var wg sync.WaitGroup

	var errorFromDeletioninPllCommentList, errorFromDeletionOfComment error

	wg.Add(1)
	go func() {
		result, err := deleteCommentFromPllCommentList(pllId, commentId, pllColl)
		errorFromDeletioninPllCommentList = err
		if result.ModifiedCount == 0{
			errorFromDeletioninPllCommentList = errors.New("unable to delete from Personal life lesson comments list")	
		}
		wg.Done()
	}()

	wg.Add(1)
	go func(){
		_, err := DeleteComment(commentId, commentColl)
		errorFromDeletionOfComment = err
		wg.Done()
	}()
	
	wg.Wait()
	
	errorMessage := ""
	if errorFromDeletionOfComment != nil{
		errorMessage += errorFromDeletionOfComment.Error()
	}
	if errorFromDeletioninPllCommentList != nil{
		errorMessage += " & " + errorFromDeletionOfComment.Error()
	}
	if errorMessage == ""{
		return nil
	}else{
		return errors.New(errorMessage) 
	}
}


/**
TODO test
*/
func addCommentToPllCommentList(pllId string, commentId string, coll *mongo.Collection)(*mongo.UpdateResult, error){
	id, err := primitive.ObjectIDFromHex(pllId)
	if err!=nil{
		return nil, err
	}
	filter := bson.M{"_id":id}
	update := bson.M{
		"$push":bson.M{"comments":commentId},
	}
	return coll.UpdateOne(context.TODO(), filter, update)
}

/**
TODO test
*/
func deleteCommentFromPllCommentList(pllId string, commentId string, coll *mongo.Collection)(*mongo.UpdateResult, error){
	id, err := primitive.ObjectIDFromHex(pllId)
	if err != nil{
		return nil, err
	}
	filter := bson.M{"_id" :id}

	update := bson.M{
		"$pull": bson.M{"comments": bson.M{"$in": bson.A{commentId}}},
	}
	return coll.UpdateMany(context.TODO(), filter, update)
}

// func DeleteCommentFromPllCommentList(pllId string, commentId string, coll *mongo.Collection)(*mongo.UpdateResult, error){
// 	id, err := primitive.ObjectIDFromHex(pllId)
// 	if err != nil{
// 		return nil, err
// 	}
// 	filter := bson.M{"_id" :id}

// 	update := bson.M{
// 		"$pull": bson.M{"comments": bson.M{"$in": bson.A{commentId}}},
// 	}
// 	return coll.UpdateMany(context.TODO(), filter, update)
// }