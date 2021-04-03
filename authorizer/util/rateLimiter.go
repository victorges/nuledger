package util

import (
	"container/list"
	"time"
)

type RateLimiter struct {
	maxEvents  int
	interval   time.Duration
	pastEvents *list.List
}

func NewRateLimiter(maxEvents int, interval time.Duration) *RateLimiter {
	return &RateLimiter{maxEvents, interval, list.New()}
}

func (l *RateLimiter) Take(event time.Time) bool {
	l.popEventsNotAfter(event.Add(-l.interval))
	if l.pastEvents.Len() >= l.maxEvents {
		return false
	}
	l.pastEvents.PushBack(event)
	return true
}

func (l *RateLimiter) Allows(event time.Time) bool {
	return l.countEventsAfter(event.Add(-l.interval)) < l.maxEvents
}

func (l *RateLimiter) popEventsNotAfter(threshold time.Time) {
	for l.pastEvents.Len() > 0 {
		elm := l.pastEvents.Front()
		value := elm.Value.(time.Time)
		if value.After(threshold) {
			break
		}
		l.pastEvents.Remove(elm)
	}
}

func (l *RateLimiter) countEventsAfter(threshold time.Time) int {
	count := 0
	for elm := l.pastEvents.Back(); elm != nil; elm = elm.Prev() {
		value := elm.Value.(time.Time)
		if !value.After(threshold) {
			break
		}
		count++
	}
	return count
}
