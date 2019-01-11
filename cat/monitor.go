package cat

import (
	"bytes"
	"encoding/xml"
	"runtime"
	"time"

	"github.com/shirou/gopsutil/cpu"

	"../message"
)

type catMonitor struct {
	ch         chan int
	collectors []Collector
}

func sleep2NextMinute() *time.Timer {
	var delta = 60 - time.Now().Second()
	return time.NewTimer(time.Duration(delta) * time.Second)
}

func (m *catMonitor) Background() {
	m.collectAndSend()
	for {
		var timer = sleep2NextMinute()
		select {
		case <-m.ch:
			break
		case <-timer.C:
			m.collectAndSend()
		}
	}
}

func (m *catMonitor) Shutdown() {
	m.ch <- 1
}

func (m *catMonitor) buildXml() *bytes.Buffer {
	type ExtensionDetail struct {
		Id    string `xml:"id,attr"`
		Value string `xml:"value,attr"`
	}

	type Extension struct {
		Id      string `xml:"id,attr"`
		Desc    string `xml:"description"`
		Details []ExtensionDetail `xml:"extensionDetail"`
	}

	type CustomInfo struct {
		Key   string `xml:"key,attr"`
		Value string `xml:"value,attr"`
	}

	type Status struct {
		XMLName     xml.Name `xml:"status"`
		Extensions []Extension `xml:"extension"`
		CustomInfos []CustomInfo `xml:"customInfo"`
	}

	status := Status{
		Extensions: make([]Extension, 0, len(m.collectors)),
		CustomInfos: make([]CustomInfo, 0, 3),
	}

	for _, collector := range m.collectors {
		extension := Extension{
			Id:      collector.GetId(),
			Desc:    collector.GetDesc(),
			Details: make([]ExtensionDetail, 0),
		}

		for k, v := range collector.GetProperties() {
			detail := ExtensionDetail{
				Id:    k,
				Value: v,
			}
			extension.Details = append(extension.Details, detail)
		}
		status.Extensions = append(status.Extensions, extension)
	}

	// add custom information.
	status.CustomInfos = append(status.CustomInfos, CustomInfo{"gocat-version", GOCAT_VERSION})
	status.CustomInfos = append(status.CustomInfos, CustomInfo{"go-version", runtime.Version()})

	buf := bytes.NewBuffer([]byte{})
	encoder := xml.NewEncoder(buf)
	encoder.Indent("", "\t")

	if err := encoder.Encode(status); err != nil {
		buf.Reset()
		buf.WriteString(err.Error())
		return buf
	}
	return buf
}

func send(m message.Messager) {
	manager.flush(m)
}

func (m *catMonitor) collectAndSend() {
	var event *message.Event
	event = &message.Event{
		Message: message.NewMessage("Cat_golang_Client_Version", GOCAT_VERSION, send),
	}
	event.Complete()

	// NOTE type & name is useless while sending a heartbeat
	heartbeat := &message.Heartbeat{
		Message: message.NewMessage("Heartbeat", config.ip, send),
	}
	heartbeat.Complete()
}

var monitor = catMonitor{
	ch: make(chan int),
	collectors: []Collector{
		&MemStatsCollector{},
		&CpuInfoCollector{
			lastTime: &cpu.TimesStat{},
			lastCPUTime: 0,
		},
	},
}

func AddMonitorCollector(collector Collector) {
	monitor.collectors = append(monitor.collectors, collector)
}
