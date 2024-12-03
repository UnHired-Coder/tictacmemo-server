package main

// import (
// 	"fmt"
// 	"game-server/common"
// 	"game-server/database"
// 	"game-server/tictacmemo"
// 	"log"
// 	"net/http"
// 	"os"
// 	"time"

// 	"github.com/gin-gonic/gin"
// 	"github.com/joho/godotenv"
// )

// func CORSMiddleware() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
// 		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
// 		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

// 		if c.Request.Method == "OPTIONS" {
// 			c.AbortWithStatus(200)
// 			return
// 		}
// 		c.Next()
// 	}
// }

// func main() {
// 	// Load the environment variables
// 	err := godotenv.Load()
// 	if err != nil {
// 		log.Fatal("Error loading .env file")
// 	}

// 	// Create a new router
// 	router := gin.Default()
// 	// Recover from any panics
// 	router.Use(gin.Recovery())
// 	router.Use(CORSMiddleware())

// 	// Get the database connection
// 	db := database.GetDatabase()

// 	// Attach the routes from common package
// 	common.AttachRoutes(router, db)

// 	// Attach the routes from tictacmemo package
// 	tictacmemo.AttachRoutes(router, db)

// 	// Start the server on the specified port
// 	server := &http.Server{
// 		Addr:           fmt.Sprintf(":%s", os.Getenv("PORT")),
// 		Handler:        router,
// 		ReadTimeout:    10 * time.Second,
// 		WriteTimeout:   10 * time.Second,
// 		MaxHeaderBytes: 1 << 20,
// 	}
// 	server.ListenAndServe()
// }
