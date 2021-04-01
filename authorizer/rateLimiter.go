package authorizer

import (
	"container/list"
	"fmt"
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

func (l *RateLimiter) Take(event time.Time) (bool, error) {
	last := l.pastEvents.Back()
	if last != nil && !event.After(last.Value.(time.Time)) {
		return false, fmt.Errorf("Events must be sent in chronological order: Received %v after %v", event, last)
	}

	l.popEventsBefore(event.Add(-l.interval))
	if l.pastEvents.Len() >= l.maxEvents {
		return false, nil
	}
	l.pastEvents.PushBack(event)
	return true, nil
}

func (l *RateLimiter) popEventsBefore(threshold time.Time) {
	for l.pastEvents.Len() > 0 {
		firstElm := l.pastEvents.Front()
		first := firstElm.Value.(time.Time)
		if !first.Before(threshold) {
			break
		}
		l.pastEvents.Remove(firstElm)
	}
}
