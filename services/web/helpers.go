package web

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"strings"

	h "github.com/dustin/go-humanize"
	log "github.com/sirupsen/logrus"
	sv "github.com/webtor-io/web-ui-v2/services"
)

func MakeJobLogURL(j *sv.Job) string {
	return fmt.Sprintf("/queue/%v/job/%v/log", j.Queue, j.ID)
}

func Log(err error) error {
	log.Error(err)
	return err
}

func ShortErr(err error) string {
	return strings.Split(err.Error(), ":")[0]
}

func BitsForHumans(b int64) string {
	return h.Bytes(uint64(b))
}

func Dev() bool {
	return gin.Mode() == "debug"
}

type ErrorData struct {
	Err error
}
type CSRFData struct {
	CSRF string
}

type JobData struct {
	Job *sv.Job
}
