package cat

import (
	"time"
)

type catMetricHelper struct {
	name string
	tags map[string]string
}

func newMetricHelper(name string) *catMetricHelper {
	return &catMetricHelper{
		name: name,
		tags: make(map[string]string),
	}
}

func (h *catMetricHelper) AddTag(key, val string) *catMetricHelper {
	h.tags[key] = val
	return h
}

func (h *catMetricHelper) Count(count int) {
	aggregator.metric.AddCount(h.name, count)
}

func (h *catMetricHelper) Duration(duration time.Duration) {
	aggregator.metric.AddDuration(h.name, duration)
}
