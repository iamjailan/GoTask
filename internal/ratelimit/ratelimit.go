package ratelimit

import (
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"gotask/internal/utils/response"
)

type clientWindow struct {
	startedAt time.Time
	requests  int
}

type Limiter struct {
	mu          sync.Mutex
	limit       int
	window      time.Duration
	clients     map[string]clientWindow
	lastCleanup time.Time
}

func New(limit int, window time.Duration) *Limiter {
	if limit < 1 {
		panic("rate limit must be positive")
	}
	if window <= 0 {
		panic("rate limit window must be positive")
	}

	return &Limiter{
		limit:   limit,
		window:  window,
		clients: make(map[string]clientWindow),
	}
}

func (l *Limiter) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if l.allow(c.ClientIP(), time.Now()) {
			c.Next()
			return
		}

		c.Header("Retry-After", strconv.Itoa(max(1, int(l.window.Seconds()))))
		c.Header("X-RateLimit-Limit", strconv.Itoa(l.limit))
		response.Error(c, http.StatusTooManyRequests, "rate limit exceeded")
		c.Abort()
	}
}

func (l *Limiter) allow(clientIP string, now time.Time) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.lastCleanup.IsZero() || now.Sub(l.lastCleanup) >= l.window {
		for client, counter := range l.clients {
			if now.Sub(counter.startedAt) >= l.window {
				delete(l.clients, client)
			}
		}
		l.lastCleanup = now
	}

	counter, exists := l.clients[clientIP]
	if !exists || now.Sub(counter.startedAt) >= l.window {
		l.clients[clientIP] = clientWindow{startedAt: now, requests: 1}
		return true
	}
	if counter.requests >= l.limit {
		return false
	}

	counter.requests++
	l.clients[clientIP] = counter
	return true
}
