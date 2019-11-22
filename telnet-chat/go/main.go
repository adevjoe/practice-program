package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
)

func main() {
	ln, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Printf("listen error: %v\n", err)
		return
	}
	log.Printf("start listening...\n")
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("%s connect fail, error: %v\n", conn.RemoteAddr(), err)
			continue
		}
		go handleConnection(NewClient(conn))
	}
}

func handleConnection(c *Client) {
	log.Printf("client %s connected.\n", c.UUID)
	c.Send("Please login. Use command '/help' to get help.")
	go func(c *Client) {
		for {
			select {
			case <-c.Done:
				return
			case msg := <-c.Message:
				_, _ = c.Conn.Write([]byte(msg))
			}
		}
	}(c)
	for {
		msg, err := c.Receive()
		if err != nil {
			if err == io.EOF {
				log.Printf("client %s disconnect.\n", c.UUID)
				return
			}
			log.Printf("client %s connect error: %v\n", c.UUID, err)
			return
		}
		if len(msg) > 0 {
			if processMsg(c, msg) == -1 { // connection closed
				return
			}
		}
	}
}

func processMsg(c *Client, msg string) int {
	msg = trimString(msg)
	command := ParseCommand(msg)
	if command.Name == "" { // chat
		c.Chat(msg)
		return 0
	}
	switch command.Name {
	case "login":
		if len(command.Args) != 2 {
			c.Send("Command error.")
			break
		}
		if checkPassword(command.Args[0], command.Args[1]) {
			c.User = getUserByUsername(command.Args[0])
			c.Send("Login successful.")
		} else {
			c.Send("Login fail.")
		}
	case "reg", "register":
		if len(command.Args) != 2 {
			c.Send("Command error.")
			break
		}
		err := register(command.Args[0], command.Args[1])
		if err != nil {
			c.Send(fmt.Sprintf("Error: %s", err.Error()))
		} else {
			c.Send("Register successful.")
		}
	case "i":
		if !c.Auth() {
			c.Send("Need login.")
			break
		}
		c.Send(fmt.Sprintf("Username: %s", c.User.Username))
	case "logout":
		if !c.Auth() {
			c.Send("Need login.")
			break
		}
		c.User = User{}
		c.Send("Logout successful.")
	case "rooms":
		if !c.Auth() {
			c.Send("Need login.")
			break
		}
		c.Send(fmt.Sprintf("rooms: %v", listRooms()))
	case "join":
		if !c.Auth() {
			c.Send("Need login.")
			break
		}
		if len(command.Args) != 1 {
			c.Send("Command error.")
			break
		}
		roomID, err := strconv.Atoi(command.Args[0])
		if err != nil {
			c.Send("room error.")
			break
		}
		room := getStoreRoomByID(roomID)
		if room.Name != "" && room.Active {
			count := getRoomUser(roomID)
			if count >= room.Limit {
				c.Send("This room is full up.")
				break
			}
			joinRoom(c.User.Username, roomID)
		}
		c.CurrentRoom = room
		c.Send(fmt.Sprintf("Welcome to room %d.", roomID))
	case "leave":
		if !c.Auth() {
			c.Send("Need login.")
			break
		}
		if c.CurrentRoom.ID == -1 {
			c.Send("You are not in the room.")
			break
		}
		leaveRoom(c.User.Username, c.CurrentRoom.ID)
		c.Send(fmt.Sprintf("Leave room %d.", c.CurrentRoom.ID))
		c.CurrentRoom = Room{ID: -1}
	case "create":
	case "del", "delete":
	case "exit":
		c.Close()
		return -1
	case "help":
		c.Send(`Chat Room.
Usage:
  /command [args...]
Commands:
   register <username> <password>    create an account
   login <username> <password>       login with your username and password
   logout                            logout your account
   i                                 get my info
   rooms                             list rooms
   join <room_id>                    join a room
   leave <room_id>                   leave a room
   exit                              exit this connection
   help                              get help tip
`)
	default:
		c.Send("Unknown command. Tap /help to get help.")
	}
	return 0
}
