package internal

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
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

type httpServer struct {
	// map[fullPath]HttpMethod
	serverMethod map[string]string
	engine       *gin.Engine
	descriptor   *httpServerDescriptor
}

func NewHttpServer(desc *httpServerDescriptor) *httpServer {
	rev := &httpServer{engine: gin.Default(), serverMethod: map[string]string{}}
	for _, item := range desc.MethodDescriptor {
		rev.serverMethod[item.Path] = item.Method
	}
	return rev
}

func (h *httpServer) MethodHandler() error {
	h.engine.NoRoute(h.NotFoundHandler)
	h.engine.NoMethod(h.NotFoundHandler)
	for path, method := range h.serverMethod {
		if method == HTTPAnyMethod {
			h.engine.Any(path, h.commonHandler)
		} else {
			h.engine.Handle(method, path, h.commonHandler)
		}
	}
	return nil
}

func (h *httpServer) NotFoundHandler(c *gin.Context) {
	c.String(http.StatusNotFound, fmt.Sprintf("not found method: %s, path: %s", c.Request.Method, c.Request.URL.Path))
}
func (h *httpServer) commonHandler(c *gin.Context) {
	//log.Println("method: ", c.Request.Method, " path: ", c.Request.URL.Path)
	if _, ok := h.serverMethod[c.FullPath()]; !ok {
		c.String(http.StatusNotFound, fmt.Sprintf("not found method: %s, path: %s", c.Request.Method, c.Request.URL.Path))
		return
	}
	log.Printf("fullPath: %s, gin.fullPath: %s", c.Request.URL.RawPath, c.FullPath())
	rCopy := CopyHttpRequest(c.Request)
	rBody, _ := io.ReadAll(rCopy.Body)
	log.Printf("body: %s", rBody)
	b, _ := c.GetRawData()
	log.Printf("body: %s", b)
	c.String(http.StatusOK, "hello world")
}

func (h *httpServer) StartWithPort(port int) error {

	return h.engine.Run(fmt.Sprintf(":%d", port))
}

func CopyHttpRequest(r *http.Request) *http.Request {
	reqCopy := new(http.Request)

	if r == nil {
		return reqCopy
	}

	*reqCopy = *r

	if r.Body != nil {
		defer r.Body.Close()

		// Buffer body data
		var bodyBuffer bytes.Buffer
		bodyBuffer2 := new(bytes.Buffer)

		io.Copy(&bodyBuffer, r.Body)
		*bodyBuffer2 = bodyBuffer

		// Create new ReadClosers so we can split output
		r.Body = ioutil.NopCloser(&bodyBuffer)
		reqCopy.Body = ioutil.NopCloser(bodyBuffer2)
	}

	return reqCopy
}
