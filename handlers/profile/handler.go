package profile

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/webtor-io/web-ui/services/template"
)

type Data struct{}

type Handler struct {
	tb template.Builder
}

func RegisterHandler(r *gin.Engine, tm *template.Manager) {
	h := &Handler{
		tb: tm.MustRegisterViews("profile/*").WithLayout("main"),
	}
	r.GET("/profile", h.get)
}

func (s *Handler) get(c *gin.Context) {
	s.tb.Build("profile/get").HTML(http.StatusOK, c, &Data{})
}
