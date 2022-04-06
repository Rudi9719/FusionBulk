package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/gorilla/mux"
	bvs "github.com/rudi9719/BulkVS2Go"
)

var (
	logger = log.Default()
	client bvs.BulkVS2GoClient
	c      Config
)

func init() {
	var u string
	var p string
	var t string

	flag.StringVar(&u, "u", "", "BulkVS User")
	flag.StringVar(&p, "p", "", "BulkVS Password")
	flag.StringVar(&t, "t", "", "Auth Token")
	flag.Parse()
	c = Config{
		BulkUser: u,
		BulkPass: p,
		Token:    t,
	}
	client = *bvs.NewClient(u, p)
}

func PanicSafe() {

}

func getSessionIdentifier(r *http.Request) string {
	defer PanicSafe()
	ipaddr := r.Header.Get("X-Real-IP")
	if ipaddr == "" {
		ipaddr = r.RemoteAddr
	}
	uri := r.URL.Path
	return fmt.Sprintf("%s:%s", ipaddr, uri)
}

func notFoundPage(w http.ResponseWriter, r *http.Request) {
	defer PanicSafe()
	logger.Printf("%s triggered notFoundPage", getSessionIdentifier(r))

	fmt.Fprint(w,
		"Sorry, a 404 error has occured. The requested page not found! <br><br>"+
			"<iframe width=\"560\" height=\"315\" src=\"https://www.youtube.com/embed/t3otBjVZzT0\" frameborder=\"0\" allow=\"accelerometer; autoplay; encrypted-media; gyroscope; picture-in-picture\" allowfullscreen></iframe>"+
			"<div class=\"error-actions\"><a href=\"/\" class=\"btn btn-primary btn-lg\"><span class=\"glyphicon glyphicon-home\"></span>Take Me Home </a>  <a href=\"mailto://rudi@nmare.net\" class=\"btn btn-default btn-lg\"><span class=\"glyphicon glyphicon-envelope\"></span> Contact Support </a></div>")

}
func bulkVSInput(w http.ResponseWriter, r *http.Request) {
	var m bvs.MessageWebhookInput
	err := json.NewDecoder(r.Body).Decode(&m)
	if err != nil {
		logger.Printf("%+v", err)
	}
	logger.Printf("%+v", m)
	if m.DeliveryReceipt {
		return
	}
	
	m.Message, err = url.QueryUnescape(m.Message)
	if err != nil {
		logger.Printf("%+v", err)
	}
	go notifyNumber(m)
}

func SendMessage(w http.ResponseWriter, r *http.Request) {
	var m MessageRequest
	var msg FusionMSG
	var ret MessageResponse

	err := json.NewDecoder(r.Body).Decode(&msg)
	if err != nil {
		logger.Printf("%+v", r.Body)
		logger.Println("Error decoding json from request")
		logger.Printf("%+v", err)
		ret.Message = fmt.Sprintf("%+v", err)
		ret.Code = 501
		fmt.Fprintf(w, "%+v", ret)
		return
	}
	username, password, ok := r.BasicAuth()
	if !ok || password != c.Token {
		logger.Printf("%+v", r.Body)
		logger.Println("Invalid Token")
		logger.Printf("%+v", err)
		ret.Message = fmt.Sprintf("%+v", err)
		ret.Code = 401
		fmt.Fprintf(w, "%+v", ret)
		return
	}
	logger.Printf("Sending message for %+v to %+v", username, m.To)
	m.From = username
	m.To = strings.Split(msg.To, ",")
	m.Message = msg.Text
	resp, err := client.PostMessageSend(&m.MessageSendRequest)
	if err != nil {
		logger.Println("Error sending message")
		logger.Printf("%+v", err)
		ret.Message = fmt.Sprintf("%+v", err)
		ret.Code = 503
		fmt.Fprintf(w, "%+v", ret)
		return
	}
	ret.MessageSendResponse = *resp
	ret.Code = 200
	ret.Message = "Message Sent."
	fmt.Fprintf(w, "%+v", ret)
}

func main() {
	defer PanicSafe()
	router := mux.NewRouter().StrictSlash(true)
	logger.Println("Adding HandleFuncs to router")
	router.NotFoundHandler = http.HandlerFunc(notFoundPage)
	router.HandleFunc("/bulkvs/webhook", bulkVSInput).Methods("POST")
	router.HandleFunc("/api/sendSMS", SendMessage).Methods("POST")
	logger.Printf("Starting server")
	go runKeybase()
	logger.Fatal(http.ListenAndServe(":8080", router))
}
