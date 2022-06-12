package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
)

type GamePlay struct {
	Message     string `json:"message"`
	PlayerName  string `json:"playerName"`
	Timestamp   int64  `json:"timestamp"`
	GameId      int    `json:"gameId"`
	Answer      uint   `json:"answer"`
	GuessResult int    `json:"guessResult"`
}

// Define Upgrader, read write buffer size
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Continue listening on connection
func reader(conn *websocket.Conn) {
	for {
		var gameplay GamePlay
		// var data map[string]interface{}
		// log.Println(data["text"].(string))

		// fmt.Printf("Message received from client: %v\n", gameplay)

		messageType, p, readErr := conn.ReadMessage()
		if readErr != nil {
			log.Println("Reader Error: ", readErr)
			return
		}

		msgdata := string(p)
		fmt.Println(msgdata)
		json.Unmarshal([]byte(msgdata), &gameplay)

		// print out the message
		fmt.Printf("Message received from client: %v\n", gameplay)
		if strings.ToLower(gameplay.Message) == "registration" && gameplay.GameId == 0 {
			if GamerCounter == 0 {
				startGame()
				GamerCounter += 1
			}
			gameplay.GameId = GamerCounter

			Players = append(Players, gameplay)
		}
		if strings.ToLower(gameplay.Message) == "guess" {
			var verifyResult = verifyAnswer(&gameplay)
			fmt.Printf("Verified Result: %v\n", verifyResult)
		}

		fmt.Printf("Message going to send to client: %v\n", gameplay)

		reqBodyBytes := new(bytes.Buffer)
		json.NewEncoder(reqBodyBytes).Encode(gameplay)

		if writeErr := conn.WriteMessage(messageType, reqBodyBytes.Bytes()); writeErr != nil {
			log.Println("Write Error: ", writeErr)
			return
		} else {
			fmt.Println(string(reqBodyBytes.Bytes()))
		}
	}
}

func wsEndPoint(w http.ResponseWriter, r *http.Request) {
	// accept allow connection
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	ws, websocketErr := upgrader.Upgrade(w, r, nil)
	if websocketErr != nil {
		log.Println("WebSocket Error: ", websocketErr)
	}

	log.Println("Client Successfully Connect...")

	reader(ws)
}
