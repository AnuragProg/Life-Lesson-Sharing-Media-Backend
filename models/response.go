package models

type ErrorResponse struct{
	Message string `json:"message"`
	Response string `json:"response,omitempty"`
}


type GeneralResponse struct{
	Message string `json:"message"`
}