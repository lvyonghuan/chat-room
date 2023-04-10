package main

import (
	"github.com/gin-gonic/gin"
)

func InitRouters() {
	r := gin.Default()
	r.GET("/chatroom", ChatRoom)
	r.Run()
}

func ChatRoom(c *gin.Context) {
	Start(c)
}
