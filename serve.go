package main

import (
	"net/http"

	"github.com/gin-contrib/multitemplate"
	"github.com/gin-gonic/gin"
	"github.com/webtor-io/web-ui-v2/services/api"
	"github.com/webtor-io/web-ui-v2/services/auth"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	cs "github.com/webtor-io/common-services"
	"github.com/webtor-io/web-ui-v2/services"
	"github.com/webtor-io/web-ui-v2/services/claims"
	"github.com/webtor-io/web-ui-v2/services/job"
	"github.com/webtor-io/web-ui-v2/services/template"
	w "github.com/webtor-io/web-ui-v2/services/web"
	wa "github.com/webtor-io/web-ui-v2/services/web/action"
	wau "github.com/webtor-io/web-ui-v2/services/web/auth"
	we "github.com/webtor-io/web-ui-v2/services/web/embed"
	wee "github.com/webtor-io/web-ui-v2/services/web/embed/example"
	wi "github.com/webtor-io/web-ui-v2/services/web/index"
	wj "github.com/webtor-io/web-ui-v2/services/web/job"
	wm "github.com/webtor-io/web-ui-v2/services/web/migration"
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
	c.Flags = api.RegisterFlags(c.Flags)
	c.Flags = w.RegisterFlags(c.Flags)
	c.Flags = services.RegisterFlags(c.Flags)
	c.Flags = auth.RegisterFlags(c.Flags)
	c.Flags = claims.RegisterFlags(c.Flags)
	c.Flags = claims.RegisterClientFlags(c.Flags)
	// c.Flags = cs.RegisterRedisClientFlags(c.Flags)
}

func serve(c *cli.Context) error {

	// Setting Redis
	// redis := cs.NewRedisClient(c)
	// defer redis.Close()

	servers := []cs.Servable{}
	// Setting Probe
	probe := cs.NewProbe(c)
	servers = append(servers, probe)
	defer probe.Close()

	// Setting template renderer
	re := multitemplate.NewRenderer()

	// Setting Gin
	r := gin.Default()
	r.HTMLRender = re

	// Setting Web
	web := w.New(c, r)
	servers = append(servers, web)
	defer web.Close()

	// Setting Migration from v1 to v2
	wm.RegisterHandler(r)

	// Setting HTTP Client
	cl := http.DefaultClient

	// Setting Api
	api := api.New(c, cl)

	// Setting Helper
	helper := w.NewHelper(c)

	// Setting TemplateManager
	tm := template.NewManager(re).WithHelper(helper).WithContextWrapper(w.NewContext)

	// Setting JobQueues
	queues := job.NewQueues()

	// Setting JobHandler
	jobs := wj.New(queues, tm, api)

	jobs.RegisterHandler(r)

	// Setting Auth
	a := auth.New(c)

	if a != nil {
		err := a.Init()
		if err != nil {
			return err
		}
		a.RegisterHandler(r)
		wau.RegisterHandler(c, r, tm)
	}

	// Setting Claims Client
	cpCl := claims.NewClient(c)
	defer cpCl.Close()

	// Setting UserClaims
	uc := claims.New(c, cpCl)
	if uc != nil {
		// Setting UserClaimsHandler
		uc.RegisterHandler(c, r)
	}

	// Setting ApiClaimsHandler
	api.RegisterHandler(c, r)

	// Setting ResourceHandler
	wr.RegisterHandler(c, r, tm, api, jobs)

	// Setting IndexHandler
	wi.RegisterHandler(c, r, tm)

	// Setting ActionHandler
	wa.RegisterHandler(c, r, tm, jobs)

	// Setting ProfileHandler
	p.RegisterHandler(c, r, tm)

	// Setting EmbedExamplesHandler
	wee.RegisterHandler(c, r, tm)

	// Setting EmbedHandler
	we.RegisterHandler(c, r, tm, jobs)

	// Render templates
	err := tm.Init()
	if err != nil {
		return err
	}

	// Setting Serve
	serve := cs.NewServe(servers...)

	// And SERVE!
	err = serve.Serve()
	if err != nil {
		log.WithError(err).Error("got server error")
	}
	return err
}
