package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/ProtonMail/gopenpgp/v3/crypto"
	"github.com/gorilla/websocket"
)

type userType struct {
	remoteAddr          string
	randomText          string
	encryptedRandomText string
	username            string
	authed              bool
}

func getAuth(wsClient *websocket.Conn) bool {
	pgp := crypto.PGP()
	randomText := rand.Text()
	var user userType
	for {
		_, messageByte, _ := wsClient.ReadMessage()
		var message string
		for _, content := range messageByte {
			message = message + string(content)
		}
		messages := strings.SplitN(message, " ", 2)
		switch messages[0] {
		case "USER":
			prodLogln("USER " + messages[1])
			user.randomText = randomText
			user.username = messages[1] // TODO: Do bounds checking before 1.0
			err := wsClient.WriteMessage(1, []byte("CHALLENGE "+randomText+" "+pgppubkeycrypto.GetHexKeyID()))
			if err != nil {
				wsClient.WriteMessage(1, []byte("500")) // dont really care if this gets to the client since there is already a error
				return false
			}	
			prodLogln(err)
			prodLogln(messages)
			err = wsClient.WriteMessage(1, []byte("MOTD "+motd))
			if err != nil {
				wsClient.WriteMessage(1, []byte("500")) // this is only during auth since it is kindof messy i can put this in a function and put it everywhere for 1.0 
				return false
			}
			prodLogln(sens(user))
		case "AUTH":
			prodLogln(sens(messages[1]))
			privateKey, _ := crypto.NewPrivateKeyFromArmored(pgpprivkey, []byte("")) // TODO: error checking before 1.0
			decHandle, _ := pgp.Decryption().DecryptionKey(privateKey).New()
			randomTextDecrypted, _ := decHandle.Decrypt([]byte(messages[1]), crypto.Armor)
			if randomTextDecrypted.String() == user.randomText {
				wsClient.WriteMessage(1, []byte("AUTHSUC"))
				user.authed = true
				return true
			} else {
				wsClient.WriteMessage(1, []byte("FAILAUTH"))
			}
			prodLogln(sens(randomTextDecrypted))
		case "REQPUBKEY":
			pubkey, _ := pgppubkeycrypto.Armor()
			pubkeybase64 := base64.StdEncoding.EncodeToString([]byte(pubkey))
			pubkeybase64 = strings.ReplaceAll(pubkeybase64, "\n", "") // just incase
			wsClient.WriteMessage(1, []byte("PUBKEY "+pubkeybase64))
		case "RESEND":
			wsClient.WriteMessage(1, []byte("CHALLENGE "+randomText+" "+pgppubkeycrypto.GetHexKeyID()))
		default:
			wsClient.WriteMessage(1, []byte("FAILAUTH"))
			continue
		}
	}
}

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
	authDone := getAuth(wsClient)
	if !authDone {
		prodLogln(httpRequest.RemoteAddr + " Failed auth")
		return
	}
	for {
		_, messageByte, err := wsClient.ReadMessage()
		if err != nil {
			if err == websocket.ErrCloseSent {
				return
			}
			fmt.Println("Error: ", err)
			wsClient.WriteMessage(1, []byte("ERROR IN WS LOOP: "+err.Error()))
			return
		}
		var message string
		for _, content := range messageByte {
			message = message + string(content)
		}
		if message == "closing" {
			fmt.Printf("Say bye to %s\n", sens(httpRequest.RemoteAddr))
			wsClient.Close()
			return
		}
		fmt.Println(sens(httpRequest.RemoteAddr), "Said:", message)
		wsClient.WriteMessage(1, []byte("Hello"))
	}
}

func handleIndex(httpClient http.ResponseWriter, httpRequest *http.Request) {
	io.WriteString(httpClient, "<h1>TODO: put frontend in here maybe</h1>")
}

func runServer(port int) {
	//pgp := crypto.PGP()
	mux := http.NewServeMux()
	mux.HandleFunc("/ws", handleWSClient)
	mux.HandleFunc("/", handleIndex)
	if windowsMentioned {
		fmt.Println("Windows style paths detected (\\) we do not do anything different other than path stuff so do not expect stuff to 100% work on windows")
	} else {
		_ = progname(os.Args[0])
	}
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(port), corsIsTheDumbestThingEver(mux)))
}
