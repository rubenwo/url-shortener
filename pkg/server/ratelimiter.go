package server

import (
	"golang.org/x/time/rate"
	"log"
	"sync"
	"time"
)

const (
	cleanupTimeout = time.Minute * 5
	sessionTimeout = 30 * 60
)

type IPRateLimiter struct {
	ips        map[string]*rate.Limiter
	lastActive map[string]int64
	sync.RWMutex
	r rate.Limit
	b int
}

// NewIPRateLimiter creates an ip rate limiter.
// r is the rate at which the bucket fills per second.
// b is the size of the bucket. Every request depletes the bucket by 1.
func NewIPRateLimiter(r rate.Limit, b int) *IPRateLimiter {
	l := &IPRateLimiter{
		ips:        make(map[string]*rate.Limiter),
		lastActive: make(map[string]int64),
		r:          r,
		b:          b,
	}
	go l.backgroundCleanup()
	return l
}

// AddIP creates a new rate limiter and adds it to the ips map,
// using the IP address as the key
func (i *IPRateLimiter) AddIp(ip string) *rate.Limiter {
	i.Lock()
	defer i.Unlock()

	limiter := rate.NewLimiter(i.r, i.b)
	i.ips[ip] = limiter
	i.lastActive[ip] = time.Now().Unix()

	return limiter
}

// GetLimiter returns the rate limiter for the provided IP address if it exists.
// Otherwise calls AddIP to add IP address to the map
func (i *IPRateLimiter) GetLimiter(ip string) *rate.Limiter {
	i.RLock()
	limiter, exists := i.ips[ip]

	if !exists {
		i.RUnlock()
		return i.AddIp(ip)
	}
	i.lastActive[ip] = time.Now().Unix()
	i.RUnlock()
	return limiter
}

// backgroundCleanup will cleanup the ips map every 5 minutes.
func (i *IPRateLimiter) backgroundCleanup() {
	timeout := time.Tick(cleanupTimeout)

	for t := range timeout {
		i.Lock()
		log.Printf("IPRateLimiter => remove old connection entries => len(ips) = %d\n", len(i.ips))
		removed := 0
		for k, v := range i.lastActive {
			if t.Unix()-v >= sessionTimeout {
				delete(i.ips, k)
				delete(i.lastActive, k)
				removed++
			}
		}
		log.Printf("IPRateLimiter => removed %d old connections => len(ips) = %d\n", removed, len(i.ips))

		i.Unlock()
	}
}
