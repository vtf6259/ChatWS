package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/websocket"
)

var wsUpgrader = websocket.Upgrader{}

func handleWSClient(httpClient http.ResponseWriter, httpRequest *http.Request) {
	fmt.Println("Got client")
	wsClient, err := wsUpgrader.Upgrade(httpClient, httpRequest, nil)
	if err != nil {
		_ = fmt.Errorf("Websocket upgrade failed %s", err.Error())
	}
	fmt.Println(wsClient)
}

func runServer(port int) {
	http.HandleFunc("/ws", handleWSClient)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(port), nil))
}
