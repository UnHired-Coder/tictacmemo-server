package handlers

import (
	"encoding/json"
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
		if err := db.Where("id = ?", playerID).First(&user).Error; err != nil {
			sendWebSocketError(conn, "User not found.")
			return
		}

		log.Printf("Client connected: PlayerID=%s, WaitlistID=%s\n", playerID, waitlistID)

		PLAYERS_WAITLIST[waitlistID] = conn

		for {
			// Read a message from the WebSocket connection
			_, msg, err := conn.ReadMessage()
			if err != nil {
				log.Println("Error reading message:", err)
				break
			}

			processMatchingWebSocketMessage(db, gameManager, conn, msg)
		}
	}
}

// Processes WebSocket messages and performs actions based on the "action" field.
func processMatchingWebSocketMessage(db *gorm.DB, gameManager *types.TicTacMemoGameManager, conn *websocket.Conn, msg []byte) {
	log.Printf("Received message from client: %s\n", msg)

	var message types.GameEvent
	if err := json.Unmarshal(msg, &message); err != nil {
		log.Println("Error unmarshaling message:", err)
		return
	}

	switch message.Action {
	case types.ActionJoinRoom:
		var joinData types.JoinRoomData
		if err := json.Unmarshal(message.Data, &joinData); err != nil {
			log.Println("Error unmarshaling join-room data:", err)
			return
		}

		gameManager.JoinRoom(db, joinData)

		joinedRoomData := gin.H{
			"event": "joined-room",
			"data": gin.H{
				"room_id":   joinData.RoomID,
				"player_id": joinData.PlayerID,
			},
		}

		conn.WriteJSON(joinedRoomData)
		conn.Close()
	default:
		log.Println("Unknown action:", message.Action)
		conn.WriteMessage(websocket.TextMessage, []byte("Unknown action"))
	}
}
