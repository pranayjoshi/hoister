package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

const (
	PORT      = "8000"
	BASE_PATH = "https://hoister.s3.us-east-2.amazonaws.com/__outputs/"
)

func main() {
	http.HandleFunc("/", handleRequest)
	log.Printf("Reverse Proxy Running on port %s", PORT)
	log.Fatal(http.ListenAndServe(":"+PORT, nil))
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	hostname := r.Host
	subdomain := strings.Split(hostname, ".")[0]

	resolvesTo := fmt.Sprintf("%s%s", BASE_PATH, subdomain)
	fmt.Println("Resolves to: ", resolvesTo)
	target, err := url.Parse(resolvesTo)
	if err != nil {
		log.Printf("Error parsing target URL: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	proxy := httputil.NewSingleHostReverseProxy(target)
	proxy.ModifyResponse = func(resp *http.Response) error {
		resp.Header.Set("Content-Disposition", "inline")
		contentType := "application/octet-stream"
		if strings.HasSuffix(resp.Request.URL.Path, ".svg") {
			contentType = "image/svg+xml"
		} else if strings.HasSuffix(resp.Request.URL.Path, ".html") {
			contentType = "text/html"
		} else if strings.HasSuffix(resp.Request.URL.Path, ".js") {
			contentType = "application/javascript"
		} else if strings.HasSuffix(resp.Request.URL.Path, ".css") {
			contentType = "text/css"
		}
		resp.Header.Set("Content-Type", contentType)
		return nil
	}
	proxy.Director = func(req *http.Request) {
		req.URL.Host = target.Host
		req.URL.Scheme = target.Scheme
		req.Host = target.Host

		// Check if the path is empty or "/"
		if req.URL.Path == "" || req.URL.Path == "/" {
			req.URL.Path = "/index.html"
		}

		originalPath := strings.TrimPrefix(req.URL.Path, "/")
		if originalPath == "favicon.ico" {
			http.NotFound(w, req)
			return
		}
		req.URL.Path = singleJoiningSlash(target.Path, originalPath)
		fmt.Println("Proxying to: ", req.URL)
	}
	fmt.Println("Proxying to: ", target)
	proxy.ServeHTTP(w, r)
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
