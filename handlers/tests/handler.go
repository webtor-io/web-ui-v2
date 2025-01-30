package tests

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/webtor-io/web-ui/services/template"
)

type Handler struct {
	tb template.Builder
}

func RegisterHandler(r *gin.Engine, tm *template.Manager) {
	h := &Handler{
		tb: tm.MustRegisterViews("tests/**/*").WithLayout("main"),
	}

	r.GET("/tests/*template", h.get)
}

type Data struct {
}

func (s *Handler) get(c *gin.Context) {
	s.tb.Build("tests"+c.Param("template")).HTML(http.StatusOK, c, &Data{})
}
