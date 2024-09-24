package legal

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/webtor-io/web-ui-v2/services/template"
)

type Handler struct {
	tb template.Builder
}

func RegisterHandler(r *gin.Engine, tm *template.Manager) {
	h := &Handler{
		tb: tm.MustRegisterViews("legal/**/*").WithLayout("main"),
	}

	r.GET("/legal/*template", h.get)
}

type Data struct {
}

func (s *Handler) get(c *gin.Context) {
	s.tb.Build("legal"+c.Param("template")).HTML(http.StatusOK, c, &Data{})
}
