package handlers

import (
	"fmt"
	cors "github.com/adhityaramadhanus/fasthttpcors"
	"github.com/fasthttp/router"
	http "github.com/valyala/fasthttp"
	"messenger/services/api/internal/db_wizard"
	helpers "messenger/services/api/pkg/helpers/http"
)

func (s *Server) FileRouter(r *router.Router, c *cors.CorsHandler) {
	r.POST("/upload", c.CorsMiddleware(s.handleUploadFiles))
}

// обработка запроса на загрузку файлов
func (s *Server) handleUploadFiles(ctx *http.RequestCtx) {
	// авторизован ли пользователь?
	token := string(ctx.Request.Header.Peek("Authorization"))
	userId := IsAuth(token)
	if userId == -1 {
		helpers.Respond(ctx, "no auth", http.StatusUnauthorized)
		return
	}

	form, err := ctx.MultipartForm()
	if err != nil {
		return
	}
	for _, v := range form.File {
		for _, header := range v {
			// загрузка файлов
			if err := http.SaveMultipartFile(header, fmt.Sprintf("/Users/sonyatolgarenko/mess/public/files/%s", header.Filename)); err != nil {
				fmt.Errorf("error save file: %s", err)
				helpers.Respond(ctx, "can't save file", http.StatusBadRequest)
				return
			}
			// добавление записи о них в БД
			err = db_wizard.SaveFile(header.Filename, userId)
			if err != nil {
				helpers.Respond(ctx, "sql error", http.StatusBadRequest)
				return
			}
		}
	}
	helpers.Respond(ctx, "ok", http.StatusOK)
	return
}
