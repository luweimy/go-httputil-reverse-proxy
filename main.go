package main

import (
	"log"
	"math/rand"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
	"time"
)

var ErrLogger = log.New(os.Stderr, "[ERROR]", log.LstdFlags|log.Lshortfile)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func singleJoiningSlash(a, b string) string {
	aslash := strings.HasSuffix(a, "/")
	bslash := strings.HasPrefix(b, "/")
	switch {
	case aslash && bslash:
		return a + b[1:]
	case !aslash && !bslash:
		return a + "/" + b
	}
	return a + b
}

func NewMultipleHostsReverseProxy(targets []*url.URL) *httputil.ReverseProxy {
	director := func(req *http.Request) {
		target := targets[rand.Int()%len(targets)]
		targetQuery := target.RawQuery

		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		req.URL.Path = singleJoiningSlash(target.Path, req.URL.Path)
		if targetQuery == "" || req.URL.RawQuery == "" {
			req.URL.RawQuery = targetQuery + req.URL.RawQuery
		} else {
			req.URL.RawQuery = targetQuery + "&" + req.URL.RawQuery
		}
	}
	proxy := &httputil.ReverseProxy{Director: director, ErrorLog: ErrLogger}
	return proxy
}

func main() {
	conf := []*url.URL{
		{
			Scheme: "http",
			Host:   "10.102.2.233:7010",
		},
		{
			Scheme: "http",
			Host:   "10.102.1.64:7012",
		},
	}
	proxy := NewMultipleHostsReverseProxy(conf)
	log.Fatal(http.ListenAndServe(":9999", proxy))
}
