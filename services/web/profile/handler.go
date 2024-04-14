package profile

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/urfave/cli"
	w "github.com/webtor-io/web-ui-v2/services/web"
)

type ProfileData struct{}

type Handler struct {
	tm *w.TemplateManager
}

func RegisterHandler(c *cli.Context, r *gin.Engine, tm *w.TemplateManager) {
	h := &Handler{
		tm: tm,
	}
	r.GET("/profile", h.get)

	tm.RegisterViews("profile/*")
}

func (s *Handler) get(c *gin.Context) {
	s.tm.MakeTemplate("profile/get").HTML(http.StatusOK, c, &ProfileData{})
}
