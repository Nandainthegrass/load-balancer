package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/gin-gonic/gin"
)

//Let's make a load balancer

type LoadBalancer struct {
	servers []*Server
}

type Server struct {
	address  string
	reqCount int
}

func NewLoadBalancer() *LoadBalancer {
	return &LoadBalancer{
		servers: make([]*Server, 0),
	}
}

func (lb *LoadBalancer) AddServer(address string) {
	lb.servers = append(lb.servers, &Server{address: address})
}

func (lb *LoadBalancer) ChooseServer() *Server {
	if len(lb.servers) == 0 {
		log.Fatal("No server to send anything to:")
	}
	nextServer := lb.servers[0]
	if len(lb.servers) != 1 {
		lb.servers = append(lb.servers[1:], lb.servers[0])
	}
	nextServer.reqCount++
	return nextServer
}

type Proxy struct {
	targetUrl string
}

func NewProxy(targetUrl string) *Proxy {
	return &Proxy{
		targetUrl: targetUrl,
	}
}

func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	targetURL, _ := url.Parse(p.targetUrl)
	r.Header.Set("Connection", "close")
	reverseProxy := httputil.NewSingleHostReverseProxy(targetURL)
	fmt.Printf("request sent to %v", p.targetUrl)
	reverseProxy.ServeHTTP(w, r)
}

func (lb *LoadBalancer) HandleRequest(c *gin.Context) {
	server := lb.ChooseServer()
	if server == nil {
		log.Fatal("No server available!")
		return
	}
	p := NewProxy(server.address)
	p.ServeHTTP(c.Writer, c.Request)
}

func (lb *LoadBalancer) MetricsHandler(c *gin.Context) {
	metrics := make(map[string]int)
	for _, server := range lb.servers {
		metrics[server.address] = server.reqCount
	}
	c.IndentedJSON(http.StatusOK, metrics)
}

func main() {
	fmt.Println("Running server!")
	lb := NewLoadBalancer()
	lb.AddServer("http://127.0.0.1:8000")
	lb.AddServer("http://127.0.0.1:8001")
	lb.AddServer("http://127.0.0.1:8002")
	r := gin.Default()
	r.GET("/", func(ctx *gin.Context) {
		lb.HandleRequest(ctx)
	})
	r.GET("/metrics", func(ctx *gin.Context) {
		lb.MetricsHandler(ctx)
	})
	r.Run()
}
