package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,

	// CheckOrigin allows connections from any origin (can be customized for security).
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func sendWebSocketError(conn *websocket.Conn, errorMsg string) error {
	errorResponse := map[string]string{"error": errorMsg}
	errorJSON, err := json.Marshal(errorResponse)
	if err != nil {
		return err
	}
	return conn.WriteMessage(websocket.TextMessage, errorJSON)
}
