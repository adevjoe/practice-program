package main

type Room struct {
	ID     int
	Name   string
	User   []string
	Active bool
	Limit  int
}

const DefaultRoomLimit int = 2
