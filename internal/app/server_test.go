package app

import (
	"bytes"
	"encoding/json"
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

			s := NewServer(tc.cidr, limit, period, wait, "")
			rec := httptest.NewRecorder()

			handler := http.HandlerFunc(s.OkHandler)
			handler.ServeHTTP(rec, req)

			if status := rec.Code; status != tc.status {
				t.Error("Wrong status")
			}
		})
	}
}

func TestServer_ResetHandler(t *testing.T) {
	limit := 100
	period := time.Second * 60
	wait := time.Second * 120
	cidr := 24

	testCases := []struct {
		name	string
		status	int
		user	User
	}{
		{
			name: 	"valid",
			status: http.StatusOK,
			user:User{
				Login:    "admin",
				Password: "123456",
			},
		},
	}

	s := NewServer(cidr, limit, period, wait, "123456")
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			body := new(bytes.Buffer)
			if err := json.NewEncoder(body).Encode(tc.user); err != nil {
				t.Error(err)
			}

			req := httptest.NewRequest("POST", "/admin/reset", body)
			rec := httptest.NewRecorder()
			handler := http.HandlerFunc(s.OkHandler)
			handler.ServeHTTP(rec, req)

			if status := rec.Code; status != tc.status {
				t.Error("Wrong status")
			}
		})
	}
}