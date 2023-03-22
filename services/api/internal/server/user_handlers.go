package server

import (
	"bytes"
	"encoding/json"
	cors "github.com/adhityaramadhanus/fasthttpcors"
	"github.com/fasthttp/router"
	http "github.com/valyala/fasthttp"
	"log"
	helpers "messenger/services/api/pkg/helpers/http"
	"messenger/services/api/pkg/helpers/models"
)

func (s *Server) UserRouter(r *router.Router, c *cors.CorsHandler) {
	r.GET("/user/", c.CorsMiddleware(s.registration))
}

func (s *Server) registration(ctx *http.RequestCtx) {
	log.Print("Add new user")
	user := models.User{}

	//TODO: do uniq email column

	if err := json.NewDecoder(bytes.NewReader(ctx.PostBody())).Decode(&user); err != nil {
		log.Print("Failed unmarshal user data. Reason: ", err.Error())

		helpers.RespondError(ctx, models.MakeErrorResponse(err.Error(), http.StatusBadRequest))
		return
	}

	helpers.Respond(ctx, []byte("done"), http.StatusOK)

}
