package resource

import (
	"github.com/gin-gonic/gin"
	"github.com/urfave/cli"
	"github.com/webtor-io/web-ui-v2/services/api"
	"github.com/webtor-io/web-ui-v2/services/template"
	j "github.com/webtor-io/web-ui-v2/services/web/job"
)

type Handler struct {
	api  *api.Api
	jobs *j.Handler
	tb   template.Builder
}

func RegisterHandler(c *cli.Context, r *gin.Engine, tm *template.Manager, api *api.Api, jobs *j.Handler) {
	helper := NewHelper(c)
	h := &Handler{
		api:  api,
		jobs: jobs,
		tb:   tm.MustRegisterViews("resource/*").WithHelper(helper).WithLayout("main"),
	}
	r.POST("/", h.post)
	r.GET("/:resource_id", h.get)
}
