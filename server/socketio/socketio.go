package socketio

import (
	"fmt"
	"log"

	"github.com/graarh/golang-socketio"
	"github.com/graarh/golang-socketio/transport"
)

type Message struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

type ClientUser struct {
	UserID   string `json:"userId"`
	Username string `json:"username"`
}

var rooms = []string{"chat", "room2", "room3"}

//在线用户
var onlineUsers = make(map[string]string)

//当前在线人数
var onlineCount = 0

// 存储channel(在socketio中实际对应的是socket) 和 用户信息的映射
var channelUserMap = make(map[string]ClientUser)

// 不准备支持 一个用户多个聊天室  只能是唯一一个当前存活的聊天室
var channelRoomMap = make(map[string]string)

func userInfoByChannel(c *gosocketio.Channel) (room string, clientUser ClientUser) {
	room = channelRoomMap[c.Id()]
	clientUser = channelUserMap[c.Id()]

	return
}

func updateRooms(c *gosocketio.Channel, newRoom string) {
	// 匿名结构体
	var eventData struct {
		Rooms       []string `json:"rooms"`
		CurrentRoom string   `json:"currentRoom"`
	}
	eventData.Rooms = rooms
	eventData.CurrentRoom = newRoom

	c.Emit("updateRooms", eventData)
}

// func main()  {
func InitHandler() *gosocketio.Server {
	//create
	server := gosocketio.NewServer(transport.GetDefaultWebsocketTransport())

	//handle connected
	server.On(gosocketio.OnConnection, func(c *gosocketio.Channel) {
		log.Println("New client connected")
		//join them to the default room
		c.Join(rooms[0])
		// 同时写入自己的存储中
		channelRoomMap[c.Id()] = rooms[0]

		updateRooms(c, rooms[0])
		/*
			// 匿名结构体
			var eventData struct {
				Rooms       []string `json:"rooms"`
				CurrentRoom string   `json:"currentRoom"`
			}
			eventData.Rooms = rooms
			eventData.CurrentRoom = rooms[0]

			c.Emit("updateRooms", eventData)
		*/
	})

	getUsersInRoom := func(room string, currentChannel *gosocketio.Channel) []ClientUser {
		var channelsInRoom []*gosocketio.Channel
		channelsInRoom = server.List(room)

		fmt.Println(" room: ", room, " channels : ", len(channelsInRoom))

		var userList []ClientUser
		// 遍历当前房间下的往昔客户channel
		for _, ch := range channelsInRoom {
			if currentChannel != nil {
				// 忽略当前的用户
				if currentChannel.Id() == ch.Id() {
					continue
				}
			}
			userList = append(userList, channelUserMap[ch.Id()])
		}
		return userList
	}

	//handle user login
	server.On("addUser", func(c *gosocketio.Channel, user ClientUser) string {
		log.Println("userInfo ： ", user)
		/*
			var channelsInRoom []*gosocketio.Channel
			channelsInRoom = server.List(rooms[0])
			var userList []ClientUser
			// 遍历当前房间下的往昔客户channel
			for _, ch := range channelsInRoom {
				userList = append(userList, channelUserMap[ch.Id()])
			}
		*/
		room := rooms[0] // 默认的房间
		userList := getUsersInRoom(room, c)
		if len(userList) > 0 {
			c.Emit("usersInRoom", userList)
		}

		channelUserMap[c.Id()] = user
		//send event to all in room
		// TODO : 对当前用户 需要发送过往所有的用户列表
		// 对已经在房间的用户 需要发送新加入了用户信息
		c.BroadcastTo(room, "addUser", Message{Type: "info", Data: user})

		// 系统通知
		msg := Message{
			Type: "sysinfo",
			Data: fmt.Sprintf("user %s joined room: %s ", user.Username, room),
		}
		//send event to all in room
		c.BroadcastTo(room, "message", msg)

		return "OK"
	})

	// handle custom event
	server.On("message", func(c *gosocketio.Channel, msg Message) string {
		log.Println("message： ", msg)
		//send event to all in room
		c.BroadcastTo("chat", "message", msg)
		return "OK"
	})

	// handle chat event
	server.On("chatMessage", func(c *gosocketio.Channel, m string) string {
		log.Println("message： ", m)

		clientUser := channelUserMap[c.Id()]
		room := channelRoomMap[c.Id()]

		msg := Message{
			Type: "chat",
			Data: fmt.Sprintf("%s say: %s ", clientUser.Username, m),
		}
		//send event to all in room
		c.BroadcastTo(room, "message", msg)
		return "OK"
	})

	server.On("switchRoom", func(c *gosocketio.Channel, r string) string {
		room, clientUser := userInfoByChannel(c)
		log.Println("switch room from ", room, " to ", r)

		// 离开旧房 换新房
		// 老房间要先通知人离开了
		c.BroadcastTo(room, "leaveRoom", clientUser)
		c.Leave(room)
		c.Join(r)
		// 房间映射更改
		channelRoomMap[c.Id()] = r
		// 刷新房间列表
		updateRooms(c, r)
		// 房间用户列表更新
		userList := getUsersInRoom(r, nil)
		if len(userList) > 0 {
			c.Emit("usersInRoom", userList)
		}

		msg := Message{
			Type: "sysinfo",
			Data: fmt.Sprintf("%s switch to : %s ", clientUser.Username, r),
		}
		//send event to all in room
		c.BroadcastTo(room, "message", msg)
		// 新房间发送信息
		c.BroadcastTo(r, "message", Message{
			Type: "sysinfo",
			Data: fmt.Sprintf("%s join in %s", clientUser.Username, r),
		})
		c.BroadcastTo(r, "addUser", Message{Type: "info", Data: clientUser})
		return "OK"
	})
	// handle system event
	server.On(gosocketio.OnDisconnection, func(c *gosocketio.Channel) string {
		clientUser := channelUserMap[c.Id()]
		room := channelRoomMap[c.Id()]

		log.Println(fmt.Sprintf("user %s leave room: %s ", clientUser.Username, room))

		c.BroadcastTo(room, "usersInRoom", getUsersInRoom(room, c))

		// 清除当前用户占用的资源
		delete(channelUserMap, c.Id())
		delete(channelRoomMap, c.Id())

		msg := Message{
			Type: "sysinfo",
			Data: fmt.Sprintf("user %s leave room: %s ", clientUser.Username, room),
		}
		//send event to all in room
		c.BroadcastTo(room, "message", msg)

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
