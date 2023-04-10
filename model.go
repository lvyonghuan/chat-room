package main

import (
	"github.com/gorilla/websocket"
	"sync"
)

type Client struct {
	Conn     *websocket.Conn
	Room     *Room
	Username string
	Send     chan []byte
}

type Room struct {
	Clients    sync.Map
	Broadcast  chan []byte
	Username   string
	RoomID     int
	Register   chan *Client
	Unregister chan *Client
}

// Data 树（大概吧）
type Data struct {
	roomMap     sync.Map
	usernameMap sync.Map
}

type Building struct {
	Rooms []*Room
	// 可能还有其他字段，比如建筑名称、地址等
}
