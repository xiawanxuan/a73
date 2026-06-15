package websocket

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func WsHandler(hub *Hub) gin.HandlerFunc {
	return func(c *gin.Context) {
		pointCloudID := c.Query("point_cloud_id")
		userID := c.Query("user_id")
		username := c.Query("username")

		if pointCloudID == "" || userID == "" || username == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "point_cloud_id, user_id and username are required"})
			return
		}

		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			return
		}

		client := NewClient(conn, hub, userID, username, pointCloudID)
		hub.Register <- client

		joinData := UserJoinData{
			UserID:   client.UserID,
			Username: client.Username,
			Color:    client.Color,
		}
		joinDataBytes, _ := json.Marshal(joinData)
		joinMsg := &WSMessage{
			Type:      MsgTypeUserJoin,
			Data:      joinDataBytes,
			Timestamp: time.Now().UnixMilli(),
		}
		hub.Broadcast <- &BroadcastMessage{
			PointCloudID:  pointCloudID,
			Message:       joinMsg,
			ExcludeClient: client,
		}

		existingClients := hub.GetRoomUsers(pointCloudID)
		onlineUsers := make([]UserJoinData, 0, len(existingClients))
		for _, ec := range existingClients {
			onlineUsers = append(onlineUsers, UserJoinData{
				UserID:   ec.UserID,
				Username: ec.Username,
				Color:    ec.Color,
			})
		}
		onlineDataBytes, _ := json.Marshal(OnlineUsersData{Users: onlineUsers})
		onlineMsg := &WSMessage{
			Type:      MsgTypeOnlineUsers,
			Data:      onlineDataBytes,
			Timestamp: time.Now().UnixMilli(),
		}
		client.SendMessage(onlineMsg)

		go client.WritePump()
		go client.ReadPump()
	}
}
