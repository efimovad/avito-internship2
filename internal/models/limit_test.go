package models

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestLimiter_Allow(t *testing.T) {
	limit := 100
	limiter := NewLimit(limit, time.Minute, 2 * time.Minute)

	testCases := []struct {
		name	string
		res		bool
		requests int
	}{
		{
			name: 	"in limit",
			requests: limit,
			res: 	true,
		},
		{
			name: 	"up than limit",
			requests: limit + 1,
			res: 	false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var res bool
			for i := 0; i < tc.requests; i++ {
				res = limiter.Allow()
			}
			assert.Equal(t, res, tc.res)
		})
	}
}

func TestLimiter_CleanUp(t *testing.T) {
	limit := 100
	limiter := NewLimit(limit, time.Minute, 2 * time.Minute)

	testCases := []struct {
		name	string
		res		bool
		requests int
	}{
		{
			name: 	"valid",
			requests: limit,
			res: 	true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			for i := 0; i < tc.requests; i++ {
				_ = limiter.Allow()
			}
			assert.NotEqual(t, 0, limiter.curr)
			limiter.CleanUp()
			assert.Equal(t, 0, limiter.curr)
		})
	}
}

func TestLimiter_GetCurr(t *testing.T) {
	limit := 100
	limiter := NewLimit(limit, time.Minute, 2 * time.Minute)

	testCases := []struct {
		name	string
		res		bool
		requests int
	}{
		{
			name: 	"valid",
			requests: limit,
			res: 	true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			for i := 0; i < tc.requests; i++ {
				_ = limiter.Allow()
			}
			assert.Equal(t, tc.requests, limiter.GetCurr())
		})
	}
}