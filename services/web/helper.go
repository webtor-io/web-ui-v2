package web

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/urfave/cli"

	h "github.com/dustin/go-humanize"
	log "github.com/sirupsen/logrus"
	"github.com/webtor-io/lazymap"
	"github.com/webtor-io/web-ui-v2/services"
	"github.com/webtor-io/web-ui-v2/services/auth"
	"github.com/webtor-io/web-ui-v2/services/claims"
	"github.com/webtor-io/web-ui-v2/services/job"
	"github.com/webtor-io/web-ui-v2/services/obfuscator"
)

func (s *Helper) MakeJobLogURL(j *job.Job) string {
	return fmt.Sprintf("/queue/%v/job/%v/log", j.Queue, j.ID)
}

func (s *Helper) Log(err error) error {
	log.Error(err)
	return err
}

func (s *Helper) ShortErr(err error) string {
	return strings.Split(err.Error(), ":")[0]
}

func (s *Helper) BitsForHumans(b int64) string {
	return h.Bytes(uint64(b))
}

func (s *Helper) Dev() bool {
	return gin.Mode() == "debug"
}

func (s *Helper) Has(obj any, fieldName string) bool {
	value := reflect.Indirect(reflect.ValueOf(obj))
	field := value.FieldByName(fieldName)
	return field.IsValid() && !field.IsNil()
}

type Helper struct {
	assetsHost string
	assetsPath string
	useAuth    bool
	domain     string
	demoMagnet string
	ah         *AssetHashes
}

func NewHelper(c *cli.Context) *Helper {
	return &Helper{
		demoMagnet: c.String(services.DemoMagnetFlag),
		assetsHost: c.String(assetsHostFlag),
		assetsPath: c.String(assetsPathFlag),
		useAuth:    c.Bool(auth.UseAuthFlag),
		domain:     c.String(services.DomainFlag),
		ah:         NewAssetHashes(c.String(assetsPathFlag)),
	}
}

func (s *Helper) HasAds(c *claims.Data) bool {
	if c == nil {
		return false
	}
	return !c.Claims.Site.NoAds
}

func (s *Helper) UseAuth() bool {
	return s.useAuth
}

func (s *Helper) Domain() string {
	return s.domain
}

func (s *Helper) DemoMagnet() string {
	return s.demoMagnet
}

func (s *Helper) IsDemoMagnet(m string) bool {
	return s.demoMagnet == m
}

func (s *Helper) Obfuscate(in string) string {
	return obfuscator.Obfuscate(in)
}

func (s *Helper) Asset(in string) string {
	path := s.assetsHost + "/assets/" + in
	if !s.Dev() {
		h, _ := s.ah.Get(in)
		path += "?" + h
	}
	return path
}

func (s *Helper) DevAsset(in string) string {
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
