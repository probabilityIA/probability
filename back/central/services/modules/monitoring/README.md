# Módulo de Monitoreo — Alertas de Servidor

## ¿Qué hace este módulo?

Recibe webhooks de **Grafana Cloud** cuando alguna métrica del servidor supera un umbral crítico, y reenvía la alerta al administrador por **WhatsApp** de forma automática.

Es un **relay puro**: no tiene base de datos ni estado propio. Solo recibe, valida y enruta.

```
Grafana Cloud
     │
     │  POST /api/v1/monitoring/alerts/grafana
     │  Header: X-Grafana-Signature-V2: sha256=<hmac>
     ▼
┌─────────────────────┐
│   WebhookGrafana    │  ← Valida HMAC-SHA256
│     (handler)       │  ← Parsea payload
└─────────┬───────────┘
          │
          ▼
┌─────────────────────┐
│  ProcessGrafanaAlert│  ← Filtra alertas "firing"
│     (use case)      │  ← Mapea alertname → tipo legible
└─────────┬───────────┘
          │
          ▼
   RabbitMQ queue
  "monitoring.alerts"
          │
          ▼
┌─────────────────────┐
│   consumeralert     │  ← Lee credenciales de env vars
│  (whatsapp module)  │  ← Construye template message
└─────────┬───────────┘
          │
          ▼
   WhatsApp API (Meta)
   → +573023406789
```

---

## Alertas configuradas

| alertname en Grafana | Tipo mostrado en mensaje | Umbral | Duración |
|----------------------|--------------------------|--------|----------|
| `RAM_ALTA`           | RAM                      | > 85%  | 2 min    |
| `CPU_ALTO`           | CPU                      | > 80%  | 5 min    |
| `DISCO_ALTO`         | Disco                    | > 80%  | 5 min    |

Si Grafana envía un `alertname` que no está en esa tabla, el módulo lo pasa tal cual como tipo de alerta (fallback transparente).

---

## Endpoint

```
POST /api/v1/monitoring/alerts/grafana
```

- **Sin JWT** — autenticado por firma HMAC-SHA256
- **Header requerido**: `X-Grafana-Signature-V2: sha256=<hex>`
- **Responde siempre 200** para evitar reintentos de Grafana (los errores internos se loguean pero no se propagan)

### Validación de firma

El handler valida la firma usando `GRAFANA_WEBHOOK_SECRET`. Si esa variable está vacía, acepta el request sin validar (modo desarrollo). En producción siempre debe estar configurada.

---

## Variables de entorno

| Variable                 | Requerida | Descripción                                              |
|--------------------------|-----------|----------------------------------------------------------|
| `GRAFANA_WEBHOOK_SECRET` | Sí (prod) | Secreto compartido para validar firma del webhook        |
| `WHATSAPP_PHONE_NUMBER_ID` | Sí      | ID del número de WhatsApp desde el que se envía la alerta |
| `WHATSAPP_TOKEN`         | Sí        | Token de acceso de la WhatsApp Cloud API                 |

---

## Cola RabbitMQ

- **Nombre**: `monitoring.alerts`
- **Durable**: sí
- **Formato del mensaje**:

```json
{
  "alert_type": "RAM",
  "summary": "87.3% - supera umbral de 85%",
  "status": "firing",
  "fired_at": "2026-02-25T10:30:00Z"
}
```

El consumer en el módulo WhatsApp ignora mensajes con `status != "firing"` para evitar spam cuando la alerta se resuelve.

---

## Plantillas de Meta que deben existir

### `alerta_servidor`

Esta es la única plantilla que usa este módulo. Debe estar **aprobada en Meta Business Manager** antes de desplegar en producción.

#### Datos de la plantilla

| Campo        | Valor            |
|--------------|------------------|
| Nombre       | `alerta_servidor` |
| Idioma       | Español (`es`)   |
| Categoría    | `UTILITY`        |
| Tipo         | Solo texto (sin botones, sin header, sin footer) |
| Variables    | 2                |

#### Variables

| Variable | Posición | Descripción                          | Ejemplo          |
|----------|----------|--------------------------------------|------------------|
| `{{1}}`  | Body     | Tipo de alerta (RAM / CPU / Disco)   | `RAM`            |
| `{{2}}`  | Body     | Descripción del valor y umbral       | `87.3% - supera umbral de 85%` |

#### Texto del body sugerido

```
⚠️ Alerta del servidor

Recurso crítico: {{1}}
Detalle: {{2}}

Por favor revisa el estado del servidor a la brevedad.
```

> El texto exacto lo decide quien configura la plantilla en Meta. Lo único que el código controla es que `{{1}}` = tipo de alerta y `{{2}}` = descripción/valor.

#### Cómo crearla en Meta Business Manager

1. Ir a **Meta Business Manager** → **WhatsApp** → **Message Templates**
2. Clic en **Create Template**
3. Configurar:
   - **Category**: Utility
   - **Name**: `alerta_servidor`
   - **Language**: Spanish
4. En el Body, escribir el texto con `{{1}}` y `{{2}}`
5. Enviar a revisión
6. Esperar aprobación (puede tardar horas o días)

> La plantilla NO puede desplegarse hasta que Meta la apruebe. Crear en Meta antes de desplegar el código en producción.

---

## Configuración en Grafana Cloud (post-despliegue)

### Reglas de alerta (Alerting → Alert Rules → New Rule)

**RAM Alta**
```
100 * container_memory_usage_bytes{name="central_reserve_prod"} / 1879048192
Condition: IS ABOVE 85
For: 2m
Labels: alertname=RAM_ALTA
Annotations: summary=<valor>% - supera umbral de 85%
```

**CPU Alto**
```
100 * rate(container_cpu_usage_seconds_total{name="central_reserve_prod"}[5m])
Condition: IS ABOVE 80
For: 5m
Labels: alertname=CPU_ALTO
Annotations: summary=<valor>% - supera umbral de 80%
```

**Disco Alto**
```
100 * container_fs_usage_bytes{id="/"} / container_fs_limit_bytes{id="/"}
Condition: IS ABOVE 80
For: 5m
Labels: alertname=DISCO_ALTO
Annotations: summary=<valor>% - supera umbral de 80%
```

### Contact Point (Alerting → Contact Points → New)

- **Type**: Webhook
- **URL**: `https://www.probabilityia.com.co/api/v1/monitoring/alerts/grafana`
- **Optional webhook secret**: mismo valor que `GRAFANA_WEBHOOK_SECRET` en el `.env`

### Notification Policy

Apuntar la política por defecto (o crear una nueva) al contact point del webhook.

---

## Verificación end-to-end

```bash
# Sin secreto (dev mode)
curl -X POST https://www.probabilityia.com.co/api/v1/monitoring/alerts/grafana \
  -H "Content-Type: application/json" \
  -d '{
    "status": "firing",
    "title": "[FIRING] RAM_ALTA",
    "alerts": [{
      "status": "firing",
      "labels": {"alertname": "RAM_ALTA"},
      "annotations": {"summary": "87.3% - supera umbral de 85%"},
      "startsAt": "2026-02-25T10:30:00Z",
      "valueString": ""
    }]
  }'
# Esperado: {"status":"received"}

# Con HMAC (producción)
SECRET="tu_secreto"
BODY='{"status":"firing","title":"[FIRING] RAM_ALTA","alerts":[{"status":"firing","labels":{"alertname":"RAM_ALTA"},"annotations":{"summary":"87.3% - supera umbral de 85%"},"startsAt":"2026-02-25T10:30:00Z","valueString":""}]}'
SIG=$(echo -n "$BODY" | openssl dgst -sha256 -hmac "$SECRET" | awk '{print $2}')
curl -X POST https://www.probabilityia.com.co/api/v1/monitoring/alerts/grafana \
  -H "Content-Type: application/json" \
  -H "X-Grafana-Signature-V2: sha256=$SIG" \
  -d "$BODY"
# Esperado: {"status":"received"} y WhatsApp llega a +573023406789
```
