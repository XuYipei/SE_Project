package net

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
	"github.com/tuplz/tuplz-be/config"
	"github.com/tuplz/tuplz-be/internal/monitor"
	"github.com/tuplz/tuplz-be/internal/pkg/tools"
)

func FindDbStatus(c *gin.Context) {
	type resp struct {
		Key   string `json:"key"`
		Name  string `json:"name"`
		Size  int    `json:"size"`
		Count int    `json:"count"`
	}

	dbInfo, names, err := monitor.FindDbStatus()
	result := make([]resp, len(dbInfo))
	for i, info := range dbInfo {
		result[i].Count = info.Count
		result[i].Size = info.Size
		result[i].Name = names[i]
		result[i].Key = fmt.Sprintf("%d", i)
	}

	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "fail", "data": result})
	} else {
		c.JSON(http.StatusOK, gin.H{"status": "success", "data": result})
	}
}

func FindDeviceStatus(c *gin.Context) {
	type resp struct {
		Name  string  `json:"name"`
		Usage float64 `json:"usage"`
	}

	cpuUsage, err0 := cpu.Percent(time.Second, false)
	memInfo, err1 := mem.VirtualMemory()
	Parts, err2 := disk.Partitions(true)
	if err0 != nil || err1 != nil || err2 != nil {
		c.JSON(http.StatusOK, gin.H{"status": "fail", "data": make([]resp, 0)})
		return
	}

	result := make([]resp, 2)
	result[0] = resp{Name: "CPU", Usage: cpuUsage[0]}
	result[1] = resp{Name: "mempry", Usage: memInfo.UsedPercent}
	for _, d := range Parts {
		usgPart, err := disk.Usage(d.Device)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"status": "fail", "data": result})
			return
		}
		result = append(result, resp{Name: "partition " + d.Device, Usage: usgPart.UsedPercent})
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": result})
}

func LockSyncDatabase(c *gin.Context) {
	err := tools.LockDb()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "fail"})
	} else {
		c.JSON(http.StatusOK, gin.H{"status": "success"})
	}
}

func StoreDatabase(c *gin.Context) {
	err := tools.LockDb()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "fail"})
		return
	}
	err = tools.SysPwsh("mongostore -d " + config.Database + " -o " + config.DatabaseDir + "--drop")
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "fail"})
		return
	}
	err = tools.UnlockDb()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "fail"})
	} else {
		c.JSON(http.StatusOK, gin.H{"status": "success"})
	}
}

func DumpDatabase(c *gin.Context) {
	err := tools.LockDb()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "fail"})
		return
	}
	err = tools.SysPwsh("mongodump -d " + config.Database + " -o " + config.DatabaseDir)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "fail"})
		return
	}
	err = tools.UnlockDb()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "fail"})
	} else {
		c.JSON(http.StatusOK, gin.H{"status": "success"})
	}
}

func UnlockDatabase(c *gin.Context) {
	err := tools.UnlockDb()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "fail"})
	} else {
		c.JSON(http.StatusOK, gin.H{"status": ""})
	}
}
