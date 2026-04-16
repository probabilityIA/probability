package usecaseintegrations

import "time"

// DateRange representa un rango de fechas para un lote de sincronización.
type DateRange struct {
	Start time.Time
	End   time.Time
}

// SplitDateRange divide un rango de fechas en chunks de chunkDays días.
// Si el rango es menor que chunkDays, retorna un solo chunk.
func SplitDateRange(start, end time.Time, chunkDays int) []DateRange {
	if chunkDays <= 0 {
		chunkDays = 7
	}

	// Normalizar: start al inicio del día, end al fin del día
	start = time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, start.Location())
	end = time.Date(end.Year(), end.Month(), end.Day(), 23, 59, 59, 999999999, end.Location())

	if !start.Before(end) {
		return []DateRange{{Start: start, End: end}}
	}

	var chunks []DateRange
	current := start

	for current.Before(end) {
		chunkEnd := current.AddDate(0, 0, chunkDays)
		// Ajustar al fin del día anterior al siguiente chunk
		chunkEnd = time.Date(chunkEnd.Year(), chunkEnd.Month(), chunkEnd.Day(), 23, 59, 59, 999999999, chunkEnd.Location())

		if chunkEnd.After(end) {
			chunkEnd = end
		}

		chunks = append(chunks, DateRange{
			Start: current,
			End:   chunkEnd,
		})

		// Siguiente chunk empieza al día siguiente
		current = time.Date(chunkEnd.Year(), chunkEnd.Month(), chunkEnd.Day(), 0, 0, 0, 0, chunkEnd.Location()).AddDate(0, 0, 1)
	}

	return chunks
}
