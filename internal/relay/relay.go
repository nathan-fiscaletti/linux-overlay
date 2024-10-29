package relay

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func NewWsHandler(outChan chan any) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			fmt.Println("Failed to set websocket upgrade:", err)
			return
		}
		defer conn.Close()

		for {
			e := <-outChan

			encodedJson, err := json.Marshal(e)
			if err != nil {
				fmt.Println("Error encoding JSON:", err)
				continue
			}

			// Send a response to the client
			err = conn.WriteMessage(websocket.TextMessage, encodedJson)
			if err != nil {
				fmt.Println("Error writing message:", err)
				break
			}
		}
	}
}
