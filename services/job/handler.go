package job

import (
	"github.com/gin-contrib/multitemplate"
	sv "github.com/webtor-io/web-ui-v2/services"
)

type Handler struct {
	api *sv.Api
	q   *sv.JobQueues
	re  multitemplate.Renderer
}

func NewHandler(re multitemplate.Renderer, api *sv.Api, queues *sv.JobQueues) *Handler {
	return &Handler{
		api: api,
		q:   queues,
		re:  re,
	}
}
