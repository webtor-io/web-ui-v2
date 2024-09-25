package claims

import (
	"context"
	"github.com/pkg/errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/urfave/cli"
	"github.com/webtor-io/lazymap"

	proto "github.com/webtor-io/claims-provider/proto"
	"github.com/webtor-io/web-ui-v2/services/auth"
)

const (
	UseFlag = "use-claims"
)

func RegisterFlags(f []cli.Flag) []cli.Flag {
	return append(f,
		cli.BoolFlag{
			Name:   UseFlag,
			Usage:  "use claims",
			EnvVar: "USE_CLAIMS",
		},
	)
}

type Claims struct {
	lazymap.LazyMap
	cl *Client
}

type Data = proto.GetResponse

func New(c *cli.Context, cl *Client) *Claims {
	if !c.Bool(UseFlag) {
		return nil
	}
	return &Claims{
		cl: cl,
		LazyMap: lazymap.New(&lazymap.Config{
			Expire:      time.Minute,
			ErrorExpire: 10 * time.Second,
		}),
	}
}

func (s *Claims) get(email string) (resp *Data, err error) {
	var cl proto.ClaimsProviderClient
	cl, err = s.cl.Get()
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	resp, err = cl.Get(ctx, &proto.GetRequest{Email: email})
	if err != nil {
		return nil, errors.WithMessage(err, "failed to get claims")
	}
	return
}

func (s *Claims) Get(email string) (*Data, error) {
	resp, err := s.LazyMap.Get(email, func() (interface{}, error) {
		return s.get(email)
	})
	if err != nil {
		return nil, err
	}
	return resp.(*Data), nil
}

func (s *Claims) MakeUserClaimsFromContext(c *gin.Context) (*Data, error) {
	u := auth.GetUserFromContext(c)
	r, err := s.Get(u.Email)
	if err != nil {
		return nil, err
	}
	return r, nil
}

type Context struct{}

func GetFromContext(c *gin.Context) *Data {
	if r := c.Request.Context().Value(Context{}); r != nil {
		return r.(*Data)
	}
	return nil
}

func (s *Claims) RegisterHandler(r *gin.Engine) {
	r.Use(func(c *gin.Context) {
		r, err := s.MakeUserClaimsFromContext(c)
		if err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), Context{}, r))
		c.Next()
	})
}
