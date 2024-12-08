package handler

import (
	"game-server/common"
	"game-server/database"
	"game-server/tictacmemo"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

var router *gin.Engine

func init() {
	// Initialize the Gin router
	router = gin.Default()
	router.Use(gin.Recovery())
	router.Use(CORSMiddleware())

	// Get the database connection
	log.Println("SETTING UP DB")
	db := database.GetDatabase()

	// Attach the routes from common package
	log.Println("ATTACHING ROUTES:common")
	common.AttachRoutes(router, db)

	// Attach the routes from tictacmemo package
	log.Println("ATTACHING ROUTES:tictacmemo")
	tictacmemo.AttachRoutes(router, db)
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
			return
		}
		c.Next()
	}
}

// Handler is the entry point for Vercel
func Handler(w http.ResponseWriter, r *http.Request) {
	log.Println("SERVING: HANDLER")
	router.ServeHTTP(w, r)
}
