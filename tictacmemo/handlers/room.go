package handlers

import (
	"fmt"
	"game-server/common/models"
	"game-server/tictacmemo/core"
	"game-server/tictacmemo/types"

	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func FindMatch(db *gorm.DB, mms *core.MatchmakingSystem) gin.HandlerFunc {
	fn := func(ctx *gin.Context) {
		// Example: Retrieve user ID from the context (maybe set by middleware)
		userId := ctx.Query("user_id")

		// Example: Retrieve a user object from the database using the user ID
		var user models.User
		if err := db.Where("id = ?", userId).First(&user).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "User not found"})
			return
		}

		waitlistId := uuid.New()
		player := models.Player{
			User:       user,
			WaitlistId: waitlistId.String(),
		}

		// Add player to the matchmaking system
		mms.AddPlayer(player)

		go startMatchMacking(ctx, mms)

		// Send a response back to the client
		ctx.JSON(http.StatusOK, gin.H{"message": "Matchmaking started!", "waitlist_id": waitlistId})
	}
	return gin.HandlerFunc(fn)
}

func startMatchMacking(ctx *gin.Context, mms *core.MatchmakingSystem) {
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

	go sendRoomId(ctx, player1, roomId.String(), room)
	go sendRoomId(ctx, player2, roomId.String(), room)
}

func sendRoomId(ctx *gin.Context, player *models.Player, roomId string, room types.Room) {
	wsURL := fmt.Sprintf("/%d/%s", player.ID, player.WaitlistId)
	log.Println("Starting socket connection on " + wsURL)

}

func JoinRoom(db *gorm.DB) gin.HandlerFunc {
	fn := func(ctx *gin.Context) {

	}
	return gin.HandlerFunc(fn)
}
