package main

import (
	"log"
	"strings"
	"encoding/json"
	"net/http"
	"github.com/gorilla/websocket"
	"github.com/nats-io/nats"
	"github.com/satori/go.uuid"
	"github.com/go-redis/redis"
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
	leave chan *wsclient
	//authed_clients map[*client]bool
	//unauthed_clients map[*client]bool
}

func new_socket_handler(n *nats.Conn) *socket_handler {

	return &socket_handler{
		nats_client: n,
		conn: make(chan []byte, 256), 
		leave: make(chan *wsclient),
	}
}

func (s *socket_handler) run() {

	for {
		select {
		case client := <-s.leave:
			rdb := redis.NewClient(&redis.Options{
				Addr:     "myredis:6379",
				Password: "", // no password set
				DB:       1,  // use default DB
			})
			
			rdb.Del(client.UUID)
		}
	}
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

	defer func() { s.leave <- client}()

	go client.listen()
	client.read()

}