package main

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"sync"
)

var GlobalUsers sync.Map
var GlobalRooms = Rooms{
	Rooms:   make([]Room, 0),
	RoomMap: sync.Map{},
	Mux:     &sync.Mutex{},
}
var GlobalClient = Clients{
	Clients: make([]*Client, 0),
	Mux:     &sync.Mutex{},
}

type Clients struct {
	Clients []*Client
	Mux     *sync.Mutex
}

type Rooms struct {
	Rooms   []Room
	RoomMap sync.Map
	Mux     *sync.Mutex
}

func createRoomStore(name string, ) (Room, error) {
	GlobalRooms.Mux.Lock()
	defer GlobalRooms.Mux.Unlock()
	count := len(GlobalRooms.Rooms)
	room := Room{
		ID:     count,
		Name:   name,
		Active: true,
		Limit:  DefaultRoomLimit,
	}
	r := getStoreRoomByName(name)
	if r.Name != "" {
		return Room{}, errors.New("room exist")
	}
	GlobalRooms.Rooms = append(GlobalRooms.Rooms, room)
	GlobalRooms.RoomMap.Store(room.Name, room)
	return room, nil
}

func getStoreRoomByName(name string) Room {
	if r, ok := GlobalRooms.RoomMap.Load(name); ok {
		if room, ok := r.(Room); ok {
			return room
		}
	}
	return Room{}
}

func getStoreRoomByID(id int) Room {
	if len(GlobalRooms.Rooms)-1 >= id {
		return GlobalRooms.Rooms[id]
	}
	return Room{}
}

func listRooms() []string {
	var list []string
	for _, value := range GlobalRooms.Rooms {
		if value.Active {
			list = append(list, fmt.Sprintf("%d-%s(%d/%d)",
				value.ID, value.Name, getRoomUser(value.ID), value.Limit))
		}
	}
	return list
}

func joinRoom(username string, roomID int) {
	GlobalRooms.Mux.Lock()
	defer GlobalRooms.Mux.Unlock()
	if len(GlobalRooms.Rooms)-1 >= roomID {
		GlobalRooms.Rooms[roomID].User = append(GlobalRooms.Rooms[roomID].User, username)
	}
}

func leaveRoom(username string, roomID int) {
	GlobalRooms.Mux.Lock()
	defer GlobalRooms.Mux.Unlock()
	if len(GlobalRooms.Rooms)-1 >= roomID {
		for key, r := range GlobalRooms.Rooms[roomID].User {
			if r == username {
				GlobalRooms.Rooms[roomID].User = append(GlobalRooms.Rooms[roomID].User[:key], GlobalRooms.Rooms[roomID].User[key+1:]...)
			}
		}
	}
}

func getRoomUser(id int) int {
	r := getStoreRoomByID(id)
	return len(r.User)
}

func appendClient(c *Client) {
	GlobalClient.Mux.Lock()
	defer GlobalClient.Mux.Unlock()
	GlobalClient.Clients = append(GlobalClient.Clients, c)
}

func removeClient(uuid uuid.UUID) {
	GlobalClient.Mux.Lock()
	defer GlobalClient.Mux.Unlock()
	for key, c := range GlobalClient.Clients {
		if c.UUID == uuid {
			GlobalClient.Clients = append(GlobalClient.Clients[:key], GlobalClient.Clients[key+1:]...)
		}
	}
}

func sendRoom(username string, roomID int, msg string) {
	r := getStoreRoomByID(roomID)
	for _, c := range GlobalClient.Clients {
		if c.CurrentRoom.ID == roomID && c.CurrentRoom.Name == r.Name {
			go func(c *Client, r Room) {
				c.Message <- fmt.Sprintf("[ (%d)%s ] %s: %s\n> ", r.ID, r.Name, username, msg)
			}(c, r)
		}
	}
}
