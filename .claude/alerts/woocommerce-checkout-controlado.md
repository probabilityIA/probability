# WooCommerce Checkout Controlado - Plan por Fases

Fecha: 2026-07-06
Contexto: hacer que el checkout de WooCommerce use nuestros datos (DANE, geocode Google,
logos de transportadoras) para garantizar direccion correcta, independiente del theme del
cliente. Se prueba primero con nuestro WooCommerce (integracion 197, business 26).

Todos los endpoints publicos usan el mismo token por integracion (X-Probability-Token) y el
rate limiter existente. La Google API key nunca se expone (geocode corre en backend).

## Fases

- [x] FASE 1 - Logos de transportadoras en el checkout (2026-07-06)
  - Backend: `logo_url` en meta_data de cada tarifa, servido desde nuestro backend
    (endpoint `/woocommerce/carrier-logo/:carrier`) para evitar mixed-content.
  - Plugin: filtro `woocommerce_cart_shipping_method_full_label` que pinta el `<img>`.

- [x] FASE 2 - Endpoints DANE (municipios) autenticados (2026-07-06)
  - `GET /woocommerce/dane/:integration_id/states`
  - `GET /woocommerce/dane/:integration_id/cities?state=<code>`
  - Datos de geozones (33 states / 1121 cities, todos con code DANE + centroid).

- [x] FASE 3 - Validacion de direccion (token) + geocoder enriquecido (2026-07-06)
  - Geocoder en shipments (replicado) que devuelve lat/lng + locality + adminArea2 +
    location_type + partial_match.
  - `POST /woocommerce/validate-address/:integration_id` -> {lat,lng,municipality,dane,
    department,confidence}. Reverse-map a DANE contra geozones.

- [x] FASE 4 - Plugin: checkout controlado, plugin v1.2.0 (2026-07-06)
  - JS: ciudad como dropdown dependiente del departamento (Fase 2); validar direccion
    (Fase 3) y setear hidden dane/lat/lng; feedback de confianza.

- [ ] FASE 5 - Robustez en la orden (follow-up, toca modulo orders)
  - Enriquecer geocode de orden con DANE + confianza; flag en panel si confianza baja.

## Criterio de cierre
Fases 1-4 desplegadas y probadas contra nuestro WooCommerce. Fase 5 puede ir despues.
