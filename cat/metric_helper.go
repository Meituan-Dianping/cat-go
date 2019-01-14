package cat

import (
	"time"
)

type MetricHelper struct {
	name string
	tags map[string]string
}

func newMetricHelper(name string) *MetricHelper {
	return &MetricHelper{
		name: name,
		tags: make(map[string]string),
	}
}

func (h *MetricHelper) AddTag(key, val string) *MetricHelper {
	h.tags[key] = val
	return h
}

func (h *MetricHelper) Count(count int) {
	aggregator.metric.AddCount(h.name, count)
}

func (h *MetricHelper) Duration(duration time.Duration) {
	aggregator.metric.AddDuration(h.name, duration)
}
