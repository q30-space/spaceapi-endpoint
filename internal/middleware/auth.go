package middleware

import (
	"log"
	"net/http"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"
)

// FailedAttempts tracks failed authentication attempts for an IP
type FailedAttempts struct {
	Count        int
	FirstAttempt time.Time
	BlockedUntil *time.Time
}

// RateLimiter manages rate limiting for failed authentication attempts
type RateLimiter struct {
	attempts map[string]*FailedAttempts
	mutex    sync.RWMutex
	stopCh   chan struct{}
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter() *RateLimiter {
	rl := &RateLimiter{
		attempts: make(map[string]*FailedAttempts),
		stopCh:   make(chan struct{}),
	}
	
	// Start cleanup goroutine only if not in test mode
	if !isTestMode() {
		go rl.cleanup()
	}
	
	return rl
}

// isTestMode checks if we're running in test mode
func isTestMode() bool {
	// Check if we're running tests by looking at the call stack
	// This is more reliable than environment variables
	for i := 0; i < 10; i++ {
		if _, file, _, ok := runtime.Caller(i); ok {
			if strings.Contains(file, "_test.go") {
				return true
			}
		}
	}
	return false
}

// Stop stops the cleanup goroutine
func (rl *RateLimiter) Stop() {
	close(rl.stopCh)
}

// cleanup removes old entries every 30 minutes
func (rl *RateLimiter) cleanup() {
	ticker := time.NewTicker(30 * time.Minute)
	defer ticker.Stop()
	
	for {
		select {
		case <-ticker.C:
			rl.mutex.Lock()
			cutoff := time.Now().Add(-2 * time.Hour)
			for ip, attempt := range rl.attempts {
				if attempt.FirstAttempt.Before(cutoff) {
					delete(rl.attempts, ip)
				}
			}
			rl.mutex.Unlock()
		case <-rl.stopCh:
			return
		}
	}
}

// isBlocked checks if an IP is currently blocked
func (rl *RateLimiter) isBlocked(ip string) bool {
	rl.mutex.RLock()
	defer rl.mutex.RUnlock()
	
	attempt, exists := rl.attempts[ip]
	if !exists {
		return false
	}
	
	if attempt.BlockedUntil != nil && time.Now().Before(*attempt.BlockedUntil) {
		return true
	}
	
	return false
}

// recordFailedAttempt records a failed authentication attempt
func (rl *RateLimiter) recordFailedAttempt(ip string) {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()
	
	now := time.Now()
	attempt, exists := rl.attempts[ip]
	
	if !exists {
		attempt = &FailedAttempts{
			Count:        1,
			FirstAttempt: now,
		}
		rl.attempts[ip] = attempt
	} else {
		attempt.Count++
	}
	
	// If this is the 5th attempt within 15 minutes, block for 1 hour
	if attempt.Count >= 5 && now.Sub(attempt.FirstAttempt) <= 15*time.Minute {
		blockedUntil := now.Add(1 * time.Hour)
		attempt.BlockedUntil = &blockedUntil
		
		log.Printf("SECURITY: IP %s blocked for 1 hour after %d failed authentication attempts", ip, attempt.Count)
	}
}

// getRetryAfter returns the time until the block expires
func (rl *RateLimiter) getRetryAfter(ip string) int {
	rl.mutex.RLock()
	defer rl.mutex.RUnlock()
	
	attempt, exists := rl.attempts[ip]
	if !exists || attempt.BlockedUntil == nil {
		return 0
	}
	
	retryAfter := int(time.Until(*attempt.BlockedUntil).Seconds())
	if retryAfter < 0 {
		return 0
	}
	
	return retryAfter
}

// Global rate limiter instance
var rateLimiter *RateLimiter
var rateLimiterOnce sync.Once

// getRateLimiter returns the global rate limiter, initializing it if needed
func getRateLimiter() *RateLimiter {
	rateLimiterOnce.Do(func() {
		rateLimiter = NewRateLimiter()
	})
	return rateLimiter
}

// AuthMiddleware validates API key and enforces rate limiting
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get client IP
		clientIP := r.RemoteAddr
		if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
			clientIP = forwarded
		}
		
		// Get rate limiter
		rl := getRateLimiter()
		
		// Check if IP is currently blocked
		if rl.isBlocked(clientIP) {
			retryAfter := rl.getRetryAfter(clientIP)
			w.Header().Set("Retry-After", string(rune(retryAfter)))
			http.Error(w, "Too many failed authentication attempts. Please try again later.", http.StatusTooManyRequests)
			return
		}
		
		// Get API key from environment
		expectedKey := os.Getenv("SPACEAPI_AUTH_KEY")
		if expectedKey == "" {
			log.Println("ERROR: SPACEAPI_AUTH_KEY environment variable not set")
			http.Error(w, "Server configuration error", http.StatusInternalServerError)
			return
		}
		
		// Extract API key from headers
		var providedKey string
		
		// Check Authorization header (Bearer token)
		if auth := r.Header.Get("Authorization"); auth != "" {
			if len(auth) > 7 && auth[:7] == "Bearer " {
				providedKey = auth[7:]
			}
		}
		
		// Check X-API-Key header
		if providedKey == "" {
			providedKey = r.Header.Get("X-API-Key")
		}
		
		// Validate API key
		if providedKey == "" {
			rl.recordFailedAttempt(clientIP)
			http.Error(w, "API key required", http.StatusUnauthorized)
			return
		}
		
		if providedKey != expectedKey {
			rl.recordFailedAttempt(clientIP)
			log.Printf("SECURITY: Invalid API key attempt from IP %s", clientIP)
			http.Error(w, "Invalid API key", http.StatusUnauthorized)
			return
		}
		
		// Authentication successful, proceed to next handler
		next.ServeHTTP(w, r)
	})
}
