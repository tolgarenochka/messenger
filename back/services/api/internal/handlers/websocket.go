package handlers

import (
	"github.com/fasthttp/websocket"
	"net/http"

	. "github.com/NGRsoftlab/ngr-logging"
)

// map with session token:websocket
var TokenWebSockets = map[string]interface{}{}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func NewWebsocket(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		Logger.Fatal(err)
	}

	defer ws.Close()

	te, p, err := ws.ReadMessage()
	if err != nil {
		Logger.Fatal(err)
	}

	Logger.Print(te, p)

	//t := Test{ChatId: 1, Msg: "TEEESTING", Attach: []string{"first file data", "secind file"}}
	//
	//err = ws.WriteJSON(t)
	//if err != nil {
	//	Logger.Fatal(err)
	//}
	//
	//te, p, err = ws.ReadMessage()
	//if err != nil {
	//	Logger.Fatal(err)
	//}
	//
	//Logger.Print(te, p)
}
