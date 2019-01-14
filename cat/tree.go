package cat

import (
	"time"

	"../message"
)

type catMessageTree struct {
}


func Instance() *catMessageTree {
	return &catMessageTree{}
}

func (t *catMessageTree) NewTransaction(mtype, name string) message.Transactor {
	return message.NewTransaction(mtype, name, manager.flush)
}

func (t *catMessageTree) NewCompletedTransactionWithDuration(mtype, name string, duration time.Duration) {
	var trans = t.NewTransaction(mtype, name)
	trans.SetDuration(duration)
	if duration > 0 && duration < 60 * time.Millisecond {
		trans.SetTime(time.Now().Add(-duration))
	}
	trans.SetStatus(message.CAT_SUCCESS)
	trans.Complete()
}

func (t *catMessageTree) NewEvent(mtype, name string) message.Messager {
	return message.NewEvent(mtype, name, manager.flush)
}

func (t *catMessageTree) LogEvent(mtype, name string, args ...string) {
	var e = t.NewEvent(mtype, name)
	if len(args) > 0 {
		e.SetStatus(args[0])
	}
	if len(args) > 1 {
		e.SetData(args[1])
	}
	e.Complete()
}

func (t *catMessageTree) LogError(err error, args ...string) {
	var category = "CAT_ERROR"

	if len(args) > 0 {
		category = args[0]
	}

	t.LogErrorWithCategory(err, category)
}

func (t *catMessageTree) LogErrorWithCategory(err error, category string) {
	var e = t.NewEvent("Error", category)
	var buf = newStacktrace(2, err)
	e.SetStatus(message.CAT_ERROR)
	e.SetData(buf.String())
	e.Complete()
}

func (t *catMessageTree) LogMetricForCount(name string, args ...int) {
	var count int
	if len(args) == 0 {
		count = 1
	} else {
		count = args[0]
	}
	aggregator.metric.AddCount(name, count)
}

func (t *catMessageTree) LogMetricForDuration(name string, duration time.Duration) {
	aggregator.metric.AddDuration(name, duration)
}

func (t *catMessageTree) NewMetricHelper(name string) *MetricHelper {
	return newMetricHelper(name)
}
