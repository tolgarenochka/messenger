package handlers

import (
	cors "github.com/adhityaramadhanus/fasthttpcors"
	"github.com/fasthttp/router"
	http "github.com/valyala/fasthttp"

	"messenger/services/api/internal/db_wizard"
	helpers "messenger/services/api/pkg/helpers/http"

	. "messenger/services/api/pkg/helpers/logger"
)

func (s *Server) DialogRouter(r *router.Router, c *cors.CorsHandler) {
	r.GET("/dialogList", c.CorsMiddleware(s.dialogList))
}

// обработка запроса на получение списка диалогов текущего пользователя
func (s *Server) dialogList(ctx *http.RequestCtx) {
	Logger.Info("Get dialog list")

	// авторизован ли пользователь?
	token := string(ctx.Request.Header.Peek("Authorization"))
	userId := IsAuth(token)
	if userId == -1 {
		helpers.Respond(ctx, "no auth", http.StatusUnauthorized)
		return
	}

	// получение списка диалогов
	dialogs, err := db_wizard.GetDialogsList(userId)
	if err != nil {
		Logger.Error("Failed to do sql req. Reason: ", err.Error())
		helpers.Respond(ctx, "sql error", http.StatusBadRequest)
		return
	}

	helpers.Respond(ctx, dialogs, http.StatusOK)
	return
}
