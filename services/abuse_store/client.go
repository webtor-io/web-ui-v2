package abuse_store

import (
	"fmt"
	"google.golang.org/grpc/credentials/insecure"
	"sync"

	"github.com/urfave/cli"
	as "github.com/webtor-io/abuse-store/proto"
	"google.golang.org/grpc"
)

const (
	HostFlag = "abuse-host"
	PortFlag = "abuse-port"
	UseFlag  = "use-abuse"
)

func RegisterFlags(f []cli.Flag) []cli.Flag {
	return append(f,
		cli.StringFlag{
			Name:   HostFlag,
			Usage:  "abuse store host",
			Value:  "",
			EnvVar: "ABUSE_STORE_SERVICE_HOST",
		},
		cli.IntFlag{
			Name:   PortFlag,
			Usage:  "port of the redis service",
			Value:  50051,
			EnvVar: "ABUSE_STORE_SERVICE_PORT",
		},
		cli.BoolFlag{
			Name:   UseFlag,
			Usage:  "use abuse",
			EnvVar: "USE_ABUSE_STORE",
		},
	)
}

type Client struct {
	once sync.Once
	cl   as.AbuseStoreClient
	err  error
	host string
	port int
	conn *grpc.ClientConn
}

func New(c *cli.Context) *Client {
	if !c.Bool(UseFlag) {
		return nil
	}
	return &Client{
		host: c.String(HostFlag),
		port: c.Int(PortFlag),
	}
}

func (s *Client) Get() (as.AbuseStoreClient, error) {
	s.once.Do(func() {
		addr := fmt.Sprintf("%s:%d", s.host, s.port)
		conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			s.err = err
			return
		}
		if err != nil {
			s.err = err
			return
		}
		s.conn = conn
		s.cl = as.NewAbuseStoreClient(conn)
	})
	return s.cl, s.err
}

func (s *Client) Close() {
	if s.conn != nil {
		_ = s.conn.Close()
	}
}
