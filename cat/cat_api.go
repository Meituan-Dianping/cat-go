package cat

import (
	"time"

	"github.com/Meituan-Dianping/cat-go/message"
)

func NewTransaction(mtype, name string) message.Transactor {
	if !IsEnabled() {
		return &message.NullTransaction{}
	}
	return message.NewTransaction(mtype, name, manager.flush)
}

func NewCompletedTransactionWithDuration(mtype, name string, duration time.Duration) {
	if !IsEnabled() {
		return
	}

	var trans = NewTransaction(mtype, name)
	trans.SetDuration(duration)
	if duration > 0 && duration < 60*time.Millisecond {
		trans.SetTime(time.Now().Add(-duration))
	}
	trans.SetStatus(message.CatSuccess)
	trans.Complete()
}

func NewEvent(mtype, name string) message.Messager {
	if !IsEnabled() {
		return &message.NullMessage{}
	}
	return message.NewEvent(mtype, name, manager.flush)
}

func LogEvent(mtype, name string, args ...string) {
	if !IsEnabled() {
		return
	}

	var e = NewEvent(mtype, name)
	if len(args) > 0 {
		e.SetStatus(args[0])
	}
	if len(args) > 1 {
		e.SetData(args[1])
	}
	e.Complete()
}

func LogError(err error, args ...string) {
	if !IsEnabled() {
		return
	}

	var category = "CAT_ERROR"

	if len(args) > 0 {
		category = args[0]
	}

	LogErrorWithCategory(err, category)
}

func LogErrorWithCategory(err error, category string) {
	if !IsEnabled() {
		return
	}

	var event = NewEvent("Error", category)
	var buf = newStacktrace(2, err)
	event.SetStatus(message.CatError)
	event.SetData(buf.String())
	event.Complete()
}

func LogMetricForCount(name string, args ...int) {
	if !IsEnabled() {
		return
	}
	if len(args) == 0 {
		aggregator.metric.AddCount(name, 1)
	} else {
		aggregator.metric.AddCount(name, args[0])
	}
}

func LogMetricForDuration(name string, duration time.Duration) {
	if !IsEnabled() {
		return
	}
	aggregator.metric.AddDuration(name, duration)
}

func NewMetricHelper(name string) MetricHelper {
	if !IsEnabled() {
		return &nullMetricHelper{}
	}
	return newMetricHelper(name)
}
