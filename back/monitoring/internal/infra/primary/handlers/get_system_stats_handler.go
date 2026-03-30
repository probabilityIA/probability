package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *handler) GetSystemStats(c *gin.Context) {
	stats, err := h.useCase.GetSystemStats(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"cpu_percent":    stats.CPUPercent,
		"cpu_cores":      stats.CPUCores,
		"memory_total":   stats.MemoryTotal,
		"memory_used":    stats.MemoryUsed,
		"memory_percent": stats.MemoryPercent,
		"disk_total":     stats.DiskTotal,
		"disk_used":      stats.DiskUsed,
		"disk_percent":   stats.DiskPercent,
	})
}
