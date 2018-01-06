package main

import (
	"log"
	"strings"
	"encoding/json"
	"net/http"
	"github.com/gorilla/websocket"
	"github.com/nats-io/nats"
	"github.com/satori/go.uuid"
)

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

var upgrader = &websocket.Upgrader{ReadBufferSize: socketBufferSize, WriteBufferSize: socketBufferSize}


type AuthInfo struct {
	UUID string `json:"uuid"`
	Token string `json:"token"`
}

type socket_handler struct {
	nats_client *nats.Conn
	conn chan []byte
	authed_clients map[*client]bool
	unauthed_clients map[*client]bool
}


func (s *socket_handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	
	log.Println("/auth handler called")
	
	socket, err := upgrader.Upgrade(w, req, nil)
	
	if err != nil {
		log.Fatal("ServeHTTP:", err)
		return
	}

	path := req.URL.Path
	token := path[strings.Index(path[1:], "/")+2:]


	a := AuthInfo{}
	
	listen_uuid := uuid.NewV4().String()
	a.UUID = listen_uuid

	a.Token = token
	send_msg, _ := json.Marshal(a)

	s.nats_client.Publish("auth", send_msg)

	client := &wsclient{
		UUID : listen_uuid,
		socket: socket,
		reply:   make(chan []byte, messageBufferSize),
		handler:   s,
	}

	go client.listen()
	client.read()

}