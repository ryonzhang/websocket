package main

import (
	"github.com/gorilla/websocket"
	"net/http"
	"time"
)

var (
	upgrader = websocket.Upgrader{
		// accept Cross Domain
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)


func main(){
	http.HandleFunc("/ws",wsHandler)
	http.ListenAndServe("0.0.0.0:7777",nil)
}

func wsHandler(w http.ResponseWriter,r *http.Request){
	var (
		wsConn *websocket.Conn
		conn *Connection
		err error
		data []byte
	)
	if wsConn,err = upgrader.Upgrade(w,r,nil);err!=nil{
		return
	}
	if conn,err = InitConnection(wsConn);err!=nil{
		goto ERR
	}

	go func(){
		for{
			if err=conn.WriteMessage([]byte("heartbeat"));err!=nil{
				return
			}
			time.Sleep(1*time.Second)
		}
	}()

	for{
		if data,err = conn.ReadMessage();err!=nil{
			goto ERR
		}
		if err=conn.WriteMessage(data);err!=nil{
			goto ERR
		}
	}
ERR:
	conn.Close()
}
