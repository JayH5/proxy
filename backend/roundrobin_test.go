package backend

import (
	"net/http/httptest"
	"net/url"
	"testing"

	. "gopkg.in/check.v1"
)

func TestRoundRobin(t *testing.T) { TestingT(t) }

type RoundRobinSuite struct{}

var _ = Suite(&RoundRobinSuite{})

// When a RoundRobin backend is created with a single server and GetServer is
// called, the server's URL should be returned.
func (s *RoundRobinSuite) TestSingleServer(c *C) {
	expectedUrl, _ := url.ParseRequestURI("http://127.0.0.1:8080")
	serverMap := map[*url.URL]int{expectedUrl: 1}

	rr, err := NewRoundRobin(serverMap)
	c.Assert(err, IsNil)

	// Check we get our URL back when we ask for a server
	request := httptest.NewRequest("GET", "/test", nil)
	actualUrl := rr.GetServer(request)
	c.Assert(actualUrl, Equals, expectedUrl)
}

// When we try to create a RoundRobin backend with no servers, an error is
// returned.
func (s *RoundRobinSuite) TestEmptyServerMaps(c *C) {
	var emptyServerMap map[*url.URL]int

	rr, err := NewRoundRobin(emptyServerMap)
	c.Assert(rr, IsNil)

	c.Assert(err, ErrorMatches, "Must specify at least one server")
}

// When we try to create a RoundRobin backend with a server with a negative
// weight, an error it returned.
func (s *RoundRobinSuite) TestNegativeWeight(c *C) {
	expectedUrl, _ := url.ParseRequestURI("http://127.0.0.1:8080")
	serverMap := map[*url.URL]int{expectedUrl: -1}

	rr, err := NewRoundRobin(serverMap)
	c.Assert(rr, IsNil)

	c.Assert(err, ErrorMatches, "Weight of server .* is less than 0")
}

// When we try to create a RoundRobin backend with all servers having a weight
// of 0, an error is returned.
func (s *RoundRobinSuite) TestAllZeroWeight(c *C) {
	url1, _ := url.ParseRequestURI("http://127.0.0.1:8080")
	url2, _ := url.ParseRequestURI("http://127.0.0.1:9090")
	serverMap := map[*url.URL]int{
		url1: 0,
		url2: 0,
	}

	rr, err := NewRoundRobin(serverMap)
	c.Assert(rr, IsNil)

	c.Assert(err, ErrorMatches, "All servers have 0 weight")
}

// When a RoundRobin backend is created with two servers with equal weight,
// consecutive calls to GetServer should toggle between the two servers.
func (s *RoundRobinSuite) TestCycleNodes(c *C) {
	url1, _ := url.ParseRequestURI("http://127.0.0.1:8080")
	url2, _ := url.ParseRequestURI("http://127.0.0.1:9090")
	serverMap := map[*url.URL]int{
		url1: 1,
		url2: 1,
	}
	urls := []*url.URL{url1, url2}

	rr, err := NewRoundRobin(serverMap)
	c.Assert(err, IsNil)

	// Get the first request to check which URL came first
	request := httptest.NewRequest("GET", "/test", nil)
	actualUrl := rr.GetServer(request)
	var nextIndex int
	if actualUrl == url1 {
		nextIndex = 1
	} else if actualUrl == url2 {
		nextIndex = 0
	} else {
		c.Fatal("Unexpected URL")
	}

	// Request a URL a few times, asserting that we toggle between the two URLs
	for i := 0; i < 5; i++ {
		actualUrl := rr.GetServer(request)
		c.Assert(actualUrl, Equals, urls[(nextIndex+i)%len(urls)])
	}
}

// When a RoundRobin backend is created with two servers -- one with double the
// weight of the other -- then the server with the larger weight should be
// returned twice as often on calls to GetServer.
func (s *RoundRobinSuite) TestWeighted(c *C) {
	url1, _ := url.ParseRequestURI("http://127.0.0.1:8080")
	url2, _ := url.ParseRequestURI("http://127.0.0.1:9090")
	serverMap := map[*url.URL]int{
		url1: 1,
		url2: 2,
	}

	rr, err := NewRoundRobin(serverMap)
	c.Assert(err, IsNil)

	request := httptest.NewRequest("GET", "/test", nil)
	url1Count := 0
	url2Count := 0
	for i := 0; i < 6; i++ {
		actualUrl := rr.GetServer(request)
		if actualUrl == url1 {
			url1Count++
		} else if actualUrl == url2 {
			url2Count++
		} else {
			c.Fatal("Unexpected URL")
		}
	}

	// Assert that url2 occurs twice as often as url1
	c.Assert(url1Count, Equals, 2)
	c.Assert(url2Count, Equals, 4)
}

// When a RoundRobin backend is created with a server with a weight of 0, that
// server's address should never be returned when GetServer is called.
func (s *RoundRobinSuite) TestZeroWeightSkipped(c *C) {
	url1, _ := url.ParseRequestURI("http://127.0.0.1:8080")
	url2, _ := url.ParseRequestURI("http://127.0.0.1:9090")
	serverMap := map[*url.URL]int{
		url1: 1,
		url2: 0,
	}

	rr, err := NewRoundRobin(serverMap)
	c.Assert(err, IsNil)

	// Assert that we never see the 0-weighted server
	request := httptest.NewRequest("GET", "/test", nil)
	for i := 0; i < 5; i++ {
		actualUrl := rr.GetServer(request)
		c.Assert(actualUrl, Equals, url1)
	}
}
