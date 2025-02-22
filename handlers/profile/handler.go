package profile

import (
	"github.com/webtor-io/web-ui/services/web"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/webtor-io/web-ui/services/template"
)

type Data struct{}

type Handler struct {
	tb template.Builder[*web.Context]
}

func RegisterHandler(r *gin.Engine, tm *template.Manager[*web.Context]) {
	h := &Handler{
		tb: tm.MustRegisterViews("profile/*").WithLayout("main"),
	}
	r.GET("/profile", h.get)
}

func (s *Handler) get(c *gin.Context) {
	s.tb.Build("profile/get").HTML(http.StatusOK, web.NewContext(c).WithData(&Data{}))
}
