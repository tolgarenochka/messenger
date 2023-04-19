package handlers

import (
	"encoding/json"
	cors "github.com/adhityaramadhanus/fasthttpcors"
	"github.com/fasthttp/router"
	http "github.com/valyala/fasthttp"
	"log"
	"messenger/services/api/internal/db_wizard"
	helpers "messenger/services/api/pkg/helpers/http"
	"messenger/services/api/pkg/helpers/models"
	"time"
)

func (s *Server) MesRouter(r *router.Router, c *cors.CorsHandler) {
	r.POST("/mesList", c.CorsMiddleware(s.mesList))
	r.POST("/sendMes", c.CorsMiddleware(s.sendMes))
}

type DialogInfo struct {
	DialogId int `json:"dialog_id"`
}

func (s *Server) mesList(ctx *http.RequestCtx) {
	log.Println("Get mes list")

	token := string(ctx.Request.Header.Peek("Authorization"))
	userId := IsAuth(token)
	if userId == -1 {
		helpers.Respond(ctx, "no auth", http.StatusUnauthorized)
		return
	}

	dialog := DialogInfo{}
	if err := json.Unmarshal(ctx.PostBody(), &dialog); err != nil {
		log.Print("Failed unmarshal user data. Reason: ", err.Error())
		log.Print(ctx.PostBody())

		helpers.Respond(ctx, "Unmarshal error", http.StatusBadRequest)
		return
	}

	mess, err := db_wizard.GetMessagesList(dialog.DialogId, userId)
	if err != nil {
		log.Print("Failed to do sql req. Reason: ", err.Error())
		helpers.Respond(ctx, "sql error", http.StatusBadRequest)
		return
	}

	helpers.Respond(ctx, mess, http.StatusOK)
	return
}

type MessageInfo struct {
	Text     string `json:"text"`
	DialogId int    `json:"dialog_id"`
}

type MesData struct {
	DialogId int       `json:"dialog_id"`
	Time     time.Time `json:"send_time"`
	Text     string    `json:"text"`
}

func (s *Server) sendMes(ctx *http.RequestCtx) {
	log.Println("Send mes")

	token := string(ctx.Request.Header.Peek("Authorization"))
	userId := IsAuth(token)
	if userId == -1 {
		helpers.Respond(ctx, "no auth", http.StatusUnauthorized)
		return
	}

	mesData := MesData{}
	mes := models.MessageDB{}
	if err := json.Unmarshal(ctx.PostBody(), &mesData); err != nil {
		log.Print("Failed unmarshal user data. Reason: ", err.Error())

		helpers.Respond(ctx, "Unmarshal error", http.StatusBadRequest)
		return
	}

	user1, user2, err := db_wizard.GetDialogParticipants(mesData.DialogId)
	if err != nil {
		log.Print("Failed to do sql req. Reason: ", err.Error())
		helpers.Respond(ctx, "sql error", http.StatusBadRequest)
		return
	}

	if user1 == userId {
		mes.Recipient = user2
	} else {
		mes.Recipient = user1
	}

	mes.Text = mesData.Text
	mes.Time = mesData.Time
	mes.Sender = userId

	mesId, err := db_wizard.PostMessage(mes, mesData.DialogId)
	if err != nil {
		log.Print("Failed to do sql req. Reason: ", err.Error())
		helpers.Respond(ctx, "sql error", http.StatusBadRequest)
		return
	}

	err = db_wizard.UpdateLastMesInDialog(mesData.DialogId, mesId, mes.Sender)
	if err != nil {
		log.Print("Failed to do sql req. Reason: ", err.Error())
		helpers.Respond(ctx, "sql error", http.StatusBadRequest)
		return
	}

	helpers.Respond(ctx, "message was send", http.StatusOK)
	return
}
