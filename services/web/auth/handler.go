package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/urfave/cli"
	"github.com/webtor-io/web-ui-v2/services/template"
)

type LoginData struct{}

type LogoutData struct{}

type VerifyData struct {
	PreAuthSessionId string
}

type Handler struct {
	tb template.Builder
}

func RegisterHandler(c *cli.Context, r *gin.Engine, tm *template.Manager) {
	h := &Handler{
		tb: tm.MustRegisterViews("auth/*").WithLayout("main"),
	}

	r.GET("/login", h.login)
	r.GET("/logout", h.logout)
	r.GET("/auth/verify", h.verify)
}

func (s *Handler) login(c *gin.Context) {
	s.tb.Build("auth/login").HTML(http.StatusOK, c, LoginData{})
}

func (s *Handler) logout(c *gin.Context) {
	s.tb.Build("auth/logout").HTML(http.StatusOK, c, LogoutData{})
}

func (s *Handler) verify(c *gin.Context) {
	s.tb.Build("auth/verify").HTML(http.StatusOK, c, &VerifyData{
		PreAuthSessionId: c.Query("preAuthSessionId"),
	})
}
