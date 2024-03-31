package job

import (
	sv "github.com/webtor-io/web-ui-v2/services"
	w "github.com/webtor-io/web-ui-v2/services/web"
)

type Handler struct {
	api *sv.Api
	q   *sv.JobQueues
	tm  *w.TemplateManager
}

func NewHandler(tm *w.TemplateManager, api *sv.Api, queues *sv.JobQueues) *Handler {
	return &Handler{
		api: api,
		q:   queues,
		tm:  tm,
	}
}
