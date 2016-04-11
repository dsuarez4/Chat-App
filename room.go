package main

import (
	"fmt"
	"net/http"
	"github.com/gorilla/websocket"
	"log"
)

type room struct {
	// forward is a channel that holds incoming messages
	// that should be forwarded to the other clients.
	forward chan []byte
	// join is a channel for clients wishing to join the room
	join chan *client
	// leave is a channel for clients wishing to leave the room
	leave chan *client
	// clients holds all current client  in this room
	clients map[*client]bool
}

// newRoom
func newRoom() *room {
	return &room{
		forward: make(chan []byte),
		join: make(chan *client),
		leave: make(chan *client),
		clients: make(map[*client] bool),
	}
}

func (r *room) run() {
	/* 	
		-Run this as a GoRoutine so that it runs in the background
		-Watch the 3 channels (join, leave, forward)
		-r.clients map will only be modified once at a time bcuz of 
		select statement
	*/
	for {
		select {
		// client is a reference
		case client := <-r.join:
			r.clients[client] = true

		case client := <-r.leave:
			delete(r.clients, client)
			close(client.send)

		// fwd message to all clients in room		
		case msg := <-r.forward:

			fmt.Println(msg)
			// TODO: Remove
			// Logic: msg is the bridge between..
			// client.send channel and room.forward channel
			for client := range r.clients {
				select {
				// push to client channel
				case client.send <- msg:

				// failed to contact client
				default:
					log.Println("client closed the connection, deleting..")
					delete(r.clients, client)
					close(client.send)
				}
			}
		}
	}
}

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)
var upgrader = &websocket.Upgrader{
	ReadBufferSize: socketBufferSize,
	WriteBufferSize: messageBufferSize,
}

func (r *room) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Fatal("ServeHTTP:", err)
		return
	}
	//create client
	client := &client {
		socket: socket,
		send: make(chan []byte, messageBufferSize),
		room: r,
	}
	r.join <- client
	defer func () {r.leave <- client} ()
	// background thread to write
	// main thread reads
	go client.write()
	client.read()

}
