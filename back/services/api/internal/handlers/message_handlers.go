package handlers

import (
	"encoding/json"
	cors "github.com/adhityaramadhanus/fasthttpcors"
	"github.com/fasthttp/router"
	http "github.com/valyala/fasthttp"
	"messenger/services/api/internal/db_wizard"
	helpers "messenger/services/api/pkg/helpers/http"
	"messenger/services/api/pkg/helpers/models"
	"time"

	. "messenger/services/api/pkg/helpers/logger"
)

func (s *Server) MesRouter(r *router.Router, c *cors.CorsHandler) {
	r.POST("/mesList", c.CorsMiddleware(s.mesList))
}

type DialogInfo struct {
	DialogId int `json:"dialog_id"`
}

// обработка запроса на получение списка сообщений
func (s *Server) mesList(ctx *http.RequestCtx) {
	Logger.Info("Get mes list")

	// авторизован ли пользователь?
	token := string(ctx.Request.Header.Peek("Authorization"))
	userId := IsAuth(token)
	if userId == -1 {
		helpers.Respond(ctx, "no auth", http.StatusUnauthorized)
		return
	}

	dialog := DialogInfo{}
	if err := json.Unmarshal(ctx.PostBody(), &dialog); err != nil {
		Logger.Error("Failed unmarshal user data. Reason: ", err.Error())
		Logger.Debug(ctx.PostBody())

		helpers.Respond(ctx, "Unmarshal error", http.StatusBadRequest)
		return
	}

	mess, err := db_wizard.GetMessagesList(dialog.DialogId, userId)
	if err != nil {
		Logger.Error("Failed to do sql req. Reason: ", err.Error())
		helpers.Respond(ctx, "sql error", http.StatusBadRequest)
		return
	}

	helpers.Respond(ctx, mess, http.StatusOK)
	return
}

type MessageInfo struct {
	FriendPhoto    string `json:"friend_photo"`
	FriendFullName string `json:"friend_full_name"`
}

type MesData struct {
	DialogId int           `json:"dialog_id"`
	FriendId int           `json:"friend_id"`
	Time     time.Time     `json:"send_time"`
	Text     string        `json:"text"`
	Files    []models.File `json:"files,omitempty"`
}
