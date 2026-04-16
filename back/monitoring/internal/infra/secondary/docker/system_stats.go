package docker

import (
	"bufio"
	"context"
	"os"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/secamc93/probability/back/monitoring/internal/domain/entities"
)

func (d *DockerClient) GetSystemStats(_ context.Context) (*entities.SystemStats, error) {
	stats := &entities.SystemStats{
		CPUCores: runtime.NumCPU(),
	}

	// CPU usage from /proc/stat
	cpuPct, err := getCPUPercent()
	if err == nil {
		stats.CPUPercent = cpuPct
	}

	// Memory from /proc/meminfo
	memTotal, memAvail, err := getMemInfo()
	if err == nil {
		stats.MemoryTotal = memTotal
		stats.MemoryUsed = memTotal - memAvail
		if memTotal > 0 {
			stats.MemoryPercent = float64(stats.MemoryUsed) / float64(memTotal) * 100
		}
	}

	// Disk usage from syscall
	var statfs syscall.Statfs_t
	if err := syscall.Statfs("/", &statfs); err == nil {
		stats.DiskTotal = statfs.Blocks * uint64(statfs.Bsize)
		stats.DiskUsed = (statfs.Blocks - statfs.Bfree) * uint64(statfs.Bsize)
		if statfs.Blocks > 0 {
			stats.DiskPercent = float64(statfs.Blocks-statfs.Bfree) / float64(statfs.Blocks) * 100
		}
	}

	return stats, nil
}

func getCPUPercent() (float64, error) {
	idle1, total1, err := readCPUStat()
	if err != nil {
		return 0, err
	}
	time.Sleep(100 * time.Millisecond)
	idle2, total2, err := readCPUStat()
	if err != nil {
		return 0, err
	}

	idleDelta := float64(idle2 - idle1)
	totalDelta := float64(total2 - total1)
	if totalDelta == 0 {
		return 0, nil
	}
	return (1.0 - idleDelta/totalDelta) * 100, nil
}

func readCPUStat() (idle, total uint64, err error) {
	f, err := os.Open("/proc/stat")
	if err != nil {
		return 0, 0, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	if scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) >= 5 && fields[0] == "cpu" {
			for i, field := range fields[1:] {
				val, _ := strconv.ParseUint(field, 10, 64)
				total += val
				if i == 3 { // idle is 4th field (index 3)
					idle = val
				}
			}
		}
	}
	return idle, total, nil
}

func getMemInfo() (total, available uint64, err error) {
	f, err := os.Open("/proc/meminfo")
	if err != nil {
		return 0, 0, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "MemTotal:") {
			total = parseMemLine(line)
		} else if strings.HasPrefix(line, "MemAvailable:") {
			available = parseMemLine(line)
		}
		if total > 0 && available > 0 {
			break
		}
	}
	return total, available, nil
}

func parseMemLine(line string) uint64 {
	fields := strings.Fields(line)
	if len(fields) >= 2 {
		val, _ := strconv.ParseUint(fields[1], 10, 64)
		return val * 1024 // kB to bytes
	}
	return 0
}
