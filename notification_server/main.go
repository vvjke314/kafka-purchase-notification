package main

import (
	"github.com/gin-gonic/gin"
	"github.com/vvjke314/kafka-purchase-notification/models"
	"log"
	"net/http"
)

func Handler(c *gin.Context) {
	req := &models.RequestMessage{}
	err := c.BindJSON(req)
	message := &models.ResponseMessage{
		Id:   req.Id,
		Name: req.Name,
		Time: req.Time,
		Note: req.Note,
	}
	if err != nil {
		log.Printf("Error occured in binding request body into message")
	}
	c.JSON(http.StatusOK, message)
}

func main() {
	r := gin.Default()
	r.POST("/notification", Handler)
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
