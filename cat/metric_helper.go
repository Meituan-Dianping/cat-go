package cat

import (
	"time"
)

type MetricHelper interface {
	AddTag(key, val string) MetricHelper
	Count(int)
	Duration(time.Duration)
}

type catMetricHelper struct {
	name string
	tags map[string]string
}

type nullMetricHelper struct {
}

func (h *nullMetricHelper) AddTag(key, val string) MetricHelper {
	return h
}

func (h *nullMetricHelper) Count(count int) {
}

func (h *nullMetricHelper) Duration(duration time.Duration) {
}

func newMetricHelper(name string) MetricHelper {
	return &catMetricHelper{
		name: name,
		tags: make(map[string]string),
	}
}

func (h *catMetricHelper) AddTag(key, val string) MetricHelper {
	h.tags[key] = val
	return h
}

func (h *catMetricHelper) Count(count int) {
	aggregator.metric.AddCount(h.name, count)
}

func (h *catMetricHelper) Duration(duration time.Duration) {
	aggregator.metric.AddDuration(h.name, duration)
}
