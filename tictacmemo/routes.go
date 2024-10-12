package tictacmemo

import (
	"game-server/common/types"
	"game-server/tictacmemo/handlers"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func AttachRoutes(router *gin.Engine, db *gorm.DB) {

	gameManager := types.NewGameManager()
	mms := InitMatchMaking()

	this := router.Group("/tictacmemo")

	{
		this.POST("/join-room/:playerID/:roomID", handlers.JoinRoom(db, gameManager))
		this.POST("/update-score", handlers.UpdateScore(db))
		this.POST("/find-match", handlers.FindMatch(db, mms, gameManager))
		this.GET("/find-match/:playerID/:waitlistID", handlers.Matching(db))
	}
}
