package handlers

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreateRoom(db *gorm.DB) gin.HandlerFunc {
	fn := func(ctx *gin.Context) {

	}
	return gin.HandlerFunc(fn)
}

func JoinRoom(db *gorm.DB) gin.HandlerFunc {
	fn := func(ctx *gin.Context) {

	}
	return gin.HandlerFunc(fn)
}
