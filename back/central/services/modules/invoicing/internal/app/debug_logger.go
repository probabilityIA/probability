package app

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// DebugLogger escribe logs detallados a archivo para debugging
type DebugLogger struct {
	filePath string
}

// NewDebugLogger crea un nuevo logger de debugging
func NewDebugLogger() *DebugLogger {
	// Obtener la ruta base del proyecto
	basePath := "/home/cam/Desktop/probability/back/central/log"

	// Crear carpeta si no existe
	os.MkdirAll(basePath, 0755)

	// Archivo con timestamp
	timestamp := time.Now().Format("2006-01-02")
	filePath := filepath.Join(basePath, fmt.Sprintf("invoicing-config-%s.log", timestamp))

	return &DebugLogger{filePath: filePath}
}

// LogConfigValidation registra la validación de configuración
func (d *DebugLogger) LogConfigValidation(orderID string, integrationID uint, config interface{}, enabled bool, autoInvoice bool) {
	d.writeLog(map[string]interface{}{
		"timestamp":      time.Now().Format(time.RFC3339),
		"event":          "config_validation",
		"order_id":       orderID,
		"integration_id": integrationID,
		"config":         config,
		"enabled":        enabled,
		"auto_invoice":   autoInvoice,
	})
}

// LogFilterValidation registra la validación de filtros
func (d *DebugLogger) LogFilterValidation(orderID string, filterName string, filterValue interface{}, orderValue interface{}, passed bool, errorMsg string) {
	d.writeLog(map[string]interface{}{
		"timestamp":    time.Now().Format(time.RFC3339),
		"event":        "filter_validation",
		"order_id":     orderID,
		"filter_name":  filterName,
		"filter_value": filterValue,
		"order_value":  orderValue,
		"passed":       passed,
		"error":        errorMsg,
	})
}

// LogFilterDetails registra los detalles de un filtro específico
func (d *DebugLogger) LogFilterDetails(orderID string, filterType string, details map[string]interface{}) {
	logData := map[string]interface{}{
		"timestamp":   time.Now().Format(time.RFC3339),
		"event":       "filter_details",
		"order_id":    orderID,
		"filter_type": filterType,
	}
	for k, v := range details {
		logData[k] = v
	}
	d.writeLog(logData)
}

// LogOrderData registra los datos completos de la orden
func (d *DebugLogger) LogOrderData(orderID string, order interface{}) {
	d.writeLog(map[string]interface{}{
		"timestamp": time.Now().Format(time.RFC3339),
		"event":     "order_data",
		"order_id":  orderID,
		"order":     order,
	})
}

// LogDecision registra la decisión final de facturación
func (d *DebugLogger) LogDecision(orderID string, shouldInvoice bool, reason string) {
	d.writeLog(map[string]interface{}{
		"timestamp":      time.Now().Format(time.RFC3339),
		"event":          "invoicing_decision",
		"order_id":       orderID,
		"should_invoice": shouldInvoice,
		"reason":         reason,
	})
}

// writeLog escribe un registro al archivo
func (d *DebugLogger) writeLog(data map[string]interface{}) {
	file, err := os.OpenFile(d.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Failed to open debug log file: %v\n", err)
		return
	}
	defer file.Close()

	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Printf("Failed to marshal log data: %v\n", err)
		return
	}

	file.WriteString(string(jsonData) + "\n")
}
