package app

import (
	"github.com/efimovad/avito-internship2/internal/models"
	"net"
	"net/http"
	"time"
)

type Middleware struct {
	mask net.IPMask
	visitors *models.Visitors
}

func NewLimitMiddleware(cidr int, limit int, period time.Duration, wait time.Duration) *Middleware {
	m, _ := GetNetIP(cidr)
	return &Middleware{
		mask:     m,
		visitors: models.NewVisitors(limit, period, wait),
	}
}

func (m *Middleware) Limit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ipStr := r.Header.Get("X-FORWARDED-FOR")
		ip := net.ParseIP(ipStr)
		netIP := ip.Mask(m.mask)

		limiter := m.visitors.GetVisitor(netIP.String())
		if limiter.Allow() == false {
			Error(w, Error429, http.StatusTooManyRequests, limiter.TimeLeft())
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (m *Middleware) Reset(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.RequestURI == "/reset" {
			ipStr := r.Header.Get("X-FORWARDED-FOR")
			ip := net.ParseIP(ipStr)
			netIP := ip.Mask(m.mask)

			limiter := m.visitors.GetVisitor(netIP.String())
			limiter.CleanUp()
			return
		}
		next.ServeHTTP(w, r)
	})
}
