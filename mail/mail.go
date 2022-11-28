package mail

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/vvjke314/kafka-purchase-notification/models"
	"log"
	"net/smtp"
	"strings"
)

func SendMessageService(mailAddress string, messageText []byte) error {
	var auth = smtp.PlainAuth("", "hopply@mail.ru", "Kneu6n4tiiGknTnkbbwr", "smtp.mail.ru")
	var conf = &tls.Config{ServerName: "smtp.mail.ru"}
	conn, err := tls.Dial("tcp", "smtp.mail.ru:465", conf)
	if err != nil {
		log.Println(err)
		return err
	}
	cl, err := smtp.NewClient(conn, "smtp.mail.ru")
	if err != nil {
		return err
	}
	err = cl.Auth(auth)
	if err != nil {
		return err
	}
	err = cl.Mail("hopply@mail.ru")
	if err != nil {
		return err
	}
	err = cl.Rcpt(strings.Replace(mailAddress, "\"", "", -1))
	if err != nil {
		return err
	}
	w, err := cl.Data()
	if err != nil {
		return err
	}
	messageBody := models.ResponseMessage{}
	err = json.Unmarshal(messageText, &messageBody)
	messageSend := "Subject: Уведомление об оплате\r\n" + fmt.Sprintf("Покупка товара: %s\n %s", messageBody.Product, messageBody.Time)
	_, err = w.Write([]byte(messageSend))
	if err != nil {
		return err
	}

	defer w.Close()

	defer cl.Quit()
	return nil
}
