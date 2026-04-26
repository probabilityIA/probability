package app

var validStatuses = map[string]bool{
	"open":           true,
	"in_review":      true,
	"in_development": true,
	"testing":        true,
	"blocked":        true,
	"resolved":       true,
	"closed":         true,
	"wont_fix":       true,
}

var validPriorities = map[string]bool{
	"low":      true,
	"medium":   true,
	"high":     true,
	"critical": true,
}

var validTypes = map[string]bool{
	"bug":         true,
	"improvement": true,
	"feature":     true,
	"data":        true,
	"integration": true,
	"support":     true,
	"complaint":   true,
	"claim":       true,
	"question":    true,
}

var validSeverities = map[string]bool{
	"low":    true,
	"medium": true,
	"high":   true,
}

var validSources = map[string]bool{
	"internal": true,
	"business": true,
}

var closedStatuses = map[string]bool{
	"resolved": true,
	"closed":   true,
	"wont_fix": true,
}
