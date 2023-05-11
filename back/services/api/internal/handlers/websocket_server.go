package handlers

import "github.com/dgrr/websocket"

const (
	RegisterUser = "RegisterUser"
	ReadDialog   = "ReadDialog"
	SendMessage  = "SendMessage"
)

type WSEvent struct {
	Token       string      `json:"token"`
	MsgMetaData MessageInfo `json:"message_info,omitempty"`
	MsgData     MesData     `json:"message_data,omitempty"`
	EventType   string      `json:"event"`
	DialogId    int         `json:"dialog_id,omitempty"`
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
