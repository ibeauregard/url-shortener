package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

type getNormalizedUrlTest struct {
	requestBodyContent map[string]any
	expectedUrl        *url.URL
	expectedError      error
	expectedHttpStatus int
}

var getNormalizedUrlTests = []getNormalizedUrlTest{
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

func TestGetNormalizedUrl(t *testing.T) {
	for _, test := range getNormalizedUrlTests {
		testName := fmt.Sprintf("getNormalizedUrl with %v as request body", test.requestBodyContent)
		t.Run(testName, func(t *testing.T) {
			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)
			ctx.Request = &http.Request{
				Header: make(http.Header),
			}
			MockJsonPost(ctx, test.requestBodyContent)
			urlOutput, err := getNormalizedUrl(ctx)
			assert.EqualValues(t, test.expectedUrl, urlOutput)
			assert.Condition(t, func() bool {
				return (err == nil) == (test.expectedError == nil)
			})
			assert.EqualValues(t, test.expectedHttpStatus, w.Code)
		})
	}
}

type getShortUrlTest struct {
	arg      string
	expected string
}

var getShortUrlTests = []getShortUrlTest{
	{"my_key", AppScheme + "://" + AppHost + "/my_key"},
}

func TestGetShortUrl(t *testing.T) {
	for _, test := range getShortUrlTests {
		testName := fmt.Sprintf("getShortUrl(%q)", test.arg)
		t.Run(testName, func(t *testing.T) {
			assert.EqualValues(t, test.expected, getShortUrl(test.arg))
		})
	}
}

type getSuccessResponseBodyTest struct {
	arg1     string
	arg2     string
	expected gin.H
}

var getSuccessResponseBodyTests = []getSuccessResponseBodyTest{
	{arg1: "Foo", arg2: "Bar", expected: gin.H{"longUrl": "Foo", "shortUrl": "Bar"}},
}

func TestGetSuccessResponseBody(t *testing.T) {
	for _, test := range getSuccessResponseBodyTests {
		testName := fmt.Sprintf("getSuccessResponseBody(%q, %q)", test.arg1, test.arg2)
		t.Run(testName, func(t *testing.T) {
			assert.EqualValues(t, test.expected, getSuccessResponseBody(test.arg1, test.arg2))
		})
	}
}

func MockJsonPost(c *gin.Context, content any) {
	c.Request.Method = "POST"
	c.Request.Header.Set("Content-Type", "application/json")

	jsonBytes, err := json.Marshal(content)
	if err != nil {
		panic(err)
	}
	c.Request.Body = io.NopCloser(bytes.NewBuffer(jsonBytes))
}
