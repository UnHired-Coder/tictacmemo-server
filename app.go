package main

import (
	"game-server/tictacmemo"
	"game-server/tictacmemo/database"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.Use(gin.Recovery())

	db := database.GetDatabase()

	tictacmemo.AttachTicTacMemoRoutes(router, db)

	router.Run() // listen and serve on 0.0.0.0:8080
}
