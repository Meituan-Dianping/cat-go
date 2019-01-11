package cat

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
	// TODO check if int is over than i64
}

func (h *MetricHelper) DurationMs(duration int) {
	// TODO check if int is over than i64
}
