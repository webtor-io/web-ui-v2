package index

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/urfave/cli"
	w "github.com/webtor-io/web-ui-v2/services/web"
)

type IndexData struct{}

type Handler struct {
	tm *w.TemplateManager
}

func RegisterHandler(c *cli.Context, r *gin.Engine, tm *w.TemplateManager) {
	h := &Handler{
		tm: tm,
	}
	r.GET("/", h.index)

	h.tm.RegisterViews("*")
}

func (s *Handler) index(c *gin.Context) {
	s.tm.MakeTemplate("index").HTML(http.StatusOK, c, &IndexData{})
}
