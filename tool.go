package main

//func serializedRoom(r *Room, c *Client) (jsonBytes []byte) {
//	// 将 *websocket.Conn 转换为字符串并添加到 ConnStrs 切片中
//	connStr := c.Conn.RemoteAddr().String()
//	r.ConnStrs = append(r.ConnStrs, connStr)
//	// 将 room 序列化为 JSON 格式
//	jsonBytes, err := json.Marshal(r)
//	if err != nil {
//		log.Println(err)
//		panic(err)
//	}
//	return
//}
//
//func deSerializedRoom(jsonBytes []byte) (room2 Room, err error) {
//	// 将房间从 JSON 格式反序列化回来
//	err = json.Unmarshal(jsonBytes, &room2)
//	if err != nil {
//		panic(err)
//	}
//
//	// 将 ConnStrs 中的字符串转换为 *websocket.Conn 类型
//	for _, connStr := range room2.ConnStrs {
//		// 使用适当的反序列化逻辑将字符串转换回 *websocket.Conn 类型
//		conn, _, err := websocket.DefaultDialer.Dial(connStr, nil)
//		if err != nil {
//			panic(err)
//		}
//		client := &Client{
//			Conn:     conn,
//			Room:     &room2,
//			Username: "",
//			Send:     make(chan []byte, 256),
//		}
//		room2.Register <- client
//	}
//	return
//}
