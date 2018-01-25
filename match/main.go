
package main

import (
	"log"
	"flag"
	"github.com/nats-io/nats"
)

var natsServer = flag.String("nats", "nats-b:4222", "NATS server URI")
var matchClient *nats.Conn
var matchlist map[string]int

type MatchInfo struct {
	UUID string `json:"uuid"`
	Cmd string `json:"cmd"`
}

func main() {
	var err error

	matchClient, err = nats.Connect("nats://" + *natsServer)

	if err != nil {
        log.Fatal(err)
	}
	

	matchClient.QueueSubscribe("match", "match_group", handleMatch)


}

func handleMatch(m *nats.Msg) {

}
