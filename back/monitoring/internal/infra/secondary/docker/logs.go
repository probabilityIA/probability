package docker

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/docker/docker/api/types/container"
	"github.com/secamc93/probability/back/monitoring/internal/domain/entities"
)

func (d *DockerClient) GetContainerLogs(ctx context.Context, id string, tail int) ([]entities.LogEntry, error) {
	opts := container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Timestamps: true,
		Tail:       fmt.Sprintf("%d", tail),
	}

	reader, err := d.cli.ContainerLogs(ctx, id, opts)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	return parseLogEntries(reader), nil
}

func (d *DockerClient) StreamContainerLogs(ctx context.Context, id string) (io.ReadCloser, error) {
	opts := container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Timestamps: true,
		Follow:     true,
		Tail:       "50",
	}

	reader, err := d.cli.ContainerLogs(ctx, id, opts)
	if err != nil {
		return nil, err
	}

	return reader, nil
}

func parseLogEntries(reader io.Reader) []entities.LogEntry {
	entries := make([]entities.LogEntry, 0)
	scanner := bufio.NewScanner(reader)
	// Docker logs can have large lines
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)

	for scanner.Scan() {
		line := scanner.Bytes()

		// Docker multiplexed stream: first 8 bytes are header
		// [stream_type(1), 0, 0, 0, size(4)]
		msg := string(line)
		if len(line) > 8 {
			streamType := line[0]
			msg = string(line[8:])

			stream := "stdout"
			if streamType == 2 {
				stream = "stderr"
			}

			// Try to split timestamp from message
			// Format: 2024-01-01T00:00:00.000000000Z message
			timestamp := ""
			if idx := strings.IndexByte(msg, ' '); idx > 20 {
				timestamp = msg[:idx]
				msg = msg[idx+1:]
			}

			entries = append(entries, entities.LogEntry{
				Timestamp: timestamp,
				Stream:    stream,
				Message:   msg,
			})
		} else if len(msg) > 0 {
			entries = append(entries, entities.LogEntry{
				Message: msg,
				Stream:  "stdout",
			})
		}
	}

	return entries
}
