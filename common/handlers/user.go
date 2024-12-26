package handlers

import (
	"fmt"
	"game-server/common/types"
	"log"
	"net/http"
	"time"

	"math/rand"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Login(db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		// Bind incoming request parameters
		var userInput struct {
			UserID   string `json:"userId" binding:"required"`
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
		var user types.User
		if err := db.Where("user_id = ?", userInput.UserID).First(&user).Error; err != nil {
			// If the user does not exist, create a new user
			user = types.User{
				UserID:   userInput.UserID,
				Username: userInput.Name,
				Email:    userInput.Email,
				AuthType: userInput.AuthType,
				Avatar:   getRandomAvatar(),
			}
			if err := db.Create(&user).Error; err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
				return
			}
		}

		// Return success response
		ctx.JSON(http.StatusOK, gin.H{
			"message": "User login successful",
			"user":    user,
		})
	}
}

func Profile(db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		// Bind incoming request parameters
		var userInput struct {
			UserID string `json:"userId" binding:"required"`
		}

		// Bind the incoming JSON payload to the userInput struct
		if err := ctx.ShouldBindJSON(&userInput); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		// Check if the user already exists in the database
		var user types.User
		if err := db.Where("user_id = ?", userInput.UserID).First(&user).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find user"})
		}

		var gameHistory []types.GameHistory
		if err := db.Where("user_id = ?", userInput.UserID).Order("created_at DESC").Find(&gameHistory).Error; err != nil {
			gameHistory = []types.GameHistory{}
		}

		var rank int
		query := `
		SELECT rank FROM (
			SELECT id, RANK() OVER (ORDER BY rating DESC) AS rank
			FROM users
		) subquery
		WHERE id = ?`

		// Get the rank of the user
		if err := db.Raw(query, user.ID).Scan(&rank).Error; err != nil {
			log.Printf("Could not update rank!")
		}

		user.Rank = rank

		// Return success response
		ctx.JSON(http.StatusOK, gin.H{
			"message":      "User ",
			"user":         user,
			"game_history": gameHistory,
		})
	}
}

func getRandomAvatar() string {
	rand.Seed(time.Now().UnixNano())                                                                                                           // Seed the random number generator with int64
	randomNumber := rand.Intn(11) + 1                                                                                                          // Generate a random number between 1 and 100
	return fmt.Sprintf("https://firebasestorage.googleapis.com/v0/b/tictacmemo-123e7.appspot.com/o/avatars%%2F%d.png?alt=media", randomNumber) // Return the avatar filename
}
