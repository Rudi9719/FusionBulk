package main

import bvs "github.com/rudi9719/BulkVS2Go"

type MessageRequest struct {
	bvs.MessageSendRequest `json:"message"`
	Token string `json:"token"`
}

type MessageResponse struct {
	bvs.MessageSendResponse
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type Config struct {
	BulkUser string `json:"BulkUser"`
	BulkPass string `json:"BulkPass"`
	Token    string `json:"Token"`
}

	
type FusionMSG struct {
	To   string `json:"to"`
	Text string `json:"text"`
}