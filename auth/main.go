package main

import (
    "encoding/json"
    "flag"
    "log"
    "fmt"
    "time"
    "reflect"
    "github.com/nats-io/nats"
    "gopkg.in/mgo.v2"    
    "gopkg.in/mgo.v2/bson"
)

type Login struct {
	Name string 
	Pass string 
}


var session *mgo.Session
 
var loginClient *nats.Conn

var natsServer = flag.String("nats", "127.0.0.1:4223", "NATS server URI")

func init() {

    flag.Parse()

}


func main() {
    var err error
    log.Println(*natsServer)
    loginClient, err = nats.Connect("nats://" + *natsServer)
    if err != nil {
        log.Fatal(err)
    }

    session, err := mgo.Dial("127.0.0.1")
    
    defer session.Close()
    defer loginClient.Close()

        
    session.SetMode(mgo.Monotonic, true)

    loginClient.QueueSubscribe("login", "auth_group", handleLogin)
    loginClient.QueueSubscribe("auth", "auth_group", handleAuth)

    for {
        time.Sleep(10*time.Second)
    }

}

func handleLogin(m *nats.Msg) {
    l := Login{}
    err := json.Unmarshal(m.Data, &l)
    if err != nil {
    log.Println("Unable to unmarshal event object")
        return
    }

    log.Println("Received login message: %v, %#v", m.Subject, l)
    
    fmt.Println(reflect.TypeOf(session))

    s, err := mgo.Dial("127.0.0.1")

    s.SetMode(mgo.Monotonic, true)
    //go processLogin(l)
    c := s.DB("bigpower").C("player")

    fmt.Println(reflect.TypeOf(c))
    result := Login{}
    err = c.Find(bson.M{"name": l.Name}).One(&result)
    if err != nil {
        log.Fatal(err)

    }

    fmt.Println(result.Pass)

}

func handleAuth(m *nats.Msg) {
     
    log.Printf("Received auth message: %v, %#v", m.Subject, m.Data)

}

//func processLogin(l login) {

//    c := session.DB("bigpower").C("player")

//    result := login{}
//    err := c.Find(bson.M{"name": "ian"}).One(&result)
//    if err != nil {
//        log.Fatal(err)

//    }

//    fmt.Println(result.Pass)

//}


