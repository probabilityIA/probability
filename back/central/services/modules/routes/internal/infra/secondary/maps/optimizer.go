package maps

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/routes/internal/domain/dtos"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/routes/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/modules/routes/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/env"
	"github.com/secamc93/probability/back/central/shared/log"
)

const computeRoutesURL = "https://routes.googleapis.com/directions/v2:computeRoutes"

const fieldMask = "routes.duration,routes.distanceMeters,routes.optimizedIntermediateWaypointIndex"

type optimizer struct {
	apiKey     string
	httpClient *http.Client
	log        log.ILogger
}

func New(cfg env.IConfig, logger log.ILogger) ports.IRouteOptimizer {
	return &optimizer{
		apiKey:     strings.TrimSpace(cfg.Get("GOOGLE_MAPS_API_KEY")),
		httpClient: &http.Client{Timeout: 20 * time.Second},
		log:        logger.WithModule("routes-optimizer"),
	}
}

func (o *optimizer) IsConfigured() bool {
	return o.apiKey != ""
}

type latLng struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type waypoint struct {
	Location struct {
		LatLng latLng `json:"latLng"`
	} `json:"location"`
}

type computeRoutesRequest struct {
	Origin                whereabouts `json:"origin"`
	Destination           whereabouts `json:"destination"`
	Intermediates         []whereabouts `json:"intermediates,omitempty"`
	TravelMode            string        `json:"travelMode"`
	OptimizeWaypointOrder bool          `json:"optimizeWaypointOrder"`
}

type whereabouts struct {
	Location struct {
		LatLng latLng `json:"latLng"`
	} `json:"location"`
}

type computeRoutesResponse struct {
	Routes []struct {
		DistanceMeters int    `json:"distanceMeters"`
		Duration       string `json:"duration"`
		OptimizedOrder []int  `json:"optimizedIntermediateWaypointIndex"`
	} `json:"routes"`
	Error *struct {
		Message string `json:"message"`
		Status  string `json:"status"`
	} `json:"error"`
}

func point(p dtos.GeoPoint) whereabouts {
	var w whereabouts
	w.Location.LatLng = latLng{Latitude: p.Lat, Longitude: p.Lng}
	return w
}

func (o *optimizer) Optimize(ctx context.Context, origin dtos.GeoPoint, stops []dtos.GeoPoint) (dtos.OptimizedRoute, error) {
	result := dtos.OptimizedRoute{}

	if !o.IsConfigured() {
		return result, domainerrors.ErrOptimizerNotConfigured
	}
	if len(stops) < 2 {
		return result, domainerrors.ErrNotEnoughStops
	}

	body := computeRoutesRequest{
		Origin:                point(origin),
		Destination:           point(stops[len(stops)-1]),
		Intermediates:         make([]whereabouts, 0, len(stops)-1),
		TravelMode:            "DRIVE",
		OptimizeWaypointOrder: true,
	}
	for _, s := range stops[:len(stops)-1] {
		body.Intermediates = append(body.Intermediates, point(s))
	}

	payload, err := json.Marshal(body)
	if err != nil {
		return result, fmt.Errorf("serializando peticion de ruta: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, computeRoutesURL, bytes.NewReader(payload))
	if err != nil {
		return result, fmt.Errorf("armando peticion de ruta: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Goog-Api-Key", o.apiKey)
	req.Header.Set("X-Goog-FieldMask", fieldMask)

	resp, err := o.httpClient.Do(req)
	if err != nil {
		return result, fmt.Errorf("llamando a Routes API: %w", err)
	}
	defer resp.Body.Close()

	raw, _ := io.ReadAll(resp.Body)

	var parsed computeRoutesResponse
	if err := json.Unmarshal(raw, &parsed); err != nil {
		return result, fmt.Errorf("respuesta ilegible de Routes API: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		msg := string(raw)
		if parsed.Error != nil {
			msg = parsed.Error.Status + ": " + parsed.Error.Message
		}
		o.log.Error(ctx).Int("status", resp.StatusCode).Str("respuesta", msg).Msg("Routes API rechazo la peticion")
		return result, fmt.Errorf("Routes API respondio %d: %s", resp.StatusCode, msg)
	}

	if len(parsed.Routes) == 0 {
		return result, domainerrors.ErrNoRouteFound
	}

	route := parsed.Routes[0]

	result.Order = make([]int, 0, len(stops))
	result.Order = append(result.Order, route.OptimizedOrder...)
	result.Order = append(result.Order, len(stops)-1)

	result.DistanceKm = float64(route.DistanceMeters) / 1000
	result.DurationMin = parseDurationMinutes(route.Duration)

	return result, nil
}

func parseDurationMinutes(d string) int {
	seconds, err := strconv.Atoi(strings.TrimSuffix(d, "s"))
	if err != nil {
		return 0
	}
	return (seconds + 59) / 60
}
