package embed

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/webtor-io/web-ui-v2/services/api"
	"github.com/webtor-io/web-ui-v2/services/embed"
	"github.com/webtor-io/web-ui-v2/services/job"
	"github.com/webtor-io/web-ui-v2/services/web/job/script"
)

type PostArgs struct {
	ID            string
	EmbedSettings *script.EmbedSettings
	Claims        *api.Claims
}

type PostData struct {
	ID             string
	EmbedSettings  *script.EmbedSettings
	DomainSettings *embed.DomainSettingsData
	Job            *job.Job
}

func (s *Handler) bindPostArgs(c *gin.Context) (*PostArgs, error) {
	rawSettings := c.PostForm("settings")
	var settings script.EmbedSettings
	err := json.Unmarshal([]byte(rawSettings), &settings)
	if err != nil {
		return nil, err
	}
	id := c.Query("id")

	return &PostArgs{
		ID:            id,
		EmbedSettings: &settings,
		Claims:        api.GetClaimsFromContext(c),
	}, nil

}

func (s *Handler) post(c *gin.Context) {
	tpl := s.tb.Build("embed/post")
	pd := PostData{}
	args, err := s.bindPostArgs(c)
	if err != nil {
		tpl.HTMLWithErr(err, http.StatusBadRequest, c, pd)
		return
	}
	pd.ID = args.ID
	u, err := url.Parse(args.EmbedSettings.Referer)
	if err != nil {
		tpl.HTMLWithErr(err, http.StatusBadRequest, c, pd)
		return
	}
	dsd, err := s.ds.Get(u.Hostname())
	if err != nil {
		tpl.HTMLWithErr(err, http.StatusBadRequest, c, pd)
		return
	}
	pd.EmbedSettings = args.EmbedSettings
	pd.DomainSettings = dsd
	job, err := s.jobs.Embed(c, s.hCl, args.Claims, args.EmbedSettings)
	if err != nil {
		tpl.HTMLWithErr(err, http.StatusBadRequest, c, pd)
		return
	}
	pd.Job = job
	tpl.HTML(http.StatusAccepted, c, pd)
}
