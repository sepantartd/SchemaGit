package cloud

import (
    "net/http"
    "sync"
    "time"
)

type rateBucket struct {
    count int
    reset time.Time
}

var (
    rateMu   sync.Mutex
    rateData = make(map[string]*rateBucket)
)

func RateLimitMiddleware(limit int, window time.Duration) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            ip := r.RemoteAddr

            rateMu.Lock()
            b, ok := rateData[ip]
            if !ok || time.Now().After(b.reset) {
                b = &rateBucket{count: 0, reset: time.Now().Add(window)}
                rateData[ip] = b
            }

            if b.count >= limit {
                rateMu.Unlock()
                http.Error(w, "rate limit exceeded", http.StatusTooManyRequests)
                return
            }

            b.count++
            rateMu.Unlock()

            next.ServeHTTP(w, r)
        })
    }
}
