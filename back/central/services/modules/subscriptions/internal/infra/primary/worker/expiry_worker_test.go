package worker

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNextRun_AntesDeLasOcho_HoyALasOcho(t *testing.T) {
	loc, _ := time.LoadLocation("America/Bogota")
	now := time.Date(2026, 9, 4, 6, 0, 0, 0, loc)

	got := nextRun(now, loc)

	assert.Equal(t, time.Date(2026, 9, 4, 8, 0, 0, 0, loc), got)
}

func TestNextRun_DespuesDeLasOcho_MananaALasOcho(t *testing.T) {
	loc, _ := time.LoadLocation("America/Bogota")
	now := time.Date(2026, 9, 4, 12, 40, 0, 0, loc)

	got := nextRun(now, loc)

	assert.Equal(t, time.Date(2026, 9, 5, 8, 0, 0, 0, loc), got)
}

func TestNextRun_ExactamenteALasOcho_MananaALasOcho(t *testing.T) {
	loc, _ := time.LoadLocation("America/Bogota")
	now := time.Date(2026, 9, 4, 8, 0, 0, 0, loc)

	got := nextRun(now, loc)

	assert.Equal(t, time.Date(2026, 9, 5, 8, 0, 0, 0, loc), got)
}
