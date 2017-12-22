package main

import (
	"log"
	"encoding/json"
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
	conn chan []byte
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


		log.Println("recv request %#v", li)

		li.UUID = uuid.NewV4().String()
		send_msg, _ := json.Marshal(li)

		log.Println("send_msg", send_msg)
		c := &client{
			handler: l,
			reply: make(chan []byte, 256),
		}
		
		l.nats_client.Publish("login", send_msg)

		a := c.wait_response(li.UUID)

		reply_msg, _ := json.Marshal(a)
		
		w.Header().Set("Content-Type", "application/json")

		w.Write(reply_msg)
		//fmt.Fprint(w, reply_msg)
	}

}