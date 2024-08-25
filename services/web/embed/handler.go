package embed

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/webtor-io/web-ui-v2/services/embed"
	"github.com/webtor-io/web-ui-v2/services/template"
	j "github.com/webtor-io/web-ui-v2/services/web/job"
)

type Handler struct {
	tb   template.Builder
	hCl  *http.Client
	jobs *j.Handler
	ds   *embed.DomainSettings
}

func RegisterHandler(hCl *http.Client, r *gin.Engine, tm *template.Manager, jobs *j.Handler, ds *embed.DomainSettings) {
	h := &Handler{
		tb:   tm.MustRegisterViews("embed/*"),
		jobs: jobs,
		ds:   ds,
		hCl:  hCl,
	}
	r.GET("/embed", h.get)
	r.POST("/embed", h.post)
}
