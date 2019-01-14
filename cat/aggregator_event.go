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
	ch      chan *message.Event
	dataMap map[string]*eventData
}

func (p *eventAggregator) send(dataMap map[string]*eventData) {
	if len(dataMap) == 0 {
		return
	}

	t := message.NewTransaction(typeSystem, nameEventAggregator, aggregator.flush)
	defer t.Complete()

	buf := newBuf()

	for _, data := range dataMap {
		event := t.NewEvent(data.mtype, data.name)
		buf.WriteRune(batchFlag)
		buf.WriteInt(data.count)
		buf.WriteRune(batchSplit)
		buf.WriteInt(data.fail)
		event.SetData(buf.String())
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

func (p *eventAggregator) BackGround() {
	var ticker = time.NewTicker(eventAggregatorInterval)
	for {
		select {
		case event := <-p.ch:
			p.getOrDefault(event).add(event)
		case <-ticker.C:
			dataMap := p.dataMap
			p.dataMap = make(map[string]*eventData)
			go p.send(dataMap)
		}
	}
}

func (p *eventAggregator) Put(event *message.Event) {
	select {
	case p.ch <- event:
	default:
		logger.Warning("Event aggregator is full")
	}
}

func (data *eventData) add(event *message.Event) {
	data.count++

	if event.GetStatus() != CAT_SUCCESS {
		data.fail++
	}
}

func newEventAggregator() *eventAggregator {
	return &eventAggregator{
		ch:      make(chan *message.Event, eventAggregatorChannelCapacity),
		dataMap: make(map[string]*eventData),
	}
}
