package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/sirupsen/logrus"
	log2 "log"
	"strconv"
	"strings"
)

func SendMessageService(TgID string, messageText string) error {
	t := strings.Replace(TgID, "\"", "", -1)
	tg, err := strconv.Atoi(t)
	if err != nil {
		log.Error(err)
		return err
	}
	bot, err := tgbotapi.NewBotAPI("5984784274:AAEJnX8wPC_xLJqZsxPDnbw4Lkqb6vLySMA")
	if err != nil {
		log.Error(err)
		return err
	}

	message := messageText[1 : len(messageText)-1]
	messageSplit := strings.Split(message, ",")

	nameSplit := strings.Split(messageSplit[1], ":")
	name1 := nameSplit[1]
	name := name1[1 : len(name1)-1]

	noteSplit := strings.Split(messageSplit[3], ":")
	note1 := noteSplit[1]
	note2 := strings.ToLower(note1)
	note := note2[1 : len(note2)-1]

	timeSplit := strings.SplitN(messageSplit[2], ":", 2)
	time1 := timeSplit[1]
	time := time1[1 : len(time1)-1]

	messageCombine := "Напоминание для Вас, " + name + ":" + "\n" + "Не забудьте " + note + " в " + time
	msg := tgbotapi.NewMessage(int64(tg), messageCombine)

	_, err = bot.Send(msg)
	if err != nil {
		log.Error(err)
		return err
	}
	log2.Println("Напоминание отправлено в телеграмм")
	return nil
}
