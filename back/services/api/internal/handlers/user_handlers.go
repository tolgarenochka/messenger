package handlers

import (
	"database/sql"
	"encoding/json"
	cors "github.com/adhityaramadhanus/fasthttpcors"
	"github.com/fasthttp/router"
	http "github.com/valyala/fasthttp"
	"log"
	"messenger/services/api/internal/db_wizard"
	helpers "messenger/services/api/pkg/helpers/http"
	"messenger/services/api/pkg/helpers/models"
)

func (s *Server) UserRouter(r *router.Router, c *cors.CorsHandler) {
	r.POST("/auth", c.CorsMiddleware(s.auth))
}

type AuthData struct {
	Mail string `json:"mail"`
	Pas  string `json:"pas"`
}

//map with session user_id:session_token
var UserToken = map[int64]string{}

func (s *Server) auth(ctx *http.RequestCtx) {
	log.Println("Auth")
	authData := AuthData{}

	if err := json.Unmarshal(ctx.PostBody(), &authData); err != nil {
		log.Print("Failed unmarshal user data. Reason: ", err.Error())

		helpers.RespondError(ctx, models.MakeErrorResponse(err.Error(), http.StatusBadRequest))
		return
	}

	user, err := db_wizard.Auth(authData.Mail, authData.Pas)
	if err != nil {
		log.Print("Failed to do sql req. Reason: ", err.Error())
		if err == sql.ErrNoRows {
			helpers.Respond(ctx, "no auth", http.StatusUnauthorized)
			return
		}
		helpers.RespondError(ctx, models.MakeErrorResponse(err.Error(), http.StatusBadRequest))
		return
	}
	//
	//token := "asdasdasdasd"
	//
	//if token, ok := UserToken[user.Id]; ok {
	//	TokenWebSockets[token].Close()
	//}
	//
	//UserToken[user.Id] = token
	//
	//TokenWebSockets[token] = NewWebsocket(ctx, ctx.Request)

	helpers.Respond(ctx, user, http.StatusOK)
}
