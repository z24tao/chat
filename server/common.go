package server

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/z24tao/chat/model"
	"log"
	"net/http"
	"strconv"
)

func init() {
	http.HandleFunc("/user", userHandler)
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

var upgrader = websocket.Upgrader{
	ReadBufferSize: 1024,
	WriteBufferSize: 1024,
}

func Upgrade(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	upgrader.CheckOrigin = func(r *http.Request) bool {return true}
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
	}
	return ws, nil
}

func send(message *model.Message) {
	usersMutex.Lock()
	defer usersMutex.Unlock()

	client, seen := users[message.To]
	if !seen {
		fmt.Println("error: to user not found")
		return
	}

	messageJson, err := json.Marshal(message)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = client.conn.WriteMessage(websocket.TextMessage, messageJson)
	if err != nil {
		fmt.Println(err)
	}
}

func listen(conn *websocket.Conn) {
	for {
		_, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		var message *model.Message
		err = json.Unmarshal(p, &message)
		if err != nil {
			log.Println(err)
			return
		}

		go send(message)
	}
}

func getUrlParam(r *http.Request, name string) string {
	return r.URL.Query()[name][0]
}

func getUrlIntParam(r *http.Request, name string) int {
	result, _ := strconv.Atoi(getUrlParam(r, name))
	return result
}

func serveStatus(w http.ResponseWriter, status int) {
	w.WriteHeader(status)
}
