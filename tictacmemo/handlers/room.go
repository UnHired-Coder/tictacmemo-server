package handlers

import (
	"game-server/tictacmemo/core"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func FindMatch(db *gorm.DB, mms *core.MatchmakingSystem) gin.HandlerFunc {
	fn := func(ctx *gin.Context) {
		// userId := ctx.GetInt("userId")
		// currentTime := time.Now()

		// Get User from DB
		// user :=

		// mms.AddPlayer(user)
		//mms.MatchPlayers(30 * time.Second)

	}
	return gin.HandlerFunc(fn)
}

func JoinRoom(db *gorm.DB) gin.HandlerFunc {
	fn := func(ctx *gin.Context) {

	}
	return gin.HandlerFunc(fn)
}
