package app

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestServer_OkHandler(t *testing.T) {
	limit := 100
	period := time.Second * 60
	wait := time.Second * 120

	testCases := []struct {
		name	string
		cidr 	int
		status	int
	}{
		{
			name: 	"valid",
			cidr: 	24,
			status: http.StatusOK,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/", nil)

			s := NewServer(tc.cidr, limit, period, wait)
			rec := httptest.NewRecorder()

			handler := http.HandlerFunc(s.OkHandler)
			handler.ServeHTTP(rec, req)

			if status := rec.Code; status != tc.status {
				t.Error("Wrong status")
			}
		})
	}
}