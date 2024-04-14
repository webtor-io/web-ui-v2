package web

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	sv "github.com/webtor-io/web-ui-v2/services"
)

type ClaimsHandler struct {
	uc *sv.UserClaims
}

func NewClaimsHandler(uc *sv.UserClaims) *ClaimsHandler {
	return &ClaimsHandler{
		uc: uc,
	}
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

func (s *ClaimsHandler) MakeClaims(c *gin.Context) (*sv.Claims, error) {
	sess, _ := c.Cookie("session")
	claims := &sv.Claims{
		SessionID:     sess,
		Domain:        "webtor.io",
		RemoteAddress: s.getRemoteAddress(c.Request),
		Agent:         c.Request.Header.Get("User-Agent"),
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(24 * 7 * time.Hour).Unix(),
		},
	}
	if s.uc != nil {
		cl, err := s.uc.GetFromContext(c)
		if err != nil {
			return nil, err
		}
		claims.Role = cl.Context.GetTier().GetName()
		rate := cl.GetClaims().GetConnection().GetRate()
		if rate > 0 {
			claims.Rate = fmt.Sprintf("%dM", rate)
		}
		h := sha1.New()
		h.Write([]byte(cl.Email))
		hashEmail := hex.EncodeToString(h.Sum(nil))
		claims.SessionID = hashEmail
	}
	return claims, nil
}
