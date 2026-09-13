package entities

import (
	"strconv"
	"strings"
)

var flowBlockedVariableSources = map[string]bool{
	"sender.name":   true,
	"campaign.name": true,
}

func IsFlowBlockedVariable(source string) bool {
	return flowBlockedVariableSources[strings.TrimSpace(source)]
}

func BuildTemplateParameters(template *WhatsappTemplate, candidate SegmentCandidate) []string {
	parameters := make([]string, 0, len(template.Variables))

	for _, variable := range template.Variables {
		value := ResolveCandidateVariable(variable, candidate)
		if value == "" {
			value = variable.Fallback
		}
		if value == "" {
			value = "-"
		}
		parameters = append(parameters, value)
	}

	return parameters
}

func ResolveCandidateVariable(variable TemplateVariable, candidate SegmentCandidate) string {
	switch variable.Source {
	case "customer.first_name":
		return FirstName(candidate.Name)
	case "customer.full_name":
		return strings.TrimSpace(candidate.Name)
	case "customer.days_inactive":
		if candidate.DaysInactive <= 0 {
			return ""
		}
		return strconv.Itoa(candidate.DaysInactive)
	case "customer.total_orders":
		if candidate.TotalOrders <= 0 {
			return ""
		}
		return strconv.Itoa(candidate.TotalOrders)
	case "customer.last_product":
		return strings.TrimSpace(candidate.LastProduct)
	case "business.name":
		return strings.TrimSpace(candidate.BusinessName)
	default:
		return ""
	}
}

func FirstName(name string) string {
	fields := strings.Fields(strings.TrimSpace(name))
	if len(fields) == 0 {
		return ""
	}
	return fields[0]
}
