package models

import (
	"sync"
	"time"
)

type Limiter struct {
	limit int
	period time.Duration
	wait time.Duration
	last time.Time
	first time.Time
	curr int
	mu sync.Mutex
	isPaused bool
}

func NewLimit(limit int, period time.Duration, wait time.Duration) *Limiter{
	return &Limiter{
		limit:  limit,
		period: period,
		wait:   wait,

		first: 	time.Now(),
		last: 	time.Now(),

		isPaused: false,
	}
}

func (l *Limiter) TimeLeft() time.Duration {
	return l.wait - time.Since(l.last)
}

func (l *Limiter) Allow() bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.curr = l.curr + 1

	if !l.isPaused && time.Since(l.first) >= l.period {
		l.first = time.Now()
		l.curr = 0
	}

	if l.isPaused && time.Since(l.last) >= l.wait {
		l.isPaused = false
		l.curr = 0
	}

	if l.curr > l.limit {
		l.isPaused = true
		return false
	}

	l.last = time.Now()

	return true
}

func (l *Limiter) CleanUp() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.curr = 0
}

func (l *Limiter) GetCurr() int {
	return l.curr
}