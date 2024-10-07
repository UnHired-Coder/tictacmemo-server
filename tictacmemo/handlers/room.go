package handlers

import (
	"game-server/common/models"
	"game-server/tictacmemo/core"
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

		// Match players with a timeout (e.g., 30 seconds)
		mms.MatchPlayers(3 * time.Second)

		// Send a response back to the client
		ctx.JSON(http.StatusOK, gin.H{"message": "Matchmaking started!", "waitlist_id": waitlistId})
	}
	return gin.HandlerFunc(fn)
}

func JoinRoom(db *gorm.DB) gin.HandlerFunc {
	fn := func(ctx *gin.Context) {

	}
	return gin.HandlerFunc(fn)
}
