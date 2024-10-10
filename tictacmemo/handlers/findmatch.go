package handlers

import (
	"fmt"
	"game-server/common/types"
	"game-server/tictacmemo/core"

	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func FindMatch(db *gorm.DB, mms *core.MatchmakingSystem) gin.HandlerFunc {
	fn := func(ctx *gin.Context) {

		//Retrive the user
		userId := ctx.Query("user_id")
		var user types.User
		if err := db.Where("id = ?", userId).First(&user).Error; err != nil {
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

		go startMatchMacking(mms)

		// Send a response back to the client
		ctx.JSON(http.StatusOK, gin.H{"message": "Matchmaking started!", "waitlist_id": waitlistId})
	}
	return gin.HandlerFunc(fn)
}

func startMatchMacking(mms *core.MatchmakingSystem) {
	// Match players with a timeout (e.g., 30 seconds)
	player1, player2, err := mms.MatchPlayers(300 * time.Second)
	if err != nil {
		log.Fatal("Error: " + err.Error())
	}

	if player1 == nil || player2 == nil {
		log.Fatal("Something went wrong!")
	}

	roomId := uuid.New()
	room := types.CreateRoom(*player1, *player2)

	go sendRoomId(player1, roomId.String(), room)
	go sendRoomId(player2, roomId.String(), room)
}

func sendRoomId(player *types.Player, roomId string, room types.Room) {
	wsURL := fmt.Sprintf("/%d/%s", player.ID, player.WaitlistId)
	log.Println("Joining room on: " + wsURL)

	roomData := map[string]any{
		"playerID": player.ID,
		"roomID":   roomId,
		"room":     room,
	}

	sendRoomDataForMatch(roomData, 1, player.WaitlistId)
}

func sendRoomDataForMatch(roomData map[string]any, attempt int, waitlistId string) {
	if attempt < 20 {
		connection, ok := PLAYERS_WAITLIST[waitlistId]
		if ok {
			// Send the structured message to the client
			err := connection.WriteJSON(roomData)
			if err != nil {
				log.Println("Failed to send message:", err)
			}
		} else {
			log.Println("Connection not established yet")
			time.Sleep(2 * time.Second)
			sendRoomDataForMatch(roomData, 1, waitlistId)
		}
	}
}
