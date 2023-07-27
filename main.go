package main

import (
	"fmt"
	"net/http/httputil"
	"net/url"
	"os"
)

type simpleServer struct {
	addr  string
	proxy *httputil.ReverseProxy
}

/**
 *
 * func to create and initialize a new server
 **/
func newSimpleServer(addr string) *simpleServer {

	//parses a raw url into a URL structure.
	serverUrl, err := url.Parse(addr)
	handleErr(err)

	return &simpleServer{
		addr:  addr,
		proxy: httputil.NewSingleHostReverseProxy(serverUrl),
	}

}

/**
 *
 * func to handle errors
 **/
func handleErr(err error) {
	if err != nil {
		fmt.Printf("error : %v\n", err)
		os.Exit(1)
	}
}

/**
 *
 *
 **/
func main() {

}
