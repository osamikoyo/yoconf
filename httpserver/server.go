package httpserver

import (
	"context"
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/osamikoyo/yoconf/config"
	"github.com/osamikoyo/yoconf/handler"
	"github.com/osamikoyo/yoconf/logger"
)

type HTTPServer struct {
	server *echo.Echo
	logger *logger.Logger
	cfg    *config.Config
	handler *handler.Handler
}

func NewHTTPServer(
	server *echo.Echo,
	logger *logger.Logger,
	cfg *config.Config,
	handler *handler.Handler,
) *HTTPServer {
	return &HTTPServer{
		server: server,
		logger: logger,
		cfg: cfg,
		handler: handler,
	}
}

func (s *HTTPServer) Close(ctx context.Context) error {
	s.logger.Info("closing http server...")
	
	return s.server.Shutdown(ctx)
}

func (s *HTTPServer) Run(ctx context.Context) error {
	s.logger.Info("starting http server...")

	s.handler.RegisterRouters(s.server)

	return s.server.Start(fmt.Sprintf("%s:%d", s.cfg.Addr, s.cfg.HTTPPort))
}