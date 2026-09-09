package mapimage

import (
	"bytes"
	"context"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/env"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/storage"
)

const staticMapURL = "https://maps.googleapis.com/maps/api/staticmap"

const carpeta = "mapas-direccion"

type generator struct {
	apiKey     string
	storage    storage.IS3Service
	httpClient *http.Client
	log        log.ILogger
}

func New(cfg env.IConfig, s3 storage.IS3Service, logger log.ILogger) ports.IMapImageGenerator {
	return &generator{
		apiKey:     strings.TrimSpace(cfg.Get("GOOGLE_MAPS_API_KEY")),
		storage:    s3,
		httpClient: &http.Client{Timeout: 20 * time.Second},
		log:        logger.WithModule("whatsapp-mapimage"),
	}
}

func (g *generator) IsConfigured() bool {
	return g.apiKey != "" && g.storage != nil
}

func (g *generator) BuildAddressMap(ctx context.Context, lat, lng float64, referencia string) (string, error) {
	if !g.IsConfigured() {
		return "", ports.ErrMapImageNotConfigured
	}

	nombre := nombreArchivo(lat, lng)

	existe, err := g.storage.ImageExists(ctx, carpeta+"/"+nombre)
	if err == nil && existe {
		return g.storage.GetImageURL(carpeta + "/" + nombre), nil
	}

	png, err := g.descargar(ctx, lat, lng)
	if err != nil {
		return "", err
	}

	urlPublica, err := g.storage.UploadFile(ctx, bytes.NewReader(png), carpeta+"/"+nombre)
	if err != nil {
		return "", fmt.Errorf("subiendo el mapa a S3: %w", err)
	}

	g.log.Info(ctx).
		Str("referencia", referencia).
		Str("archivo", nombre).
		Msg("mapa de la direccion generado")

	return urlPublica, nil
}

func (g *generator) descargar(ctx context.Context, lat, lng float64) ([]byte, error) {
	centro := fmt.Sprintf("%.6f,%.6f", lat, lng)

	params := url.Values{}
	params.Set("size", "640x400")
	params.Set("zoom", "16")
	params.Set("scale", "2")
	params.Set("maptype", "roadmap")
	params.Set("language", "es")
	params.Set("region", "co")
	params.Set("center", centro)
	params.Set("markers", "color:red|"+centro)
	params.Set("key", g.apiKey)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, staticMapURL+"?"+params.Encode(), nil)
	if err != nil {
		return nil, fmt.Errorf("armando peticion del mapa: %w", err)
	}

	resp, err := g.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("descargando el mapa: %w", err)
	}
	defer resp.Body.Close()

	cuerpo, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("leyendo el mapa: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("Static Maps respondio %d: %s", resp.StatusCode, string(cuerpo))
	}

	return cuerpo, nil
}

func nombreArchivo(lat, lng float64) string {
	suma := sha1.Sum([]byte(fmt.Sprintf("%.6f,%.6f", lat, lng)))
	return hex.EncodeToString(suma[:]) + ".png"
}
