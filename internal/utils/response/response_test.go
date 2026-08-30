package response

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestErrorUsesTopLevelErrorFormat(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/", func(c *gin.Context) {
		Error(c, http.StatusBadRequest, "title is required and must be 1-255 characters")
	})

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/", nil))

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusBadRequest)
	}
	const want = `{"success":false,"statusCode":400,"error":"title is required and must be 1-255 characters"}`
	if got := recorder.Body.String(); got != want {
		t.Errorf("response = %s, want %s", got, want)
	}
}
