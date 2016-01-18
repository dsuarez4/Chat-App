package main

// Page 18
import (
	"github.com/gorilla/websocket"
)

// client object will represent a single chatting user
type client struct {

	// socket for this client
	socket *websocket.Conn
	// send is a channel on which messages are sent
	send chan []byte
	// room for the client
	room *room
}

func (c *client) read() {
	for {
		if _, msg, err := c.socket.ReadMessage(); err == nil {
			c.room.forward <- msg
		} else {
			break
		}
	}

	c.socket.Close()
}

// EDU: research range functionality
func (c *client) write() {
	for msg := range c.send {
		if err := c.socket.WriteMessage(websocket.TextMessage, msg); err != nil {
			break
		}
	}
	c.socket.Close()
}