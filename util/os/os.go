package ml_os

import (
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"
)

func DiskUsagePercentage() int {
	if runtime.GOOS != "linux" {
		return 0
	}

	cmd := exec.Command("df", "-h", "/")
	output, err := cmd.Output()
	if err != nil {
		return 0
	}

	lines := strings.Split(string(output), "\n")

	if len(lines) > 1 {
		fields := strings.Fields(lines[1])
		if len(fields) > 4 {
			usage := strings.TrimSuffix(fields[4], "%")
			//fmt.Println("Процент использования диска:", usage)
			u, _ := strconv.Atoi(usage)
			return u
		}
	}

	return 0
}

func DiskUsageScheduleNotificator(thr int, checkEach time.Duration, fn func(v int)) {
	for {
		v := DiskUsagePercentage()
		if DiskUsagePercentage() > thr {
			fn(v)
		}
		time.Sleep(checkEach)
	}
}
