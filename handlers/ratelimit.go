package handlers

import (
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

// IPRateLimiter stores a rate limiter for each IP address.
type IPRateLimiter struct {
	ips   map[string]*rate.Limiter
	mu    *sync.RWMutex
	rate  rate.Limit
	burst int
}

// NewIPRateLimiter creates a new IP-based rate limiter.
// It allows 'r' requests per second with a maximum burst of 'b' requests.
func NewIPRateLimiter(r rate.Limit, b int) *IPRateLimiter {
	limiter := &IPRateLimiter{
		ips:   make(map[string]*rate.Limiter),
		mu:    &sync.RWMutex{},
		rate:  r,
		burst: b,
	}

	// Periodically clean up old entries from the IP map to prevent memory leaks.
	go limiter.cleanupVisitors()

	return limiter
}

// getVisitorLimiter returns the rate limiter for a given IP address.
func (i *IPRateLimiter) getVisitorLimiter(ip string) *rate.Limiter {
	i.mu.RLock()
	limiter, exists := i.ips[ip]
	i.mu.RUnlock()

	if !exists {
		i.mu.Lock()
		// Double-check after acquiring write lock
		if limiter, exists = i.ips[ip]; !exists {
			limiter = rate.NewLimiter(i.rate, i.burst)
			i.ips[ip] = limiter
		}
		i.mu.Unlock()
	}

	return limiter
}

// cleanupVisitors removes old entries from the ips map.
func (i *IPRateLimiter) cleanupVisitors() {
	for {
		// Wait for a minute before cleaning up.
		time.Sleep(time.Minute)

		i.mu.Lock()
		// In a real-world scenario, you'd track last access time.
		// For this implementation, we'll keep it simple and assume a fixed cleanup interval is sufficient.
		// A more advanced implementation would use a library like `patrickmn/go-cache`.
		if len(i.ips) > 1000 { // Only clean if the map is getting large
			// A simple strategy: remove a fraction of the entries.
			// This is not perfect but prevents unbounded growth.
			count := 0
			for ip := range i.ips {
				if count > 500 {
					break
				}
				delete(i.ips, ip)
				count++
			}
		}
		i.mu.Unlock()
	}
}

// RateLimitMiddleware creates a middleware that limits requests based on IP address.
func RateLimitMiddleware(next http.Handler) http.Handler {
	// Allow 10 requests per second with a burst of 20.
	// These values should ideally be configurable.
	limiter := NewIPRateLimiter(10, 20)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			// If we can't get the IP, we can choose to block or allow.
			// Allowing is safer to not block legitimate but malformed requests (e.g., behind proxies).
			log.Printf("Could not parse IP address from %s", r.RemoteAddr)
			ip = r.RemoteAddr // Fallback
		}

		// Check if the IP is allowed to make a request.
		if !limiter.getVisitorLimiter(ip).Allow() {
			http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}
