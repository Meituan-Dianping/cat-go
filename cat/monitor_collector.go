package cat

import (
	"fmt"
	"runtime"
	"strconv"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/load"
)

type Collector interface {
	GetId() string
	GetDesc() string
	GetProperties() map[string]string
}

func b2kbstr(b uint64) string {
	return strconv.Itoa(int(b / 1024))
}

func f642str(b float64) string {
	return fmt.Sprintf("%f", b)
}

type memStatsCollector struct {
	m runtime.MemStats

	alloc,
	mallocs,
	lookups,
	frees uint64
}

func (c *memStatsCollector) GetId() string {
	return "mem.runtime"
}

func (c *memStatsCollector) GetDesc() string {
	return "mem.runtime"
}

func (c *memStatsCollector) GetProperties() map[string]string {
	runtime.ReadMemStats(&c.m)

	m := map[string]string{
		"mem.sys": b2kbstr(c.m.Sys),

		// heap
		"mem.heap.alloc":    b2kbstr(c.m.HeapAlloc),
		"mem.heap.sys":      b2kbstr(c.m.HeapSys),
		"mem.heap.idle":     b2kbstr(c.m.HeapIdle),
		"mem.heap.inuse":    b2kbstr(c.m.HeapInuse),
		"mem.heap.released": b2kbstr(c.m.HeapReleased),
		"mem.heap.objects":  strconv.Itoa(int(c.m.HeapObjects)),

		// stack
		"mem.stack.inuse": b2kbstr(c.m.StackInuse),
		"mem.stack.sys":   b2kbstr(c.m.StackSys),
	}

	if c.alloc > 0 {
		m["mem.alloc"] = b2kbstr(c.m.TotalAlloc - c.alloc)
		m["mem.mallocs"] = strconv.Itoa(int(c.m.Mallocs - c.mallocs))
		m["mem.lookups"] = strconv.Itoa(int(c.m.Lookups - c.lookups))
		m["mem.frees"] = strconv.Itoa(int(c.m.Frees - c.frees))
	}
	c.alloc = c.m.TotalAlloc
	c.mallocs = c.m.Mallocs
	c.lookups = c.m.Lookups
	c.frees = c.m.Frees

	return m
}

type cpuInfoCollector struct {
	lastTime    *cpu.TimesStat
	lastCPUTime float64
}

func (c *cpuInfoCollector) GetId() string {
	return "cpu"
}

func (c *cpuInfoCollector) GetDesc() string {
	return "cpu"
}

func (c *cpuInfoCollector) GetProperties() map[string]string {

	m := make(map[string]string)

	if avg, err := load.Avg(); err == nil {
		m["load.1min"] = f642str(avg.Load1)
		m["load.5min"] = f642str(avg.Load5)
		m["load.15min"] = f642str(avg.Load15)
		m["system.load.average"] = m["load.1min"]
	}

	if times, err := cpu.Times(false); err == nil {
		if len(times) > 0 {
			currentTime := times[0]

			currentCpuTime := 0.0 +
				currentTime.User +
				currentTime.System +
				currentTime.Idle +
				currentTime.Nice +
				currentTime.Iowait +
				currentTime.Irq +
				currentTime.Softirq +
				currentTime.Steal +
				currentTime.Guest +
				currentTime.GuestNice

			if c.lastCPUTime > 0 {
				cpuTime := currentCpuTime - c.lastCPUTime

				if cpuTime > 0.0 {
					user := currentTime.User - c.lastTime.User
					system := currentTime.System - c.lastTime.System
					nice := currentTime.Nice - c.lastTime.Nice
					idle := currentTime.Idle - c.lastTime.Idle
					iowait := currentTime.Iowait - c.lastTime.Iowait
					softirq := currentTime.Softirq - c.lastTime.Softirq
					irq := currentTime.Irq - c.lastTime.Irq
					steal := currentTime.Steal - c.lastTime.Steal

					m["cpu.user"] = f642str(user)
					m["cpu.sys"] = f642str(system)
					m["cpu.nice"] = f642str(nice)
					m["cpu.idle"] = f642str(idle)
					m["cpu.iowait"] = f642str(iowait)
					m["cpu.softirq"] = f642str(softirq)
					m["cpu.irq"] = f642str(irq)
					m["cpu.steal"] = f642str(steal)

					m["cpu.user.percent"] = f642str(user / cpuTime * 100)
					m["cpu.sys.percent"] = f642str(system / cpuTime * 100)
					m["cpu.nice.percent"] = f642str(nice / cpuTime * 100)
					m["cpu.idle.percent"] = f642str(idle / cpuTime * 100)
					m["cpu.iowait.percent"] = f642str(iowait / cpuTime * 100)
					m["cpu.softirq.percent"] = f642str(softirq / cpuTime * 100)
					m["cpu.irq.percent"] = f642str(irq / cpuTime * 100)
					m["cpu.steal.percent"] = f642str(steal / cpuTime * 100)
				}
			}
			c.lastCPUTime = currentCpuTime
			c.lastTime = &currentTime
		}
	}

	// TODO process status
	// if processes, err := process.Processes(); err == nil {
	// 	for _, p := range processes {
	// 	}
	// }

	return m
}
