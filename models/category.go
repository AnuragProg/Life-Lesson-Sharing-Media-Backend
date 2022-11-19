package models

type CategoryRequest struct{
	Title string `json:"title" bson:"title"`
	Description string `json:"description" bson:"description"`
}


type Category struct{
	ID string `json:"_id" bson:"_id"`
	Title string `json:"title" bson:"title"`
	Description string `json:"description" bson:"description"`
}