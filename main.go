package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

type Server interface {
	Address() string
	IsAlive() bool
	Serve(rw http.ResponseWriter, r *http.Request)
}
type simpleServer struct {
	addr  string
	proxy *httputil.ReverseProxy
}

type LoadBalancer struct {
	port            string
	roundRobinCount int
	servers         []Server
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
 * func to create and initialize a new Load Balancer
 **/
func NewLoadBalancer(port string, servers []Server) *LoadBalancer {
	return &LoadBalancer{
		port:            port,
		roundRobinCount: 0,
		servers:         servers,
	}
}

/**
 *
 * func to return simple server address
 **/
func (s *simpleServer) Address() string {
	return s.addr
}

/**
 *
 * func to check if  server is alive
 **/
func (s *simpleServer) IsAlive() bool {
	return true
}

/**
 *
 * func to serve server
 **/
func (s *simpleServer) Serve(rw http.ResponseWriter, req *http.Request) {
	s.proxy.ServeHTTP(rw, req)
}

/**
 *
 * func to get next available server
 **/
func (lb *LoadBalancer) getNextAvailableServer() Server {
	server := lb.servers[lb.roundRobinCount%len(lb.servers)]
	for !server.IsAlive() {
		lb.roundRobinCount++
		server = lb.servers[lb.roundRobinCount%len(lb.servers)]
	}
	lb.roundRobinCount++
	return server
}

/**
 *
 * func to serve the next available server throught proxy
 **/
func (lb *LoadBalancer) serveProxy(rw http.ResponseWriter, r *http.Request) {
	targetServer := lb.getNextAvailableServer()
	fmt.Printf("forwading request to address %q\n", targetServer.Address())
	targetServer.Serve(rw, r)
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

	servers := []Server{
		newSimpleServer("http://www.facebook.com"),
		newSimpleServer("http://www.bing.com"),
		newSimpleServer("http://www.google.com"),
	}

	lb := NewLoadBalancer("8000", servers)

	handleRedirect := func(rw http.ResponseWriter, req *http.Request) {
		lb.serveProxy(rw, req)
	}
	http.HandleFunc("/", handleRedirect)

	fmt.Printf("serving request at 'localhost:%s'\n", lb.port)
	http.ListenAndServe(":"+lb.port, nil)
}
