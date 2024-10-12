package handlers

import (
	commonTypes "game-server/common/types"
	"game-server/tictacmemo/types"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func JoinRoom(db *gorm.DB, gameManager *types.TicTacMemoGameManager) gin.HandlerFunc {
	fn := func(ctx *gin.Context) {
		playerID := ctx.Param("playerID")
		// Parse and validate the UUID
		roomID, err := uuid.Parse(ctx.Param("roomID"))
		if err != nil {
			// Return a JSON error response if the UUID is invalid
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid room ID: " + roomID.String()})
			return
		}

		log.Println("Player " + playerID + " requested to join the room")

		var user commonTypes.User
		result := db.First(&user, playerID) // Find user by ID
		if result.Error != nil {
			log.Fatal("Error fetching user:", result.Error)
		}

		// Check if the room exists in the gameManager

		joinFunc := func(room *types.TicTacMemoRoom, player *commonTypes.User) error {
			return room.JoinRoom(player)
		}

		err = gameManager.JoinRoom(&user, roomID, joinFunc)
		if err != nil {
			log.Printf("Room with ID %d not found", roomID)
			ctx.JSON(http.StatusNotFound, gin.H{"error": err})
			return
		}
	}
	return gin.HandlerFunc(fn)
}
