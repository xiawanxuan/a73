package websocket

import (
	"encoding/json"
	"sync"
	"time"
)

type Hub struct {
	Rooms      map[string]map[*Client]bool
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan *BroadcastMessage
	BroadcastAll chan *WSMessage
	mu         sync.RWMutex
}

type BroadcastMessage struct {
	PointCloudID  string
	Message       *WSMessage
	ExcludeClient *Client
}

func NewHub() *Hub {
	return &Hub{
		Rooms:      make(map[string]map[*Client]bool),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan *BroadcastMessage),
		BroadcastAll: make(chan *WSMessage),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.JoinRoom(client.PointCloudID, client)
		case client := <-h.Unregister:
			pointCloudID := client.PointCloudID
			h.LeaveRoom(pointCloudID, client)
			leaveData := UserLeaveData{
				UserID: client.UserID,
			}
			leaveDataBytes, _ := json.Marshal(leaveData)
			leaveMsg := &WSMessage{
				Type:      MsgTypeUserLeave,
				Data:      leaveDataBytes,
				Timestamp: time.Now().UnixMilli(),
			}
			h.Broadcast <- &BroadcastMessage{
				PointCloudID:  pointCloudID,
				Message:       leaveMsg,
				ExcludeClient: nil,
			}
		case msg := <-h.Broadcast:
			h.mu.RLock()
			room, ok := h.Rooms[msg.PointCloudID]
			h.mu.RUnlock()
			if !ok {
				continue
			}
			for client := range room {
				if msg.ExcludeClient != nil && client == msg.ExcludeClient {
					continue
				}
				select {
				case client.Send <- msg.Message:
				default:
					h.LeaveRoom(client.PointCloudID, client)
				}
			}
		case msg := <-h.BroadcastAll:
			h.mu.RLock()
			allClients := make([]*Client, 0)
			for _, room := range h.Rooms {
				for client := range room {
					allClients = append(allClients, client)
				}
			}
			h.mu.RUnlock()
			for _, client := range allClients {
				select {
				case client.Send <- msg:
				default:
					h.LeaveRoom(client.PointCloudID, client)
				}
			}
		}
	}
}

func (h *Hub) JoinRoom(pointCloudID string, client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if _, ok := h.Rooms[pointCloudID]; !ok {
		h.Rooms[pointCloudID] = make(map[*Client]bool)
	}
	h.Rooms[pointCloudID][client] = true
	client.PointCloudID = pointCloudID
}

func (h *Hub) LeaveRoom(pointCloudID string, client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if room, ok := h.Rooms[pointCloudID]; ok {
		if _, ok := room[client]; ok {
			delete(room, client)
			close(client.Send)
			if len(room) == 0 {
				delete(h.Rooms, pointCloudID)
			}
		}
	}
}

func (h *Hub) GetRoomUsers(pointCloudID string) []*Client {
	h.mu.RLock()
	defer h.mu.RUnlock()
	room, ok := h.Rooms[pointCloudID]
	if !ok {
		return []*Client{}
	}
	users := make([]*Client, 0, len(room))
	for client := range room {
		users = append(users, client)
	}
	return users
}
