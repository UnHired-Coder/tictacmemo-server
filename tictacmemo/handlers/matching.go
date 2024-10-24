package handlers

import (
	commonTypes "game-server/common/types"
	"game-server/tictacmemo/types"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"
)

var PLAYERS_WAITLIST = make(map[string]*websocket.Conn)

func Matching(db *gorm.DB, gameManager *types.TicTacMemoGameManager) gin.HandlerFunc {
	return func(c *gin.Context) {

		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Println("Failed to upgrade to websocket:", err)
			return
		}
		defer conn.Close()

		playerID := c.Param("playerID")
		waitlistID := c.Param("waitlistID")

		var user commonTypes.User
		if err := db.Where("user_id = ?", playerID).First(&user).Error; err != nil {
			sendWebSocketError(conn, "User not found.")
			return
		}

		log.Printf("Client connected: PlayerID=%s, WaitlistID=%s\n", playerID, waitlistID)

		PLAYERS_WAITLIST[waitlistID] = conn

		for {
			// Read a message from the WebSocket connection
			_, _, err := conn.ReadMessage()
			if err != nil {
				log.Println("Error reading message:", err)
				break
			}
		}
	}
}
