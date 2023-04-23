package handlers

import (
	"encoding/binary"
	"encoding/json"
	"github.com/dgrr/websocket"
	"time"

	. "messenger/services/api/pkg/helpers/logger"
)

// register a connection when it's open
func RegisterConn(c *websocket.Conn) {
	UnauthorizedWebsockets[c.ID()] = c
	Logger.Debugf("Client %s connected", c.RemoteAddr())
}

// remove the connection when receiving the close
func RemoveConn(c *websocket.Conn, err error) {
	if err != nil {
		Logger.Info("Close websocket connection. Reason: ", err.Error())
	}

	var token string
	for _token, ws := range TokenWebSockets {
		if ws.ID() == c.ID() {
			token = _token
		}
	}

	if token == "" {
		Logger.Error("Bad close websocket connection: connection not found.")
		return
	}

	delete(TokenWebSockets, token)
	Logger.Debugf("Client %s disconnected with token %s", c.RemoteAddr(), token)

	_ = c.Close()
}

// handle the pong message
func OnPong(c *websocket.Conn, data []byte) {
	if len(data) == 8 {
		n := binary.BigEndian.Uint64(data)
		ts := time.Unix(0, int64(n))

		Logger.Debugf("RTT with %s is %s\n", c.RemoteAddr(), time.Now().Sub(ts))
	}
}


func OnMessage(c *websocket.Conn, isBinary bool, data []byte) {
	var newEvent WSEvent

	err := json.Unmarshal(data, &newEvent)
	if err != nil {
		Logger.Error("Failed get new websocket event. Reason: ", err)
		return
	}

	Logger.Debug(newEvent.EventType)
	switch newEvent.EventType {
	case RegisterUser:
		if newEvent.Token == "" {
			Logger.Error("Bad user token.")
			return
		}

		if _, ok := TokenWebSockets[newEvent.Token];!ok {

			//TODO: uncomment
			//if _, ok2 := UserToken[token]; !ok2 {
			//	Logger.Error()
			//	return
			//}

			if  _, ok3 := UnauthorizedWebsockets[c.ID()]; !ok3 {
				Logger.Error("Websocket connection not found!")
				return
			}

			TokenWebSockets[newEvent.Token] = c
			delete(UnauthorizedWebsockets, c.ID())

			Logger.Debug("Register user")
		} else {
			Logger.Debug("User already registered!")
		}
	case MessageRead:
		// get message and user id from newEvent
		// send event "msg read" to author of message
		// update database
	}



	Logger.Infof("Received data from %s: %s", c.RemoteAddr(), data)
}
