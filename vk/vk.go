package vk

import (
	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/api/params"
	log "github.com/sirupsen/logrus"
)

func SendMessageService(VkID int, messageText string) error {
	vk := api.NewVK("vk1.a.r00BYA7KZ-1VhhdjD-Vj8RHICtGTV1KWfVyxBrhnuoaE03f4L1WT9xabHzM-12h3eeYk3N3H7ns4ai3M96XG_3UDMEaIkKDAZRIL3rBZrstHHgUIojKT47llOej2QE8TYdadHx70ngH9YhGu4PMvGHs40xGWeW0suE21xWq9MXgQxbsafg9HEPW2k1pjAN43BD9dWcTvhAmCaZjnKHUy8Q")
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	//b.Message(messageText[1 : len(messageText)-1])
	b.Message(messageText)
	b.PeerID(VkID)
	_, err := vk.MessagesSend(b.Params)
	if err != nil {
		log.Println("Failed to get record")
		log.Error(err)
		return err
	}
	return nil
}
