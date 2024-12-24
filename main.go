package main

import (
	"fmt"
	"game-server/common"
	"game-server/database"
	"game-server/tictacmemo"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

var router *gin.Engine

func main() {
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

	for _, route := range router.Routes() {
		log.Printf("Registered Route: %s %s", route.Method, route.Path)
	}

	// Start the server on the specified port
	server := &http.Server{
		Addr:           fmt.Sprintf(":%s", os.Getenv("PORT")),
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	server.ListenAndServe()
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Connection, Upgrade, Sec-WebSocket-Key, Sec-WebSocket-Version, Sec-WebSocket-Extensions")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length, Connection, Upgrade")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		// Allow preflight requests to pass
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
			return
		}

		c.Next()
	}
}
