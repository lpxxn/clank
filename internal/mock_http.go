package internal

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tidwall/sjson"
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
	getResponse  func(methodName string, jBody string) (string, error)
}

func NewHttpServer(desc *httpServerDescriptor) *httpServer {
	rev := &httpServer{engine: gin.Default(), serverMethod: map[string]string{},
		getResponse: desc.GetResponse}
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
			h.engine.Any(path, metadataHandler, h.commonHandler)
		} else {
			h.engine.Handle(method, path, metadataHandler, h.commonHandler)
		}
	}
	return nil
}

func metadataHandler(c *gin.Context) {
	c.Request.ParseForm()
	jBody := ``
	var err error
	param := map[string]interface{}{}
	for _, item := range c.Params {
		param[item.Key] = item.Value
	}

	query := map[string]interface{}{}
	for key, value := range c.Request.URL.Query() {
		query[key] = value[0]
	}
	body, _ := io.ReadAll(CopyHttpRequest(c.Request).Body)
	form := map[string]interface{}{}
	for k, v := range c.Request.Form {
		form[k] = v[0]
	}
	if jBody, err = sjson.Set(jBody, "param", param); err != nil {
		c.Writer.WriteString(err.Error())
	}
	if jBody, err = sjson.Set(jBody, "query", query); err != nil {
		c.Writer.WriteString(err.Error())
	}
	if jBody, err = sjson.Set(jBody, "body", string(body)); err != nil {
		c.Writer.WriteString(err.Error())
	}
	if jBody, err = sjson.Set(jBody, "form", form); err != nil {
		c.Writer.WriteString(err.Error())
	}
	log.Printf("jBody: %s", jBody)
	c.Set("metadata", jBody)
	c.Next()
}
func MustGetJBody(c *gin.Context) string {
	jBody, ok := c.Get("metadata")
	if !ok {
		return ``
	}
	return jBody.(string)
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
	resp, err := h.getResponse(c.FullPath(), MustGetJBody(c))
	if err != nil {
		c.String(http.StatusExpectationFailed, err.Error())
		return
	}
	c.String(http.StatusOK, resp)
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
