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
    "github.com/satori/go.uuid"
)
type UserInfo struct {
    Id   bson.ObjectId `json:"id" bson:"_id,omitempty"`
    Name string `json:"name"`
    Pass string `json:"pass"`
    Token string `json:"token"`
}

type AuthInfo struct {
	UUID string `json:"uuid"`
	Token string `json:"token"`
}

type AuthReply struct {
    Error int `json:"error"`
    Cmd string `json:"cmd"`
    Name string `json:"name"`
}

type Login struct {
	UUID string `json:"uuid"`
	Name string `json:"name"`
	Pass string `json:"pass"`
}



type reply_msg struct {
	Token string `json:"token"`
	Error int `json:"error"`
}

var session *mgo.Session
 
var loginClient *nats.Conn

//var natsServer = flag.String("nats", "127.0.0.1:4223", "NATS server URI")
var natsServer = flag.String("nats", "nats-b:4222", "NATS server URI")

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

    session, err := mgo.Dial("mymongo")
    
    
    fmt.Println("session addr", session)

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
    fmt.Println("session addr", session)
    
    fmt.Println(reflect.TypeOf(session))
    s, err := mgo.Dial("mymongo")

    s.SetMode(mgo.Monotonic, true)
    //go processLogin(l)
    c := s.DB("bigpower").C("player")
    
    fmt.Println(reflect.TypeOf(c))

    result := UserInfo{}
    err = c.Find(bson.M{"name": l.Name}).One(&result)
    if err != nil {
        fmt.Println(err)

    }

    fmt.Println(result.Pass)

    r := reply_msg{}

    if (l.Pass == result.Pass) {
        token := uuid.NewV4().String()
        r.Token = token
        r.Error = 0
        result.Token = token

        selector := bson.M{"_id": result.Id}
        data := bson.M{"$set": bson.M{"token": token}}
         
        err = c.Update(selector, data)

    } else {
        r.Error = -1
    }
    send_msg, _ := json.Marshal(r)

    loginClient.Publish(l.UUID, send_msg)

}

func handleAuth(m *nats.Msg) {
     
    log.Printf("Received auth message: %v, %#v", m.Subject, m.Data)

    a := AuthInfo{}
    err := json.Unmarshal(m.Data, &a)

    if err != nil {
        log.Println("Unable to unmarshal event object")
            return
        }
    
        log.Println("Received auth message: %v, %#v", m.Subject, a)
        fmt.Println("session addr", session)
        
        fmt.Println(reflect.TypeOf(session))
        s, err := mgo.Dial("mymongo")
    
        s.SetMode(mgo.Monotonic, true)
        //go processLogin(l)
        c := s.DB("bigpower").C("player")

        result := UserInfo{}
        err = c.Find(bson.M{"token": a.Token}).One(&result)
        
        r := AuthReply{}
        r.Cmd = "Auth"
        
        if (a.Token == result.Token) {
            r.Error = 0 
            r.Name = result.Name

        } else {
            r.Error = -1
        }

        send_msg, _ := json.Marshal(r)

        loginClient.Publish(a.UUID, send_msg)
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


