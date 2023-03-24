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

func (s *Server) auth(ctx *http.RequestCtx) {
	log.Println("Auth")
	authData := AuthData{}

	if err := json.Unmarshal(ctx.PostBody(), &authData); err != nil {
		log.Print("Failed unmarshal user data. Reason: ", err.Error())

		helpers.RespondError(ctx, models.MakeErrorResponse(err.Error(), http.StatusBadRequest))
		return
	}

	store, err := db_wizard.NewConnect()
	if err != nil {
		log.Print("Failed connect to db. Reason: ", err.Error())
		return
	}

	user := models.User{}
	user, err = store.Auth(authData.Mail, authData.Pas)
	if err != nil {
		log.Print("Failed to do sql req. Reason: ", err.Error())
		if err == sql.ErrNoRows {
			helpers.Respond(ctx, []byte("no auth"), http.StatusUnauthorized)
			return
		}
		helpers.RespondError(ctx, models.MakeErrorResponse(err.Error(), http.StatusBadRequest))
		return
	}

	data := make([]byte, 0)

	data, err = json.Marshal(&user)
	if err != nil {
		log.Print("Failed me. Reason: ", err.Error())
		helpers.RespondError(ctx, models.MakeErrorResponse(err.Error(), http.StatusBadRequest))
		return
	}

	_, err = ctx.Write(data)
	if err != nil {
		log.Print("Failed to write data to resp. Reason: ", err.Error())
		helpers.RespondError(ctx, models.MakeErrorResponse(err.Error(), http.StatusBadRequest))
		return
	}

}
