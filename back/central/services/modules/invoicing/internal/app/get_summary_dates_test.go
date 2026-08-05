package app

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func mustTime(t *testing.T, value string) time.Time {
	t.Helper()
	parsed, err := time.Parse("2006-01-02 15:04:05", value)
	require.NoError(t, err)
	return parsed
}

func TestStartOfDay(t *testing.T) {
	got := startOfDay(mustTime(t, "2026-08-05 17:42:31"))

	assert.Equal(t, 2026, got.Year())
	assert.Equal(t, time.August, got.Month())
	assert.Equal(t, 5, got.Day())
	assert.Equal(t, 0, got.Hour())
	assert.Equal(t, 0, got.Minute())
	assert.Equal(t, 0, got.Second())
	assert.Equal(t, 0, got.Nanosecond())
}

func TestEndOfDay(t *testing.T) {
	got := endOfDay(mustTime(t, "2026-08-05 01:02:03"))

	assert.Equal(t, 5, got.Day())
	assert.Equal(t, 23, got.Hour())
	assert.Equal(t, 59, got.Minute())
	assert.Equal(t, 59, got.Second())
	assert.Equal(t, 999999999, got.Nanosecond())
}

func TestStartOfDay_YaEsMedianoche_NoCambia(t *testing.T) {
	original := mustTime(t, "2026-08-05 00:00:00")

	assert.True(t, startOfDay(original).Equal(original))
}

func TestStartOfWeek_LunesEsElInicio(t *testing.T) {
	cases := []struct {
		name    string
		fecha   string
		wantDia int
	}{
		{name: "lunes", fecha: "2026-08-03 10:00:00", wantDia: 3},
		{name: "martes", fecha: "2026-08-04 10:00:00", wantDia: 3},
		{name: "miercoles", fecha: "2026-08-05 10:00:00", wantDia: 3},
		{name: "sabado", fecha: "2026-08-08 10:00:00", wantDia: 3},
		{name: "domingo", fecha: "2026-08-09 10:00:00", wantDia: 3},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := startOfWeek(mustTime(t, tc.fecha))

			assert.Equal(t, time.Monday, got.Weekday())
			assert.Equal(t, tc.wantDia, got.Day())
			assert.Equal(t, 0, got.Hour())
		})
	}
}

func TestStartOfWeek_DomingoNoSaltaALaSemanaSiguiente(t *testing.T) {
	domingo := mustTime(t, "2026-08-09 23:59:00")

	got := startOfWeek(domingo)

	assert.Equal(t, time.Monday, got.Weekday())
	assert.Equal(t, 3, got.Day())
	assert.True(t, got.Before(domingo))
}

func TestEndOfWeek_EsDomingoALas2359(t *testing.T) {
	got := endOfWeek(mustTime(t, "2026-08-05 10:00:00"))

	assert.Equal(t, time.Sunday, got.Weekday())
	assert.Equal(t, 9, got.Day())
	assert.Equal(t, 23, got.Hour())
	assert.Equal(t, 59, got.Minute())
	assert.Equal(t, 59, got.Second())
}

func TestSemanaCompleta_SonSieteDias(t *testing.T) {
	base := mustTime(t, "2026-08-05 10:00:00")

	inicio := startOfWeek(base)
	fin := endOfWeek(base)

	assert.True(t, inicio.Before(fin))
	assert.InDelta(t, 7*24, fin.Sub(inicio).Hours(), 0.01)
}

func TestStartOfMonth(t *testing.T) {
	got := startOfMonth(mustTime(t, "2026-08-19 13:00:00"))

	assert.Equal(t, 1, got.Day())
	assert.Equal(t, time.August, got.Month())
	assert.Equal(t, 0, got.Hour())
}

func TestEndOfMonth(t *testing.T) {
	cases := []struct {
		name    string
		fecha   string
		wantDia int
		wantMes time.Month
	}{
		{name: "agosto tiene 31", fecha: "2026-08-05 10:00:00", wantDia: 31, wantMes: time.August},
		{name: "abril tiene 30", fecha: "2026-04-15 10:00:00", wantDia: 30, wantMes: time.April},
		{name: "febrero no bisiesto", fecha: "2026-02-10 10:00:00", wantDia: 28, wantMes: time.February},
		{name: "febrero bisiesto", fecha: "2028-02-10 10:00:00", wantDia: 29, wantMes: time.February},
		{name: "diciembre", fecha: "2026-12-01 10:00:00", wantDia: 31, wantMes: time.December},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := endOfMonth(mustTime(t, tc.fecha))

			assert.Equal(t, tc.wantDia, got.Day())
			assert.Equal(t, tc.wantMes, got.Month())
			assert.Equal(t, 23, got.Hour())
			assert.Equal(t, 59, got.Minute())
			assert.Equal(t, 59, got.Second())
		})
	}
}

func TestStartOfYear(t *testing.T) {
	got := startOfYear(mustTime(t, "2026-08-05 10:00:00"))

	assert.Equal(t, 2026, got.Year())
	assert.Equal(t, time.January, got.Month())
	assert.Equal(t, 1, got.Day())
	assert.Equal(t, 0, got.Hour())
}

func TestEndOfYear(t *testing.T) {
	got := endOfYear(mustTime(t, "2026-08-05 10:00:00"))

	assert.Equal(t, 2026, got.Year())
	assert.Equal(t, time.December, got.Month())
	assert.Equal(t, 31, got.Day())
	assert.Equal(t, 23, got.Hour())
	assert.Equal(t, 59, got.Second())
}

func TestCalculatePeriod(t *testing.T) {
	ahora := time.Now()

	cases := []struct {
		periodo string
		check   func(t *testing.T, desde, hasta time.Time)
	}{
		{
			periodo: "today",
			check: func(t *testing.T, desde, hasta time.Time) {
				assert.Equal(t, ahora.Day(), desde.Day())
				assert.Equal(t, 0, desde.Hour())
				assert.Equal(t, 23, hasta.Hour())
			},
		},
		{
			periodo: "week",
			check: func(t *testing.T, desde, hasta time.Time) {
				assert.Equal(t, time.Monday, desde.Weekday())
				assert.Equal(t, time.Sunday, hasta.Weekday())
			},
		},
		{
			periodo: "year",
			check: func(t *testing.T, desde, hasta time.Time) {
				assert.Equal(t, time.January, desde.Month())
				assert.Equal(t, 1, desde.Day())
				assert.Equal(t, time.December, hasta.Month())
				assert.Equal(t, 31, hasta.Day())
			},
		},
		{
			periodo: "all",
			check: func(t *testing.T, desde, hasta time.Time) {
				assert.Equal(t, ahora.Year()-10, desde.Year())
				assert.WithinDuration(t, ahora, hasta, 5*time.Second)
			},
		},
		{
			periodo: "month",
			check: func(t *testing.T, desde, hasta time.Time) {
				assert.Equal(t, 1, desde.Day())
				assert.Equal(t, ahora.Month(), desde.Month())
			},
		},
		{
			periodo: "valor-desconocido",
			check: func(t *testing.T, desde, hasta time.Time) {
				assert.Equal(t, 1, desde.Day())
				assert.Equal(t, ahora.Month(), desde.Month())
			},
		},
		{
			periodo: "",
			check: func(t *testing.T, desde, hasta time.Time) {
				assert.Equal(t, 1, desde.Day())
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.periodo, func(t *testing.T) {
			desde, hasta := calculatePeriod(tc.periodo)

			assert.True(t, desde.Before(hasta), "el inicio debe ser anterior al fin")
			tc.check(t, desde, hasta)
		})
	}
}
