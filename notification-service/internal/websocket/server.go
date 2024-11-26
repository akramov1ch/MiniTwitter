package websocket

import (
	"net/http"
	"github.com/gorilla/websocket"
	"log"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Client struct {
	conn *websocket.Conn
}

var clients = make(map[*Client]bool)

func StartWebSocketServer() {
	http.HandleFunc("/ws", handleConnections)
	log.Println("WebSocket server is running on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	client := &Client{conn: conn}
	clients[client] = true

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			delete(clients, client)
			break
		}
	}
}

func BroadcastMessage(message []byte) {
	for client := range clients {
		err := client.conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			log.Println(err)
			client.conn.Close()
			delete(clients, client)
		}
	}
}