package services

import (
	"fmt"
	"sync"

	"github.com/urfave/cli"
	proto "github.com/webtor-io/claims-provider/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	claimsProviderHostFlag = "claims-provider-host"
	claimsProviderPortFlag = "claims-provider-port"
)

func RegisterClaimsProviderClientFlags(f []cli.Flag) []cli.Flag {
	return append(f,
		cli.StringFlag{
			Name:   claimsProviderHostFlag,
			Usage:  "claims provider host",
			Value:  "",
			EnvVar: "CLAIMS_PROVIDER_SERVICE_HOST",
		},
		cli.IntFlag{
			Name:   claimsProviderPortFlag,
			Usage:  "claims provider port",
			Value:  50051,
			EnvVar: "CLAIMS_PROVIDER_SERVICE_PORT",
		},
	)
}

type ClaimsProviderClient struct {
	once sync.Once
	cl   proto.ClaimsProviderClient
	err  error
	host string
	port int
	conn *grpc.ClientConn
}

func NewClaimsProviderClient(c *cli.Context) *ClaimsProviderClient {
	return &ClaimsProviderClient{
		host: c.String(claimsProviderHostFlag),
		port: c.Int(claimsProviderPortFlag),
	}
}

func (s *ClaimsProviderClient) Get() (proto.ClaimsProviderClient, error) {
	s.once.Do(func() {
		addr := fmt.Sprintf("%s:%d", s.host, s.port)
		conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			s.err = err
			return
		}
		s.conn = conn
		s.cl = proto.NewClaimsProviderClient(conn)
	})
	return s.cl, s.err
}

func (s *ClaimsProviderClient) Close() {
	if s.conn != nil {
		_ = s.conn.Close()
	}
}
