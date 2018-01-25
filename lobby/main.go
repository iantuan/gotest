package main

import (
	"flag"
	"github.com/nats-io/nats"
)

var natsServer = flag.String("nats", "nats-c:4222", "NATS server URI")

func main() {

	client, _ := nats.Connect("nats://" + *natsServer)

	client.QueueSubscribe("lobby", "lobby_group", handle_lobby)

}


func handle_lobby(m *nats.Msg) {


}