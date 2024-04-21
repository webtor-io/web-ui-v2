package web

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"html/template"
	"io"
	"os"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/urfave/cli"
	csrf "github.com/utrack/gin-csrf"

	h "github.com/dustin/go-humanize"
	log "github.com/sirupsen/logrus"
	"github.com/webtor-io/lazymap"
	"github.com/webtor-io/web-ui-v2/services"
	"github.com/webtor-io/web-ui-v2/services/auth"
	"github.com/webtor-io/web-ui-v2/services/claims"
	"github.com/webtor-io/web-ui-v2/services/job"
	st "github.com/webtor-io/web-ui-v2/services/template"
)

func MakeJobLogURL(j *job.Job) string {
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

func Has(obj any, fieldName string) bool {
	value := reflect.Indirect(reflect.ValueOf(obj))
	field := value.FieldByName(fieldName)
	return field.IsValid() && !field.IsNil()
}

type JobData struct {
	Job *job.Job
}
type Helper struct {
	assetsHost string
	assetsPath string
	useAuth    bool
	domain     string
	ah         *AssetHashes
}

func NewHelper(c *cli.Context) *Helper {
	return &Helper{
		assetsHost: c.String(assetsHostFlag),
		assetsPath: c.String(assetsPathFlag),
		useAuth:    c.Bool(auth.UseAuthFlag),
		domain:     c.String(services.DomainFlag),
		ah:         NewAssetHashes(c.String(assetsPathFlag)),
	}
}

func (s *Helper) GetFuncs() template.FuncMap {
	return template.FuncMap{
		"asset":         s.MakeAsset,
		"devAsset":      s.MakeDevAsset,
		"makeJobLogURL": MakeJobLogURL,
		"bitsForHumans": BitsForHumans,
		"log":           Log,
		"shortErr":      ShortErr,
		"dev":           Dev,
		"useAuth":       s.UseAuth,
		"domain":        s.Domain,
		"has":           Has,
	}
}

func (s *Helper) UseAuth() bool {
	return s.useAuth
}

func (s *Helper) Domain() string {
	return s.domain
}

func (s *Helper) MakeAsset(in string) string {
	path := s.assetsHost + "/assets/" + in
	if !Dev() {
		h, _ := s.ah.Get(in)
		path += "?" + h
	}
	return path
}

func (s *Helper) MakeDevAsset(in string) string {
	return s.assetsHost + "/assets/dev/" + in
}

type AssetHashes struct {
	lazymap.LazyMap
	path string
}

func (s *AssetHashes) get(name string) (hash string, err error) {
	f, err := os.Open(s.path + "/" + name)
	if err != nil {
		return "", err
	}
	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

func (s *AssetHashes) Get(name string) (string, error) {
	resp, err := s.LazyMap.Get(name, func() (interface{}, error) {
		return s.get(name)
	})
	if err != nil {
		return "", err
	}
	return resp.(string), nil
}

func NewAssetHashes(path string) *AssetHashes {
	return &AssetHashes{
		LazyMap: lazymap.New(&lazymap.Config{}),
		path:    path,
	}
}

func (s *Helper) WrapContext(c *gin.Context, obj any, err error) any {
	return NewContext(c, obj, err)
}

type Context struct {
	Data   any
	CSRF   string
	Err    error
	User   *auth.User
	Claims *claims.Data
}

func NewContext(c *gin.Context, obj any, err error) any {
	user := auth.GetUserFromContext(c)
	cl := claims.GetFromContext(c)

	return &Context{
		Data:   obj,
		CSRF:   csrf.GetToken(c),
		Err:    err,
		User:   user,
		Claims: cl,
	}
}

var _ st.Args = (*Helper)(nil)
