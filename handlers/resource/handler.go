package resource

import (
	"github.com/gin-gonic/gin"
	j "github.com/webtor-io/web-ui-v2/handlers/job"
	"github.com/webtor-io/web-ui-v2/services/api"
	"github.com/webtor-io/web-ui-v2/services/template"
	"strings"
)

type Handler struct {
	api  *api.Api
	jobs *j.Handler
	tb   template.Builder
}

func RegisterHandler(r *gin.Engine, tm *template.Manager, api *api.Api, jobs *j.Handler) {
	helper := NewHelper()
	h := &Handler{
		api:  api,
		jobs: jobs,
		tb:   tm.MustRegisterViews("resource/*").WithHelper(helper).WithLayout("main"),
	}
	r.POST("/", h.post)
	r.GET("/:resource_id", func(c *gin.Context) {
		if strings.HasPrefix(c.Param("resource_id"), "magnet") {
			h.post(c)
			return
		}
		h.get(c)
	})
}
