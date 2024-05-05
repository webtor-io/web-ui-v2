package web

import (
	"github.com/gin-gonic/gin"
	csrf "github.com/utrack/gin-csrf"
	"github.com/webtor-io/web-ui-v2/services/auth"
	"github.com/webtor-io/web-ui-v2/services/claims"
)

type Context struct {
	Data   any
	CSRF   string
	Err    error
	User   *auth.User
	Claims *claims.Data
}

func NewContext(c *gin.Context, obj any, err error) any {
	user := auth.GetUserFromContext(c)
	cl := claims.GetFromContext(c)

	return &Context{
		Data:   obj,
		CSRF:   csrf.GetToken(c),
		Err:    err,
		User:   user,
		Claims: cl,
	}
}
