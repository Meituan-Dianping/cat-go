package main

import (
	"errors"
	"time"

	gocat "../cat"
)

var cat = gocat.Instance()

const TestType = "foo"

func init() {
	gocat.Init("gocat.v2")
}

// send transaction
func case1() {
	t := cat.NewTransaction(TestType, "test")
	defer t.Complete()
	t.AddData("foo", "bar")
	t.SetStatus(gocat.FAIL)
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
	cat.LogEvent(TestType, "event-2")
	cat.LogEvent(TestType, "event-3", gocat.FAIL)
	cat.LogEvent(TestType, "event-4", gocat.FAIL, "foobar")
}

// send error with backtrace
func case4() {
	err := errors.New("error")
	cat.LogError(err)
}

// send metric
func case5() {
	cat.LogMetricForCount("metric-1")
	cat.LogMetricForCount("metric-2", 3)
	cat.LogMetricForDuration("metric-3", 150*time.Millisecond)
}

func run(f func()) {
	for {
		f()
		time.Sleep(time.Millisecond)
	}
}

func main() {
	// go run(case1)
	// go run(case2)
	// go run(case3)
	// go run(case4)
	// go run(case5)

	// // wait until main process has been killed
	// var ch chan int
	// <-ch
	gocat.DebugOn()

	case1()
	gocat.Shutdown()
}
