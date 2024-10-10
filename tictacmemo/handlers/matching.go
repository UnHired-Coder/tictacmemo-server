package handlers

import (
	"encoding/json"
	"game-server/common/types"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// CheckOrigin allows connections from any origin (can be customized for security).
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func Matching(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Upgrade the HTTP request to a WebSocket connection
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Println("Failed to upgrade to websocket:", err)
			return
		}
		defer conn.Close()

		// Extract playerID and waitlistID from the URL parameters
		playerID := c.Param("playerID")
		waitlistID := c.Param("waitlistID")

		var user types.User
		if err := db.Where("id = ?", playerID).First(&user).Error; err != nil {
			sendWebSocketError(conn, "User not found.")
			return
		}

		log.Printf("Client connected: PlayerID=%s, WaitlistID=%s\n", playerID, waitlistID)

		PLAYERS_WAITLIST[waitlistID] = conn

		// TODO:: On receiving the room details the client
		// will send back a message and the match will start.

		// Handle incoming messages from the client
		for {
			// Read a message from the WebSocket connection
			_, msg, err := conn.ReadMessage()
			if err != nil {
				log.Println("Error reading message:", err)
				break
			}
			log.Printf("Received message from client: %s\n", msg)

			// Optionally echo the message back to the client
			err = conn.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				log.Println("Error sending message:", err)
				break
			}
		}
	}
}

func sendWebSocketError(conn *websocket.Conn, errorMsg string) error {
	// Create a map similar to gin.H to hold the error message
	errorResponse := map[string]string{"error": errorMsg}

	// Marshal the map to a JSON string
	errorJSON, err := json.Marshal(errorResponse)
	if err != nil {
		return err
	}

	// Write the error message to the WebSocket connection
	return conn.WriteMessage(websocket.TextMessage, errorJSON)
}
