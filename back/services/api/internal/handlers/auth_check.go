package handlers

func IsAuth(token string) int {
	id, ok := UserToken[token]
	if ok {
		return id
	}
	return -1
}

//func IsWsOpen(userId int) *websocket.Conn {
//
//	ws, ok := TokenWebSockets[token]
//	if ok {
//		return id
//	}
//	return ws
//}
