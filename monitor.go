package vpanel

import (
	"fmt"
	"github.com/cloudfoundry/gosigar"
	"math"
	"os"
	"time"
)

type DiskStats struct {
	Utilization float64 `json:"utilization"`
	Free        string  `json:"free"`
	Total       string  `json:"total"`
}

type MemoryStats struct {
	Utilization float64 `json:"utilization"`
	Free        string  `json:"free"`
	Total       string  `json:"total"`
}

type HostStats struct {
	Hostname string               `json:"hostname"`
	Cpu      float64              `json:"cpu"`
	Memory   MemoryStats          `json:"memory"`
	Disk     map[string]DiskStats `json:"disk"`
	Uptime   string               `json:"uptime"`
}

type HostMonitor struct {
	interval time.Duration
	readCh   chan chan HostStats
	stopCh   chan bool
	manager  *Manager
	stats    *HostStats
}

func NewHostMonitor(m *Manager) *HostMonitor {
	var interval time.Duration
	Config.Get("hostMonitorInterval", &interval)
	monitor := HostMonitor{
		interval: interval,
		manager:  m,
	}
	return &monitor
}

func round(value float64) float64 {
	return math.Floor(value*100+0.5) / 100
}

func utilize(value, max uint64) float64 {
	return round(100 * float64(value) / float64(max))
}

func humanize(value uint64) string {
	var megabyte = 1024.0
	var gigabyte = megabyte * 1024.0
	var terabyte = gigabyte * 1024.0

	v := float64(value)
	if v >= terabyte {
		return fmt.Sprintf("%.1fT", v/terabyte)
	} else if v >= gigabyte {
		return fmt.Sprintf("%.1fG", v/gigabyte)
	} else if v >= megabyte {
		return fmt.Sprintf("%.1fH", v/megabyte)
	} else {
		return fmt.Sprintf("%.1f", value)
	}
}

func (m *HostMonitor) HostStats() HostStats {
	ch := make(chan HostStats)
	m.readCh <- ch
	return <-ch
}

func (m *HostMonitor) monitorLoop() {
	cpu := sigar.Cpu{}
	cpu.Get()

	memory := sigar.Mem{}
	fslist := sigar.FileSystemList{}
	uptime := sigar.Uptime{}

	ticker := time.NewTicker(m.interval)

	for {
		select {
		case <-ticker.C:
			oldCpu := cpu
			cpu.Get()
			deltaCpu := cpu.Delta(oldCpu)

			memory.Get()
			uptime.Get()
			fslist.Get()

			m.stats.Hostname, _ = os.Hostname()
			m.stats.Cpu = utilize(deltaCpu.Total()-deltaCpu.Idle, deltaCpu.Total())
			m.stats.Memory.Free = humanize(memory.Free)
			m.stats.Memory.Total = humanize(memory.Total)
			m.stats.Memory.Utilization = utilize(memory.Total-memory.Free, memory.Total)
			m.stats.Uptime = uptime.Format()

			m.stats.Disk = make(map[string]DiskStats)
			for _, fs := range fslist.List {
				usage := sigar.FileSystemUsage{}
				usage.Get(fs.DirName)
				m.stats.Disk[fs.DirName] = DiskStats{
					Utilization: utilize(usage.Total-usage.Free, usage.Total),
					Free:        fmt.Sprintf("%d", usage.Free),
					Total:       fmt.Sprintf("%d", usage.Total),
				}
			}
		case ch := <-m.readCh:
			ch <- *m.stats
		case <-m.stopCh:
			ticker.Stop()
			m.stopCh <- true
			return
		}
	}
}

func (m *HostMonitor) Start() {
	m.stats = new(HostStats)
	m.stopCh = make(chan bool)
	m.readCh = make(chan chan HostStats)
	go m.monitorLoop()
}

func (m *HostMonitor) Stop() {
	m.stopCh <- true
	<-m.stopCh
	close(m.stopCh)
}
