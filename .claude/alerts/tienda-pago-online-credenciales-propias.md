# Tienda publica: pago en linea con credenciales propias del negocio

**Fecha:** 2026-07-30

## Contexto

El checkout de la tienda publica (`publicsite`) firmaba pagos Bold con las
**credenciales de plataforma de Probability** (`GetCachedPlatformCredentials`),
es decir, la plata de las ventas de los negocios entraba a la cuenta Bold de
Probability. Se elimino ese camino: hoy el unico metodo de pago de la tienda es
**"acordar con el vendedor"** (`payment_method: "agree"`), que crea la orden por
el pipeline canonico con pago pendiente (`original_status: "agreed"`), referencia
`SFA-...`.

## Urgente

- Ninguno (el camino inseguro quedo cerrado; `payment_method: "online"` devuelve
  409 `ErrOnlinePayNotReady`).

## Importante (para habilitar pago en linea por negocio)

1. Guardar credenciales Bold propias en la integracion Bold del negocio
   (`integrations.credentials`, pasan por `IEncryptionService` del core de
   integraciones — hay que desencriptar via integration core, no leer crudo).
2. Nueva firma en `pay`: variante de `BoldGenerateSignatureForReference` que use
   las credenciales del negocio (no `GetCachedPlatformCredentials`).
3. Webhook: los eventos de la cuenta Bold del negocio deben llegar a nuestro
   endpoint (configurar webhook por cuenta) y rutearse por referencia `SFO-`.
4. `publicsite`: endpoint `payment-options` (online true solo si el negocio tiene
   credenciales propias) + el front ya envia `payment_method`; reactivar la rama
   online del `CartWidget` (quedo el codigo del widget Bold intacto).

## Criterio de cierre

Un negocio con su propia cuenta Bold configurada puede cobrar en linea en su
tienda y el dinero entra a SU cuenta; los negocios sin pasarela siguen solo con
"acordar con el vendedor".
