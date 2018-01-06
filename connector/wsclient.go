package main


import (
	"log"
	"github.com/nats-io/nats"
	"github.com/gorilla/websocket"
)


type wsclient struct {
	UUID string
	reply chan []byte
	socket *websocket.Conn
	handler *socket_handler
}


func (wsc *wsclient) read() {
	
	defer wsc.socket.Close()

	for {
		_, msg, err := wsc.socket.ReadMessage()
		if err != nil {
			return
		}
		log.Println(msg)
	}
}

func (wsc *wsclient) reply_handler(m *nats.Msg) {

	wsc.reply <- m.Data
	
	log.Println("wsccleint response_handler", m.Data)
}

func (wsc *wsclient) listen() {

	defer wsc.socket.Close()
	wsc.handler.nats_client.QueueSubscribe(wsc.UUID, "wsc_group",  wsc.reply_handler)

	for {
		select {
		case msg := <-wsc.reply:
			log.Println(msg)
			wsc.socket.WriteMessage(websocket.TextMessage, msg)
		}

	}

	//for msg := range wsc.send {
	//	err := wsc.socket.WriteMessage(websocket.TextMessage, msg)
	//	if err != nil {
	//		return
	//	}
	//}
}




