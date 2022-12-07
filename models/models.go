package models

type RequestMessage struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Time string `json:"time"`
	Note string `json:"note"`
}

type ResponseMessage struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Time string `json:"time"`
	Note string `json:"note"`
}
