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
		if m.Content.TypeName != "text" {
			logger.Printf("%+v", m.Content.TypeName)
			return
		}
		if !strings.HasPrefix(m.Channel.Name, "voipkjongsys.") {
			logger.Printf("%+v", m.Channel.Name)
			return
		}
		msg := bvs.MessageSendRequest {
			To: strings.Split( m.Channel.TopicName, ","),
			From: strings.Replace(m.Channel.Name, "voipkjongsys.", "", -1),
			Message: m.Content.Text.Body,
		}
		resp, err := client.PostMessageSend(&msg)
		if err != nil {
			log.Printf("Error posing message from Keybase: %+v", err)
		}
		for _, v := range(resp.Results) {
			if v.Status != "SUCCESS" {
				k.ReactByConvID(m.ConvID, m.Id, "-1")
				return
			}
		}
		k.ReactByConvID(m.ConvID, m.Id, "+1")


}

func logError(e error) {
	log.Printf("%+v", e)
}

func notifyNumber(m bvs.MessageWebhookInput) {
	msg := m.Message
	k.SendMessageByChannel(chat1.ChatChannel{
		Name: fmt.Sprintf("voipkjongsys.%+v", m.To[0]),
		TopicName: m.From,
	}, msg)
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