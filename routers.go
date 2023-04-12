package main

import (
	"github.com/gin-gonic/gin"
)

func InitRouters() {
	r := gin.Default()
	r.GET("/chatroom/:username/:roomID", NewChatRoom)
	r.Run()
}

func NewChatRoom(c *gin.Context) {
	IntoRoom(c)
}
