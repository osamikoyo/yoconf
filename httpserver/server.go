package httpserver

import (
	"github.com/labstack/echo/v4"
	"github.com/osamikoyo/yoconf/config"
	"github.com/osamikoyo/yoconf/logger"
)

type HTTPServer struct {
	server *echo.Echo
	logger *logger.Logger
	cfg    *config.Config
}

func NewHTTPServer(
	server *echo.Echo,
	logger *logger.Logger,
	cfg *config.Config,
) *HTTPServer {
}
