package index

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/webtor-io/web-ui/services/template"
)

type Data struct {
	Instruction string
}

type Handler struct {
	tb template.Builder
}

func RegisterHandler(r *gin.Engine, tm *template.Manager) {
	h := &Handler{
		tb: tm.MustRegisterViews("*").WithLayout("main"),
	}
	r.GET("/", h.index)
	r.GET("/torrent-to-ddl", h.index)
	r.GET("/torrent-to-zip", h.index)
	r.GET("/magnet-to-ddl", h.index)
	r.GET("/magnet-to-torrent", h.index)
}

func (s *Handler) index(c *gin.Context) {
	s.tb.Build("index").HTML(http.StatusOK, c, &Data{
		Instruction: strings.TrimPrefix(c.Request.URL.Path, "/"),
	})
}
