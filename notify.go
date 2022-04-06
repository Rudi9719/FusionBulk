package main

import (
	"fmt"
	"log"
	"strings"

	bvs "github.com/rudi9719/BulkVS2Go"
	"samhofi.us/x/keybase/v2"
	"samhofi.us/x/keybase/v2/types/chat1"
)

var (
	k = keybase.NewKeybase()
)


func routeMessage(m chat1.MsgSummary) {
		if m.Sender.Username == k.Username {
			return
		}
		if m.Content.TypeName != "text" {
			logger.Printf("%+v is not text, ignoring.", m.Content.TypeName)
			return
		}
		if !strings.HasPrefix(m.Channel.Name, "voipkjongsys.") {
			logger.Printf("%+v is not prefixed with voipkjongsys, ignoring.", m.Channel.Name)
			return
		}
		logger.Printf("Converting Keybase message to SMS: %+v", m)
		msg := bvs.MessageSendRequest {
			To: strings.Split(m.Channel.TopicName, ","),
			From: strings.Replace(m.Channel.Name, "voipkjongsys.", "", -1),
			Message: m.Content.Text.Body,
		}
		resp, err := client.PostMessageSend(&msg)
		if err != nil {
			log.Printf("Error posing message from Keybase: %+v", err)
			k.ReactByConvID(m.ConvID, m.Id, ":-1:")
			return
		}
		for _, v := range(resp.Results) {
			if v.Status != "SUCCESS" {
				k.ReactByConvID(m.ConvID, m.Id, ":-1:")
				return
			}
		}
		k.ReactByConvID(m.ConvID, m.Id, ":+1:")


}

func logError(e error) {
	log.Printf("%+v", e)
}

func notifyNumber(m bvs.MessageWebhookInput) {
	msg := m.Message
	for i, _ := range(m.To) {
		k.SendMessageByChannel(chat1.ChatChannel{
			Name: fmt.Sprintf("voipkjongsys.%+v", m.To[i]),
			TopicName: m.From,
			TopicType: keybase.TEAM,
		}, msg)
	}
}

func runKeybase() {
	logger.Printf("Starting Keybase!")
	chat := routeMessage
	err := logError

	handlers := keybase.Handlers{
		ChatHandler:  &chat,
		ErrorHandler: &err,
	}
	k.Run(handlers, &keybase.RunOptions{})
}