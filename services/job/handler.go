package job

import (
	sv "github.com/webtor-io/web-ui-v2/services"
)

type Handler struct {
	api *sv.Api
	q   *sv.JobQueues
}

func NewHandler(api *sv.Api, queues *sv.JobQueues) *Handler {
	return &Handler{
		api: api,
		q:   queues,
	}
}
