package embed

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type GetData struct {
	CheckScript string
	CheckHash   string
	ID          string
}

func (s *Handler) get(c *gin.Context) {
	id := c.Query("id")
	code := uuid.New().String()
	h := sha1.New()
	h.Write([]byte(id + code))
	gd := GetData{
		CheckHash:   hex.EncodeToString(h.Sum(nil)),
		CheckScript: s.generateCheckScript(code, id),
		ID:          id,
	}
	s.tb.Build("embed/get").HTML(http.StatusOK, c, gd)
}

func (s *Handler) generateCheckScript(code string, id string) string {
	return fmt.Sprintf(`
		var found = false;
		var scripts = document.getElementsByTagName('script');
			for (var i = scripts.length; i--;) {
				if (
					scripts[i].src.includes('https://cdn.jsdelivr.net/npm/@webtor/') ||
					scripts[i].src.includes('http://localhost:9009/')
				) {
					found = '%v';
				}
			}
		var f = window.frames['webtor-%v'];
		f.contentWindow.postMessage({id: '%v', name: 'check', data: found}, '*');
	`, code, id, id)
}
