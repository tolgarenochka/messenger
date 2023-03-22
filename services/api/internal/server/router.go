package server

import (
	cors "github.com/adhityaramadhanus/fasthttpcors"
	"github.com/fasthttp/router"
)

func (s *Server) initRouter() *router.Router {
	c := cors.DefaultHandler()
	r := router.New()

	s.UserRouter(r, c)

	return r
}
