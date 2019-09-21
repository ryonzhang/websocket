package main

import (
	"errors"
	"github.com/gorilla/websocket"
	"sync"
)

type Connection struct{
	Conn *websocket.Conn
	inChan chan []byte
	outChan chan []byte
	closeChan chan byte
	mutex sync.Mutex
	isClosed bool
}

func InitConnection(Conn *websocket.Conn)(conn *Connection,err error){
	conn = &Connection{
		Conn:Conn,
		inChan:make(chan []byte,1000),
		outChan:make(chan []byte,1000),
		closeChan:make(chan byte,1),
	}
	go conn.readLoop()
	go conn.writeLoop()
	return
}

func (Conn *Connection) ReadMessage()(data []byte,err error){
	select{
	case data=<-Conn.inChan:
	case <-Conn.closeChan:
		err = errors.New("Connection is closed")
	}
	data=<-Conn.inChan
	return
}

func (Conn *Connection) WriteMessage(data []byte) (err error){
	select {
	case Conn.outChan <- data:
	case <-Conn.closeChan:
		err = errors.New("Connection is closed")
	}
	return
}

func (Conn *Connection) Close(){
	// repeatable
	Conn.Close()

	//only execute once
	Conn.mutex.Lock()
	if !Conn.isClosed{
		close(Conn.closeChan)
		Conn.isClosed = true
	}
	Conn.mutex.Unlock()

}

func (Conn *Connection) readLoop(){
	var (
		data []byte
		err error
	)
	for {
		if _,data,err = Conn.Conn.ReadMessage();err!=nil{
			goto ERR
		}
		select{
		case Conn.inChan<-data:
		case <-Conn.closeChan:
			goto ERR
		}
	}
ERR:
	Conn.Close()
}

func (Conn *Connection) writeLoop(){
	var (
		data []byte
		err error
	)
	for {

		if err = Conn.Conn.WriteMessage(websocket.TextMessage,data);err!=nil{
			goto ERR
		}
		select{
		case data=<-Conn.outChan:
		case <-Conn.closeChan:
			goto ERR
		}
	}
ERR:
	Conn.Close()
}
