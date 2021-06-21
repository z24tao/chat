package server

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
)

type messageClient struct {
	userId int
	conn   *websocket.Conn
}

var users = map[int]*messageClient{}
var usersMutex = new(sync.Mutex)

func userHandler(w http.ResponseWriter, r *http.Request) {
	usersMutex.Lock()
	defer usersMutex.Unlock()

	userId := getUrlIntParam(r, "id")
	fmt.Println("attempting to connect with", userId)

	if _, seen := users[userId]; seen {
		serveStatus(w, http.StatusBadRequest)
		return
	}

	conn, err := Upgrade(w, r)
	if err != nil {
		fmt.Println(err)
		serveStatus(w, http.StatusInternalServerError)
		return
	}

	users[userId] = &messageClient{
		userId: userId,
		conn:   conn,
	}

	userId ++
	go listen(conn)
}
