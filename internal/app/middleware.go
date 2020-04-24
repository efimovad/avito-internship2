package app

import (
	"github.com/efimovad/avito-internship2/internal/models"
	"golang.org/x/crypto/bcrypt"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Middleware struct {
	mask net.IPMask
	visitors *models.Visitors
	limit int
	hash string
}

func NewLimitMiddleware(cidr int, limit int, period time.Duration, wait time.Duration, pass string) *Middleware {
	m, _ := GetNetIP(cidr)
	passHash, _ := HashPassword(pass)

	return &Middleware{
		mask:     m,
		limit:	  limit,
		visitors: models.NewVisitors(limit, period, wait),
		hash: passHash,
	}
}

func (m *Middleware) Limit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ipStr := r.Header.Get("X-FORWARDED-FOR")
		ip := net.ParseIP(ipStr)
		if ip == nil {
			body := getResponseText(WRONG_IP_ERROR, "X-Forwarded-For header include wrong ID")
			Respond(w, body, http.StatusForbidden)
			return
		}
		netIP := ip.Mask(m.mask)

		limiter := m.visitors.GetVisitor(netIP.String())
		if limiter.Allow() == false {
			text := "I only allow " + strconv.Itoa(m.limit) + " requests per period to this Web site per net. Try again soon."
			body := getResponseText(TOO_MANY_REQ_ERROR, text)
			Respond(w, body, http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (m *Middleware) Admin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasPrefix(r.URL.Path, "/admin") {
			next.ServeHTTP(w, r)
			return
		}

		login := r.URL.Query().Get("login")
		password := r.URL.Query().Get("password")

		if !m.isAdmin(login, password) {
			body := getResponseText(NO_ACCESS_ERROR, "Enter admin login and password")
			Respond(w, body, http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (m *Middleware) isAdmin(login string, password string) bool {
	if login == "admin" && CheckPasswordHash(password, m.hash) {
		return true
	}
	return false
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 8)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}