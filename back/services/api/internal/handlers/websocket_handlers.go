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

type TimeStruct struct {
	Event string `json:"event"`
	Id    int    `json:"id"`
}

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
	case ReadDialog:
		Logger.Info("Read dialog")

		err := db_wizard.ReadDialog(newEvent.DialogId)
		if err != nil {
			Logger.Error("Failed to do sql req. Reason: ", err.Error())
			return
		}

	case SendMessage:
		Logger.Info("Send mes")

		var dId int
		var flag = false
		var a WSEvent

		if newEvent.MsgData.DialogId == 0 {
			flag = true
			dId, err = db_wizard.CreateDialog(userId, newEvent.MsgData.FriendId)
			if err != nil {
				Logger.Error("Failed to send mes to ws. Reason: ", err.Error())
				return
			}

			newEvent.MsgData.DialogId = dId

			userInf, err := db_wizard.GetUserInfoById(userId)
			if err != nil {
				Logger.Error("Failed to do sql req. Reason: ", err.Error())
				return
			}
			newEvent.MsgMetaData.FriendPhoto = userInf.Photo
			newEvent.MsgMetaData.FriendFullName = userInf.FullName

			a := TimeStruct{Event: "dialogId", Id: dId}

			jsn, err := json.Marshal(a)
			_, err = c.Write(jsn)
			if err != nil {
				Logger.Error("Failed to send mes to ws. Reason: ", err.Error())
				return
			}
		} else {
			dId = newEvent.MsgData.DialogId
		}

		user1, user2, err := db_wizard.GetDialogParticipants(dId)
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
		mesId, err := db_wizard.PostMessage(mes, dId)
		if err != nil {
			Logger.Error("Failed to do sql req. Reason: ", err.Error())
			return
		}

		// update last message in db
		err = db_wizard.UpdateLastMesInDialog(dId, mesId, mes.Sender)
		if err != nil {
			Logger.Error("Failed to do sql req. Reason: ", err.Error())
			return
		}

		for token, id := range UserToken {
			if id == newEvent.MsgData.FriendId {
				ws, ok := TokenWebSockets[token]
				if ok {
					if flag == false {
						a = WSEvent{Token: "", EventType: "NewMessage", MsgData: newEvent.MsgData}
					} else {
						newEvent.MsgData.FriendId = userId
						a = WSEvent{Token: "", EventType: "NewMessage", MsgData: newEvent.MsgData, MsgMetaData: newEvent.MsgMetaData}
					}

					jsn, err := json.Marshal(a)
					_, err = ws.Write(jsn)
					if err != nil {
						Logger.Error("Failed to send mes to ws. Reason: ", err.Error())
						return
					}
				}
			}
		}

		//c.Write([]byte("seiuthkd"))

	}

	Logger.Infof("Received data from %s: %s", c.RemoteAddr(), data)
}
