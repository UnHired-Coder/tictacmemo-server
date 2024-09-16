package tictacmemo

import (
	"game-server/tictacmemo/handlers"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func AttachTicTacMemoRoutes(router *gin.Engine, db *gorm.DB) {
	this := router.Group("/tictacmemo")

	{
		this.POST("/create-room", handlers.CreateRoom(db))
		this.POST("/join-room", handlers.JoinRoom(db))
		this.POST("/update-score", handlers.UpdateScore(db))

	}
}
