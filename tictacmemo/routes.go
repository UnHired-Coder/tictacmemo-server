package tictacmemo

import (
	"game-server/tictacmemo/handlers"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func AttachRoutes(router *gin.Engine, db *gorm.DB) {

	mms := InitMatchMaking()

	this := router.Group("/tictacmemo")

	{
		this.POST("/join-room", handlers.JoinRoom(db))
		this.POST("/update-score", handlers.UpdateScore(db))
		this.POST("/find-match", handlers.FindMatch(db, mms))
		this.GET("/find-match/:playerID/:waitlistID", handlers.Matching(db))
	}
}
