package internal

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	HTTPAnyMethod    = "ANY"
	HTTPGETMethod    = http.MethodGet
	HTTPPOSTMethod   = http.MethodPost
	HTTPDELETEMethod = http.MethodDelete
	HTTPPATCHMethod  = http.MethodPatch
	HTTPPUTMethod    = http.MethodPut
)

var methodMap map[string]string = map[string]string{
	"ANY":    HTTPAnyMethod,
	"GET":    HTTPGETMethod,
	"POST":   HTTPPOSTMethod,
	"DELETE": HTTPDELETEMethod,
	"PATCH":  HTTPPATCHMethod,
	"PUT":    HTTPPUTMethod,
}

type httpServerDescriptor struct {
	Names            string
	methodDescriptor []*httpMethodDescriptor
}

type httpMethodDescriptor struct {
	Name     string
	FullPath string
	Method   string
}

type httpServer struct {
	// map[serverName]map[fullPath]
	serverMethod map[string]map[string]struct{}
	engine       *gin.Engine
}

func (h *httpServer) StartWithPort(port int) error {

	return h.engine.Run(fmt.Sprintf(":%d", port))
}
