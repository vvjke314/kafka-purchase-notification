package models

type RequestMessage struct {
	Email   string `json:"email"`
	Product string `json:"product"`
}

type ResponseMessage struct {
	Email   string `json:"email"`
	Product string `json:"product"`
	Time    string `json:"time"`
}
