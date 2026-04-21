# Requisitos del Backend - TopDaysChart

## Cambios Requeridos en `orders_by_date`

El componente **TopDaysChart** fue refactorizado para que **todos los cálculos se hagan en el backend** y optimizar rendimiento en alta concurrencia.

### Estructura Actual (Frontend calcula)
```json
{
  "orders_by_date": [
    { "date": "2026-03-20", "count": 150 },
    { "date": "2026-03-19", "count": 120 }
  ]
}
```

### Estructura Requerida (Backend calcula)
```json
{
  "orders_by_date": [
    {
      "date": "2026-03-20",
      "count": 150,
      "rank": 1,
      "heightPercent": 100,
      "opacity": 1,
      "maxCount": 150,
      "dayName": "viernes",
      "dayShort": "vie",
      "weekdayCount": 3
    },
    {
      "date": "2026-03-19",
      "count": 120,
      "rank": 2,
      "heightPercent": 80,
      "opacity": 0.82,
      "maxCount": 150,
      "dayName": "jueves",
      "dayShort": "jue",
      "weekdayCount": 1
    }
  ]
}
```

## Cálculos que Debe Hacer el Backend

### 1. **Ordenamiento y Top 5**
- Ordenar por `count` descendente
- Retornar solo los 5 primeros

### 2. **heightPercent** (altura proporcional)
```go
heightPercent = (count / maxCount) * 100
```

### 3. **opacity** (opacidad decreciente)
```go
opacityMap := []float64{1, 0.82, 0.66, 0.5, 0.36}
opacity = opacityMap[rank - 1]  // rank es 1-5
```

### 4. **dayName** y **dayShort** (nombres del día)
```go
dayName = date.Format("Monday")  // en español: "viernes", "lunes", etc.
dayShort = date.Format("Mon")     // en español: "vie", "lun", etc.
// Usar tiempo en zona horaria local del negocio
```

### 5. **maxCount**
- El valor máximo del período (el `count` más alto)
- Necesario para cálculos en frontend si falta `heightPercent`

### 6. **weekdayCount** (patrón de días)
- Contar cuántas veces aparece cada día de la semana en el top 5
- Ejemplo: si 3 de los 5 días son viernes → `weekdayCount: 3`
- Solo incluir si es > 1 (para badges)

### 7. **rank** (posición)
- 1 para el primero
- 2 para el segundo
- Hasta 5 para el último

## Ejemplo de Implementación (Go)

```go
type OrderByDate struct {
    Date           string    `json:"date"`
    Count          int       `json:"count"`
    Rank           int       `json:"rank"`
    HeightPercent  float64   `json:"heightPercent"`
    Opacity        float64   `json:"opacity"`
    MaxCount       int       `json:"maxCount"`
    DayName        string    `json:"dayName"`
    DayShort       string    `json:"dayShort"`
    WeekdayCount   int       `json:"weekdayCount,omitempty"`
}

func (r *Repository) GetTopOrdersByDate(ctx context.Context, businessID uint) ([]OrderByDate, error) {
    // 1. Obtener órdenes por fecha, ordenadas descendente
    // 2. Limitar a top 5
    // 3. Iterar y calcular campos enriquecidos
    // 4. Retornar
    
    maxCount := results[0].Count
    opacityMap := []float64{1, 0.82, 0.66, 0.5, 0.36}
    dayCountMap := make(map[string]int)
    
    for i, order := range results {
        // Calcular opacidad
        order.Opacity = opacityMap[i]
        order.Rank = i + 1
        order.HeightPercent = (float64(order.Count) / float64(maxCount)) * 100
        order.MaxCount = maxCount
        
        // Obtener nombre del día
        date := time.Parse("2006-01-02", order.Date)
        order.DayName = date.Format("Monday")  // en locale es-CO
        order.DayShort = date.Format("Mon")    // en locale es-CO
        
        dayCountMap[order.DayName]++
    }
    
    // Agregar weekdayCount si se repite
    for i := range results {
        if count := dayCountMap[results[i].DayName]; count > 1 {
            results[i].WeekdayCount = count
        }
    }
    
    return results, nil
}
```

## Beneficios

✅ **Rendimiento**: Cálculos una sola vez en el servidor  
✅ **Escalabilidad**: Sin overhead en clientes bajo alta concurrencia  
✅ **Lógica centralizada**: Backend es fuente de verdad  
✅ **Frontend agnóstico**: TopDaysChart solo renderiza sin cálculos  
✅ **Fallback seguro**: Frontend funciona si faltan campos (usa valores por defecto)

## Compatibilidad

El frontend es **backward compatible**. Si el backend aún no enriquece los datos, TopDaysChart sigue funcionando usando cálculos locales de fallback.

Sin embargo, **se recomienda implementar esto en el backend cuanto antes** para máximo rendimiento.

## Próximos Pasos

1. Implementar cálculos en el endpoint de `dashboard/stats`
2. Actualizar el modelo `OrdersByDate` en Go
3. Testear con datos reales
4. Monitorear rendimiento en producción
