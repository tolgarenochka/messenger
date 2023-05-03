package handlers

import "github.com/dgrr/websocket"

const (
	RegisterUser = "RegisterUser"
	MessageRead  = "ReadDialog"
	SendMessage  = "SendMessage"
)

type WSEvent struct {
	Token       string      `json:"token"`
	MsgMetaData MessageInfo `json:"message_info,omitempty"` //TODO: it is needed in msg id. Check front.
	MsgData     MesData     `json:"message_data,omitempty"`
	EventType   string      `json:"event"`
}

type WS struct {
	WsID uint64
	Conn *websocket.Conn
}

// map with session token:websocket
var TokenWebSockets = map[string]*websocket.Conn{}

var UnauthorizedWebsockets = map[uint64]*websocket.Conn{}

func InitWebSocketServer() websocket.Server {
	ws := websocket.Server{}

	ws.HandleOpen(RegisterConn)
	ws.HandleClose(RemoveConn)
	ws.HandlePong(OnPong)
	ws.HandleData(OnMessage)

	return ws
}
