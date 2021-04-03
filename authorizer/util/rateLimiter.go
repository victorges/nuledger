// Package util defines some more generic utilities that can be used in the core
// business logic of the authorizer/ledger components.
package util

import (
	"container/list"
	"time"
)

// RateLimiter is a utility that can be used to limit the maximum rate that any
// event could happen, defined by its consumers. All it needs to receive is the
// timestamp of the next event, which must be sent in ascending order. It will
// handle any internal state necessary for rate limiting the events, given the
// configuration of maximum number of events allowed in an interval. The zero
// value for RateLimiter is one that never allows any event.
//
// Notice that, since it is used for a financial system, it uses an expensive
// algorithm optimized to give the exact answer given the past recent events. It
// is not appropriate for handling a very large number of events, say for
// rate-limiting requests on a web server. In that other case, we could use some
// more efficient algorithm here like a periodically resetting counter or a
// leaky bucket, which would save a lot of CPU and memory but allow for several
// edge cases with the wrong response regarding exceeding the rate limit or not.
type RateLimiter struct {
	// MaxEvents specifies the maximum amount of events allowed in the
	// configured interval. If left zero, no events will ever be allowed.
	MaxEvents int
	// Interval specifies what is the interval to be analyzed for determining
	// the rate of events and either allowing them or not. Works together with
	// the MaxEvents configured above. If left zero, all events will be allowed,
	// which can be thought of as an infinite frequency.
	Interval time.Duration

	pastEvents list.List
}

// Allows function checks whether the given event is allowed to happen without
// making any changes to the internal state of the rate limiter. It returns true
// if the event is allowed.
//
// It is useful in case there are other reasons for an event not to be allowed,
// in which case we don't want to consider them in the rate limiter in case it
// ends up not being allowed for other reasons.
func (l *RateLimiter) Allows(event time.Time) bool {
	return l.countEventsAfter(event.Add(-l.Interval)) < l.MaxEvents
}

// Take actually updates the internal state of the rate limiter in order to
// consider the given as event as having happened. It still checks whether the
// event is actually allowed and returns whether the event was actually taken or
// not.
//
// If the same event is sent to the Allows function and no other change has been
// made to the rate limiter, Take is guaranteed to return the same result. In
// other words, if Allows has allowed an event to happen, Take will certainly
// accept it unless other events have been sent to the rate limiter in between.
func (l *RateLimiter) Take(event time.Time) bool {
	l.popEventsNotAfter(event.Add(-l.Interval))
	if l.pastEvents.Len() >= l.MaxEvents {
		return false
	}
	l.pastEvents.PushBack(event)
	return true
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
