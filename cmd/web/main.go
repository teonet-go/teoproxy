package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

func main() {

	// Start web server
	http.Handle("/", http.FileServer(http.Dir("./")))

	// Websocket handler
	http.HandleFunc("/ws", handleWebSocket)

	log.Fatal(http.ListenAndServe(":8093", nil))

}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {

	log.Println("ws client start connection", r.RemoteAddr)

	var upgrader = websocket.Upgrader{} // use default options

	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()

	log.Println("ws client connected", c.RemoteAddr())

	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", message)
		err = c.WriteMessage(mt, message)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}
