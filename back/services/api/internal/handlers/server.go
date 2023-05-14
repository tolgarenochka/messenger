package handlers

import (
	"context"
	"fmt"
	"github.com/fasthttp/router"
	http "github.com/valyala/fasthttp"
	"golang.org/x/sync/errgroup"
	"messenger/services/api/pkg/helpers/models"
)

type Server struct {
	eg  *errgroup.Group
	ctx context.Context

	router        *router.Router
	server        *HTTPServer
	configuration models.Configuration
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

func NewServer(configuration models.Configuration) *Server {
	return &Server{
		server: &HTTPServer{
			serverHTTP: &http.Server{
				Name: "http handlers",
			},
			serverHTTPS: &http.Server{
				Name: "https handlers",
			},
		},
		configuration: configuration,
	}
}

func (s *Server) Run(ctx context.Context) error {
	s.eg, s.ctx = errgroup.WithContext(ctx)
	s.eg.Go(func() error {
		return s.server.serverHTTP.ListenAndServe(fmt.Sprintf("%s:%s", s.configuration.IP, s.configuration.HTTPPort))
	})
	s.eg.Go(func() error {
		if err := s.server.serverHTTPS.AppendCert(s.configuration.PathToServerCrt, s.configuration.PathToServerKey); err != nil {
			return err
		}
		return s.server.serverHTTPS.ListenAndServeTLS(fmt.Sprintf("%s:%s", s.configuration.IP, s.configuration.HTTPSPort), "", "")
	})

	ws := InitWebSocketServer()
	s.eg.Go(func() error {
		return http.ListenAndServe(fmt.Sprintf("%s:%s", s.configuration.IP, s.configuration.WSPort), ws.Upgrade)
	})
	s.eg.Go(func() error {
		if err := s.server.serverHTTPS.AppendCert(s.configuration.PathToServerCrt, s.configuration.PathToServerKey); err != nil {
			return err
		}
		return http.ListenAndServeTLS(fmt.Sprintf("%s:%s", s.configuration.IP, s.configuration.WSSPort), "", "", ws.Upgrade)
	})

	return s.eg.Wait()
}
