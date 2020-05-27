package tokenbucket

import (
	"math"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTake(t *testing.T) {
	should := assert.New(t)

	tb := NewBucket(250*time.Millisecond, 10)
	d := tb.Take(10)
	should.Equal(time.Duration(0), d)

	tb.TakeAvailable(10)
	should.Equal(time.Duration(0), d)

	tb = NewBucket(250*time.Millisecond, 10)
	d, ok := tb.TakeMaxDuration(1, 250*time.Millisecond)
	if should.Equal(true, ok) {
		should.Equal(time.Duration(0), d)
	}

	tb = NewBucket(250*time.Millisecond, 10)
	tb.Wait(10)
	ok = tb.WaitMaxDuration(1, 250*time.Millisecond)
	should.Equal(true, ok)
	should.Equal(int64(10), tb.Capacity())

	tb = NewBucket(250*time.Millisecond, 10)
	should.Equal(true, tb.TakeOneAvailable())
}

type takeReq struct {
	time       time.Duration
	count      int64
	expectWait time.Duration
}

var takeTests = []struct {
	about        string
	fillInterval time.Duration
	capacity     int64
	reqs         []takeReq
}{{
	about:        "serial requests",
	fillInterval: 250 * time.Millisecond,
	capacity:     10,
	reqs: []takeReq{{
		time:       0,
		count:      0,
		expectWait: 0,
	}, {
		time:       0,
		count:      10,
		expectWait: 0,
	}, {
		time:       0,
		count:      1,
		expectWait: 250 * time.Millisecond,
	}, {
		time:       250 * time.Millisecond,
		count:      1,
		expectWait: 250 * time.Millisecond,
	}},
}, {
	about:        "concurrent requests",
	fillInterval: 250 * time.Millisecond,
	capacity:     10,
	reqs: []takeReq{{
		time:       0,
		count:      10,
		expectWait: 0,
	}, {
		time:       0,
		count:      2,
		expectWait: 500 * time.Millisecond,
	}, {
		time:       0,
		count:      2,
		expectWait: 1000 * time.Millisecond,
	}, {
		time:       0,
		count:      1,
		expectWait: 1250 * time.Millisecond,
	}},
}, {
	about:        "more than capacity",
	fillInterval: 1 * time.Millisecond,
	capacity:     10,
	reqs: []takeReq{{
		time:       0,
		count:      10,
		expectWait: 0,
	}, {
		time:       20 * time.Millisecond,
		count:      15,
		expectWait: 5 * time.Millisecond,
	}},
}, {
	about:        "sub-quantum time",
	fillInterval: 10 * time.Millisecond,
	capacity:     10,
	reqs: []takeReq{{
		time:       0,
		count:      10,
		expectWait: 0,
	}, {
		time:       7 * time.Millisecond,
		count:      1,
		expectWait: 3 * time.Millisecond,
	}, {
		time:       8 * time.Millisecond,
		count:      1,
		expectWait: 12 * time.Millisecond,
	}},
}, {
	about:        "within capacity",
	fillInterval: 10 * time.Millisecond,
	capacity:     5,
	reqs: []takeReq{{
		time:       0,
		count:      5,
		expectWait: 0,
	}, {
		time:       60 * time.Millisecond,
		count:      5,
		expectWait: 0,
	}, {
		time:       60 * time.Millisecond,
		count:      1,
		expectWait: 10 * time.Millisecond,
	}, {
		time:       80 * time.Millisecond,
		count:      2,
		expectWait: 10 * time.Millisecond,
	}},
}}

func TestTask(t *testing.T) {
	should := assert.New(t)
	for i, test := range takeTests {
		tb := NewBucket(test.fillInterval, test.capacity)
		for j, req := range test.reqs {
			d, ok := tb.take(tb.startTime.Add(req.time), req.count, infinityDuration)
			if should.Equal(ok, true) {
				should.Equal(req.expectWait, d, "test %d.%d, %s, got %v want %v", i, j, test.about, d, req.expectWait)
			}
		}
	}
}

func TestTakeMaxDuration(t *testing.T) {
	should := assert.New(t)

	for i, test := range takeTests {
		tb := NewBucket(test.fillInterval, test.capacity)
		for j, req := range test.reqs {
			if req.expectWait > 0 {
				d, ok := tb.take(tb.startTime.Add(req.time), req.count, req.expectWait-1)
				if should.Equal(ok, false) {
					should.Equal(d, time.Duration(0))
				}
			}
			d, ok := tb.take(tb.startTime.Add(req.time), req.count, req.expectWait)
			if should.Equal(ok, true) {
				should.Equal(req.expectWait, d, "test %d.%d, %s, got %v want %v", i, j, test.about, d, req.expectWait)
			}

		}
	}
}

type takeAvailableReq struct {
	time   time.Duration
	count  int64
	expect int64
}

var takeAvailableTests = []struct {
	about        string
	fillInterval time.Duration
	capacity     int64
	reqs         []takeAvailableReq
}{{
	about:        "serial requests",
	fillInterval: 250 * time.Millisecond,
	capacity:     10,
	reqs: []takeAvailableReq{{
		time:   0,
		count:  0,
		expect: 0,
	}, {
		time:   0,
		count:  10,
		expect: 10,
	}, {
		time:   0,
		count:  1,
		expect: 0,
	}, {
		time:   250 * time.Millisecond,
		count:  1,
		expect: 1,
	}},
}, {
	about:        "concurrent requests",
	fillInterval: 250 * time.Millisecond,
	capacity:     10,
	reqs: []takeAvailableReq{{
		time:   0,
		count:  5,
		expect: 5,
	}, {
		time:   0,
		count:  2,
		expect: 2,
	}, {
		time:   0,
		count:  5,
		expect: 3,
	}, {
		time:   0,
		count:  1,
		expect: 0,
	}},
}, {
	about:        "more than capacity",
	fillInterval: 1 * time.Millisecond,
	capacity:     10,
	reqs: []takeAvailableReq{{
		time:   0,
		count:  10,
		expect: 10,
	}, {
		time:   20 * time.Millisecond,
		count:  15,
		expect: 10,
	}},
}, {
	about:        "within capacity",
	fillInterval: 10 * time.Millisecond,
	capacity:     5,
	reqs: []takeAvailableReq{{
		time:   0,
		count:  5,
		expect: 5,
	}, {
		time:   60 * time.Millisecond,
		count:  5,
		expect: 5,
	}, {
		time:   70 * time.Millisecond,
		count:  1,
		expect: 1,
	}},
}}

func TestTakeAvailable(t *testing.T) {
	should := assert.New(t)

	for i, test := range takeAvailableTests {
		tb := NewBucket(test.fillInterval, test.capacity)
		for j, req := range test.reqs {
			d := tb.takeAvailable(tb.startTime.Add(req.time), req.count)
			should.Equal(req.expect, d, "test %d.%d, %s, got %v want %v", i, j, test.about, d, req.expect)
		}
	}
}

func TestPanics(t *testing.T) {
	should := assert.New(t)
	should.Panics(func() { NewBucket(0, 1) }, "token bucket fill interval is not > 0")
	should.Panics(func() { NewBucket(-2, 1) }, "token bucket fill interval is not > 0")
	should.Panics(func() { NewBucket(1, 0) }, "token bucket capacity is not > 0")
	should.Panics(func() { NewBucket(1, -2) }, "token bucket capacity is not > 0")
}

func isCloseTo(x, y, tolerance float64) bool {
	return math.Abs(x-y)/y < tolerance
}

func TestRate(t *testing.T) {
	should := assert.New(t)

	tb := NewBucket(1, 1)
	should.Equal(true, isCloseTo(tb.Rate(), 1e9, 0.00001), "got %v want 1e9", tb.Rate())

	tb = NewBucket(2*time.Second, 1)
	should.Equal(true, isCloseTo(tb.Rate(), 0.5, 0.00001), "got %v want 0.5", tb.Rate())

	tb = NewBucketWithQuantum(100*time.Millisecond, 1, 5)
	should.Equal(true, isCloseTo(tb.Rate(), 50, 0.00001), "got %v want 50", tb.Rate())
}

func checkRate(t *testing.T, rate float64) {
	should := assert.New(t)

	tb := NewBucketWithRate(rate, 1<<62)
	should.Equal(true, isCloseTo(tb.Rate(), rate, rateMargin), "got %g want %v", tb.Rate(), rate)

	d, ok := tb.take(tb.startTime, 1<<62, infinityDuration)
	should.Equal(true, ok)
	should.Equal(time.Duration(0), d)

	// Check that the actual rate is as expected by
	// asking for a not-quite multiple of the bucket's
	// quantum and checking that the wait time
	// correct.
	d, ok = tb.take(tb.startTime, tb.quantum*2-tb.quantum/2, infinityDuration)
	should.Equal(true, ok)

	expectTime := 1e9 * float64(tb.quantum) * 2 / rate
	should.Equal(true, isCloseTo(float64(d), expectTime, rateMargin), "rate %g: got %g want %v", rate, float64(d), expectTime)
}

func TestNewBucketWithRate(t *testing.T) {
	for rate := float64(1); rate < 1e6; rate += 7 {
		checkRate(t, rate)
	}
	for _, rate := range []float64{
		1024 * 1024 * 1024,
		1e-5,
		0.9e-5,
		0.5,
		0.9,
		0.9e8,
		3e12,
		4e18,
		float64(1<<63 - 1),
	} {
		checkRate(t, rate)
		checkRate(t, rate/3)
		checkRate(t, rate*1.3)
	}
}

var availTests = []struct {
	about        string
	capacity     int64
	fillInterval time.Duration
	take         int64
	sleep        time.Duration

	expectCountAfterTake  int64
	expectCountAfterSleep int64
}{{
	about:                 "should fill tokens after interval",
	capacity:              5,
	fillInterval:          time.Second,
	take:                  5,
	sleep:                 time.Second,
	expectCountAfterTake:  0,
	expectCountAfterSleep: 1,
}, {
	about:                 "should fill tokens plus existing count",
	capacity:              2,
	fillInterval:          time.Second,
	take:                  1,
	sleep:                 time.Second,
	expectCountAfterTake:  1,
	expectCountAfterSleep: 2,
}, {
	about:                 "shouldn't fill before interval",
	capacity:              2,
	fillInterval:          2 * time.Second,
	take:                  1,
	sleep:                 time.Second,
	expectCountAfterTake:  1,
	expectCountAfterSleep: 1,
}, {
	about:                 "should fill only once after 1*interval before 2*interval",
	capacity:              2,
	fillInterval:          2 * time.Second,
	take:                  1,
	sleep:                 3 * time.Second,
	expectCountAfterTake:  1,
	expectCountAfterSleep: 2,
}}

func TestAvailable(t *testing.T) {
	should := assert.New(t)

	for i, tt := range availTests {
		tb := NewBucket(tt.fillInterval, tt.capacity)
		c := tb.takeAvailable(tb.startTime, tt.take)
		should.Equal(tt.take, c, "#%d: %s, take = %d, want = %d", i, tt.about, c, tt.take)

		c = tb.available(tb.startTime)
		should.Equal(tt.expectCountAfterTake, c, "#%d: %s, after take, available = %d, want = %d", i, tt.about, c, tt.expectCountAfterTake)

		c = tb.available(tb.startTime.Add(tt.sleep))
		should.Equal(tt.expectCountAfterSleep, c, "#%d: %s, after some time it should fill in new tokens, available = %d, want = %d",
			i, tt.about, c, tt.expectCountAfterSleep)
	}
}

func TestNoBonusTokenAfterBucketIsFull(t *testing.T) {
	tb := NewBucketWithQuantum(time.Second*1, 100, 20)
	curAvail := tb.Available()
	if curAvail != 100 {
		t.Fatalf("initially: actual available = %d, expected = %d", curAvail, 100)
	}

	time.Sleep(time.Second * 5)

	curAvail = tb.Available()
	if curAvail != 100 {
		t.Fatalf("after pause: actual available = %d, expected = %d", curAvail, 100)
	}

	cnt := tb.TakeAvailable(100)
	if cnt != 100 {
		t.Fatalf("taking: actual taken count = %d, expected = %d", cnt, 100)
	}

	curAvail = tb.Available()
	if curAvail != 0 {
		t.Fatalf("after taken: actual available = %d, expected = %d", curAvail, 0)
	}
}

func BenchmarkWait(b *testing.B) {
	tb := NewBucket(1, 16*1024)
	for i := b.N - 1; i >= 0; i-- {
		tb.Wait(1)
	}
}

func BenchmarkNewBucket(b *testing.B) {
	for i := b.N - 1; i >= 0; i-- {
		NewBucketWithRate(4e18, 1<<62)
	}
}
