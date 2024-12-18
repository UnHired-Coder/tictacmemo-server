package handlers

import (
	"encoding/json"
	"game-server/tictacmemo/types"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"
)

func Room(db *gorm.DB, gameManager *types.TicTacMemoGameManager) gin.HandlerFunc {
	return func(c *gin.Context) {

		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Println("Failed to upgrade to websocket:", err)
			return
		}
		defer conn.Close()

		roomID := c.Param("roomID")

		roomIDuuid, err := uuid.Parse(roomID)
		if err != nil {
			log.Println("Invalid room ID: Parsing", err)
		}

		room, exists := gameManager.Rooms[roomIDuuid]

		if !exists {
			error := gin.H{
				"event": "error",
				"data": gin.H{
					"message": "Room does not exist!",
					"room_id": roomIDuuid,
				},
			}
			conn.WriteJSON(error)
		}

		room.Clients = append(room.Clients, conn)

		for {
			// Read a message from the WebSocket connection
			_, msg, err := conn.ReadMessage()
			if err != nil {
				log.Println("Error reading message:", err)
				break
			}

			processRoomWebSocketMessage(db, gameManager, conn, msg, roomIDuuid)
		}
	}
}

// Processes WebSocket messages and performs actions based on the "action" field.
func processRoomWebSocketMessage(db *gorm.DB, gameManager *types.TicTacMemoGameManager, conn *websocket.Conn, msg []byte, roomID uuid.UUID) {
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

		joinedRoomData := gin.H{
			"event": "joined-room",
			"data": gin.H{
				"room_id":   joinData.RoomID,
				"player_id": joinData.PlayerID,
			},
		}
		conn.WriteJSON(joinedRoomData)

		gameManager.JoinRoom(db, joinData)

	case types.ActionMakeMove:
		var makeMoveData types.MakeMoveData
		if err := json.Unmarshal(message.Data, &makeMoveData); err != nil {
			log.Println("Error unmarshaling make-move data:", err)
			return
		}

		room := gameManager.Rooms[roomID]
		room.MakeMove(db, makeMoveData, makeMoveData.PlayerID)

		gameMoveEvent := gin.H{
			"event": "make-move",
			"data":  room.GameState,
		}

		room.BroadcastGameState(gameMoveEvent)

	case types.ActionUpdateScore:
		var updateScoreData types.UpdateScoreData
		if err := json.Unmarshal(message.Data, &updateScoreData); err != nil {
			log.Println("Error unmarshaling make-move data:", err)
			return
		}

		room := gameManager.Rooms[roomID]
		room.UpdateScore(db, updateScoreData)

	default:
		log.Println("Unknown action:", message.Action)
		conn.WriteMessage(websocket.TextMessage, []byte("Unknown action"))
	}
}
