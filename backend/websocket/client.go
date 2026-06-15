package websocket

import (
	"encoding/json"
	"errors"
	"math/rand"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/google/uuid"
)

var ErrSendBufferFull = errors.New("send buffer is full")

var colorPalette = []string{
	"#FF6B6B", "#4ECDC4", "#45B7D1", "#96CEB4", "#FFEAA7",
	"#DDA0DD", "#98D8C8", "#F7DC6F", "#BB8FCE", "#85C1E9",
}

func generateColor() string {
	src := rand.NewSource(time.Now().UnixNano())
	r := rand.New(src)
	return colorPalette[r.Intn(len(colorPalette))]
}

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 8192
)

type Client struct {
	ID           string
	UserID       string
	Username     string
	Color        string
	PointCloudID string
	Conn         *websocket.Conn
	Send         chan *WSMessage
	Hub          *Hub
	mu           sync.Mutex
}

func NewClient(conn *websocket.Conn, hub *Hub, userID, username, pointCloudID string) *Client {
	return &Client{
		ID:           uuid.New().String(),
		UserID:       userID,
		Username:     username,
		Color:        generateColor(),
		PointCloudID: pointCloudID,
		Conn:         conn,
		Send:         make(chan *WSMessage, 256),
		Hub:          hub,
	}
}

func (c *Client) ReadPump() {
	defer func() {
		c.Hub.Unregister <- c
		c.Conn.Close()
	}()

	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, rawMsg, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
			}
			break
		}

		var msg WSMessage
		if err := json.Unmarshal(rawMsg, &msg); err != nil {
			continue
		}

		switch msg.Type {
		case MsgTypePing:
			pongMsg := &WSMessage{
				Type:      MsgTypePong,
				Timestamp: time.Now().UnixMilli(),
			}
			c.SendMessage(pongMsg)
		default:
			c.Hub.Broadcast <- &BroadcastMessage{
				PointCloudID:  c.PointCloudID,
				Message:       &msg,
				ExcludeClient: nil,
			}
		}
	}
}

func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case msg, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}

			data, _ := json.Marshal(msg)
			w.Write(data)

			n := len(c.Send)
			for i := 0; i < n; i++ {
				extra := <-c.Send
				extraData, _ := json.Marshal(extra)
				w.Write([]byte{'\n'})
				w.Write(extraData)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (c *Client) SendMessage(msg *WSMessage) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	select {
	case c.Send <- msg:
		return nil
	default:
		return ErrSendBufferFull
	}
}
