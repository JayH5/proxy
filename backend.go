package proxy

import (
	"fmt"
	"net/http"
	"net/url"
	"sync"
)

type Backend interface {
	GetServer(req *http.Request) *url.URL
}

type SingleServer struct {
	server *url.URL
}

func (s *SingleServer) GetServer(req *http.Request) *url.URL {
	return s.server
}

type RoundRobin struct {
	// List of servers
	servers []*server

	// index and currentWeight are protected by the mutex
	index         int
	currentWeight int
	mutex         *sync.Mutex

	// Pre-calculated values
	maxWeight int
	weightGcd int
}

func NewRoundRobin(serverMap map[*url.URL]int) (*RoundRobin, error) {
	servers, maxWeight, weightGcd, err := prepareRoundRobinServers(serverMap)
	if err != nil {
		return nil, err
	}

	return &RoundRobin{
		servers: servers,

		index:         -1,
		currentWeight: 0,
		mutex:         &sync.Mutex{},

		maxWeight: maxWeight,
		weightGcd: weightGcd,
	}, nil
}

func prepareRoundRobinServers(serverMap map[*url.URL]int) (servers []*server, maxWeight int, weightGcd int, err error) {
	if len(serverMap) == 0 {
		return nil, -1, -1, fmt.Errorf("Must specify at least one server")
	}

	servers = make([]*server, 0, len(serverMap))
	maxWeight = -1
	weightGcd = -1
	for url, weight := range serverMap {
		if weight < 0 {
			return nil, -1, -1, fmt.Errorf("Weight of server '%s' is less than 0", url.String())
		}

		servers = append(servers, &server{url: url, weight: weight})

		// Find the max weight
		if weight > maxWeight {
			maxWeight = weight
		}

		// Calculate the greatest common divisor
		if weightGcd == -1 {
			weightGcd = weight
		} else {
			weightGcd = gcd(weightGcd, weight)
		}
	}

	if maxWeight == 0 {
		return nil, -1, -1, fmt.Errorf("All servers have 0 weight")
	}

	return
}

type server struct {
	url    *url.URL
	weight int
}

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func (r *RoundRobin) GetServer(req *http.Request) *url.URL {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	// http://kb.linuxvirtualserver.org/wiki/Weighted_Round-Robin_Scheduling
	for {
		r.index = (r.index + 1) % len(r.servers)
		if r.index == 0 {
			r.currentWeight = r.currentWeight - r.weightGcd
			if r.currentWeight <= 0 {
				r.currentWeight = r.maxWeight
			}
		}
		srv := r.servers[r.index]
		if srv.weight >= r.currentWeight {
			return srv.url
		}
	}
}
