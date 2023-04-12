package main

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"
)

var Upgrade = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func NewHub(c Client, roomID int) *User {
	room := &User{
		Broadcast:  c.Send,
		Clients:    sync.Map{},
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Username:   c.Username,
		RoomID:     roomID,
		Bucket:     10,
	}
	room.Clients.Store(c, true)
	return room
}

func IntoRoom(ctx *gin.Context) {
	username := ctx.Param("username")
	roomIDString := ctx.Param("roomID")
	roomId, err := strconv.Atoi(roomIDString)
	if err != nil {
		log.Println(err)
		return
	}
	conn, err := Upgrade.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Println("upgrade req failed, err:", err)
		ctx.JSON(http.StatusInternalServerError, "upgrade failed")
		return
	}
	c := Client{
		Conn:     conn,
		Username: username,
		Send:     make(chan []byte, 1024),
	}
	// 防止并发
	lock.Lock()
	hub := NewHub(c, roomId)
	c.Room = hub
	data.AddRoom(hub.RoomID, hub)
	go hub.run()
	lock.Unlock()
	hub.Register <- &c
	// 开俩协程读写消息
	go Read(&c)
	go Write(&c)
	go pingPong(c.Conn)
}

func (h *User) run() {
	var bucketLock sync.Mutex
	go func() {
		for {
			time.Sleep(60 * time.Second)
			bucketLock.Lock()
			h.Bucket = 10 //每隔60秒，桶中可用的发言机会重置
			bucketLock.Unlock()
		}
	}()
	for {
		select {
		case client := <-h.Register:
			h.Clients.Store(client, true)
		case client := <-h.Unregister:
			if _, ok := h.Clients.Load(client); ok {
				data.DeleteRoom(h.Username)
				h.Clients.Delete(client)
				close(client.Send)
			}
		case message := <-h.Broadcast:
			if message == nil {
				continue
			}
			h.Clients.Range(func(key, value interface{}) bool {
				client, ok := key.(Client)
				if !ok {
					log.Println("Failed to convert key to Client:", key)
				}
				client.Send <- message
				return true
			})
		}
	}
}
