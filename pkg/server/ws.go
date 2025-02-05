package server

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/wafi04/backend/pkg/types"
)

func WebSocketHandler(c *gin.Context) {
	ws, err := types.Upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}
	defer ws.Close()
	types.Clients[ws] = true

	for {
		_, msg, err := ws.ReadMessage()
		if err != nil {
			log.Println("WebSocket read error:", err)
			delete(types.Clients, ws)
			break
		}
		log.Printf("Received message: %s", msg)
	}
}

func BroadcastMessages() {
	for {
		msg := <-types.Broadcast
		for client := range types.Clients {
			err := client.WriteMessage(websocket.TextMessage, []byte(msg))
			if err != nil {
				log.Println("WebSocket write error:", err)
				client.Close()
				delete(types.Clients, client)
			}
		}
	}
}
