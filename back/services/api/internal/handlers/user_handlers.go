package handlers

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	cors "github.com/adhityaramadhanus/fasthttpcors"
	"github.com/fasthttp/router"
	http "github.com/valyala/fasthttp"
	"log"
	"messenger/services/api/internal/db_wizard"
	helpers "messenger/services/api/pkg/helpers/http"
)

func (s *Server) UserRouter(r *router.Router, c *cors.CorsHandler) {
	r.POST("/auth", c.CorsMiddleware(s.auth))
	r.POST("/updatePhoto", c.CorsMiddleware(s.updatePhoto))
	r.GET("/usersList", c.CorsMiddleware(s.usersList))
	r.POST("/logout", c.CorsMiddleware(s.logout))
}

type AuthData struct {
	Mail string `json:"mail"`
	Pas  string `json:"pas"`
}

type UserInfo struct {
	Token    string `json:"token"`
	FullName string `json:"full_name"`
	Photo    string `json:"photo"`
}

func GenerateSecureToken(length int) string {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return hex.EncodeToString(b)
}

// UserToken map with session user_id:session_token
var UserToken = map[string]int{}

func (s *Server) auth(ctx *http.RequestCtx) {
	log.Println("Auth")
	authData := AuthData{}

	if err := json.Unmarshal(ctx.PostBody(), &authData); err != nil {
		log.Print("Failed unmarshal user data. Reason: ", err.Error())

		helpers.Respond(ctx, "Unmarshal error", http.StatusBadRequest)
		return
	}

	user, err := db_wizard.Auth(authData.Mail, authData.Pas)
	if err != nil {
		log.Print("Failed to do sql req. Reason: ", err.Error())
		if err == sql.ErrNoRows {
			helpers.Respond(ctx, "not valid mail or pas", http.StatusUnauthorized)
			return
		}
		helpers.Respond(ctx, "sql error", http.StatusBadRequest)
		return
	}

	token := GenerateSecureToken(20)

	// TODO: add ws
	//if token, ok := UserToken[user.Id]; ok {
	//	TokenWebSockets[token].Close()
	//}
	//
	UserToken[token] = user.Id
	//TokenWebSockets[token] = NewWebsocket(ctx, ctx.Request)

	usr := UserInfo{
		Token:    token,
		FullName: user.SecondName + " " + user.FirstName + " " + user.ThirdName,
		Photo:    user.Photo,
	}
	helpers.Respond(ctx, usr, http.StatusOK)
}

func (s *Server) logout(ctx *http.RequestCtx) {
	log.Println("Log Out")

	token := string(ctx.Request.Header.Peek("Authorization"))
	userId := IsAuth(token)
	if userId == -1 {
		helpers.Respond(ctx, "no auth", http.StatusUnauthorized)
		return
	}

	delete(UserToken, token)
	helpers.Respond(ctx, "un auth", http.StatusOK)
}

// NewPhoto TODO: front gives us id? or token for getting id from UserToken?
type NewPhoto struct {
	Photo string `json:"photo"`
	Id    int    `json:"id"`
}

func (s *Server) updatePhoto(ctx *http.RequestCtx) {
	log.Println("Update photo")

	token := string(ctx.Request.Header.Peek("Authorization"))
	userId := IsAuth(token)
	if userId == -1 {
		helpers.Respond(ctx, "no auth", http.StatusUnauthorized)
		return
	}

	newPhoto := NewPhoto{}

	if err := json.Unmarshal(ctx.PostBody(), &newPhoto); err != nil {
		log.Print("Failed unmarshal user data. Reason: ", err.Error())

		helpers.Respond(ctx, "Unmarshal error", http.StatusBadRequest)
		return
	}

	count, err := db_wizard.UpdatePhoto(newPhoto.Photo, newPhoto.Id)
	if err != nil {
		log.Print("Failed to do sql req. Reason: ", err.Error())
		helpers.Respond(ctx, "sql error", http.StatusBadRequest)
		return
	}
	if count == 0 {
		helpers.Respond(ctx, "can't update photo", http.StatusBadRequest)
		return
	}

	helpers.Respond(ctx, "photo was successfully updated", http.StatusOK)
	return
}

func (s *Server) usersList(ctx *http.RequestCtx) {
	log.Println("Get users list")

	token := string(ctx.Request.Header.Peek("Authorization"))
	userId := IsAuth(token)
	if userId == -1 {
		helpers.Respond(ctx, "no auth", http.StatusUnauthorized)
		return
	}

	users, err := db_wizard.GetUsersList()
	if err != nil {
		log.Print("Failed to do sql req. Reason: ", err.Error())
		helpers.Respond(ctx, "sql error", http.StatusBadRequest)
		return
	}

	helpers.Respond(ctx, users, http.StatusOK)
	return
}
