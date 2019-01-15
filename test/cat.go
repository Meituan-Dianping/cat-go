package main

import (
	"errors"
	"math/rand"
	"sync"
	"time"

	"../cat"
)

const TestType = "foo"

var wg = sync.WaitGroup{}

func init() {
	cat.DebugOn()
	cat.Init("gocat.v2")
}

// send transaction
func case1() {
	t := cat.NewTransaction(TestType, "test")
	defer t.Complete()
	t.AddData("foo", "bar")

	if rand.Int31n(1000) == 0 {
		t.SetStatus(cat.FAIL)
	}

	t.SetDurationStart(time.Now().Add(-5 * time.Second))
	t.SetTime(time.Now().Add(-5 * time.Second))
	t.SetDuration(time.Millisecond * 500)
}

// send completed transaction with duration
func case2() {
	cat.NewCompletedTransactionWithDuration(TestType, "completed", time.Second*24)
	cat.NewCompletedTransactionWithDuration(TestType, "completed-over-60s", time.Second*65)
}

// send event
func case3() {
	// way 1
	e := cat.NewEvent(TestType, "event-1")
	e.Complete()
	// way 2

	if rand.Int31n(1000) == 0 {
		cat.LogEvent(TestType, "event-2")
	} else {
		cat.LogEvent(TestType, "event-2", cat.FAIL)
	}
	cat.LogEvent(TestType, "event-3", cat.SUCCESS, "foobar")
}

// send error with backtrace
func case4() {
	if rand.Int31n(1000) == 0 {
		err := errors.New("error")
		cat.LogError(err)
	}
}

// send metric
func case5() {
	cat.LogMetricForCount("metric-1")
	cat.LogMetricForCount("metric-2", 3)
	cat.LogMetricForDuration("metric-3", 150*time.Millisecond)
	cat.NewMetricHelper("metric-4").Count(7)
	cat.NewMetricHelper("metric-5").Duration(time.Second)
}

func run(f func()) {
	defer wg.Done()

	for i := 0; i < 10000; i++ {
		f()
		time.Sleep(time.Microsecond * 10)
	}
}

func start(f func()) {
	wg.Add(1)
	go run(f)
}

func main() {
	start(case1)
	start(case2)
	start(case3)
	start(case4)
	start(case5)

	wg.Wait()

	cat.Shutdown()
}
