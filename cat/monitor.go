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
	signalsMixin
	collectors []Collector
}

func (m *catMonitor) GetName() string {
	return "Monitor"
}

func sleep2NextMinute() *time.Timer {
	var delta = 60 - time.Now().Second()
	return time.NewTimer(time.Duration(delta) * time.Second)
}

func (m *catMonitor) Background() {
	// Report monitor info at very beginning.
	m.collectAndSend()

	Instance().LogEvent(typeSystem, nameReboot)

	for m.isAlive {
		var timer = sleep2NextMinute()
		select {
		case signal := <-m.signals:
			if signal == signalShutdown {
				timer.Stop()
				m.stop()
			}
		case <-timer.C:
			m.collectAndSend()
		}
	}

	m.exit()
}

func (m *catMonitor) buildXml() *bytes.Buffer {
	type ExtensionDetail struct {
		Id    string `xml:"id,attr"`
		Value string `xml:"value,attr"`
	}

	type Extension struct {
		Id      string            `xml:"id,attr"`
		Desc    string            `xml:"description"`
		Details []ExtensionDetail `xml:"extensionDetail"`
	}

	type CustomInfo struct {
		Key   string `xml:"key,attr"`
		Value string `xml:"value,attr"`
	}

	type Status struct {
		XMLName     xml.Name     `xml:"status"`
		Extensions  []Extension  `xml:"extension"`
		CustomInfos []CustomInfo `xml:"customInfo"`
	}

	status := Status{
		Extensions:  make([]Extension, 0, len(m.collectors)),
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
	status.CustomInfos = append(status.CustomInfos, CustomInfo{"gocat-version", GoCatVersion})
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

func (m *catMonitor) collectAndSend() {
	var trans = message.NewTransaction(typeSystem, "Status", manager.flush)
	defer trans.Complete()

	trans.LogEvent("Cat_golang_Client_Version", GoCatVersion)

	// NOTE type & name is useless while sending a heartbeat
	heartbeat := message.NewHeartbeat("Heartbeat", config.ip, nil)
	heartbeat.SetData(m.buildXml().String())
	heartbeat.Complete()

	trans.AddChild(heartbeat)
}

var monitor = catMonitor{
	signalsMixin: makeSignalsMixedIn(signalMonitorExit),
	collectors: []Collector{
		&MemStatsCollector{},
		&CpuInfoCollector{
			lastTime:    &cpu.TimesStat{},
			lastCPUTime: 0,
		},
	},
}

func AddMonitorCollector(collector Collector) {
	monitor.collectors = append(monitor.collectors, collector)
}
