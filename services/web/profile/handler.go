package profile

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/urfave/cli"
	"github.com/webtor-io/web-ui-v2/services/template"
)

type ProfileData struct{}

type Handler struct {
	tm *template.Manager
}

func RegisterHandler(c *cli.Context, r *gin.Engine, tm *template.Manager) {
	h := &Handler{
		tm: tm,
	}
	r.GET("/profile", h.get)

	tm.RegisterViews("profile/*")
}

func (s *Handler) get(c *gin.Context) {
	s.tm.MakeTemplate("profile/get").HTML(http.StatusOK, c, &ProfileData{})
}
