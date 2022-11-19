package models

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