package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

const (
	BasePath = "https://us-east-2.console.aws.amazon.com/s3/buckets/hoister/__outputs/"
	Port     = "8000"
)

func main() {
	http.HandleFunc("/", handler)
	log.Printf("Reverse Proxy Running on Port %s", Port)
	log.Fatal(http.ListenAndServe(":"+Port, nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	hostname := r.Host
	subdomain := strings.Split(hostname, ".")[0]

	// Custom Domain - DB Query

	resolvesTo := BasePath + "/" + subdomain

	target, _ := url.Parse(resolvesTo)
	proxy := httputil.NewSingleHostReverseProxy(target)

	r.URL.Path = modifyPath(r.URL.Path)
	proxy.ServeHTTP(w, r)
}

func modifyPath(path string) string {
	if path == "/" {
		return "/index.html"
	}
	return path
}
