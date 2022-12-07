package mail

import (
	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/api/params"
	log "github.com/sirupsen/logrus"
	"strconv"
	"strings"
)

func SendMessageService(VkID string, messageText string) error {
	Vk := strings.Replace(VkID, "\"", "", -1)
	log.Println(Vk)
	vkId, err := strconv.Atoi(Vk)
	if err != nil {
		return err
	}
	vk := api.NewVK("vk1.a.xor1NCLq_CEWCtDdFOTsqYpSHtXWK9DCDqV2YbF8to4ZMw0bNgXyziEiOPz10AV6VmGNyIdKfP3b-PXcaJ-dfKTTwAhXlDLHl0mNffeXW3cFudt_eOHLfD9t6iKa3jV1io8018u4anxgd8SPFIld4qAyowXcoDO7jjOTML9elPGPosr7u092BdzE1wOTS-8w")
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	message := strings.Split(messageText, ":")
	status := strings.Split(message[2], ",")
	car := strings.Split(message[3], ",")
	b.Message(status[0] + "\nНомер вашей машины: " + car[0] + "\nВремя: " + message[4] + ":" + message[5][:len(message[5])-1])
	b.PeerID(vkId)
	_, err = vk.MessagesSend(b.Params)
	if err != nil {
		log.Println("Failed to get record")
		log.Error(err)
		return err
	}
	return nil
}
