/*
WebSocket example
https://www.youtube.com/watch?v=dniVs0xKYKk

WebSocket is upgraded HTTP connection
Lives until connection is killed by client or server
Duplex communication using a single connection
1 TCP connection for all communication > reduce network overhead

Download websocket
1. go get github.com/gorilla/websocket
	https://www.youtube.com/watch?v=n3BQLHtsrkM
2. Socket.IO would be another choice
	go get github.com/googellee/go-socket.io
	https://www.youtube.com/watch?v=ycgCMOWPgiw
Install live-server for test
go to extensions > search live server > install
Run
go run main.go

Test:
https://chrome.google.com/webstore/detail/websocket-test-client/fgponpodhbmadfljofbimhhlengambbn
*/

package main

import (
	"fmt"
	"log"
	"net/http"
	// "github.com/Tinprosword/go-websockets/game"
	"math/rand"
	"time"
)

var GamerCounter = 0
var Players = []GamePlay {}
var SecretAnswer uint
var MinAnswerRange uint = 1
var MaxAnswerRange uint = 500

func setupRoutes() {
	// http.HandleFunc("/", homePage)
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/", fs)
	http.HandleFunc("/ws", wsEndPoint)
}

func main() {
	fmt.Println("Go WebSockets")
	setupRoutes()
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func startGame() {
	// Generate the random secret answer
	var answerRange = int(MaxAnswerRange - MinAnswerRange)
	rand.Seed(time.Now().UTC().UnixNano())
	SecretAnswer = uint(rand.Intn(answerRange)) + MinAnswerRange
}

func verifyAnswer(gamePlay *GamePlay) *GamePlay {
	// func verifyAnswer(gamePlay GamePlay) GamePlay {
	var message string
	var answer = gamePlay.Answer
	var player = gamePlay.PlayerName
	if answer == SecretAnswer {
		message = fmt.Sprintf("Congratulation to %v, the answer is %v\n", player, answer)
	} else if answer > SecretAnswer {
		message = "too large"
	} else {
		message = "too small"
	}
	gamePlay.Message = message
	return gamePlay
}
