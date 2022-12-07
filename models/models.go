package models

type RequestMessage struct {
	ID     int    `json:"ID"`
	Status string `json:"Status"`
}

type ResponseMessage struct {
	ID     int    `json:"ID"`
	Number string `json:"Number"`
	Status string `json:"Status"`
}
