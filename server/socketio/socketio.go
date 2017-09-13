package socketio

import (
	"log"

	"github.com/graarh/golang-socketio"
	"github.com/graarh/golang-socketio/transport"
)

var rooms = []string{"chat", "room2", "room3"}

// func main()  {
func InitHandler() *gosocketio.Server {
	//create
	server := gosocketio.NewServer(transport.GetDefaultWebsocketTransport())

	//handle connected
	server.On(gosocketio.OnConnection, func(c *gosocketio.Channel) {
		log.Println("New client connected")
		//join them to the default room
		c.Join(rooms[0])

		// 匿名结构体
		var eventData struct {
			Rooms       []string `json:"rooms"`
			CurrentRoom string   `json:"currentRoom"`
		}
		eventData.Rooms = rooms
		eventData.CurrentRoom = rooms[0]

		c.Emit("updateRooms", eventData)
	})

	type Message struct {
		Name    string `json:"name"`
		Message string `json:"message"`
	}

	//handle custom event
	server.On("message", func(c *gosocketio.Channel, msg Message) string {
		log.Println("message： ", msg)
		//send event to all in room
		c.BroadcastTo("chat", "message", msg)
		return "OK"
	})

	/*
		//setup http server
		serveMux := http.NewServeMux()
		serveMux.Handle("/socket.io/", server)
		// log.Panic(http.ListenAndServe(":80", serveMux))

		return serveMux
	*/
	return server

}
