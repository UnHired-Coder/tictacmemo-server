package handlers

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreateUser(db *gorm.DB) gin.HandlerFunc {
	fn := func(ctx *gin.Context) {

	}
	return gin.HandlerFunc(fn)
}

func Login(db *gorm.DB) gin.HandlerFunc {
	fn := func(ctx *gin.Context) {

	}
	return gin.HandlerFunc(fn)
}
