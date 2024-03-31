package auth

import (
	"context"
	defaultErrors "errors"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/recipe/session/errors"
	"github.com/supertokens/supertokens-golang/recipe/session/sessmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
	"github.com/urfave/cli"
	w "github.com/webtor-io/web-ui-v2/services/web"
)

type LoginData struct{}

type LogoutData struct{}

type VerifyData struct {
	PreAuthSessionId string
}

type Handler struct {
	tm *w.TemplateManager
}

type AuthError struct{}

func myVerifySession(options *sessmodels.VerifySessionOptions, otherHandler http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := session.GetSession(r, w, options)
		if err != nil {
			if defaultErrors.As(err, &errors.TryRefreshTokenError{}) {
				if r.Header.Get("X-Requested-With") != "XMLHttpRequest" {
					ctx := context.WithValue(r.Context(), AuthError{}, err)
					otherHandler(w, r.WithContext(ctx))
					return
				}
				// This means that the session exists, but the access token
				// has expired.

				// You can handle this in a custom way by sending a 401.
				// Or you can call the errorHandler middleware as shown below
			} else if defaultErrors.As(err, &errors.UnauthorizedError{}) {
				otherHandler(w, r)
				return
				// This means that the session does not exist anymore.

				// You can handle this in a custom way by sending a 401.
				// Or you can call the errorHandler middleware as shown below
			} else if defaultErrors.As(err, &errors.InvalidClaimError{}) {
				otherHandler(w, r)
				return
				// The user is missing some required claim.
				// You can pass the missing claims to the frontend and handle it there
			}

			// OR you can use this errorHandler which will
			// handle all of the above errors in the default way
			err = supertokens.ErrorHandler(err, r, w)
			if err != nil {
				// TODO: send a 500 error to the frontend
			}
			return
		}
		if session != nil {
			ctx := context.WithValue(r.Context(), sessmodels.SessionContext, session)
			otherHandler(w, r.WithContext(ctx))
		} else {
			otherHandler(w, r)
		}
	})
}

func verifySession(options *sessmodels.VerifySessionOptions) gin.HandlerFunc {
	return func(c *gin.Context) {
		myVerifySession(options, func(rw http.ResponseWriter, r *http.Request) {
			c.Request = c.Request.WithContext(r.Context())
			c.Next()
		})(c.Writer, c.Request)
		// we call Abort so that the next handler in the chain is not called, unless we call Next explicitly
		c.Abort()
	}
}

func RegisterHandler(c *cli.Context, r *gin.Engine, tm *w.TemplateManager) error {
	h := &Handler{
		tm: tm,
	}
	// CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "DELETE", "PUT", "OPTIONS"},
		AllowHeaders:     append([]string{"content-type"}, supertokens.GetAllCORSHeaders()...),
		MaxAge:           1 * time.Minute,
		AllowCredentials: true,
	}))
	r.Use(func(c *gin.Context) {
		supertokens.Middleware(http.HandlerFunc(
			func(rw http.ResponseWriter, r *http.Request) {
				c.Next()
			})).ServeHTTP(c.Writer, c.Request)
		// we call Abort so that the next handler in the chain is not called, unless we call Next explicitly
		c.Abort()
	})
	sessionRequired := false
	r.Use(verifySession(&sessmodels.VerifySessionOptions{
		SessionRequired: &sessionRequired,
	}))
	r.GET("/login", h.login)
	r.GET("/logout", h.logout)
	r.GET("/auth/verify", h.verify)

	h.tm.RegisterViews("auth/*")

	return nil
}

func (s *Handler) login(c *gin.Context) {
	s.tm.MakeTemplate("auth/login").HTML(http.StatusOK, c, LoginData{})
}

func (s *Handler) logout(c *gin.Context) {
	s.tm.MakeTemplate("auth/logout").HTML(http.StatusOK, c, LogoutData{})
}

func (s *Handler) verify(c *gin.Context) {
	s.tm.MakeTemplate("auth/verify").HTML(http.StatusOK, c, &VerifyData{
		PreAuthSessionId: c.Query("preAuthSessionId"),
	})
}
