package handlers

// функция проверки авторизован ли пользователь
func IsAuth(token string) int {
	id, ok := UserToken[token]
	if ok {
		return id
	}
	return -1
}
