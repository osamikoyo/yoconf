package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/osamikoyo/yoconf/core"
)

type Handler struct {
	core *core.Core
}

func NewHandler(core *core.Core) *Handler {
	return &Handler{
		core: core,
	}
}

func (h *Handler) RegisterRouters(e *echo.Echo) {
	e.Use(middleware.Logger())

	e.GET("/get/:project", h.GetChunkHandler)
}

func (h *Handler) GetChunkHandler(c echo.Context) error {
	project := c.Param("project")

	chunk, err := h.core.GetConfig(project)
	if err != nil{
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, chunk)
}
