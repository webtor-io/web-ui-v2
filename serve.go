package main

import (
	"net/http"

	"github.com/gin-contrib/multitemplate"
	"github.com/gin-gonic/gin"
	"github.com/webtor-io/web-ui-v2/services/web/auth"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	cs "github.com/webtor-io/common-services"
	s "github.com/webtor-io/web-ui-v2/services"
	j "github.com/webtor-io/web-ui-v2/services/job"
	w "github.com/webtor-io/web-ui-v2/services/web"
	wa "github.com/webtor-io/web-ui-v2/services/web/action"
	wi "github.com/webtor-io/web-ui-v2/services/web/index"
	wj "github.com/webtor-io/web-ui-v2/services/web/job"
	p "github.com/webtor-io/web-ui-v2/services/web/profile"
	wr "github.com/webtor-io/web-ui-v2/services/web/resource"
)

func makeServeCMD() cli.Command {
	serveCmd := cli.Command{
		Name:    "serve",
		Aliases: []string{"s"},
		Usage:   "Serves web server",
		Action:  serve,
	}
	configureServe(&serveCmd)
	return serveCmd
}

func configureServe(c *cli.Command) {
	c.Flags = cs.RegisterProbeFlags(c.Flags)
	c.Flags = s.RegisterWebFlags(c.Flags)
	c.Flags = s.RegisterApiFlags(c.Flags)
	c.Flags = w.RegisterTemplateHandlerFlags(c.Flags)
	c.Flags = s.RegisterCommonFlags(c.Flags)
	c.Flags = s.RegisterAuthFlags(c.Flags)
	c.Flags = s.RegisterClaimsProviderClientFlags(c.Flags)
	c.Flags = s.RegisterClaimsFlags(c.Flags)
}

func serve(c *cli.Context) error {

	servers := []cs.Servable{}
	// Setting Probe
	probe := cs.NewProbe(c)
	servers = append(servers, probe)
	defer probe.Close()

	// Setting HTTP Client
	cl := http.DefaultClient

	// Setting Api
	api := s.NewApi(c, cl)

	// Setting template renderer
	re := multitemplate.NewRenderer()

	// Setting TemplateManager
	tm := w.NewMyTemplateManager(c, re)

	// Setting JobQueues
	queues := s.NewJobQueues()

	// Setting JobHandler
	jobs := j.NewHandler(tm, api, queues)

	// Setting ClaimsProviderClient
	cpCl := s.NewClaimsProviderClient(c)
	defer cpCl.Close()

	// Setting UserClaims
	claims := s.NewUserClaims(c, cpCl)

	// Setting Gin
	r := gin.Default()
	r.HTMLRender = re

	// Setting Web
	web := s.NewWeb(c, r)
	servers = append(servers, web)
	defer web.Close()

	// Setting Auth
	a := s.NewAuth(c)

	if a != nil {
		err := a.Init()
		if err != nil {
			return err
		}
		auth.RegisterHandler(c, r, tm)
	}

	// Setting ResourceHandler
	wr.RegisterHandler(c, r, tm, api, jobs, claims)

	// Setting JobHandler
	wj.RegisterHandler(r, queues)

	// Setting IndexHandler
	wi.RegisterHandler(c, r, tm)

	// Setting ActionHandler
	wa.RegisterHandler(c, r, tm, jobs, claims)

	// Setting ProfileHandler
	p.RegisterHandler(c, r, tm)

	// Render templates
	tm.Init()

	// Setting Serve
	serve := cs.NewServe(servers...)

	// And SERVE!
	err := serve.Serve()
	if err != nil {
		log.WithError(err).Error("got server error")
	}
	return err
}
