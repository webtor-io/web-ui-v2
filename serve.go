package main

import (
	"github.com/gin-contrib/multitemplate"
	"github.com/gin-gonic/gin"
	"net/http"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	cs "github.com/webtor-io/common-services"
	s "github.com/webtor-io/web-ui-v2/services"
	j "github.com/webtor-io/web-ui-v2/services/job"
	w "github.com/webtor-io/web-ui-v2/services/web"
	wa "github.com/webtor-io/web-ui-v2/services/web/action"
	wi "github.com/webtor-io/web-ui-v2/services/web/index"
	wj "github.com/webtor-io/web-ui-v2/services/web/job"
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
}

func serve(c *cli.Context) error {
	// Setting Probe
	probe := cs.NewProbe(c)
	defer probe.Close()

	// Setting HTTP Client
	cl := http.DefaultClient

	// Setting Api
	api := s.NewApi(c, cl)

	// Setting template renderer
	re := multitemplate.NewRenderer()

	// Setting JobQueues
	queues := s.NewJobQueues()

	// Setting JobHandler
	jobs := j.NewHandler(re, api, queues)

	// Setting Gin
	r := gin.Default()
	r.HTMLRender = re

	// Setting Web
	web := s.NewWeb(c, r)
	defer web.Close()

	// Setting ResourceHandler
	wr.RegisterHandler(c, r, re, api, jobs)

	// Setting JobHandler
	wj.RegisterHandler(r, queues)

	// Setting IndexHandler
	wi.RegisterHandler(c, r, re)

	// Setting ActionHandler
	wa.RegisterHandler(c, r, re, jobs)

	// Setting Serve
	serve := cs.NewServe(probe, web)

	// And SERVE!
	err := serve.Serve()
	if err != nil {
		log.WithError(err).Error("got server error")
	}
	return err
}
