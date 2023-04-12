package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"sync"
	"time"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

var (
	lock = sync.Mutex{}
	mu   sync.Mutex
)

//func pingPong(c *websocket.Conn) {
//	// 定义最后一次收到 pong 消息的时间戳
//	lastPong := time.Now()
//	// 设置 Pong 处理函数，更新最后一次收到 pong 消息的时间戳
//	c.SetPongHandler(func(string) error {
//		lastPong = time.Now()
//		return nil
//	})
//	// 启动一个新的协程，定期向客户端发送 Ping 消息
//	go func() {
//		ticker := time.NewTicker(pingPeriod)
//		defer ticker.Stop()
//		for {
//			select {
//			case <-ticker.C:
//				if err := c.WriteMessage(websocket.PingMessage, nil); err != nil {
//					log.Println("ping error", err)
//					return
//				}
//				// 检查最后一次收到 pong 消息的时间戳
//				if time.Since(lastPong) > pongWait {
//					log.Printf("pong timeout, connection lost to %s\n", c.RemoteAddr())
//					_ = c.Close()
//					return
//				}
//			}
//		}
//	}()
//}

func Read(c *Client) {
	defer func() {
		log.Printf("user %s exit the room\n", c.Username)
		c.Room.Unregister <- c
		err := c.Conn.Close()
		if err != nil {
			log.Println(err)
			return
		}
	}()
	c.Conn.SetReadLimit(maxMessageSize)
	//err := c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	//if err != nil {
	//	log.Println(err)
	//	return
	//}
	c.Conn.SetPongHandler(func(string) error {
		//err := c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		//if err != nil {
		//	log.Println(err)
		//	return err
		//}
		return nil
	})
	for {
		select {
		case msg := <-c.Send:
			err := c.Conn.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				log.Println("write msg failed,err:", err)
				return
			}
		}
	}
}

func Write(c *Client) {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		log.Printf("user %s exit the room\n", c.Username)
		ticker.Stop()
		c.Room.Unregister <- c
		err := c.Conn.Close()
		if err != nil {
			log.Println(err)
			return
		}
	}()
	go func() {
		for {
			select {
			case <-ticker.C:
				//err := c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
				//if err != nil {
				//	log.Println(err)
				//	return
				//}
				if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
					log.Println(err)
					return
				}
			}
		}
	}()
	for {
		msgType, msgByte, err := c.Conn.ReadMessage()
		if err != nil {
			// 这里遇到错误一般是断开websocket链接，不管怎样，咱们关闭链接就是了
			log.Println("read msg failed, err:", err)
			break
		}
		//发言频次检测
		//log.Println(c.Room.Bucket)
		if c.Room.Bucket == 0 {
			c.Room.Broadcast <- []byte("时间段内发言次数已达到上限")
			continue
		}
		// 这里只处理一个消息类型
		switch msgType {
		case websocket.TextMessage:
			msg := []byte(fmt.Sprintf("%s %s说:%s", time.Now().Format("01/02 03:04"), c.Username, string(msgByte)))
			RoomList := data.GetBuilding(c.Room.RoomID)
			if err != nil {
				log.Println("110:", err)
			}
			c.Room.Bucket--
			for i := range RoomList {
				RoomList[i].Broadcast <- msg
			}
		default:
			log.Println("receive don't know msg type is ", msgType)
			continue
		}
	}
}
