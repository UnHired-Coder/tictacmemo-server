package handlers

import (
	commonTypes "game-server/common/types"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func JoinRoom(db *gorm.DB, gameManager *commonTypes.GameManager) gin.HandlerFunc {
	fn := func(ctx *gin.Context) {
		playerID := ctx.Param("playerID")
		roomID, err := strconv.Atoi(ctx.Param("roomID"))

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Invalid Player id: " + playerID})
			return
		}

		log.Println("Player " + playerID + " requested to join the room")

		var user commonTypes.User
		result := db.First(&user, playerID) // Find user by ID
		if result.Error != nil {
			log.Fatal("Error fetching user:", result.Error)
		}

		// Check if the room exists in the gameManager
		err = gameManager.JoinRoom(&user, roomID, 2)
		if err != nil {
			log.Printf("Room with ID %d not found", roomID)
			ctx.JSON(http.StatusNotFound, gin.H{"error": err})
			return
		}
	}
	return gin.HandlerFunc(fn)
}
