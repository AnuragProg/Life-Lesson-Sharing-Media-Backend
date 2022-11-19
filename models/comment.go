package models

import(
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CommentRequest struct{
	UserId string `json:"userId" bson:"userId"`
	Username string `json:"username" bson:"username"`
	Comment string `json:"comment" bson:"comment"`
	CommentedOn uint32 `json:"commentedOn" bson:"commentedOn"`
}


type Comment struct{
	ID string `json:"_id" bson:"_id"`
	UserId string `json:"userId" bson:"userId"`
	Username string `json:"username" bson:"username"`
	Comment string `json:"comment" bson:"comment"`
	CommentedOn uint32 `json:"commentedOn" bson:"commentedOn"`
}

func (comment *CommentRequest) AddComment(coll *mongo.Collection)(*mongo.InsertOneResult, error){
	return coll.InsertOne(context.TODO(), comment)	
}

func (comment *Comment) UpdateComment(coll *mongo.Collection)(*mongo.UpdateResult, error){
	id, err := primitive.ObjectIDFromHex(comment.ID)
	if err!=nil{
		return nil, err
	}
	filter := bson.M{ "_id" : id, }
	update := bson.M{
		"$set": bson.M{
			"username":comment.Username,
			"comment" :comment.Comment,
			"commentedOn":comment.CommentedOn,
		},
	}
	return coll.UpdateOne(context.TODO(), filter, update)
}

func DeleteComment(commentId string, coll *mongo.Collection) (*mongo.DeleteResult, error){
	id, err := primitive.ObjectIDFromHex(commentId)
	if err != nil{
		return nil, err
	}
	filter := bson.M{"_id": id }
	return coll.DeleteOne(context.TODO(), filter)
}

func GetComments(commentIds []string, coll *mongo.Collection) ([]Comment, error){
	comments := make([]Comment, 0)

	for _, commentId := range commentIds{
		id, err := primitive.ObjectIDFromHex(commentId)
		if err != nil{
			continue
		}
		var comment Comment
		filter := bson.M{"_id": id}
		result := coll.FindOne(context.TODO(), filter)
		err = result.Decode(&comment)
		if err != nil{
			continue
		}
		comments = append(comments, comment)
	}
	return comments, nil
}