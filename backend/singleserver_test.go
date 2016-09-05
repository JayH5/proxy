package backend

import (
	"net/http/httptest"
	"net/url"
	"testing"

	. "gopkg.in/check.v1"
)

func TestSingleServer(t *testing.T) { TestingT(t) }

type SingleServerSuite struct{}

var _ = Suite(&SingleServerSuite{})

func (s *SingleServerSuite) TestSingleServer(c *C) {
	expectedUrl, _ := url.ParseRequestURI("http://127.0.0.1:8080")

	ss := &SingleServer{server: expectedUrl}

	// Check we get our URL back when we ask for a server
	request := httptest.NewRequest("GET", "/test", nil)
	actualUrl := ss.GetServer(request)
	c.Assert(actualUrl, Equals, expectedUrl)
}
