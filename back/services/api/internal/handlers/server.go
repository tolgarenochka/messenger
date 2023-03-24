package handlers

import (
	"context"
	"github.com/fasthttp/router"
	http "github.com/valyala/fasthttp"
	"golang.org/x/sync/errgroup"
)

type Server struct {
	eg  *errgroup.Group
	ctx context.Context

	router *router.Router
	server *HTTPServer
}

func (s *Server) Init() {
	s.router = s.initRouter()
	s.server.serverHTTP.Handler = s.router.Handler
	s.server.serverHTTPS.Handler = s.router.Handler
}

type HTTPServer struct {
	serverHTTP  *http.Server
	serverHTTPS *http.Server
}

func NewServer() *Server {
	return &Server{
		server: &HTTPServer{
			serverHTTP: &http.Server{
				Name: "http handlers",
			},
			serverHTTPS: &http.Server{
				Name: "https handlers",
			},
		},
	}
}

func (s *Server) Run(ctx context.Context, cert, key string) error {
	s.eg, s.ctx = errgroup.WithContext(ctx)

	s.eg.Go(func() error {
		return s.server.serverHTTP.ListenAndServe("localhost:8080")
	})
	s.eg.Go(func() error {
		if err := s.server.serverHTTPS.AppendCert(cert, key); err != nil {
			return err
		}
		return s.server.serverHTTPS.ListenAndServeTLS("localhost:8080", "", "")
	})

	return s.eg.Wait()
}
