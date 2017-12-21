package main

import (
	"log"
)

type clinet struct {
	reply chan string
	handler *login_handler
}


func (c* client) response_handler(m *nats.Msg) {


	c.reply <- m.Data
	
	log.Println("cleint response_handler", m.Data)
}

func (c *client) wait_response(uuid string) reply_msg {

	 c.handler.nats_client.QueueSubscribe(uuid, "auth_group",  c.response_handler)
	 
	 msg := ""

	 for msg = range c.reply {
		 log.Println("client get reply", msg)
		 break
	 }

	 return msg
}
