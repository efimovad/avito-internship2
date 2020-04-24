package app

import (
	"fmt"
	"github.com/gorilla/mux"
	"net"
	"net/http"
	"strconv"
	"time"
)

type Server struct {
	cidr int
	router *mux.Router
	middleware *Middleware
}

func NewServer(cidr int, limit int, period time.Duration, wait time.Duration, pass string) Server {
	return Server{
		router:		mux.NewRouter(),
		middleware: NewLimitMiddleware(cidr, limit, period, wait, pass),
	}
}

func (s *Server) Start(port int) error {
	s.configure()
	fmt.Println("Running on", port, "port")
	return http.ListenAndServe(":" + strconv.Itoa(port), s)
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *Server) configure() {
	s.router.Use(s.middleware.Admin, s.middleware.Limit)
	s.router.HandleFunc("/", s.OkHandler).Methods(http.MethodGet)
	s.router.HandleFunc("/admin/reset", s.ResetHandler).Methods(http.MethodPost)
}

func (s *Server) OkHandler(w http.ResponseWriter, r *http.Request) {
	Respond(w, getResponseText("Hello","Hello from server"), http.StatusOK)
}

func (s *Server) ResetHandler(w http.ResponseWriter, r *http.Request) {
	ipStr := r.Header.Get("X-FORWARDED-FOR")
	ip := net.ParseIP(ipStr)
	netIP := ip.Mask(s.middleware.mask)

	limiter := s.middleware.visitors.GetVisitor(netIP.String())
	limiter.ResetLimit()

	Respond(w, getResponseText("Reset", "Limit was reset"), http.StatusOK)
}

