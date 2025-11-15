// apps/backend/internal/infrastructure/websocket/message.go

package websocket

import "time"

// MessageType نوع پیام WebSocket را مشخص می‌کند.
type MessageType string

const (
	// پیام‌های ارسالی از کلاینت
	TypeJoinRoom     MessageType = "join_room"
	TypeLeaveRoom    MessageType = "leave_room"
	TypeSendMessage  MessageType = "send_message"
	
	// پیام‌های ارسالی از سرور
	TypeNewMessage   MessageType = "new_message"
	TypeError        MessageType = "error"
	TypeNotification MessageType = "notification"
)

// InboundMessage ساختار پیام‌های ورودی از کلاینت است.
type InboundMessage struct {
	Type    MessageType     `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

// OutboundMessage ساختار پیام‌های خروجی به کلاینت است.
type OutboundMessage struct {
	Type    MessageType `json:"type"`
	Payload interface{} `json:"payload"`
}

// Payloads
type JoinRoomPayload struct {
	RoomID string `json:"roomId"`
}

type SendMessagePayload struct {
	RoomID string `json:"roomId"`
	Text   string `json:"text"`
}

type NewMessagePayload struct {
	ID        string    `json:"id"`
	RoomID    string    `json:"roomId"`
	UserID    string    `json:"userId"`
	Username  string    `json:"username"`
	Text      string    `json:"text"`
	Timestamp time.Time `json:"timestamp"`
}