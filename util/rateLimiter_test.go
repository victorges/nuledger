package util_test

import (
	"nuledger/util"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

var startTime = time.Date(2021, time.April, 1, 12, 26, 47, 0, time.Local)

func TestRateLimiter(t *testing.T) {
	Convey("Given a RateLimiter", t, func() {
		limiter := util.RateLimiter{MaxEvents: 5, Interval: 1 * time.Second}

		Convey("Allows should not affect state", func() {
			allAllowed := true
			for i := 0; i < 100; i++ {
				allAllowed = allAllowed && limiter.Allows(startTime)
			}
			So(allAllowed, ShouldBeTrue)
		})

		testTake := func(time time.Time) bool {
			allows := limiter.Allows(time)
			taken := limiter.Take(time)
			So(allows, ShouldEqual, taken)
			return allows
		}

		Convey("When events come periodically", func() {
			Convey("It should allow if at least of minimum period", func() {
				currTime := startTime
				period := limiter.Interval / time.Duration(limiter.MaxEvents)
				for i := 0; i < 10*limiter.MaxEvents; i++ {
					currTime = currTime.Add(period)
					So(testTake(currTime), ShouldBeTrue)
				}
			})

			Convey("It should NOT allow if period just a bit too short", func() {
				currTime := startTime
				period := (limiter.Interval - 1) / time.Duration(limiter.MaxEvents)
				for i := 0; i < limiter.MaxEvents; i++ {
					currTime = currTime.Add(period)
					So(testTake(currTime), ShouldBeTrue)
				}
				So(testTake(currTime), ShouldBeFalse)
			})
		})

		Convey("When events come in bursts", func() {
			Convey("In the beginning of the interval", func() {
				for i := 0; i < limiter.MaxEvents; i++ {
					So(testTake(startTime), ShouldBeTrue)
				}

				Convey("It should only allow new events after the end of the interval", func() {
					So(testTake(startTime.Add(limiter.Interval/2)), ShouldBeFalse)
					So(testTake(startTime.Add(limiter.Interval-1*time.Millisecond)), ShouldBeFalse)
					So(testTake(startTime.Add(limiter.Interval-1)), ShouldBeFalse)
					So(testTake(startTime.Add(limiter.Interval)), ShouldBeTrue)
				})
			})

			Convey("In the end of the interval", func() {
				So(testTake(startTime), ShouldBeTrue)

				currTime := startTime.Add(limiter.Interval - 1)
				for i := 1; i < limiter.MaxEvents; i++ {
					So(testTake(currTime), ShouldBeTrue)
				}
				So(testTake(currTime), ShouldBeFalse)

				Convey("It should allow 1 event as soon as the first one expires", func() {
					currTime = startTime.Add(limiter.Interval)
					So(testTake(currTime), ShouldBeTrue)

					Convey("Then only by the end of the next interval", func() {
						So(testTake(currTime.Add(limiter.Interval/2)), ShouldBeFalse)
						So(testTake(currTime.Add(limiter.Interval-1*time.Millisecond)), ShouldBeFalse)
						So(testTake(currTime.Add(limiter.Interval-1)), ShouldBeTrue)
					})
				})
			})
		})

		Convey("Corner cases", func() {
			Convey("Its zero value should never take any event", func() {
				limiter = util.RateLimiter{}
				So(testTake(startTime), ShouldBeFalse)
			})
			Convey("Non-zero interval with zero max events should never take any event", func() {
				limiter.MaxEvents = 0
				So(testTake(startTime), ShouldBeFalse)
			})
			Convey("Non-zero max events with zero interval should take all events", func() {
				limiter.Interval = 0
				for i := 0; i < 10; i++ {
					So(testTake(startTime), ShouldBeTrue)
				}
			})
		})
	})
}
