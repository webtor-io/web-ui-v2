package resource

import (
	"html/template"

	"github.com/gin-gonic/gin"
	"github.com/urfave/cli"
	sv "github.com/webtor-io/web-ui-v2/services"
	j "github.com/webtor-io/web-ui-v2/services/job"
	w "github.com/webtor-io/web-ui-v2/services/web"
)

type Handler struct {
	*w.ClaimsHandler
	api  *sv.Api
	jobs *j.Handler
	tm   *w.TemplateManager
}

func RegisterHandler(c *cli.Context, r *gin.Engine, tm *w.TemplateManager, api *sv.Api, jobs *j.Handler, uc *sv.UserClaims) {
	h := &Handler{
		ClaimsHandler: w.NewClaimsHandler(uc),
		api:           api,
		jobs:          jobs,
		tm:            tm,
	}
	r.POST("/", h.post)
	r.GET("/:resource_id", h.get)

	tm.RegisterViewsWithFuncs("resource/*", template.FuncMap{
		"makeBreadcrumbs":  MakeBreadcrumbs,
		"hasBreadcrumbs":   HasBreadcrumbs,
		"hasPagination":    HasPagination,
		"makePagination":   MakePagination,
		"makeFileDownload": MakeFileDownload,
		"makeDirDownload":  MakeDirDownload,
		"makeImage":        MakeImage,
		"makeAudio":        MakeAudio,
		"makeVideo":        MakeVideo,
	})
}
