package action

import (
	"github.com/webtor-io/web-ui-v2/services"
)

func GetDurationSec(mp *services.MediaProbe) string {
	return mp.Format.Duration
}
