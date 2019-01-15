package cat

import (
	"fmt"
	"time"

	"../message"
)

type eventData struct {
	mtype, name string

	count, fail int
}

type eventAggregator struct {
	scheduleMixin
	ch      chan *message.Event
	dataMap map[string]*eventData
	ticker *time.Ticker
}

func (p *eventAggregator) GetName() string {
	return "EventAggregator"
}

func (p *eventAggregator) collectAndSend() {
	dataMap := p.dataMap
	p.dataMap = make(map[string]*eventData)
	p.send(dataMap)
}

func (p *eventAggregator) send(dataMap map[string]*eventData) {
	if len(dataMap) == 0 {
		return
	}

	t := message.NewTransaction(typeSystem, nameEventAggregator, aggregator.flush)
	defer t.Complete()

	for _, data := range dataMap {
		event := t.NewEvent(data.mtype, data.name)
		event.SetData(fmt.Sprintf("%c%d%c%d", batchFlag, data.count, batchSplit, data.fail))
	}
}

func (p *eventAggregator) getOrDefault(event *message.Event) *eventData {
	key := fmt.Sprintf("%s,%s", event.Type, event.Name)

	if data, ok := p.dataMap[key]; ok {
		return data
	} else {
		p.dataMap[key] = &eventData{
			mtype: event.Type,
			name:  event.Name,
			count: 0,
			fail:  0,
		}
		return p.dataMap[key]
	}
}

func (p *eventAggregator) afterStart() {
	p.ticker = time.NewTicker(eventAggregatorInterval)
}

func (p *eventAggregator) beforeStop() {
	close(p.ch)

	for event := range p.ch {
		p.getOrDefault(event).add(event)
	}
	p.collectAndSend()

	p.ticker.Stop()
}

func (p *eventAggregator) process() {
	select {
	case sig := <-p.signals:
		p.handle(sig)
	case event := <-p.ch:
		p.getOrDefault(event).add(event)
	case <-p.ticker.C:
		p.collectAndSend()
	}
}

func (p *eventAggregator) Put(event *message.Event) {
	if !IsEnabled() {
		return
	}

	select {
	case p.ch <- event:
	default:
		logger.Warning("Event aggregator is full")
	}
}

func (data *eventData) add(event *message.Event) {
	data.count++

	if event.GetStatus() != SUCCESS {
		data.fail++
	}
}

func newEventAggregator() *eventAggregator {
	return &eventAggregator{
		scheduleMixin: makeScheduleMixedIn(signalEventAggregatorExit),
		ch:            make(chan *message.Event, eventAggregatorChannelCapacity),
		dataMap:       make(map[string]*eventData),
	}
}
