package main

import (
	"errors"
	"fmt"
	"net/url"
	"testing"
)

type normalizeTest struct {
	arg       string
	expected1 *url.URL
	expected2 error
}

var normalizeTests = []normalizeTest{
	{"[[[]]]", nil, errors.New("")},
	{"http://google.ca:80", &url.URL{
		Scheme: "http",
		Host:   "google.ca",
	}, nil},
	{"https://google.ca:443", &url.URL{
		Scheme: "https",
		Host:   "google.ca",
	}, nil},
	{"ftp://google.ca:21", &url.URL{
		Scheme: "ftp",
		Host:   "google.ca",
	}, nil},
	{"http://gooGlE.cA", &url.URL{
		Scheme: "http",
		Host:   "google.ca",
	}, nil},
	{"http://google.ca/", &url.URL{
		Scheme: "http",
		Host:   "google.ca",
	}, nil},
	{"http://google.ca//search", &url.URL{
		Scheme: "http",
		Host:   "google.ca",
		Path:   "/search",
	}, nil},
	{"http://google.ca/search//results", &url.URL{
		Scheme: "http",
		Host:   "google.ca",
		Path:   "/search/results",
	}, nil},
	{"http://google.ca/search/results//", &url.URL{
		Scheme: "http",
		Host:   "google.ca",
		Path:   "/search/results",
	}, nil},
	{"http://google.ca/search//results/", &url.URL{
		Scheme: "http",
		Host:   "google.ca",
		Path:   "/search/results",
	}, nil},
}

func TestNormalize(t *testing.T) {
	for _, test := range normalizeTests {
		testName := fmt.Sprintf("normalize(%q)", test.arg)
		t.Run(testName, func(t *testing.T) {
			output1, output2 := normalize(test.arg)
			if (output1 == nil) != (test.expected1 == nil) ||
				(output1 != nil && test.expected1 != nil && *output1 != *test.expected1) ||
				(output2 == nil) != (test.expected2 == nil) {
				t.Errorf("got (%+v, %v), expected (%+v, %v)", *output1, output2, *test.expected1, test.expected2)
			}
		})
	}
}

type performBasicNormalizationTest struct {
	arg       string
	expected1 string
	expected2 error
}

var performBasicNormalizationTests = []performBasicNormalizationTest{
	{"http://google.ca", "http://google.ca", nil},
	{" http://google.ca ", "http://google.ca", nil},
	{"google.ca", "http://google.ca", nil},
	{"https://google.ca", "https://google.ca", nil},
	{"http://", "", errors.New("")},
	{"http://?my_query", "", errors.New("")},
	{"http://$", "", errors.New("")},
	{"/my_path", "", errors.New("")},
	{"[[[]]]", "", errors.New("")},
	{"http:www.google.ca", "", errors.New("")},
	{"http:/www.google.ca", "", errors.New("")},
	{"http:///www.google.ca", "", errors.New("")},
	{"://www.google.ca", "", errors.New("")},
	{"//www.google.ca", "", errors.New("")},
}

func TestPerformBasicNormalization(t *testing.T) {
	for _, test := range performBasicNormalizationTests {
		testName := fmt.Sprintf("performBasicNormalization(%q)", test.arg)
		t.Run(testName, func(t *testing.T) {
			output1, output2 := performBasicNormalization(test.arg)
			if output1 != test.expected1 || (output2 == nil) != (test.expected2 == nil) {
				t.Errorf("got (%q, %v), expected (%q, %v)", output1, output2, test.expected1, test.expected2)
			}
		})
	}
}

type getUrlStringFromMapTest struct {
	arg      map[string]string
	expected string
}

var getUrlStringFromMapTests = []getUrlStringFromMapTest{
	{map[string]string{"scheme": "",
		"host":     "google.ca",
		"port":     "",
		"path":     "",
		"query":    "",
		"fragment": ""}, "http://google.ca"},
	{map[string]string{"scheme": "",
		"host":     "www.google.ca",
		"port":     "",
		"path":     "",
		"query":    "",
		"fragment": ""}, "http://www.google.ca"},
	{map[string]string{"scheme": "http",
		"host":     "google.ca",
		"port":     "",
		"path":     "",
		"query":    "",
		"fragment": ""}, "http://google.ca"},
	{map[string]string{"scheme": "https",
		"host":     "google.ca",
		"port":     "",
		"path":     "",
		"query":    "",
		"fragment": ""}, "https://google.ca"},
	{map[string]string{"scheme": "https",
		"host":     "google.ca",
		"port":     "80",
		"path":     "",
		"query":    "",
		"fragment": ""}, "https://google.ca:80"},
	{map[string]string{"scheme": "https",
		"host":     "google.ca",
		"port":     "80",
		"path":     "/search",
		"query":    "",
		"fragment": ""}, "https://google.ca:80/search"},
	{map[string]string{"scheme": "https",
		"host":     "google.ca",
		"port":     "80",
		"path":     "/search",
		"query":    "?q=url%20shortening",
		"fragment": ""}, "https://google.ca:80/search?q=url%20shortening"},
	{map[string]string{"scheme": "https",
		"host":     "google.ca",
		"port":     "80",
		"path":     "/search",
		"query":    "?q=url%20shortening",
		"fragment": "#hello"}, "https://google.ca:80/search?q=url%20shortening#hello"},
}

func TestGetUrlStringFromMap(t *testing.T) {
	for _, test := range getUrlStringFromMapTests {
		testName := fmt.Sprintf("getUrlStringFromMap(%v)", test.arg)
		t.Run(testName, func(t *testing.T) {
			if output := getUrlStringFromMap(test.arg); output != test.expected {
				t.Errorf("got %q, expected %q", output, test.expected)
			}
		})
	}
}