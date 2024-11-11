package handlers

import (
	"game-server/common/types"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Leaderboard(db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var leaderboard []types.User
		if err := db.Where("").Order("rating DESC").Find(&leaderboard).Error; err != nil {
			leaderboard = []types.User{}
		}

		for index := range leaderboard {
			leaderboard[index].Rank = index + 1
		}

		// Return success response
		ctx.JSON(http.StatusOK, gin.H{
			"message":     "Leaderboard fetched successfully",
			"leaderboard": leaderboard,
		})
	}
}
