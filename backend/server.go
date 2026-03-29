package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/websocket"
)

func corsIsTheDumbestThingEver(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Handle preflight requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

var wsUpgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // i hate cors
	},
}

func handleWSClient(httpClient http.ResponseWriter, httpRequest *http.Request) {
	fmt.Printf("Got client: %s\n", sens(httpRequest.RemoteAddr))
	wsClient, err := wsUpgrader.Upgrade(httpClient, httpRequest, nil)
	if err != nil {
		fmt.Printf("Websocket upgrade failed: %s\n", err.Error())
		return
	}
	for {
		_, messageByte, err := wsClient.ReadMessage()
		if err != nil {
			fmt.Println("Error: ", err)
			continue
		}
		var message string
		for _, content := range messageByte {
			message = message + string(content)
		}
		fmt.Println(httpRequest.RemoteAddr, "Said:", message)
		wsClient.WriteMessage(1, []byte("Hello"))
	}
}

func runServer(port int) {
	mux := http.NewServeMux()
	mux.HandleFunc("/ws", handleWSClient)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(port), corsIsTheDumbestThingEver(mux)))
}
