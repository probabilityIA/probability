package client

import (
	"strings"
)

func parseLinkHeader(header string) string {
	if header == "" {
		return ""
	}
	links := strings.Split(header, ",")
	for _, link := range links {
		parts := strings.Split(link, ";")
		if len(parts) < 2 {
			continue
		}
		if strings.Contains(parts[1], `rel="next"`) {
			url := strings.Trim(parts[0], " <>")
			return url
		}
	}
	return ""
}
