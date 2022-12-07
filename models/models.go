package models

type RequestMessage struct {
	Plate  string `json:"plate"`
	Id     string `json:"id"`
	Status string `json:"status"`
}

type ResponseMessage struct {
	Id     string `json:"id"`
	Status string `json:"status"`
	Plate  string `json:"plate"`
	Time   string `json:"time"`
}
