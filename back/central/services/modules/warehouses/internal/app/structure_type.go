package app

import "strings"

const (
	StructureSimple = "simple"
	StructureZones  = "zones"
	StructureWMS    = "wms"
)

func normalizeStructureType(value string) string {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case StructureZones:
		return StructureZones
	case StructureWMS:
		return StructureWMS
	default:
		return StructureSimple
	}
}

func structureRank(value string) int {
	switch normalizeStructureType(value) {
	case StructureWMS:
		return 2
	case StructureZones:
		return 1
	default:
		return 0
	}
}
