package vpanel

import (
	"fmt"
	"github.com/cloudfoundry/gosigar"
	"math"
	"os"
	"time"
)

type host struct{}

var Host host

type MemoryStats struct {
	Utilization float64 `json:"utilization"`
	Free        string  `json:"free"`
	Total       string  `json:"total"`
}

type DiskStats struct {
	Utilization float64 `json:"utilization"`
	Free        string  `json:"free"`
	Total       string  `json:"total"`
}

type HostStats struct {
	Hostname string      `json:"hostname"`
	Cpu      float64     `json:"cpu"`
	Memory   MemoryStats `json:"memory"`
	Disk     DiskStats   `json:"disk"`
	Uptime   string      `json:"uptime"`
}

var hostStats HostStats

var megabyte = 1024.0
var gigabyte = megabyte * 1024.0
var terabyte = gigabyte * 1024.0

func round(value float64) float64 {
	return math.Floor(value*100+0.5) / 100
}

func utilize(value, max uint64) float64 {
	return round(100 * float64(value) / float64(max))
}

func humanize(value uint64) string {
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

func (h host) Stats() (HostStats, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return hostStats, err
	}
	hostStats.Hostname = hostname

	concreteSigar := sigar.ConcreteSigar{}

	c, stop := concreteSigar.CollectCpuStats(500 * time.Millisecond)
	cpu := <-c
	cpu = <-c
	stop <- struct{}{}
	hostStats.Cpu = utilize(cpu.Total()-cpu.Idle, cpu.Total())

	mem := sigar.Mem{}
	mem.Get()
	hostStats.Memory.Free = humanize(mem.Free)
	hostStats.Memory.Total = humanize(mem.Total)
	hostStats.Memory.Utilization = utilize(mem.Total-mem.Free, mem.Total)

	fslist := sigar.FileSystemList{}
	fslist.Get()
	var diskTotal, diskFree uint64
	for _, fs := range fslist.List {
		dir_name := fs.DirName

		usage := sigar.FileSystemUsage{}
		usage.Get(dir_name)

		diskTotal += usage.Total
		diskFree += usage.Free
	}
	hostStats.Disk.Utilization = utilize(diskTotal-diskFree, diskTotal)
	hostStats.Disk.Free = humanize(diskFree)
	hostStats.Disk.Total = humanize(diskTotal)

	uptime := sigar.Uptime{}
	uptime.Get()
	hostStats.Uptime = uptime.Format()

	return hostStats, nil
}
