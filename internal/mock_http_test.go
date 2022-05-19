package internal

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"regexp"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/lpxxn/clank/internal/clanklog"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

var testEngine *gin.Engine

func initGin() {
	gin.SetMode(gin.DebugMode)
	testEngine = gin.Default()
	testEngine.NoRoute(func(c *gin.Context) {
		clanklog.Infof("NoRoute: %s", c.Request.RequestURI)
		c.String(http.StatusOK, "NoRoute")
	})
	testEngine.NoMethod(func(c *gin.Context) {
		clanklog.Infof("NoMethod: %s", c.Request.RequestURI)
		c.String(http.StatusOK, "NoMethod")
	})
}

func TestSchema1(t *testing.T) {
	initGin()
	httpDescriptor := &httpServerDescriptor{MethodDescriptor: []*httpMethodDescriptor{
		{
			Name:   "testApi",
			Path:   "/test",
			Method: HTTPPOSTMethod,
			DefaultResponse: `{"code": 0,"message": "success",
				"data": {"name": "Jerry","age": 18}
			}`,
		},
		{
			Name:   "testApi2",
			Path:   "/user/:userID/order/:orderNo",
			Method: HTTPPOSTMethod,
			DefaultResponse: `{
				"code": 0, "message": "success",
				"data": {
					"orderNo": "$param.orderNo",
					"userID": $param.userID,
					"desc": "{{RandString 5 20}}"
				}
			}`,
		},
		{
			Name:   "testApi2",
			Path:   "/user/:userID/createOrder",
			Method: HTTPPOSTMethod,
			DefaultResponse: `{
				"code": 0, "message": "success",
				"data": {
					"orderNo": "OrderNo{{RandString 5 10}}",
					"userID": $param.userID,
					"desc": "{{RandString 5 20}}"
				}
			}`,
			Conditions: ResponseConditionDescriptionList{
				{
					Condition: "$query.userID == 1",
					Response: `{
						"code": 0, "message": "success",
						"data": {
							"orderNo": "OrderNo{{RandString 5 10}}",
							"userID": $query.userID,
							"desc": "query.userID == 1"
						}
					}`,
				},
				{
					Condition: "$param.userID == 1",
					Response: `{
						"code": 0, "message": "success",
						"data": {
							"orderNo": "OrderNo{{RandString 5 10}}",
							"userID": $param.userID,
							"desc": "param.userID == 1"
						}
					}`,
				},
				{
					Condition: "$body.userID == 1 && $query.userID == 2",
					Response: `{
						"code": 0, "message": "success",
						"data": {
							"orderNo": "OrderNo{{RandString 5 10}}",
							"userID": $body.userID,
							"queryUserID": $query.userID,
							"desc": "body.userID == 1&& query.userID == 2"
						}
					}`,
				},
				{
					Condition: "$body.userID == 1",
					Response: `{
						"code": 0, "message": "success",
						"data": {
							"orderNo": "OrderNo{{RandString 5 10}}",
							"userID": $param.userID,
							"desc": "$body.userID == 1""
						}
					}`,
				},
			},
		},
	}}
	assert.Nil(t, httpDescriptor.Validate())

	serv := NewHttpServer(httpDescriptor)
	assert.NotNil(t, serv)
	assert.Nil(t, serv.MethodHandler())
	form := url.Values{"name": {"Jerry"}, "age": {"18"}}
	r, _ := http.NewRequest(HTTPPOSTMethod, "/user/1233/order/13", strings.NewReader(form.Encode()))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	testHTTPResponse(t, serv.engine, r, func(w *httptest.ResponseRecorder) bool {
		t.Log(w.Body.String())
		return true
	})

	form = url.Values{"name": {"Jerry"}, "age": {"18"}, "userID": {"1"}}
	r, _ = http.NewRequest(HTTPPOSTMethod, "/user/1233/createOrder?userID=2", strings.NewReader(form.Encode()))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	testHTTPResponse(t, serv.engine, r, func(w *httptest.ResponseRecorder) bool {
		t.Log(w.Body.String())
		return true
	})

	r, _ = http.NewRequest(HTTPPOSTMethod, "/user/1233/createOrder?userID=1", nil)
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	testHTTPResponse(t, serv.engine, r, func(w *httptest.ResponseRecorder) bool {
		t.Log(w.Body.String())
		return true
	})

	r, _ = http.NewRequest(HTTPPOSTMethod, "/user/1233/createOrder?userID=2", strings.NewReader(`{"name": "Jerry", "age": 18, "userID": 1}`))
	r.Header.Set("Content-Type", "application/json")
	testHTTPResponse(t, serv.engine, r, func(w *httptest.ResponseRecorder) bool {
		t.Log(w.Body.String())
		return true
	})
}

func TestParam(t *testing.T) {
	initGin()

	orderPath1 := "/restaurant/:id/order/:orderNo"
	testEngine.Any(orderPath1, func(c *gin.Context) {
		t.Log(c.Params)
		t.Logf("full path: %s path: %s rawPath: %s", c.FullPath(), c.Request.URL.Path, c.Request.URL.RawPath)
		body, err := c.GetRawData()
		if err != nil {
			t.Error(err)
		}
		t.Logf("body: %s", string(body))
		t.Logf("query: %s", c.Request.URL.RawQuery)
	})
	testEngine.Handle(HTTPAnyMethod, "/testAny", func(c *gin.Context) {
		c.Writer.WriteString("testAny")
	})
	r, _ := http.NewRequest("GET", orderPath1, strings.NewReader(`{"name": "manu"}`))
	testHTTPResponse(t, testEngine, r, func(w *httptest.ResponseRecorder) bool {
		t.Log(w.Body.String())
		return true
	})

	r = httptest.NewRequest("GET", "/restaurant/1/order/2?a=v1&b=v2", nil)
	testHTTPResponse(t, testEngine, r, func(w *httptest.ResponseRecorder) bool {
		t.Log(w.Body.String())
		return true
	})
	rTestPath := httptest.NewRequest("GET", "/testAny", nil)
	testHTTPResponse(t, testEngine, rTestPath, func(w *httptest.ResponseRecorder) bool {
		t.Log(w.Body.String())
		return true
	})

}

// Helper function to process a request and test its response
func testHTTPResponse(t *testing.T, r *gin.Engine, req *http.Request, f func(w *httptest.ResponseRecorder) bool) {
	// Create a response recorder
	w := httptest.NewRecorder()
	// Create the service and process the above request.
	r.ServeHTTP(w, req)

	if f != nil && !f(w) {
		t.Fail()
	}
}

func TestNoRouter(t *testing.T) {
	initGin()

	path1 := "/restaurant/:id/order/:orderNo"
	path2 := "/restaurant/:id/:action"
	httpServ := &httpServer{
		serverMethod: map[string]string{
			path1: HTTPAnyMethod,
			path2: HTTPPOSTMethod,
		},
		engine: gin.Default(),
	}
	httpServ.MethodHandler()
	r, _ := http.NewRequest(HTTPPOSTMethod, "/restaurant/1/order/2?a=v1&b=v2", strings.NewReader(`{"id": 1, "name": "Tom"}`))
	testHTTPResponse(t, httpServ.engine, r, func(w *httptest.ResponseRecorder) bool {
		t.Log(w.Body.String())
		return true
	})

	r, _ = http.NewRequest(HTTPPOSTMethod, "/restaurant/1/2?a=v1&b=v2", strings.NewReader(`{"id": 1, "name": "Tom"}`))
	testHTTPResponse(t, httpServ.engine, r, func(w *httptest.ResponseRecorder) bool {
		t.Log(w.Body.String())
		return true
	})
}

func TestHttpRegex(t *testing.T) {
	initGin()

	regex := regexp.MustCompile(`\$(?P<parameter>(param|body|query|form)\.\w+[.\w]*)`)

	str := `$body.name=$param.id || "$abcdef.eeeee" = 1334`
	match := regex.FindAllStringSubmatch(str, -1)
	idx := regex.SubexpIndex("parameter")
	for _, matchItem := range match {
		t.Log(matchItem[idx])
	}
}

func TestHttpYAML(t *testing.T) {
	b, err := yaml.Marshal(httpSchema{Server: httpServerDescriptor{
		MethodDescriptor: []*httpMethodDescriptor{
			&httpMethodDescriptor{
				Name:            "abcafe",
				Path:            "/index",
				Method:          "GET",
				DefaultResponse: "hello",
			},
			&httpMethodDescriptor{
				Name:            "bbb",
				Path:            "/bbb",
				Method:          "GET",
				DefaultResponse: "hello",
			},
		},
	}})
	assert.Nil(t, err)
	t.Log(string(b))
}
