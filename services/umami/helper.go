package umami

import (
	"github.com/urfave/cli"
)

var (
	UseFlag       = "umami-use"
	SrcFlag       = "umami-src"
	WebsiteIDFlag = "umami-website-id"
	HostUrlFlag   = "umami-host-url"
)

func RegisterFlags(f []cli.Flag) []cli.Flag {
	return append(f,
		cli.BoolFlag{
			Name:   UseFlag,
			Usage:  "use umami",
			EnvVar: "USE_UMAMI",
		},
		cli.StringFlag{
			Name:   WebsiteIDFlag,
			Usage:  "umami website-id",
			EnvVar: "UMAMI_WEBSITE_ID",
		},
		cli.StringFlag{
			Name:   HostUrlFlag,
			Usage:  "umami host url",
			EnvVar: "UMAMI_HOST_URL",
		},
	)
}

type Helper struct {
	use       bool
	WebsiteID string `json:"website_id,omitempty"`
	HostURL   string `json:"host_url,omitempty"`
}

func NewHelper(cli *cli.Context) *Helper {
	return &Helper{
		use:       cli.Bool(UseFlag),
		WebsiteID: cli.String(WebsiteIDFlag),
		HostURL:   cli.String(HostUrlFlag),
	}
}

func (s *Helper) UseUmami() bool {
	return s.use
}

func (s *Helper) UmamiConfig() *Helper {
	return s
}
