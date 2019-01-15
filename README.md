# Gocat.v2

`gocat.v2` is the client sdk for the [CAT(Central Application Tracking)](https://github.com/dianping/cat) for the Go programming language. 

See the [**DIFFERENCES**](DIFFERENCES.md) between [gocat](https://github.com/dianping/cat/tree/master/lib) and `gocat.v2`

The sdk requires a minimum version of Go 1.8

The current latest version is 2.0.0

Please check [**CHANGELOG**](./CHANGELOG.md) for information about our latest updates to the sdk.

## Installation 

You can simply get started with our sdk by using the following command:
```
go get github.com/dianping/gocat.v2
```

## Quickstart 

This example shows how you can send tracking data to CAT using `gocat.v2`.

```go
package main

import (
	"errors"
	"math/rand"
	"sync"
	"time"

	"./cat" 
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

	if rand.Int31n(100) == 0 {
		t.SetStatus(cat.FAIL)
	}

	t.AddData("foo", "bar")

	t.NewEvent(TestType, "event-1")
	t.Complete()

	if rand.Int31n(100) == 0 {
		t.LogEvent(TestType, "event-2", cat.FAIL)
	} else {
		t.LogEvent(TestType, "event-2")
	}
	t.LogEvent(TestType, "event-3", cat.SUCCESS, "k=v")

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
	e := cat.NewEvent(TestType, "event-4")
	e.Complete()
	// way 2

	if rand.Int31n(100) == 0 {
		cat.LogEvent(TestType, "event-5", cat.FAIL)
	} else {
		cat.LogEvent(TestType, "event-5")
	}
	cat.LogEvent(TestType, "event-6", cat.SUCCESS, "foobar")
}

// send error with backtrace
func case4() {
	if rand.Int31n(100) == 0 {
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

	for i := 0; i < 100000000; i++ {
		f()
		time.Sleep(time.Microsecond * 100)
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

```

## License

This SDK is distributed under the [Apache License, Version 2.0](http://www.apache.org/licenses/LICENSE-2.0),
