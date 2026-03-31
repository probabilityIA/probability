package handlers

import (
	"bufio"
	"fmt"
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
)

var ansiRegex = regexp.MustCompile(`\x1b\[[0-9;]*[a-zA-Z]`)

func stripANSI(s string) string {
	return ansiRegex.ReplaceAllString(s, "")
}

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
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)

	for scanner.Scan() {
		select {
		case <-c.Request.Context().Done():
			return
		default:
			line := scanner.Text()
			// Docker log lines have 8-byte header, skip it
			if len(line) > 8 {
				line = line[8:]
			}
			// Strip ANSI escape codes
			line = stripANSI(line)
			fmt.Fprintf(c.Writer, "data: %s\n\n", line)
			flusher.Flush()
		}
	}
}
