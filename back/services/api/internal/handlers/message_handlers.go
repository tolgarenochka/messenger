package handlers

import (
	"encoding/json"
	cors "github.com/adhityaramadhanus/fasthttpcors"
	"github.com/fasthttp/router"
	http "github.com/valyala/fasthttp"
	"log"
	"messenger/services/api/internal/db_wizard"
	helpers "messenger/services/api/pkg/helpers/http"
)

func (s *Server) MesRouter(r *router.Router, c *cors.CorsHandler) {
	r.GET("/mesList", c.CorsMiddleware(s.mesList))
}

type DialogInfo struct {
	DialogId int `json:"dialog_id"`
}

func (s *Server) mesList(ctx *http.RequestCtx) {
	log.Println("Get mes list")
	dialog := DialogInfo{}
	if err := json.Unmarshal(ctx.PostBody(), &dialog); err != nil {
		log.Print("Failed unmarshal user data. Reason: ", err.Error())

		helpers.Respond(ctx, "Unmarshal error", http.StatusBadRequest)
		return
	}

	mess, err := db_wizard.GetMessagesList(dialog.DialogId)
	if err != nil {
		log.Print("Failed to do sql req. Reason: ", err.Error())
		helpers.Respond(ctx, "sql error", http.StatusBadRequest)
		return
	}

	helpers.Respond(ctx, mess, http.StatusOK)
	return
}

//type Message struct {
//	Id        int64         `json:"id" db:"id"`
//	Time      time.Time     `json:"time" db:"time"`
//	Text      string        `json:"text" db:"text"`
//	Sender    string        `json:"sender" db:"sender"`
//	Recipient string        `json:"recipient" db:"recipient"`
//	File      []models.File `json:"file"`
//}
