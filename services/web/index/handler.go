package index

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/urfave/cli"
	"github.com/webtor-io/web-ui-v2/services/template"
)

type IndexData struct {
}

type Handler struct {
	tb template.Builder
}

func RegisterHandler(c *cli.Context, r *gin.Engine, tm *template.Manager) {
	h := &Handler{
		tb: tm.MustRegisterViews("*").WithLayout("main"),
	}
	r.GET("/", h.index)
}

func (s *Handler) index(c *gin.Context) {
	s.tb.Build("index").HTML(http.StatusOK, c, &IndexData{})
}
