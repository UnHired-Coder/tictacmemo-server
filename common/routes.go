package common

import (
	"game-server/common/handlers"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func AttachRoutes(router *gin.Engine, db *gorm.DB) {
	this := router.Group("/common")
	{
		this.POST("/login", handlers.Login(db))
	}
}
