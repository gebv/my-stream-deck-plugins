package main

import (
	"testing"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
)

func zzzTestCPU(t *testing.T) {
	res, _ := cpu.Percent(time.Second, true)
	t.Log("res:", res)
	hec, hpc := coreAverages(res, 2, 8)
	t.Log("res:", hec, hpc)
}
