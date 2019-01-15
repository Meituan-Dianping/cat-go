package cat

import (
	"time"

	"../message"
)

func NewTransaction(mtype, name string) message.Transactor {
	return message.NewTransaction(mtype, name, manager.flush)
}

func NewCompletedTransactionWithDuration(mtype, name string, duration time.Duration) {
	var trans = NewTransaction(mtype, name)
	trans.SetDuration(duration)
	if duration > 0 && duration < 60*time.Millisecond {
		trans.SetTime(time.Now().Add(-duration))
	}
	trans.SetStatus(message.CatSuccess)
	trans.Complete()
}

func NewEvent(mtype, name string) message.Messager {
	return message.NewEvent(mtype, name, manager.flush)
}

func LogEvent(mtype, name string, args ...string) {
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
	var category = "CAT_ERROR"

	if len(args) > 0 {
		category = args[0]
	}

	LogErrorWithCategory(err, category)
}

func LogErrorWithCategory(err error, category string) {
	var event = NewEvent("Error", category)
	var buf = newStacktrace(2, err)
	event.SetStatus(message.CatError)
	event.SetData(buf.String())
	event.Complete()
}

func LogMetricForCount(name string, args ...int) {
	var count int
	if len(args) == 0 {
		count = 1
	} else {
		count = args[0]
	}
	aggregator.metric.AddCount(name, count)
}

func LogMetricForDuration(name string, duration time.Duration) {
	aggregator.metric.AddDuration(name, duration)
}

func NewMetricHelper(name string) *catMetricHelper {
	return newMetricHelper(name)
}
