package handlers

import (
	cors "github.com/adhityaramadhanus/fasthttpcors"
	"github.com/fasthttp/router"
	http "github.com/valyala/fasthttp"
	"log"
	"messenger/services/api/internal/db_wizard"
	helpers "messenger/services/api/pkg/helpers/http"
)

func (s *Server) DialogRouter(r *router.Router, c *cors.CorsHandler) {
	r.GET("/dialogList", c.CorsMiddleware(s.dialogList))
}

func (s *Server) dialogList(ctx *http.RequestCtx) {
	log.Println("Get dialog list")

	token := string(ctx.Request.Header.Peek("Authorization"))
	userId := IsAuth(token)
	if userId == -1 {
		helpers.Respond(ctx, "no auth", http.StatusUnauthorized)
		return
	}

	dialogs, err := db_wizard.GetDialogsList(userId)
	if err != nil {
		log.Print("Failed to do sql req. Reason: ", err.Error())
		helpers.Respond(ctx, "sql error", http.StatusBadRequest)
		return
	}

	helpers.Respond(ctx, dialogs, http.StatusOK)
	return
}
