package internal

import (
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

var testEngine *gin.Engine

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	testEngine = gin.Default()
	testEngine.NoRoute(func(c *gin.Context) {
		log.Printf("NoRoute: %s", c.Request.RequestURI)
		c.String(http.StatusOK, "NoRoute")
	})
	testEngine.NoMethod(func(c *gin.Context) {
		log.Printf("NoMethod: %s", c.Request.RequestURI)
		c.String(http.StatusOK, "NoMethod")
	})
	os.Exit(m.Run())
}

func TestParam(t *testing.T) {
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
