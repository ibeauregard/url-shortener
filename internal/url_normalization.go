// Adapted from https://github.com/sekimura/go-normalize-url/blob/master/normalizeurl.go
// The original author decided to always strip the leading www, which I believe is not advisable

package main

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"golang.org/x/net/idna"
)

var (
	DefaultPorts = map[string]int{
		"http":  80,
		"https": 443,
		"ftp":   21,
	}
)

func normalize(s string) (*url.URL, error) {
	s, err := performBasicNormalization(s)
	if err != nil {
		return nil, err
	}
	u, err := url.ParseRequestURI(s)
	if err != nil {
		return nil, err
	}

	p, ok := DefaultPorts[u.Scheme]
	if ok {
		u.Host = strings.TrimSuffix(u.Host, fmt.Sprintf(":%d", p))
	}

	got, err := idna.ToUnicode(u.Host)
	if err != nil {
		return nil, err
	} else {
		u.Host = strings.ToLower(got)
	}

	multipleSlashRegex := regexp.MustCompile(`//+`)
	u.Path = multipleSlashRegex.ReplaceAllString(u.Path, "/")
	u.Path = strings.TrimSuffix(u.Path, "/")

	v := u.Query()
	u.RawQuery = v.Encode()
	u.RawQuery, _ = url.QueryUnescape(u.RawQuery)

	return u, nil
}

var urlRegex = regexp.MustCompile(`^((?P<scheme>[\w-.]+)://)?(?P<host>[a-zA-Z\d-]+(\.[a-zA-Z\d-]+)*)(:(?P<port>\d{1,5}))?(?P<path>/([\w.\-~!$&'()*+,;=:@/]|%[\da-fA-F]{2})*)?(?P<query>\?(&?[^=&#]*=[^=&#]*)*)?(?P<fragment>#([\w?/:@\-.~!$&'()*+,;=]|%[\da-fA-F]{2})*)?$`)

func performBasicNormalization(s string) (string, error) {
	s = strings.TrimSpace(s)
	match := urlRegex.FindStringSubmatch(s)
	if match == nil {
		return "", fmt.Errorf("unable to match against URL regex")
	}
	result := make(map[string]string)
	for i, name := range urlRegex.SubexpNames() {
		if i != 0 && name != "" {
			result[name] = match[i]
		}
	}
	return getUrlStringFromMap(result), nil
}

func getUrlStringFromMap(urlMap map[string]string) string {
	var stringBuilder strings.Builder
	if urlMap["scheme"] != "" {
		stringBuilder.WriteString(strings.ToLower(urlMap["scheme"]))
	} else {
		stringBuilder.WriteString("http")
	}
	stringBuilder.WriteString("://")
	stringBuilder.WriteString(strings.ToLower(urlMap["host"]))
	if urlMap["port"] != "" {
		stringBuilder.WriteString(":")
	}
	stringBuilder.WriteString(urlMap["port"])
	stringBuilder.WriteString(urlMap["path"])
	stringBuilder.WriteString(urlMap["query"])
	stringBuilder.WriteString(urlMap["fragment"])
	return stringBuilder.String()
}
