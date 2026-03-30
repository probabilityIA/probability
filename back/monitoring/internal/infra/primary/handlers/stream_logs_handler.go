package handlers

import (
	"bufio"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *handler) StreamLogs(c *gin.Context) {
	id := c.Param("id")

	reader, err := h.useCase.StreamContainerLogs(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer reader.Close()

	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("X-Accel-Buffering", "no")

	flusher, ok := c.Writer.(http.Flusher)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "streaming not supported"})
		return
	}

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		select {
		case <-c.Request.Context().Done():
			return
		default:
			line := scanner.Text()
			// Docker log lines have 8-byte header, skip it if present
			if len(line) > 8 {
				line = line[8:]
			}
			fmt.Fprintf(c.Writer, "data: %s\n\n", line)
			flusher.Flush()
		}
	}
}
