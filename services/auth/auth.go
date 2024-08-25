package auth

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/supertokens/supertokens-golang/ingredients/emaildelivery"
	"github.com/supertokens/supertokens-golang/recipe/dashboard"
	"github.com/supertokens/supertokens-golang/recipe/passwordless"
	"github.com/supertokens/supertokens-golang/recipe/passwordless/plessmodels"
	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/recipe/session/errors"
	"github.com/supertokens/supertokens-golang/recipe/session/sessmodels"
	"github.com/supertokens/supertokens-golang/recipe/usermetadata"
	"github.com/supertokens/supertokens-golang/recipe/userroles"
	"github.com/supertokens/supertokens-golang/supertokens"
	"github.com/urfave/cli"
	sv "github.com/webtor-io/web-ui-v2/services"

	defaultErrors "errors"
)

const (
	supertokensHostFlag = "supertokens-host"
	supertokensPortFlag = "supertokens-port"
	UseAuthFlag         = "use-auth"
)

func RegisterFlags(f []cli.Flag) []cli.Flag {
	return append(f,
		cli.StringFlag{
			Name:   supertokensHostFlag,
			Usage:  "supertokens host",
			Value:  "",
			EnvVar: "SUPERTOKENS_SERVICE_HOST",
		},
		cli.IntFlag{
			Name:   supertokensPortFlag,
			Usage:  "supertokens port",
			EnvVar: "SUPERTOKENS_SERVICE_PORT",
		},
		cli.BoolFlag{
			Name:   UseAuthFlag,
			Usage:  "use auth",
			EnvVar: "USE_AUTH",
		},
	)
}

type Auth struct {
	url        string
	smtpUser   string
	smtpPass   string
	smtpSecure bool
	smtpHost   string
	smtpPort   int
	domain     string
}

func New(c *cli.Context) *Auth {
	if !c.Bool(UseAuthFlag) {
		return nil
	}
	return &Auth{
		url:        c.String(supertokensHostFlag) + ":" + c.String(supertokensPortFlag),
		smtpUser:   c.String(sv.SMTPUserFlag),
		smtpPass:   c.String(sv.SMTPPassFlag),
		smtpHost:   c.String(sv.SMTPHostFlag),
		smtpSecure: c.BoolT(sv.SMTPSecureFlag),
		smtpPort:   c.Int(sv.SMTPPortFlag),
		domain:     c.String(sv.DomainFlag),
	}
}

func (s *Auth) Init() error {
	smtpSettings := emaildelivery.SMTPSettings{
		Host: s.smtpHost,
		From: emaildelivery.SMTPFrom{
			Name:  "Webtor",
			Email: s.smtpUser,
		},
		Username: &s.smtpUser,
		Port:     s.smtpPort,
		Password: s.smtpPass,
		Secure:   s.smtpSecure,
		TLSConfig: &tls.Config{
			InsecureSkipVerify: false,
			ServerName:         s.smtpHost,
		},
	}
	apiBasePath := "/auth"
	websiteBasePath := "/auth"
	return supertokens.Init(supertokens.TypeInput{
		// Debug: true,
		Supertokens: &supertokens.ConnectionInfo{
			// https://try.supertokens.com is for demo purposes. Replace this with the address of your core instance (sign up on supertokens.com), or self host a core.
			ConnectionURI: s.url,
			// APIKey: <API_KEY(if configured)>,
		},
		AppInfo: supertokens.AppInfo{
			AppName:         "webtor",
			APIDomain:       s.domain,
			WebsiteDomain:   s.domain,
			APIBasePath:     &apiBasePath,
			WebsiteBasePath: &websiteBasePath,
		},
		RecipeList: []supertokens.Recipe{
			passwordless.Init(plessmodels.TypeInput{
				FlowType: "MAGIC_LINK",
				ContactMethodEmail: plessmodels.ContactMethodEmailConfig{
					Enabled: true,
				},
				EmailDelivery: &emaildelivery.TypeInput{
					Service: passwordless.MakeSMTPService(emaildelivery.SMTPServiceConfig{
						Settings: smtpSettings,
						Override: func(originalImplementation emaildelivery.SMTPInterface) emaildelivery.SMTPInterface {
							*originalImplementation.GetContent = func(input emaildelivery.EmailType, userContext supertokens.UserContext) (emaildelivery.EmailContent, error) {

								email := input.PasswordlessLogin.Email

								// magic link
								urlWithLinkCode := *input.PasswordlessLogin.UrlWithLinkCode
								body := fmt.Sprintf("<a href=\"%v\">Login to your account!</a>", urlWithLinkCode)

								// send some custom email content
								return emaildelivery.EmailContent{
									Body:    body,
									IsHtml:  true,
									Subject: "Login to your account!",
									ToEmail: email,
								}, nil

							}

							return originalImplementation
						},
					}),
				},
			}),
			session.Init(nil), // initializes session features
			dashboard.Init(nil),
			usermetadata.Init(nil),
			userroles.Init(nil),
		},
	})
}

type User struct {
	ID      string
	Email   string
	Expired bool
}

func GetUserFromContext(c *gin.Context) *User {
	u := &User{}
	if sessionContainer := session.GetSessionFromRequestContext(c.Request.Context()); sessionContainer != nil {
		userID := sessionContainer.GetUserID()
		userInfo, err := passwordless.GetUserByID(userID)
		if err == nil && userInfo != nil {
			u.ID = userInfo.ID
			u.Email = *userInfo.Email
		}
	}
	if err := c.Request.Context().Value(ErrorContext{}); err != nil {
		if defaultErrors.As(err.(error), &errors.TryRefreshTokenError{}) {
			u.Expired = true
		}
	}
	return u
}

type ErrorContext struct{}

func myVerifySession(options *sessmodels.VerifySessionOptions, otherHandler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sess, err := session.GetSession(r, w, options)
		if err != nil {
			ctx := context.WithValue(r.Context(), ErrorContext{}, err)
			r := r.WithContext(ctx)
			if defaultErrors.As(err, &errors.TryRefreshTokenError{}) {
				if r.Header.Get("X-Requested-With") != "XMLHttpRequest" {
					otherHandler(w, r)
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
		if sess != nil {
			ctx := context.WithValue(r.Context(), sessmodels.SessionContext, sess)
			otherHandler(w, r.WithContext(ctx))
		} else {
			otherHandler(w, r)
		}
	}
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

func (s *Auth) RegisterHandler(r *gin.Engine) {
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
}
