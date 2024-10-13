package handlers

import (
	"encoding/json"
	"fmt"
	commonTypes "game-server/common/types"
	"game-server/tictacmemo/types"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"
)

var PLAYERS_WAITLIST = make(map[string]*websocket.Conn)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,

	// CheckOrigin allows connections from any origin (can be customized for security).
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

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

			processWebSocketMessage(db, gameManager, conn, msg)
		}
	}
}

// Processes WebSocket messages and performs actions based on the "action" field.
func processWebSocketMessage(db *gorm.DB, gameManager *types.TicTacMemoGameManager, conn *websocket.Conn, msg []byte) {
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

		response := fmt.Sprintf("User %d successfully joined room %s", joinData.PlayerID, joinData.RoomID)
		conn.WriteMessage(websocket.TextMessage, []byte(response))

	case types.ActionMakeMove:
		var makeMoveData types.MakeMoveData
		if err := json.Unmarshal(message.Data, &makeMoveData); err != nil {
			log.Println("Error unmarshaling make-move data:", err)
			return
		}

		gameManager.Rooms[makeMoveData.RoomID].MakeMove(db, makeMoveData)

		response := fmt.Sprintf("User successfully joined room %s", gameManager.Rooms[makeMoveData.RoomID].ID)
		conn.WriteMessage(websocket.TextMessage, []byte(response))

	default:
		log.Println("Unknown action:", message.Action)
		conn.WriteMessage(websocket.TextMessage, []byte("Unknown action"))
	}
}

func sendWebSocketError(conn *websocket.Conn, errorMsg string) error {
	errorResponse := map[string]string{"error": errorMsg}
	errorJSON, err := json.Marshal(errorResponse)
	if err != nil {
		return err
	}
	return conn.WriteMessage(websocket.TextMessage, errorJSON)
}
