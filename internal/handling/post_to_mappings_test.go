package handling

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestHandlePostToMappings(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/test", strings.NewReader("{}"))
	r := gin.Default()
	r.POST("/test", HandlePostToMappings(&repoAdapter{}))
	r.ServeHTTP(w, req)
	assert.EqualValues(t, http.StatusBadRequest, w.Code)
}

func TestHandlePostToMappingsBadUserInput(t *testing.T) {
	badInputs := []map[string]any{
		{},
		{"foo": "bar"},
		{"longUrl": ""},
		{"longUrl": ";[*"},
	}
	for _, badInput := range badInputs {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		MockJsonPost(ctx, badInput)
		t.Run(fmt.Sprintf("POST %v", badInput), func(t *testing.T) {
			(&concretePostHandler{ctx}).handle(&repoAdapter{})
			assert.Condition(t, func() bool { return 400 <= w.Code && w.Code < 500 })
		})
	}
}

func TestHandlePostToMappingsBlacklistedDomain(t *testing.T) {
	AppHost = "apphost"
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	MockJsonPost(ctx, map[string]any{"longUrl": AppHost})
	(&concretePostHandler{ctx}).handle(&repoAdapter{})
	assert.EqualValues(t, http.StatusUnprocessableEntity, w.Code)
}

func (m *repoAdapterMock) getKey(_ string) (string, bool) {
	return m.outputStr, m.outputFoundStatus
}

func TestHandlePostToMappingsLongUrlFound(t *testing.T) {
	longUrl := "http://foobar.com"
	key := "my_key"
	AppHost = "apphost"
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	MockJsonPost(ctx, map[string]any{"longUrl": longUrl})
	(&concretePostHandler{ctx}).handle(&repoAdapterMock{
		outputStr:         key,
		outputFoundStatus: true,
	})
	assert.EqualValues(t, http.StatusOK, w.Code)
	validateResponseBody(t, w, longUrl, key)
}

func (m *repoAdapterMock) addMapping(_ string) (string, error) {
	return m.outputStr, m.outputError
}

func TestHandlePostToMappingsErrorWhileAdding(t *testing.T) {
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	MockJsonPost(ctx, map[string]any{"longUrl": "http://foobar.com"})
	(&concretePostHandler{ctx}).handle(&repoAdapterMock{
		outputFoundStatus: false,
		outputError:       errors.New(""),
	})
	assert.EqualValues(t, http.StatusInternalServerError, w.Code)
}

func TestHandlePostToMappingsAddSuccess(t *testing.T) {
	longUrl := "http://foobar.com"
	key := "my_key"
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	MockJsonPost(ctx, map[string]any{"longUrl": longUrl})
	(&concretePostHandler{ctx}).handle(&repoAdapterMock{
		outputFoundStatus: false,
		outputStr:         key,
		outputError:       nil,
	})
	assert.EqualValues(t, http.StatusCreated, w.Code)
	validateResponseBody(t, w, longUrl, key)
}

func TestGetNormalizedUrl(t *testing.T) {
	getNormalizedUrlTests := []struct {
		requestBodyContent map[string]any
		expectedUrl        *url.URL
		expectedError      error
		expectedHttpStatus int
	}{
		{
			requestBodyContent: map[string]any{},
			expectedUrl:        nil,
			expectedError:      errors.New(""),
			expectedHttpStatus: http.StatusBadRequest,
		},
		{
			requestBodyContent: map[string]any{"foo": "bar"},
			expectedUrl:        nil,
			expectedError:      errors.New(""),
			expectedHttpStatus: http.StatusBadRequest,
		},
		{
			requestBodyContent: map[string]any{"longUrl": ""},
			expectedUrl:        nil,
			expectedError:      errors.New(""),
			expectedHttpStatus: http.StatusBadRequest,
		},
		{
			requestBodyContent: map[string]any{"longUrl": ";[*"},
			expectedUrl:        nil,
			expectedError:      errors.New(""),
			expectedHttpStatus: http.StatusUnprocessableEntity,
		},
		{
			requestBodyContent: map[string]any{"longUrl": "http://gooGlE.ca:80/search//results/"},
			expectedUrl: &url.URL{
				Scheme: "http",
				Host:   "google.ca",
				Path:   "/search/results",
			},
			expectedError:      nil,
			expectedHttpStatus: http.StatusOK,
		},
	}

	for _, test := range getNormalizedUrlTests {
		testName := fmt.Sprintf("getNormalizedUrl with %v as request body", test.requestBodyContent)
		t.Run(testName, func(t *testing.T) {
			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)
			MockJsonPost(ctx, test.requestBodyContent)
			urlOutput, err := (&concretePostHandler{ctx}).getNormalizedUrl()
			assert.EqualValues(t, test.expectedUrl, urlOutput)
			assert.Condition(t, func() bool {
				return (err == nil) == (test.expectedError == nil)
			})
			assert.EqualValues(t, test.expectedHttpStatus, w.Code)
		})
	}
}

func TestGetSuccessResponseBody(t *testing.T) {
	getSuccessResponseBodyTests := []struct {
		arg1     string
		arg2     string
		expected gin.H
	}{
		{arg1: "Foo", arg2: "Bar", expected: gin.H{"longUrl": "Foo", "shortUrl": "Bar"}},
	}

	for _, test := range getSuccessResponseBodyTests {
		testName := fmt.Sprintf("getSuccessResponseBody(%q, %q)", test.arg1, test.arg2)
		t.Run(testName, func(t *testing.T) {
			assert.EqualValues(t, test.expected, getSuccessResponseBody(test.arg1, test.arg2))
		})
	}
}

func TestGetShortUrl(t *testing.T) {
	AppHost = "apphost"
	getShortUrlTests := []struct {
		arg      string
		expected string
	}{
		{"my_key", AppScheme + "://" + AppHost + "/my_key"},
	}

	for _, test := range getShortUrlTests {
		testName := fmt.Sprintf("getShortUrl(%q)", test.arg)
		t.Run(testName, func(t *testing.T) {
			assert.EqualValues(t, test.expected, getShortUrl(test.arg))
		})
	}
}

func MockJsonPost(c *gin.Context, content any) {
	jsonBytes, err := json.Marshal(content)
	if err != nil {
		panic(err)
	}
	c.Request = &http.Request{
		Method: "POST",
		Header: http.Header{"Content-Type": []string{"application/json"}},
		Body:   io.NopCloser(bytes.NewBuffer(jsonBytes)),
	}
}

func validateResponseBody(t *testing.T, w *httptest.ResponseRecorder, longUrl string, key string) {
	type responseBody struct {
		LongUrl  string
		ShortUrl string
	}
	expectedResponseBody := responseBody{longUrl, getShortUrl(key)}
	rawResponseBody, _ := ioutil.ReadAll(w.Result().Body)
	var parsedResponseBody responseBody
	_ = json.Unmarshal(rawResponseBody, &parsedResponseBody)
	assert.EqualValues(t, expectedResponseBody, parsedResponseBody)
}
