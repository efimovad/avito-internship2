package models

import (
	"errors"
	"sync"
	"time"
)

type visitor struct {
	limiter *Limiter
	lastSeen time.Time
}

type Visitors struct {
	visitors	map[string]*visitor
	mu			sync.Mutex
	limit		int
	period		time.Duration
	wait		time.Duration
}

func NewVisitors(limit int, period time.Duration, wait time.Duration) *Visitors{
	return &Visitors{
		visitors:	make(map[string]*visitor),
		limit:		limit,
		period:		period,
		wait:		wait,
	}
}

func (vis *Visitors) GetVisitor(ip string) *Limiter {
	vis.mu.Lock()
	defer vis.mu.Unlock()

	v, exists := vis.visitors[ip]
	if !exists {
		limiter := NewLimit(vis.limit, vis.period, vis.wait)
		vis.visitors[ip] = &visitor{limiter, time.Now()}
		return limiter
	}

	v.lastSeen = time.Now()
	return v.limiter
}

func (vis *Visitors) GetLimitCurr(ip string) (int, error) {
	v, exists := vis.visitors[ip]
	if !exists {
		return 0, errors.New("not found")
	}
	return v.limiter.GetCurr(), nil
}

