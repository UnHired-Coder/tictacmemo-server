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

	"game-server/common/websocketserver"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func matchMakingWebsocketHandler(ctx *gin.Context) {
	// Extract player ID and waitlist ID from the URL parameters
	playerID := ctx.Param("playerID")
	waitlistID := ctx.Param("waitlistID")

	// Logging connection details
	log.Printf("WebSocket connection for playerID: %s, waitlistID: %s", playerID, waitlistID)

	// Pass context to the WebSocket handler
	websocketserver.HandleWebSocket(ctx)
}

func main() {
	// Load the environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Create a new router
	router := gin.Default()
	// Recover from any panics
	router.Use(gin.Recovery())

	// Get the database connection
	db := database.GetDatabase()

	// Attach the routes from common package
	common.AttachRoutes(router, db)

	// Attach the routes from tictacmemo package
	tictacmemo.AttachRoutes(router, db)

	// General WebSocket route for handling dynamic player connections
	router.GET("/ws/:playerID/:waitlistID", matchMakingWebsocketHandler)

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
