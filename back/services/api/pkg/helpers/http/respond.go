package http

import (
	"encoding/json"
	"github.com/valyala/fasthttp"
	"log"
	"messenger/services/api/pkg/helpers/models"
	"net/http"
)

func Respond(ctx *fasthttp.RequestCtx, data interface{}, statusCode int) {
	bytes, err := json.Marshal(data)
	if err != nil {
		log.Print("Failed to marshal error data. Reason:", err.Error())
	}

	ctx.Response.Header.Set("Content-Type", "application/json")

	_, err = ctx.Write(bytes)
	if err != nil {
		log.Print("Failed to marshal error data. Reason:", err.Error())
	}

	ctx.SetStatusCode(statusCode)
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
