package index

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/urfave/cli"
	"github.com/webtor-io/web-ui-v2/services/template"
)

type IndexData struct{}

type Handler struct {
	tm *template.Manager
}

func RegisterHandler(c *cli.Context, r *gin.Engine, tm *template.Manager) {
	h := &Handler{
		tm: tm,
	}
	r.GET("/", h.index)

	h.tm.RegisterViews("*")
}

func (s *Handler) index(c *gin.Context) {
	s.tm.MakeTemplate("index").HTML(http.StatusOK, c, &IndexData{})
}
