package handlers

import (
	"fmt"
	commonTypes "game-server/common/types"
	"game-server/tictacmemo/types"

	"game-server/tictacmemo/core"

	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func FindMatch(db *gorm.DB, mms *core.MatchmakingSystem, gameManager *types.TicTacMemoGameManager) gin.HandlerFunc {
	fn := func(ctx *gin.Context) {

		userId := ctx.Query("user_id")
		var user commonTypes.User
		if err := db.Where("user_id = ?", userId).First(&user).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "User not found"})
			return
		}

		waitlistId := uuid.New()
		player := types.Player{
			User:       user,
			WaitlistId: waitlistId.String(),
		}

		// Add player to the matchmaking system
		mms.AddPlayer(player)

		go startMatchMacking(mms, gameManager)

		ctx.JSON(http.StatusOK, gin.H{
			"event": "matching-started",
			"data": gin.H{
				"waitlist_id": waitlistId,
			},
		})

	}
	return gin.HandlerFunc(fn)
}

func startMatchMacking(mms *core.MatchmakingSystem, gameManager *types.TicTacMemoGameManager) {
	// Match players with a timeout (e.g., 30 seconds)
	player1, player2, err := mms.MatchPlayers(300 * time.Second)
	if err != nil {
		log.Fatal("Error: " + err.Error())
	}

	if player1 == nil || player2 == nil {
		log.Fatal("Something went wrong!")
	}

	MAX_PLAYERS_TIC_TAC_MEMEO := 2
	roomId, room, err := gameManager.CreateRoom(MAX_PLAYERS_TIC_TAC_MEMEO, player1.UserID, player2.UserID)

	if err != nil {
		log.Println("Failed to Create Room:", err)
	}

	player1.InitialGameData = &types.InitialGameData{
		AssignedLable: "X",
	}

	player2.InitialGameData = &types.InitialGameData{
		AssignedLable: "O",
	}

	// Now we have enough players, emit room Id over the websocket
	go sendRoomId(player1, roomId, room)
	go sendRoomId(player2, roomId, room)
}

func sendRoomId(player *types.Player, roomId uuid.UUID, room *types.TicTacMemoRoom) {
	wsURL := fmt.Sprintf("/%s/%s", player.UserID, player.WaitlistId)
	log.Println("Player added to waitlist: " + wsURL)

	roomData := gin.H{
		"event": "player-matched",
		"data": gin.H{
			"room_id":         roomId,
			"room":            room,
			"InitialGameData": player.InitialGameData,
		},
	}

	sendRoomDataForMatch(roomData, 1, player.WaitlistId)
}

func sendRoomDataForMatch(roomData map[string]any, attempt int, waitlistId string) {
	if attempt < 20 {
		connection, ok := PLAYERS_WAITLIST[waitlistId]
		if ok {

			log.Println("Match found, join the room")

			err := connection.WriteJSON(roomData)
			if err != nil {
				log.Println("Failed to send message:", err)
			}
		} else {

			log.Println("Waiting for other player to join...")
			time.Sleep(2 * time.Second)
			sendRoomDataForMatch(roomData, 1, waitlistId)
		}
	}
}
