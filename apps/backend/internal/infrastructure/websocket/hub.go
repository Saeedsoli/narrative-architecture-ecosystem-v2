// apps/backend/internal/infrastructure/websocket/hub.go

package websocket

import "log"

// Hub تمام کلاینت‌های فعال و اتاق‌ها را مدیریت می‌کند.
type Hub struct {
	clients    map[*Client]bool
	rooms      map[string]map[*Client]bool
	broadcast  chan *OutboundMessage
	register   chan *Client
	unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		rooms:      make(map[string]map[*Client]bool),
		broadcast:  make(chan *OutboundMessage),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

// Run حلقه اصلی Hub را برای مدیریت رویدادها اجرا می‌کند.
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
			log.Println("INFO: New client connected. Total clients:", len(h.clients))
		
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
				h.removeClientFromAllRooms(client)
				log.Println("INFO: Client disconnected. Total clients:", len(h.clients))
			}
		
		case message := <-h.broadcast:
			// این پیام به تمام کلاینت‌های یک اتاق خاص ارسال می‌شود.
			if roomID, ok := message.Payload.(map[string]interface{})["roomId"].(string); ok {
				if room, ok := h.rooms[roomID]; ok {
					for client := range room {
						select {
						case client.send <- message:
						default:
							close(client.send)
							delete(h.clients, client)
							delete(room, client)
						}
					}
				}
			}
		}
	}
}

func (h *Hub) removeClientFromAllRooms(client *Client) {
	for roomID, room := range h.rooms {
		if _, ok := room[client]; ok {
			delete(room, client)
			if len(room) == 0 {
				delete(h.rooms, roomID)
			}
		}
	}
}