package embed

import (
	"github.com/gin-gonic/gin"
	"github.com/urfave/cli"
	"github.com/webtor-io/web-ui-v2/services/template"
	j "github.com/webtor-io/web-ui-v2/services/web/job"
)

type Handler struct {
	tb   template.Builder
	jobs *j.Handler
}

func RegisterHandler(c *cli.Context, r *gin.Engine, tm *template.Manager, jobs *j.Handler) {
	h := &Handler{
		tb:   tm.MustRegisterViews("embed/*"),
		jobs: jobs,
	}
	r.GET("/embed", h.get)
	r.POST("/embed", h.post)
}
