package models

import (
	"time"
	"sync"
	"errors"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Request body from user when updating comment
type CommentUpdateRequest struct{
	ID string `json:"_id" bson:"_id"`
	Comment string `json:"comment" bson:"comment"`
}

// Request body from user when adding new comment
type CommentRequest struct{
	PllId string `json:"pllId" bson:"pllId"`
	Comment string `json:"comment" bson:"comment"`
}

// Internally to be used
// Actual data that will be added to the db
type CommentRequestIntermediate struct{
	PllId string `json:"pllId" bson:"pllId"`
	UserId string `json:"userId" bson:"userId"`
	Username string `json:"username" bson:"username"`
	Comment string `json:"comment" bson:"comment"`
	CommentedOn time.Time `json:"commentedOn" bson:"commentedOn"`
}

// Full data that is stored in db
type Comment struct{
	ID string `json:"_id" bson:"_id"`
	PllId string `json:"pllId" bson:"pllId"`
	UserId string `json:"userId" bson:"userId"`
	Username string `json:"username" bson:"username"`
	Comment string `json:"comment" bson:"comment"`
	CommentedOn time.Time `json:"commentedOn" bson:"commentedOn"`
}

func (comment *CommentRequest)ToCommentRequestIntermediate(userId, username string) *CommentRequestIntermediate{
	return &CommentRequestIntermediate{
		PllId: comment.PllId,
		UserId: userId,
		Username: username,
		Comment: comment.Comment,
		CommentedOn: time.Now(),
	}
}


func (comment *CommentRequestIntermediate) AddComment(pllColl, commentColl *mongo.Collection)(*mongo.InsertOneResult, error){
	//Check if given pll exists
	pllId, err := primitive.ObjectIDFromHex(comment.PllId)
	if err != nil{
		return nil, err
	}
	var pll PersonalLifeLesson
	filter := bson.M{"_id": pllId}
	result := pllColl.FindOne(context.TODO(), filter)
	err = result.Decode(&pll)
	if err != nil{
		return nil, errors.New("no such personal life lesson post exist") 
	}

	// Inserting comment
	commentInsertResult, err := commentColl.InsertOne(context.TODO(), comment)	
	if err != nil{
		return nil, err 
	}
	commentId := commentInsertResult.InsertedID.(primitive.ObjectID).Hex()
	return nil, addCommentToPllCommentsTable(pll.ID, commentId, pllColl)
}


func (comment *CommentUpdateRequest) UpdateComment(coll *mongo.Collection)(*mongo.UpdateResult, error){
	id, err := primitive.ObjectIDFromHex(comment.ID)
	if err!=nil{
		return nil, err
	}
	filter := bson.M{ "_id" : id }
	update := bson.M{
		"$set": bson.M{
			"comment" :comment.Comment,
			"commentedOn":time.Now().Unix(),
		},
	}
	return coll.UpdateOne(context.TODO(), filter, update)
}


func DeleteComment(commentId string, pllColl, commentColl *mongo.Collection) error{
	id, err := primitive.ObjectIDFromHex(commentId)
	if err != nil{
		return err
	}
	var comment Comment
	filter := bson.M{"_id":id}
	result:= commentColl.FindOne(context.TODO(), filter)

	err = result.Decode(&comment)
	if err != nil{
		return err
	}
	
	err = deleteCommentFromPllCommentsTable(comment.PllId, commentId, pllColl)
	if err != nil{
		return err
	}
	_, err = commentColl.DeleteOne(context.TODO(), filter)
	
	return err
}


func GetComments(commentIds []string, coll *mongo.Collection) []Comment{
	comments := make([]Comment, 0)
	for _, commentId := range commentIds{
		id, err := primitive.ObjectIDFromHex(commentId)
		if err != nil{
			continue
		}
		filter := bson.M{"_id":id}
		result := coll.FindOne(context.TODO(), filter)
		var comment Comment
		err = result.Decode(&comment)
		if err==nil{
			comments = append(comments, comment)
		}
	}
	return comments
}

func deleteCommentFromPllCommentsTable(pllId, commentId string, coll *mongo.Collection) error{
	id, err := primitive.ObjectIDFromHex(pllId)
	if err != nil{
		return err
	}

	filter := bson.M{"_id":id}
	update := bson.M{
		"$pull":bson.M{
			"comments": bson.M{"$in": bson.A{commentId}},
		},
	}
	_, err = coll.UpdateOne(context.TODO(), filter, update)
	return err
}


func addCommentToPllCommentsTable(pllId, commentId string, coll *mongo.Collection) error{
	id, err := primitive.ObjectIDFromHex(pllId)
	if err != nil{
		return err
	}
	filter := bson.M{"_id":id}
	update := bson.M{
		"$push":bson.M{
			"comments": commentId,
		},
	}
	_, err = coll.UpdateOne(context.TODO(), filter, update)
	return err
}

// Only for admin
func DeleteComments(commentIds []string, coll *mongo.Collection) error{
	var wg sync.WaitGroup
	for _, commentId := range commentIds{
		commentId, err := primitive.ObjectIDFromHex(commentId)
		if err != nil{
			continue
		}
		filter := bson.M{"_id":commentId}
		wg.Add(1)
		go func() {
			coll.DeleteOne(context.TODO(), filter)
			wg.Done()
		}()
	}
	wg.Wait()
	return nil
}