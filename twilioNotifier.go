package main

import (
	bvs "github.com/rudi9719/BulkVS2Go"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// Set to your endpoint for twilio
var endpoint = "https://voipkjongsys.com/app/sms/hook/sms_hook_twilio.php"

// Notifier to mimic Twilio SMS Hook Output
var twilioNotifier = NotifyListener{
	To:  nil,
	Run: notifyTwilio,
}

func notifyTwilio(m bvs.MessageWebhookInput) {
	if m.DeliveryReceipt {
		return
	}
	for _, nmbr := range m.To {
		data := url.Values{}
		data.Set("From", m.From)
		data.Set("To", nmbr)
		data.Set("Body", m.Message)
		client := &http.Client{}
		r, err := http.NewRequest("POST", endpoint, strings.NewReader(data.Encode())) // URL-encoded payload
		if err != nil {
			log.Println(err)
			continue
		}
		r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

		res, err := client.Do(r)
		if err != nil {
			log.Println(err)
			continue
		}
		log.Println(res.Status)
		defer res.Body.Close()
		_, err = ioutil.ReadAll(res.Body)
		if err != nil {
			log.Println(err)
			continue
		}
	}
}
