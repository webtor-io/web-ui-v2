package embed

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ExampleData struct {
}

func (s *Handler) example(c *gin.Context) {
	s.tb.Build("embed/example").HTML(http.StatusOK, c, &ExampleData{})
}
