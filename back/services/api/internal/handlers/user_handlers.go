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

var (
	corsAllowHeaders     = "authorization"
	corsAllowMethods     = "HEAD,GET,POST,PUT,DELETE,OPTIONS"
	corsAllowOrigin      = "*"
	corsAllowCredentials = "true"
)

func CORS(next http.RequestHandler) http.RequestHandler {
	return func(ctx *http.RequestCtx) {

		ctx.Response.Header.Set("Access-Control-Allow-Credentials", corsAllowCredentials)
		ctx.Response.Header.Set("Access-Control-Allow-Headers", corsAllowHeaders)
		ctx.Response.Header.Set("Access-Control-Allow-Methods", corsAllowMethods)
		ctx.Response.Header.Set("Access-Control-Allow-Origin", corsAllowOrigin)

		next(ctx)
	}
}

func (s *Server) UserRouter(r *router.Router, c *cors.CorsHandler) {
	r.POST("/auth", CORS(s.auth))
	r.POST("/updatePhoto", c.CorsMiddleware(s.updatePhoto))
	r.GET("/usersList", c.CorsMiddleware(s.usersList))
}

type AuthData struct {
	Mail string `json:"mail"`
	Pas  string `json:"pas"`
}

func GenerateSecureToken(length int) string {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return hex.EncodeToString(b)
}

// UserToken map with session user_id:session_token
var UserToken = map[int64]string{}

func (s *Server) auth(ctx *http.RequestCtx) {
	ctx.Response.Header.Set("Access-Control-Allow-Credentials", "true")
	ctx.Response.Header.SetBytesV("Access-Control-Allow-Origin", ctx.Request.Header.Peek("Origin"))

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
	UserToken[user.Id] = token
	//TokenWebSockets[token] = NewWebsocket(ctx, ctx.Request)

	helpers.Respond(ctx, token, http.StatusOK)
}

// NewPhoto TODO: front gives us id? or token for getting id from UserToken?
type NewPhoto struct {
	Photo string `json:"photo"`
	Id    int    `json:"id"`
}

func (s *Server) updatePhoto(ctx *http.RequestCtx) {
	log.Println("Update photo")
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

	users, err := db_wizard.GetUsersList()
	if err != nil {
		log.Print("Failed to do sql req. Reason: ", err.Error())
		helpers.Respond(ctx, "sql error", http.StatusBadRequest)
		return
	}

	helpers.Respond(ctx, users, http.StatusOK)
	return
}