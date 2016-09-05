package backend

import (
	"net/http"
	"net/url"
)

// The simplest possible backend. Just a single server address for all requests.
// Super fast.
type SingleServer struct {
	server *url.URL
}

func (s *SingleServer) GetServer(req *http.Request) *url.URL {
	return s.server
}
