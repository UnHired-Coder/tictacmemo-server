package handler

import (
	"game-server/common"
	"game-server/database"
	"game-server/tictacmemo"
	"net/http"

	"github.com/gin-gonic/gin"
)

var router *gin.Engine

func init() {
	// Load the environment variables
	/*err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}*/

	// Initialize the Gin router
	router = gin.Default()
	router.Use(gin.Recovery())
	router.Use(CORSMiddleware())

	// Get the database connection
	db := database.GetDatabase()

	// Attach the routes from common package
	common.AttachRoutes(router, db)

	// Attach the routes from tictacmemo package
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
	router.ServeHTTP(w, r)
}
