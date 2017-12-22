package main

import (
	"log"
	"encoding/json"
	"github.com/nats-io/nats"
)

type auth_token struct {
	Token string `json:"token"`
	Error int `json:"error"`
}

type client struct {
	reply chan []byte
	handler *login_handler
}


func (c* client) response_handler(m *nats.Msg) {

	c.reply <- m.Data
	
	log.Println("cleint response_handler", m.Data)
}

func (c *client) wait_response(uuid string) auth_token {

	 c.handler.nats_client.QueueSubscribe(uuid, "auth_group",  c.response_handler)
	 
	 t := auth_token{}

	 for data := range c.reply {

		 err := json.Unmarshal(data, &t)
		 if err != nil {
			 log.Println("Unable to unmarshal event")
			 break
		 }
		 log.Println("client get reply", t)
		 break
	 }

	 return t
}
