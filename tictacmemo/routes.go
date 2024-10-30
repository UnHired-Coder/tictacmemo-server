package tictacmemo

import (
	"game-server/tictacmemo/handlers"
	"game-server/tictacmemo/types"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func AttachRoutes(router *gin.Engine, db *gorm.DB) {

	gameManager := types.NewTicTacMemoGameManager()
	mms := InitMatchMaking()

	this := router.Group("/tictacmemo")

	{
		this.POST("/find-match", handlers.FindMatch(db, mms, gameManager))
		this.GET("/find-match/:playerID/:waitlistID", handlers.Matching(db, gameManager))
		this.GET("/play-game/:roomID", handlers.Room(db, gameManager))
	}
}
