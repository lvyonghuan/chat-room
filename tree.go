package main

var data = &Data{}

func (d *Data) AddRoom(roomID int, room *Room) {
	roomDataInterface, _ := d.roomMap.LoadOrStore(roomID, &Building{Rooms: []*Room{}})
	roomData := roomDataInterface.(*Building)
	d.usernameMap.Store(room.Username, room)
	roomData.Rooms = append(roomData.Rooms, room)
}

func (d *Data) GetBuilding(building int) []*Room {
	buildingDataInterface, ok := d.roomMap.Load(building)
	if ok {
		buildingData := buildingDataInterface.(*Building)
		return buildingData.Rooms
	}
	return nil
}

func (d *Data) GetRoom(roomID int) *Room {
	var room *Room
	d.roomMap.Range(func(key, value interface{}) bool {
		buildingData := value.(*Building)
		for _, r := range buildingData.Rooms {
			if r.RoomID == roomID {
				room = r
				return false
			}
		}
		return true
	})
	return room
}

func (d *Data) DeleteRoom(username string) {
	roomInterface, ok := d.usernameMap.Load(username)
	if ok {
		room := roomInterface.(*Room)
		d.usernameMap.Delete(username)
		buildingDataInterface, ok := d.roomMap.Load(room.RoomID)
		if ok {
			buildingData := buildingDataInterface.(*Building)
			for i, r := range buildingData.Rooms {
				if r == room {
					buildingData.Rooms = append(buildingData.Rooms[:i], buildingData.Rooms[i+1:]...)
					break
				}
			}
		}
	}
}

//怀旧作
//import (
//	"github.com/gin-gonic/gin"
//	"github.com/redis/go-redis/v9"
//	"log"
//	"strconv"
//	"time"
//)
//
////用redis写的一个储存用户和房间信息的数据库
//
//var rds *redis.Client
//
//func InitRedis() {
//	rds = redis.NewClient(&redis.Options{
//		Addr:     "localhost:6379",
//		DB:       0,
//		Password: "",
//		PoolSize: 10,
//	})
//}
//
////房间相关操作
//
//// StoreRoomIDAndUser 匹配用户信息和房间id
//func StoreRoomIDAndUser(c *gin.Context, user string, roomID int) (err error) {
//	err = rds.Set(c, user, roomID, 3600*24*time.Second).Err()
//	return err
//}
//
//// StoreRoomAndRoomID 储存房间id下的所有房间
//func StoreRoomAndRoomID(c *gin.Context, roomID int, serializedRoom []byte) (err error) {
//	//serializedRoom, err := json.Marshal(room) //序列化房间信息
//	if err != nil {
//		return err
//	}
//	err = rds.SAdd(c, strconv.Itoa(roomID), serializedRoom).Err() //储存ID和房间信息到一个集合里
//	if err != nil {
//		return err
//	}
//	err = rds.Expire(c, strconv.Itoa(roomID), 3600*24*time.Second).Err() //一天之后过期
//	return err
//}
//
//// StoreRoomAndUser 存储用户名和房间信息
//func StoreRoomAndUser(c *gin.Context, user string, serializedRoom []byte) (err error) {
//	//serializedRoom := encode(room)                                    //序列化房间信息
//	err = rds.Set(c, user, serializedRoom, 3600*24*time.Second).Err() //储存用户名和房间信息,一天后过期
//	//log.Println(serializedRoom)
//	log.Println(serializedRoom)
//	return err
//}
//
//// SearchRoomIDByUser 根据用户名索引房间ID
//func SearchRoomIDByUser(c *gin.Context, user string) (roomID int, err error) {
//	roomID, err = rds.Get(c, user).Int()
//	return roomID, err
//}
//
//// SearchRoomByID 根据房间id索引房间
//func SearchRoomByID(c *gin.Context, roomID int) (room []Room, err error) {
//	serializedRoomSlice, err := rds.SMembers(c, strconv.Itoa(roomID)).Result()
//	log.Println(serializedRoomSlice)
//	if err != nil && err != redis.Nil {
//		return nil, err
//	}
//	var rooms []Room
//	for _, serializedRoom := range serializedRoomSlice {
//		var room *Room
//		*room, err = deSerializedRoom([]byte(serializedRoom))
//		log.Println(*room)
//		if err != nil {
//			return nil, err
//		}
//		if room.RoomID == roomID {
//			rooms = append(rooms, *room)
//		}
//	}
//	return rooms, err
//}
//
//func SearchRoomByUsername(c *gin.Context, username string) (serializedRoom string, err error) {
//	serializedRoomByte, err := rds.Get(c, username).Bytes()
//	return string(serializedRoomByte), err
//}
//
//func DeleteRoom(c *gin.Context, user string) (err error) {
//	roomID, err := SearchRoomIDByUser(c, user)
//	if err != nil {
//		return err
//	}
//	err = rds.Del(c, user).Err()
//	if err != nil {
//		return err
//	}
//	serializedRoom, err := SearchRoomByUsername(c, user)
//	if err != nil {
//		return err
//	}
//	err = rds.SRem(c, strconv.Itoa(roomID), serializedRoom).Err()
//	return err
//}

//真不能处。摆了，留作纪念
