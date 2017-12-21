package main

import (
    "encoding/json"
    "flag"
    "io/ioutil"
    "log"
    "net/http"
    "github.com/nats-io/nats"
)

var natsServer = flag.String("nats", "nats-a:4222", "NATS server URI")

type login struct {
	Name string `json:"name"`
	Pass  string `json:"pass"`
}


func init() {
	
    flag.Parse()
}



func main() {
	
    var err error
    log.Println(*natsServer)
	
    natsClient, err = nats.Connect("nats://" + *natsServer)
	
    if err != nil {
        log.Fatal(err)	
    }
	
    defer natsClient.Close()

    l := &login_handler{
        nats_client: natsClient,
        conn: make(chan string),
    }

    http.Handle("/login", l)
    //http.DefaultServeMux.HandleFunc("/login", loginHandler)
    go l.run()

    if err := http.ListenAndServe(*addr, nil); err != nil {
        log.Fatal("ListenAndServe:", err)
    }

    log.Println("Starting product write service on port 8080") 
    log.Fatal(http.ListenAndServe(":8080", http.DefaultServeMux))
}



