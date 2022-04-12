package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	bvs "github.com/rudi9719/BulkVS2Go"
	"samhofi.us/x/keybase/v2"
	"samhofi.us/x/keybase/v2/types/chat1"
)

var (
	k               = keybase.NewKeybase()
	keybaseListener = NotifyListener{
		To:  nil,
		Run: notifyNumber,
	}
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
	if m.Channel.TopicName == "general" {
		logger.Printf("%+v is general, ignoring.", m.Channel.TopicName)
		return
	}
	logger.Printf("Converting Keybase message to SMS: %+v", m)

	msg := bvs.MessageSendRequest{
		To:      strings.Split(m.Channel.TopicName, ","),
		From:    strings.Replace(m.Channel.Name, "voipkjongsys.", "", -1),
		Message: m.Content.Text.Body,
	}
	resp, err := SendMessage(&msg)
	if err != nil {
		log.Printf("Error posing message from Keybase: %+v", err)
		k.ReactByConvID(m.ConvID, m.Id, ":-1:")
		return
	}
	for _, v := range resp.Results {
		if v.Status != "SUCCESS" {
			k.ReactByConvID(m.ConvID, m.Id, ":-1:")
			return
		}
	}
	k.KVPut(&m.Channel.Name, m.Channel.TopicName, resp.RefID, fmt.Sprintf("%+v", m.Id))
	k.ReactByConvID(m.ConvID, m.Id, ":+1:")

}

func logError(e error) {
	log.Printf("KEYBASE: %+v", e)
}

func notifyNumber(m bvs.MessageWebhookInput) {
	if m.DeliveryReceipt {
		for i := range m.To {
			team := fmt.Sprintf("voipkjongsys.%+v", m.To[i])
			test, err := k.KVGet(&team, m.To[i], m.RefID)
			if err != nil {
				continue
			}
			mid, err := strconv.ParseUint(test.EntryValue, 10, 32)
			if err != nil {
				continue
			}
			k.ReactByChannel(chat1.ChatChannel{
				Name:        fmt.Sprintf("voipkjongsys.%+v", m.To[i]),
				TopicName:   m.From,
				MembersType: keybase.TEAM,
			}, chat1.MessageID(mid), ":white_check_mark:")
			defer k.KVDelete(&team, m.To[i], m.RefID)
		}
	}
	msg := m.Message
	for i := range m.To {
		_, err := k.SendMessageByChannel(chat1.ChatChannel{
			Name:        fmt.Sprintf("voipkjongsys.%+v", m.To[i]),
			TopicName:   m.From,
			MembersType: keybase.TEAM,
		}, msg)
		if err != nil {
			_, err := k.SendMessageByChannel(chat1.ChatChannel{
				Name:        fmt.Sprintf("voipkjongsys.%+v", m.To[i]),
				TopicName:   "general",
				MembersType: keybase.TEAM,
			}, fmt.Sprintf("%+v", msg))
			if err != nil {
				_, err := k.SendMessageByChannel(chat1.ChatChannel{
					Name:        "voipkjongsys",
					TopicName:   "general",
					MembersType: keybase.TEAM,
				}, fmt.Sprintf("%+v tried to send %+v: %+v", m.From, m.To[i], msg))
				if err != nil {
					logger.Printf("Unable to pass message %+v along to voipkjongsys.%+v, #%+v", msg, m.To[i], m.From)
				}
			}
		}
	}
}

func runNotifier() {
	listeners = append(listeners, keybaseListener)
	logger.Printf("Starting Keybase!")
	chat := routeMessage
	err := logError

	handlers := keybase.Handlers{
		ChatHandler:  &chat,
		ErrorHandler: &err,
	}
	k.Run(handlers, &keybase.RunOptions{})
}