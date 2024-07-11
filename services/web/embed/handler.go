package embed

import (
	"github.com/gin-gonic/gin"
	"github.com/urfave/cli"
	"github.com/webtor-io/web-ui-v2/services/embed"
	"github.com/webtor-io/web-ui-v2/services/template"
	j "github.com/webtor-io/web-ui-v2/services/web/job"
)

type Handler struct {
	tb   template.Builder
	jobs *j.Handler
	ds   *embed.DomainSettings
}

func RegisterHandler(c *cli.Context, r *gin.Engine, tm *template.Manager, jobs *j.Handler, ds *embed.DomainSettings) {
	h := &Handler{
		tb:   tm.MustRegisterViews("embed/*"),
		jobs: jobs,
		ds:   ds,
	}
	r.GET("/embed", h.get)
	r.POST("/embed", h.post)
}
