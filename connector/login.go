package main

import (
	"log"
	"flag"
	"io/ioutil"
	"net/http"
	"github.com/nats-io/nats"
	"github.com/satori/go.uuid"
)



type login_info struct {
	UUID string `json:"uuid"`
	Name string `json:"name"`
	Pass string `json:"pass"`
}

type login_handler struct {
	nats_client *nats.Conn
	conn chan string
}

func (l *login_handler) run() {
	
	l.nats_client.QueueSubscribe("conn-1", "auth_group", l.handlePublish)

	for {
		select {
		case msg := <-l.conn:
			log.Println("get connector msg", msg)
		}


	}
}

func (l *login_handler) handlePublish(m *nats.Msg) {

	l.conn <- m.Data

}

func (l *login_handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {
		log.Println("/login handler called")
		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		log.Println("data: %v", data)

		li := login_info{}
		err1 := json.Unmarshal(data, &li)

		if err1 != nil {
			log.Println("Unable to unmarshal event object")
			return
		}

		li.UUID = uuid.NewV4().String()
		send_msg, _ := json.Marshal(li)

		c := &client{
			resp
		}
		
		l.nats_client.Publish("login", send_msg)

		c.wait_response(li.UUID)
	}

}