// apps/backend/internal/infrastructure/websocket/client.go

package websocket

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
	"github.com/oklog/ulid/v2"
)

type Client struct {
	hub    *Hub
	conn   *websocket.Conn
	send   chan *OutboundMessage
	UserID string // شناسه کاربر پس از احراز هویت
}

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

// readPump پیام‌ها را از کلاینت می‌خواند و پردازش می‌کند.
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("ERROR: WebSocket read error: %v", err)
			}
			break
		}

		var msg InboundMessage
		if err := json.Unmarshal(message, &msg); err != nil {
			log.Printf("ERROR: Failed to unmarshal message: %v", err)
			continue
		}
		
		c.handleMessage(msg)
	}
}

// writePump پیام‌ها را از Hub دریافت و به کلاینت ارسال می‌کند.
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.conn.WriteJSON(message); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// handleMessage منطق پردازش پیام‌های ورودی را پیاده‌سازی می‌کند.
func (c *Client) handleMessage(msg InboundMessage) {
	switch msg.Type {
	case TypeJoinRoom:
		var p JoinRoomPayload
		if err := json.Unmarshal(msg.Payload, &p); err == nil {
			if c.hub.rooms[p.RoomID] == nil {
				c.hub.rooms[p.RoomID] = make(map[*Client]bool)
			}
			c.hub.rooms[p.RoomID][c] = true
			log.Printf("INFO: Client %s joined room %s", c.UserID, p.RoomID)
		}

	case TypeSendMessage:
		var p SendMessagePayload
		if err := json.Unmarshal(msg.Payload, &p); err == nil {
			// TODO: ذخیره پیام در دیتابیس
			
			// ایجاد پیام خروجی برای ارسال به همه اعضای اتاق
			outMsg := &OutboundMessage{
				Type: TypeNewMessage,
				Payload: NewMessagePayload{
					ID:       ulid.New().String(),
					RoomID:   p.RoomID,
					UserID:   c.UserID,
					Username: "some_username", // باید از اطلاعات کاربر گرفته شود
					Text:     p.Text,
					Timestamp: time.Now(),
				},
			}
			c.hub.broadcast <- outMsg
		}
	}
}