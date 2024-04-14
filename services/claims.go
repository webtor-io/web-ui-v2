package services

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/urfave/cli"
	"github.com/webtor-io/lazymap"

	proto "github.com/webtor-io/claims-provider/proto"
)

const (
	UseClaimsFlag = "use-claims"
)

func RegisterClaimsFlags(f []cli.Flag) []cli.Flag {
	return append(f,
		cli.BoolFlag{
			Name:   UseClaimsFlag,
			Usage:  "use claims",
			EnvVar: "USE_CLAIMS",
		},
	)
}

type UserClaims struct {
	lazymap.LazyMap
	cl *ClaimsProviderClient
}

func NewUserClaims(c *cli.Context, cl *ClaimsProviderClient) *UserClaims {
	if !c.Bool(UseClaimsFlag) || !c.Bool(UseAuthFlag) {
		return nil
	}
	return &UserClaims{
		cl: cl,
		LazyMap: lazymap.New(&lazymap.Config{
			Expire:      time.Minute,
			ErrorExpire: 10 * time.Second,
		}),
	}
}

func (s *UserClaims) get(email string) (resp *proto.GetResponse, err error) {
	var cl proto.ClaimsProviderClient
	cl, err = s.cl.Get()
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	resp, err = cl.Get(ctx, &proto.GetRequest{Email: email})
	if err != nil {
		return nil, err
	}
	return
}

func (s *UserClaims) Get(email string) (*proto.GetResponse, error) {
	resp, err := s.LazyMap.Get(email, func() (interface{}, error) {
		return s.get(email)
	})
	if err != nil {
		return nil, err
	}
	return resp.(*proto.GetResponse), nil
}

type UserWithClaims struct {
	*User
	*proto.GetResponse
}

func (s *UserClaims) GetFromContext(c *gin.Context) (*UserWithClaims, error) {
	u := GetUserFromContext(c)
	r, err := s.Get(u.Email)
	if err != nil {
		return nil, err
	}
	return &UserWithClaims{u, r}, nil
}
