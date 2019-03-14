package cat

import (
	"fmt"
	"strconv"
	"time"

	"github.com/Meituan-Dianping/cat-go/message"
)

type metricData struct {
	name     string
	count    int
	duration time.Duration
}

type metricAggregator struct {
	scheduleMixin
	ch      chan *metricData
	dataMap map[string]*metricData
	ticker  *time.Ticker
}

func (p *metricAggregator) GetName() string {
	return "MetricAggregator"
}

func (p *metricAggregator) afterStart() {
	p.ticker = time.NewTicker(metricAggregatorInterval)
}

func (p *metricAggregator) beforeStop() {
	close(p.ch)

	for event := range p.ch {
		p.putOrMerge(event)
	}
	p.collectAndSend()

	p.ticker.Stop()
}

func (p *metricAggregator) process() {
	select {
	case sig := <-p.signals:
		p.handle(sig)
	case data := <-p.ch:
		p.putOrMerge(data)
	case <-p.ticker.C:
		p.collectAndSend()
	}
}

func (p *metricAggregator) collectAndSend() {
	dataMap := p.dataMap
	p.dataMap = make(map[string]*metricData)
	p.send(dataMap)
}

func (p *metricAggregator) send(dataMap map[string]*metricData) {
	if len(dataMap) == 0 {
		return
	}

	t := message.NewTransaction(typeSystem, nameMetricAggregator, aggregator.flush)
	defer t.Complete()

	for _, data := range dataMap {
		metric := message.NewMetric("", data.name, nil)

		if data.duration > 0 {
			metric.SetStatus("S,C")
			duration := data.duration.Nanoseconds() / time.Millisecond.Nanoseconds()
			metric.SetData(fmt.Sprintf("%d,%d", data.count, duration))
		} else {
			metric.SetStatus("C")
			metric.SetData(strconv.Itoa(data.count))
		}

		t.AddChild(metric)
	}
}

func (p *metricAggregator) putOrMerge(data *metricData) {
	if item, ok := p.dataMap[data.name]; ok {
		item.count += data.count
		item.duration += data.duration
	} else {
		p.dataMap[data.name] = data
	}
}

func newMetricAggregator() *metricAggregator {
	return &metricAggregator{
		scheduleMixin: makeScheduleMixedIn(signalMetricAggregatorExit),
		ch:            make(chan *metricData, metricAggregatorChannelCapacity),
		dataMap:       make(map[string]*metricData),
	}
}

func (p *metricAggregator) AddDuration(name string, duration time.Duration) {
	select {
	case p.ch <- &metricData{
		name:     name,
		count:    1,
		duration: duration,
	}:
	default:
		logger.Warning("Metric aggregator is full")
	}
}

func (p *metricAggregator) AddCount(name string, count int) {
	select {
	case p.ch <- &metricData{
		name:     name,
		count:    count,
		duration: 0,
	}:
	default:
		logger.Warning("Metric aggregator is full")
	}
}
