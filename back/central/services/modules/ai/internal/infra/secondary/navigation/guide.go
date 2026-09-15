package navigation

import (
	_ "embed"
	"strings"
)

//go:embed guide.md
var guideSource string

func parseGuide(source string) map[string]string {
	guides := make(map[string]string)
	var key string
	var body strings.Builder

	flush := func() {
		if key != "" {
			guides[key] = strings.TrimSpace(body.String())
		}
		body.Reset()
	}

	for _, line := range strings.Split(source, "\n") {
		if strings.HasPrefix(line, "## ") {
			flush()
			key = strings.TrimSpace(strings.TrimPrefix(line, "## "))
			continue
		}
		body.WriteString(line)
		body.WriteString("\n")
	}
	flush()
	return guides
}
