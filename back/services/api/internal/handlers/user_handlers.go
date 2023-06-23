package handlers

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	cors "github.com/adhityaramadhanus/fasthttpcors"
	"github.com/fasthttp/router"
	http "github.com/valyala/fasthttp"
	"messenger/services/api/internal/db_wizard"
	helpers "messenger/services/api/pkg/helpers/http"
	. "messenger/services/api/pkg/helpers/logger"
	"messenger/services/api/pkg/helpers/utils"
	"strconv"
)

// UserToken map with session session_token:user_id
var UserToken = map[string]int{}

type AuthData struct {
	Mail string `json:"mail"`
	Pas  string `json:"pas"`
}

type UserInfo struct {
	Token    string `json:"token"`
	FullName string `json:"full_name"`
	Photo    string `json:"photo"`
}

// функция генерации токена
func GenerateSecureToken(length int) string {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return hex.EncodeToString(b)
}

func (s *Server) UserRouter(r *router.Router, c *cors.CorsHandler) {
	r.POST("/auth", c.CorsMiddleware(s.auth))
	r.POST("/updatePhoto", c.CorsMiddleware(s.updatePhoto))
	r.GET("/usersList", c.CorsMiddleware(s.usersList))
	r.POST("/logout", c.CorsMiddleware(s.logout))
}

// обработка запроса на авторизацию
func (s *Server) auth(ctx *http.RequestCtx) {
	Logger.Info("Auth")
	authData := AuthData{}

	if err := json.Unmarshal(ctx.PostBody(), &authData); err != nil {
		Logger.Error("Failed unmarshal user data. Reason: ", err.Error())

		helpers.Respond(ctx, "Unmarshal error", http.StatusBadRequest)
		return
	}

	user, err := db_wizard.Auth(authData.Mail, utils.GetSHA256(authData.Pas, s.configuration.Salt))
	if err != nil {
		Logger.Error("Failed to do sql req. Reason: ", err.Error())
		if err == sql.ErrNoRows {
			helpers.Respond(ctx, "not valid mail or pas", http.StatusUnauthorized)
			return
		}
		helpers.Respond(ctx, "sql error", http.StatusBadRequest)
		return
	}

	// генерация токена пользователю
	token := GenerateSecureToken(20)

	UserToken[token] = user.Id

	usr := UserInfo{
		Token:    token,
		FullName: user.SecondName + " " + user.FirstName + " " + user.ThirdName,
		Photo:    user.Photo,
	}

	helpers.Respond(ctx, usr, http.StatusOK)
}

// обработка запроса на выход из системы
func (s *Server) logout(ctx *http.RequestCtx) {
	Logger.Info("Log Out")

	token := string(ctx.Request.Header.Peek("Authorization"))
	userId := IsAuth(token)
	if userId == -1 {
		helpers.Respond(ctx, "no auth", http.StatusUnauthorized)
		return
	}

	// удаление токена пользователя
	delete(UserToken, token)
	helpers.Respond(ctx, "un auth", http.StatusOK)
}

type NewPhoto struct {
	Photo string `json:"photo"`
}

// обработка запроса на обновление фото профиля
func (s *Server) updatePhoto(ctx *http.RequestCtx) {
	Logger.Info("Update photo")

	token := string(ctx.Request.Header.Peek("Authorization"))
	userId := IsAuth(token)
	if userId == -1 {
		helpers.Respond(ctx, "no auth", http.StatusUnauthorized)
		return
	}

	fh, err := ctx.FormFile("photo")
	if err != nil {
		fmt.Errorf("error get file: %s", err)
		helpers.Respond(ctx, "can't get file", http.StatusBadRequest)
		return
	}

	if err := http.SaveMultipartFile(fh, fmt.Sprintf("/Users/sonyatolgarenko/mess/public/user_photo/user%s.jpeg",
		strconv.Itoa(userId))); err != nil {
		fmt.Errorf("error save file: %s", err)
		helpers.Respond(ctx, "can't save file", http.StatusBadRequest)
		return
	}

	helpers.Respond(ctx, "ok", http.StatusOK)
	return
}

// обработка запроса на получение списка пользователей
func (s *Server) usersList(ctx *http.RequestCtx) {
	Logger.Info("Get users list")

	token := string(ctx.Request.Header.Peek("Authorization"))
	userId := IsAuth(token)
	if userId == -1 {
		helpers.Respond(ctx, "no auth", http.StatusUnauthorized)
		return
	}

	users, err := db_wizard.GetUsersList(userId)
	if err != nil {
		Logger.Error("Failed to do sql req. Reason: ", err.Error())
		helpers.Respond(ctx, "sql error", http.StatusBadRequest)
		return
	}

	helpers.Respond(ctx, users, http.StatusOK)
	return
}
