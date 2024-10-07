package handlers

import (
	"game-server/common/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// PlayerReadyToPlay adds a player to the DB when they are ready to play
func PlayerReadyToPlay(db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var input struct {
			UserID string `json:"user_id" binding:"required"`
		}

		// Bind the incoming JSON data
		if err := ctx.ShouldBindJSON(&input); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		// Check if the user exists in the User table
		var user models.User
		if err := db.Where("id = ?", input.UserID).First(&user).Error; err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}

		// Check if the player is already in the Player table
		var player models.Player
		if err := db.Where("user_id = ?", input.UserID).First(&player).Error; err == nil {
			ctx.JSON(http.StatusConflict, gin.H{"error": "Player is already in the active list"})
			return
		}

		// Add the player to the Player table
		player = models.Player{
			UserID:    input.UserID,
			CreatedAt: time.Now(),
		}

		if err := db.Create(&player).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add player to the active list"})
			return
		}

		// Respond with success
		ctx.JSON(http.StatusOK, gin.H{"message": "Player is ready to play", "player": player})
	}
}

func Login(db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		// Bind incoming request parameters
		var userInput struct {
			ID       int    `json:"id" binding:"required"`
			Name     string `json:"name" binding:"required"`
			Email    string `json:"email" binding:"required,email"`
			AuthType string `json:"authType" binding:"required"`
		}

		// Bind the incoming JSON payload to the userInput struct
		if err := ctx.ShouldBindJSON(&userInput); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		// Check if the user already exists in the database
		var user models.User
		if err := db.Where("id = ?", userInput.ID).First(&user).Error; err != nil {
			// If the user does not exist, create a new user
			user = models.User{
				ID:       userInput.ID,
				Username: userInput.Name,
				Email:    userInput.Email,
				AuthType: userInput.AuthType,
			}
			if err := db.Create(&user).Error; err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
				return
			}
		} else {
			// If user exists, update their information
			user.Username = userInput.Name
			user.Email = userInput.Email
			db.Save(&user)
		}

		// Return success response
		ctx.JSON(http.StatusOK, gin.H{
			"message": "User login successful",
			"user":    user,
		})
	}
}
