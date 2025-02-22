package resource

import (
	"github.com/gin-gonic/gin"
	j "github.com/webtor-io/web-ui/handlers/job"
	"github.com/webtor-io/web-ui/services/api"
	"github.com/webtor-io/web-ui/services/template"
	"github.com/webtor-io/web-ui/services/web"
	"strings"
)

type Handler struct {
	api  *api.Api
	jobs *j.Handler
	tb   template.Builder[*web.Context]
}

func RegisterHandler(r *gin.Engine, tm *template.Manager[*web.Context], api *api.Api, jobs *j.Handler) {
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
