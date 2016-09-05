package backend

import (
	"net/http"
	"net/url"
)

// Defines a simple backend that maps requests to the addresses for the servers
// that the request should be proxied to. Backends in general can have state but
// should *not* be mutable. Any changes to the servers in the backend should
// result in a new backend instance being created. Thus, backends should depend
// on only consistent distribution methods.
type Backend interface {
	// GetServer should return the URL for upstream server that should be
	// delivered the given request. This method must be safe for concurrent
	// access. The given request should not be mutated in any way.
	GetServer(req *http.Request) *url.URL
}
