package app

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestMiddleware_Limit(t *testing.T) {
	limit := 100
	period := time.Second * 60
	wait := time.Second * 120

	testCases := []struct {
		name		string
		cidr 		int
		status		int
		requests	int
		ip 			[]string
	}{
		{
			name: 	"valid: 100 request per minute",
			requests: 100,
			cidr: 	24,
			status: http.StatusOK,
			ip:		[]string{"169.89.1.0"},
		},
		{
			name: 	"invalid: 101 request per minute",
			requests: 101,
			cidr: 	24,
			status: http.StatusTooManyRequests,
			ip:		[]string{"169.89.1.0"},
		},
		{
			name: 	"valid: 160 request per minute from diff nets",
			requests: 160,
			cidr: 	24,
			status: http.StatusOK,
			ip:		[]string{"169.89.1.0", "169.89.2.0"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			s := NewServer(tc.cidr, limit, period, wait)

			testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// before any request
			})

			rr := httptest.NewRecorder()
			handler := s.middleware.Limit(testHandler)

			for i := 0; i < tc.requests; i++ {
				req := httptest.NewRequest("GET", "/", nil)
				req.Header.Set("X-FORWARDED-FOR", tc.ip[i % len(tc.ip)])
				handler.ServeHTTP(rr, req)
			}

			assert.Equal(t, rr.Code, tc.status)
		})
	}
}