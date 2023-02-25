package web

import (
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	sv "github.com/webtor-io/web-ui-v2/services"
)

type ClaimsHandler struct {
}

func NewClaimsHandler() *ClaimsHandler {
	return &ClaimsHandler{}
}

func (s *ClaimsHandler) getRemoteAddress(r *http.Request) string {
	forwarded := r.Header.Get("X-FORWARDED-FOR")
	if forwarded != "" {
		return strings.Split(forwarded, ",")[0]
	}
	ip, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr))
	if err != nil {
		return ""
	}
	return ip
}

func (s *ClaimsHandler) MakeClaims(c *gin.Context) *sv.Claims {
	sess, _ := c.Cookie("session")
	claims := &sv.Claims{
		SessionID:     sess,
		Role:          "nobody",
		Rate:          "10M",
		Domain:        "webtor.io",
		RemoteAddress: s.getRemoteAddress(c.Request),
		Agent:         c.Request.Header.Get("User-Agent"),
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(24 * 7 * time.Hour).Unix(),
		},
	}
	return claims
}
