package app

import (
	"encoding/json"
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
			http.Error(w, body, http.StatusForbidden)
			return
		}
		netIP := ip.Mask(m.mask)

		limiter := m.visitors.GetVisitor(netIP.String())
		if limiter.Allow() == false {
			text := "I only allow " + strconv.Itoa(m.limit) + " requests per period to this Web site per net. Try again soon."
			body := getResponseText(TOO_MANY_REQ_ERROR, text)
			http.Error(w, body, http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}

type User struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (m *Middleware) Admin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasPrefix(r.URL.Path, "/admin") {
			next.ServeHTTP(w, r)
			return
		}

		var user User
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if !m.isAdmin(user) {
			body := getResponseText(NO_ACCESS_ERROR, "Enter admin login and password")
			http.Error(w, body, http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (m *Middleware) isAdmin(user User) bool {
	if user.Login == "admin" && CheckPasswordHash(user.Password, m.hash) {
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