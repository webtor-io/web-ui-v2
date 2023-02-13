package main

import (
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

	// Setting JobQueues
	queues := s.NewJobQueues()

	// Setting JobHandler
	jobs := j.NewHandler(api, queues)

	// Setting Web
	web := s.NewWeb(c)
	defer web.Close()

	// Setting ResourceHandler
	resourceH := wr.NewHandler(c, api, jobs)
	web.RegisterHandler(resourceH)

	// Setting JobHandler
	jobH := wj.NewHandler(queues)
	web.RegisterHandler(jobH)

	// Setting IndexHandler
	indexH := wi.NewHandler(c)
	web.RegisterHandler(indexH)

	// Setting ActionHandler
	actionH := wa.NewHandler(c, jobs)
	web.RegisterHandler(actionH)

	// Setting Serve
	serve := cs.NewServe(probe, web)

	// And SERVE!
	err := serve.Serve()
	if err != nil {
		log.WithError(err).Error("got server error")
	}
	return err
}
