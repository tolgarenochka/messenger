package http

import (
	"github.com/valyala/fasthttp"
	"log"
	"messenger/services/api/pkg/helpers/models"
	"net/http"
)

func Respond(ctx *fasthttp.RequestCtx, message []byte, code int) {
	ctx.Response.Header.SetContentLength(len(message))
	ctx.SetContentType("application/json")

	_, err := ctx.Write(message)
	if err != nil {
		log.Print("Failed send respond. Reason: ", err.Error())

		ctx.SetStatusCode(http.StatusInternalServerError)
	}

	ctx.SetStatusCode(code)
}

func RespondError(ctx *fasthttp.RequestCtx, errorObject *models.Error) {
	bytes, err := errorObject.Marshal()
	if err != nil {
		log.Print("Failed marshal error response. Reason: ", err.Error())

		ctx.SetStatusCode(http.StatusInternalServerError)
	}

	ctx.Response.Header.SetContentLength(len(bytes))
	ctx.SetContentType("application/json")

	_, err = ctx.Write(bytes)
	if err != nil {
		log.Print("Failed sent error response. Reason: ", err.Error())

		ctx.SetStatusCode(http.StatusInternalServerError)
	}

	ctx.SetStatusCode(int(errorObject.Code()))
}
