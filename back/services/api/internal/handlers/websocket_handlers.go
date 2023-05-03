package handlers

import (
	"encoding/binary"
	"encoding/json"
	"github.com/dgrr/websocket"
	"messenger/services/api/internal/db_wizard"
	"messenger/services/api/pkg/helpers/models"
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
	userId := IsAuth(newEvent.Token)

	Logger.Debug(newEvent.EventType)
	switch newEvent.EventType {
	case RegisterUser:
		if newEvent.Token == "" {
			Logger.Error("Bad user token.")
			return
		}

		if _, ok := TokenWebSockets[newEvent.Token]; !ok {

			if _, ok2 := UserToken[newEvent.Token]; !ok2 {
				Logger.Error()
				return
			}

			if _, ok3 := UnauthorizedWebsockets[c.ID()]; !ok3 {
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

	case SendMessage:
		Logger.Info("Send mes")

		user1, user2, err := db_wizard.GetDialogParticipants(newEvent.MsgData.DialogId)
		if err != nil {
			Logger.Error("Failed to do sql req. Reason: ", err.Error())
			return
		}

		mes := models.MessageDB{}

		if user1 == userId {
			mes.Recipient = user2
		} else {
			mes.Recipient = user1
		}

		mes.Text = newEvent.MsgData.Text
		mes.Time = newEvent.MsgData.Time
		mes.Sender = userId

		// write message to db
		mesId, err := db_wizard.PostMessage(mes, newEvent.MsgData.DialogId)
		if err != nil {
			Logger.Error("Failed to do sql req. Reason: ", err.Error())
			return
		}

		// update last message in db
		err = db_wizard.UpdateLastMesInDialog(newEvent.MsgData.DialogId, mesId, mes.Sender)
		if err != nil {
			Logger.Error("Failed to do sql req. Reason: ", err.Error())
			return
		}

		//friendC =

		//c.Write([]byte("seiuthkd"))

	}

	Logger.Infof("Received data from %s: %s", c.RemoteAddr(), data)
}
