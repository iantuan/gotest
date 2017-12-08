package main

import (
    "encoding/json"
    "flag"
    "io/ioutil"
    "log"
    "net/http"
    "github.com/nats-io/nats"
)

type login struct {
	Name string `json:"name"`
	Pass  string `json:"pass"`
}



//var natsServer = flag.String("nats", "127.0.0.1:4224", "NATS server URI")
var natsServer = flag.String("nats", "nats-a:4222", "NATS server URI")

var natsClient *nats.Conn

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

    http.DefaultServeMux.HandleFunc("/login", loginHandler)
    

    log.Println("Starting product write service on port 8080") 
    log.Fatal(http.ListenAndServe(":8080", http.DefaultServeMux))
}

func loginHandler(rw http.ResponseWriter, r *http.Request) {

    if r.Method == "POST" {
             
        log.Println("/login handler called")
        data, err := ioutil.ReadAll(r.Body)
        if err != nil {
            rw.WriteHeader(http.StatusBadRequest)
            return
        }
        defer r.Body.Close()

        log.Println("data: %v", data)        	

        l := login{}
        err1 := json.Unmarshal(data, &l)

        if err1 != nil {
        log.Println("Unable to unmarshal event object")
            return
        }

        log.Printf("Received login message: %#v", l)
        natsClient.Publish("login", data)
    }

}


