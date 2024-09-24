package web

import (
	"github.com/gin-gonic/gin"
	"github.com/webtor-io/web-ui-v2/handlers/session"
	"github.com/webtor-io/web-ui-v2/services/auth"
	"github.com/webtor-io/web-ui-v2/services/claims"
)

type Context struct {
	Data      any
	CSRF      string
	SessionID string
	Err       error
	User      *auth.User
	Claims    *claims.Data
}

func NewContext(c *gin.Context, obj any, err error) any {
	user := auth.GetUserFromContext(c)
	cl := claims.GetFromContext(c)
	sess := session.GetFromContext(c)

	return &Context{
		Data:      obj,
		CSRF:      sess.CSRF,
		Err:       err,
		User:      user,
		Claims:    cl,
		SessionID: sess.ID,
	}
}
