package main

import (
	wa "github.com/webtor-io/web-ui-v2/handlers/action"
	wau "github.com/webtor-io/web-ui-v2/handlers/auth"
	"github.com/webtor-io/web-ui-v2/handlers/donate"
	we "github.com/webtor-io/web-ui-v2/handlers/embed"
	wee "github.com/webtor-io/web-ui-v2/handlers/embed/example"
	"github.com/webtor-io/web-ui-v2/handlers/ext"
	wi "github.com/webtor-io/web-ui-v2/handlers/index"
	wj "github.com/webtor-io/web-ui-v2/handlers/job"
	"github.com/webtor-io/web-ui-v2/handlers/legal"
	wm "github.com/webtor-io/web-ui-v2/handlers/migration"
	p "github.com/webtor-io/web-ui-v2/handlers/profile"
	wr "github.com/webtor-io/web-ui-v2/handlers/resource"
	sess "github.com/webtor-io/web-ui-v2/handlers/session"
	sta "github.com/webtor-io/web-ui-v2/handlers/static"
	"github.com/webtor-io/web-ui-v2/handlers/support"
	"github.com/webtor-io/web-ui-v2/handlers/tests"
	as "github.com/webtor-io/web-ui-v2/services/abuse_store"
	"github.com/webtor-io/web-ui-v2/services/umami"
	"net/http"

	"github.com/gin-contrib/multitemplate"
	"github.com/gin-gonic/gin"
	"github.com/webtor-io/web-ui-v2/services/api"
	"github.com/webtor-io/web-ui-v2/services/auth"
	"github.com/webtor-io/web-ui-v2/services/embed"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	cs "github.com/webtor-io/common-services"
	"github.com/webtor-io/web-ui-v2/services"
	"github.com/webtor-io/web-ui-v2/services/claims"
	"github.com/webtor-io/web-ui-v2/services/job"
	"github.com/webtor-io/web-ui-v2/services/template"
	w "github.com/webtor-io/web-ui-v2/services/web"
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
	c.Flags = cs.RegisterPGFlags(c.Flags)
	c.Flags = cs.RegisterProbeFlags(c.Flags)
	c.Flags = api.RegisterFlags(c.Flags)
	c.Flags = w.RegisterFlags(c.Flags)
	c.Flags = services.RegisterFlags(c.Flags)
	c.Flags = auth.RegisterFlags(c.Flags)
	c.Flags = claims.RegisterFlags(c.Flags)
	c.Flags = claims.RegisterClientFlags(c.Flags)
	c.Flags = sess.RegisterFlags(c.Flags)
	c.Flags = sta.RegisterFlags(c.Flags)
	c.Flags = cs.RegisterRedisClientFlags(c.Flags)
	c.Flags = as.RegisterFlags(c.Flags)
	c.Flags = cs.RegisterPprofFlags(c.Flags)
	c.Flags = umami.RegisterFlags(c.Flags)
}

func serve(c *cli.Context) error {
	// Setting DB
	pg := cs.NewPG(c)
	defer pg.Close()

	// Setting Migrations
	m := cs.NewPGMigration(pg)
	err := m.Run()
	if err != nil {
		return err
	}

	// Setting template renderer
	re := multitemplate.NewRenderer()

	// Setting Helper
	helper := w.NewHelper(c)

	// Setting Helper
	umamiHelper := umami.NewHelper(c)

	// Setting TemplateManager
	tm := template.NewManager(re).
		WithHelper(helper).
		WithHelper(umamiHelper).
		WithContextWrapper(w.NewContext)

	var servers []cs.Servable
	// Setting Probe
	probe := cs.NewProbe(c)
	if probe != nil {
		servers = append(servers, probe)
		defer probe.Close()
	}

	// Setting Pprof
	pprof := cs.NewPprof(c)
	if pprof != nil {
		servers = append(servers, pprof)
		defer pprof.Close()
	}
	// Setting Gin
	r := gin.Default()
	r.HTMLRender = re

	// Setting Web
	web, err := w.New(c, r)
	if err != nil {
		return err
	}
	servers = append(servers, web)
	defer web.Close()

	err = sess.RegisterHandler(c, r)
	if err != nil {
		return err
	}

	// Setting Auth
	a := auth.New(c)

	if a != nil {
		err := a.Init()
		if err != nil {
			return err
		}
		a.RegisterHandler(r)
		wau.RegisterHandler(r, tm)
	}

	// Setting Claims Client
	cpCl := claims.NewClient(c)
	defer cpCl.Close()

	// Setting UserClaims
	uc := claims.New(c, cpCl)
	if uc != nil {
		// Setting UserClaimsHandler
		uc.RegisterHandler(r)
	}

	// Setting HTTP Client
	cl := http.DefaultClient

	// Setting Api
	sapi := api.New(c, cl)

	// Setting ApiClaimsHandler
	sapi.RegisterHandler(r)

	err = sta.RegisterHandler(c, r)
	if err != nil {
		return err
	}

	// Setting Migration from v1 to v2
	wm.RegisterHandler(r)

	// Setting Redis
	redis := cs.NewRedisClient(c)
	defer redis.Close()

	// Setting JobQueues
	queues := job.NewQueues(job.NewStorage(redis, gin.Mode()))

	// Setting JobHandler
	jobs := wj.New(queues, tm, sapi)

	jobs.RegisterHandler(r)

	// Setting AbuseStore
	asc := as.New(c)

	if asc != nil {
		defer asc.Close()
		// Setting Support
		support.RegisterHandler(r, tm, asc)

		// Setting Legal
		legal.RegisterHandler(r, tm)
	}

	// Setting DomainSettings
	ds := embed.NewDomainSettings(pg, uc)

	// Setting ResourceHandler
	wr.RegisterHandler(r, tm, sapi, jobs)

	// Setting IndexHandler
	wi.RegisterHandler(r, tm)

	// Setting ActionHandler
	wa.RegisterHandler(r, tm, jobs)

	// Setting ProfileHandler
	p.RegisterHandler(r, tm)

	// Setting EmbedExamplesHandler
	wee.RegisterHandler(r, tm)

	// Setting EmbedHandler
	we.RegisterHandler(cl, r, tm, jobs, ds)

	// Setting ExtHandler
	ext.RegisterHandler(r, tm)

	// Setting Donate
	donate.RegisterHandler(r)

	// Setting Tests
	tests.RegisterHandler(r, tm)

	// Render templates
	err = tm.Init()
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
