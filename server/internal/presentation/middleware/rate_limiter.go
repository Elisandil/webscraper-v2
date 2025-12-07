package middleware

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

type RateLimiter struct {
	visitors     map[string]*visitor
	userVisitors map[int64]*visitor
	mu           sync.RWMutex
	limit        rate.Limit
	burst        int
	trackUsers   bool
	ctx          context.Context
	cancel       context.CancelFunc
}

type visitor struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

func NewRateLimiter(requestsPerSecond float64, burst int) *RateLimiter {
	return NewConfigurableRateLimiter(requestsPerSecond, burst, false)
}

func NewConfigurableRateLimiter(requestsPerSecond float64, burst int, trackUsers bool) *RateLimiter {
	ctx, cancel := context.WithCancel(context.Background())
	rl := &RateLimiter{
		visitors:     make(map[string]*visitor),
		userVisitors: make(map[int64]*visitor),
		limit:        rate.Limit(requestsPerSecond),
		burst:        burst,
		trackUsers:   trackUsers,
		ctx:          ctx,
		cancel:       cancel,
	}
	go rl.cleanupVisitors()
	return rl
}

func (rl *RateLimiter) Shutdown() {
	rl.cancel()
}

func (rl *RateLimiter) Limit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := getClientIP(r)

		ipLimiter := rl.getVisitor(ip)
		if !ipLimiter.Allow() {
			rl.sendRateLimitError(w)
			return
		}

		if rl.trackUsers {
			user := GetUserFromContext(r.Context())
			if user != nil {
				userLimiter := rl.getUserVisitor(user.ID)
				if !userLimiter.Allow() {
					rl.sendRateLimitError(w)
					return
				}
			}
		}

		next.ServeHTTP(w, r)
	})
}

func (rl *RateLimiter) getVisitor(ip string) *rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	v, exists := rl.visitors[ip]
	if !exists {
		limiter := rate.NewLimiter(rl.limit, rl.burst)
		rl.visitors[ip] = &visitor{limiter, time.Now()}
		return limiter
	}

	v.lastSeen = time.Now()
	return v.limiter
}

func (rl *RateLimiter) getUserVisitor(userID int64) *rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	v, exists := rl.userVisitors[userID]
	if !exists {
		limiter := rate.NewLimiter(rl.limit, rl.burst)
		rl.userVisitors[userID] = &visitor{limiter, time.Now()}
		return limiter
	}

	v.lastSeen = time.Now()
	return v.limiter
}

func (rl *RateLimiter) cleanupVisitors() {
	ticker := time.NewTicker(3 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			rl.mu.Lock()
			for ip, v := range rl.visitors {
				if time.Since(v.lastSeen) > 5*time.Minute {
					delete(rl.visitors, ip)
				}
			}
			for userID, v := range rl.userVisitors {
				if time.Since(v.lastSeen) > 5*time.Minute {
					delete(rl.userVisitors, userID)
				}
			}
			rl.mu.Unlock()
		case <-rl.ctx.Done():
			return
		}
	}
}

func (rl *RateLimiter) sendRateLimitError(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Retry-After", "60")
	w.Header().Set("X-RateLimit-Limit", fmt.Sprintf("%d", rl.burst))
	w.Header().Set("X-RateLimit-Window", "60s")
	w.WriteHeader(http.StatusTooManyRequests)
	sendJSONResponse(w, ErrorResponse{
		Error:   "Rate limit exceeded",
		Message: "Too many requests. Please try again later. Check Retry-After header.",
		Code:    http.StatusTooManyRequests,
	})
}

func getClientIP(r *http.Request) string {
	ip := r.Header.Get("X-Real-IP")
	if ip == "" {
		ip = r.Header.Get("X-Forwarded-For")
	}
	if ip == "" {
		ip = r.RemoteAddr
	}
	return ip
}

func NewStrictRateLimiter() *RateLimiter {
	return NewRateLimiter(5.0/60.0, 5)
}

func NewModerateRateLimiter() *RateLimiter {
	return NewRateLimiter(10.0/60.0, 10)
}

func NewGeneralRateLimiter() *RateLimiter {
	return NewRateLimiter(1, 60)
}

func NewPublicRateLimiter() *RateLimiter {
	return NewRateLimiter(3.0/60.0, 3)
}
