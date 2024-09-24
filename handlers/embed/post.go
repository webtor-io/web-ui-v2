package embed

import (
	"encoding/json"
	"github.com/webtor-io/web-ui-v2/services/models"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/webtor-io/web-ui-v2/services/api"
	"github.com/webtor-io/web-ui-v2/services/embed"
	"github.com/webtor-io/web-ui-v2/services/job"
)

type PostArgs struct {
	ID            string
	EmbedSettings *models.EmbedSettings
	Claims        *api.Claims
}

type PostData struct {
	ID             string
	EmbedSettings  *models.EmbedSettings
	DomainSettings *embed.DomainSettingsData
	Job            *job.Job
}

func (s *Handler) bindPostArgs(c *gin.Context) (*PostArgs, error) {
	rawSettings := c.PostForm("settings")
	var settings models.EmbedSettings
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
	embedJob, err := s.jobs.Embed(c, s.hCl, args.Claims, args.EmbedSettings)
	if err != nil {
		tpl.HTMLWithErr(err, http.StatusBadRequest, c, pd)
		return
	}
	pd.Job = embedJob
	tpl.HTML(http.StatusAccepted, c, pd)
}
