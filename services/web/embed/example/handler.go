package example

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/urfave/cli"
	"github.com/webtor-io/web-ui-v2/services/template"
)

type Handler struct {
	tb template.Builder
}

func RegisterHandler(c *cli.Context, r *gin.Engine, tm *template.Manager) {
	h := &Handler{
		tb: tm.MustRegisterViews("embed/example/*").WithLayout("embed/example"),
	}

	r.GET("/embed/example/:name", h.get)
}

type ExampleData struct {
}

func (s *Handler) get(c *gin.Context) {
	s.tb.Build("embed/example/"+c.Param("name")).HTML(http.StatusOK, c, &ExampleData{})
}
