package web

import (
	"github.com/gin-gonic/gin"
	"github.com/webtor-io/web-ui/handlers/geo"
	"github.com/webtor-io/web-ui/handlers/session"
	"github.com/webtor-io/web-ui/services/auth"
	"github.com/webtor-io/web-ui/services/claims"
	"github.com/webtor-io/web-ui/services/geoip"
)

type Context struct {
	Data      any
	CSRF      string
	SessionID string
	Err       error
	User      *auth.User
	Claims    *claims.Data
	Geo       *geoip.Data
}

func NewContext(c *gin.Context, obj any, err error) any {
	user := auth.GetUserFromContext(c)
	cl := claims.GetFromContext(c)
	sess := session.GetFromContext(c)
	geoData := geo.GetFromContext(c)

	return &Context{
		Data:      obj,
		CSRF:      sess.CSRF,
		Err:       err,
		User:      user,
		Claims:    cl,
		SessionID: sess.ID,
		Geo:       geoData,
	}
}
