package main

import (
	"github.com/gorilla/websocket"
	"sync"
)

type Client struct {
	Conn     *websocket.Conn
	Room     *User
	Username string
	Send     chan []byte
}

type User struct {
	Clients    sync.Map
	Broadcast  chan []byte
	Username   string
	RoomID     int
	Register   chan *Client
	Unregister chan *Client
	Bucket     int //桶实现发言频次的限制
}

// Data 树（大概吧）
type Data struct {
	roomMap     sync.Map
	usernameMap sync.Map
}

type Room struct {
	Users []*User
	// 可能还有其他字段，比如建筑名称、地址等
}
