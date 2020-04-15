package app

import (
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

type Server struct {
	cidr int
	router *mux.Router
	middleware *Middleware
}

func NewServer(cidr int, limit int, period time.Duration, wait time.Duration) Server {
	return Server{
		router:		mux.NewRouter(),
		middleware: NewLimitMiddleware(cidr, limit, period, wait),
	}
}

func (s *Server) Start() error {
	s.configure()
	return http.ListenAndServe(":8080", s)
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *Server) configure() {
	s.router.Use(s.middleware.Reset, s.middleware.Limit)
	s.router.HandleFunc("/", s.OkHandler)
}

func (s *Server) OkHandler(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte("OK"))
}