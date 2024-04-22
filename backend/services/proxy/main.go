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

	// Custom Domain - DB Query
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
		return nil
	}
	proxy.ModifyResponse = func(resp *http.Response) error {
		if strings.HasSuffix(resp.Request.URL.Path, ".html") {
			resp.Header.Set("Content-Type", "text/html")
		}
		return nil
	}
	proxy.Director = func(req *http.Request) {
		req.URL.Host = target.Host
		req.URL.Scheme = target.Scheme
		req.Host = target.Host
		// req.URL.Path += ""
		req.URL.Path = strings.TrimPrefix(req.URL.Path, "/")
		req.URL.Path = "/__outputs/" + subdomain + "/index.html"

		fmt.Println("Proxyaaaing to: ", req.URL)
	}
	fmt.Println("Proxying to: ", target)
	proxy.ServeHTTP(w, r)
}
