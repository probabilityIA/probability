# Plan Variantes WMS - Modulo Products - Estado Actual

> Rama: `feat/inventory-wms-phases` | sesion 2026-04-23 | backend compile OK en `back/central` y `back/migration`

---

## Objetivo

Separar el concepto de catalogo (`producto padre` o `familia`) del concepto operativo de inventario (`SKU` o `variante`) sin romper el CRUD actual de `products`.

La decision aplicada fue:

- `ProductFamily` = producto padre o familia
- `Product` = variante inventariable / SKU real del WMS
- el stock, picking, barcode, precio operativo e integraciones siguen viviendo en `Product`

---

## Ejecutado

## Fase 1 - Base de variantes sin romper el modulo actual

### Dominio

Se extendio el dominio de `products` para soportar variantes:

- `ProductFamily` como entidad nueva
- `Product` ahora soporta:
  - `family_id`
  - `barcode`
  - `variant_label`
  - `variant_attributes`
  - `variant_signature`
- `ProductBusinessIntegration` ahora soporta:
  - `external_variant_id`
  - `external_sku`
  - `external_barcode`

### Reglas de negocio aplicadas

- el `SKU` sigue siendo unico por negocio
- una familia puede tener multiples variantes
- no se permiten dos variantes con la misma combinacion de atributos dentro de la misma familia
- los atributos de variante se normalizan de forma canonica antes de validar

### Persistencia y migraciones

Se agrego:

- tabla `product_families`
- nuevos campos en `products`
- nuevos campos en `product_business_integrations`
- migracion `migrate_product_variants.go`

### API de productos

El CRUD actual de `products` ya acepta:

- `family_id`
- `family`
- `barcode`
- `variant_label`
- `variant_attributes`

Y responde con:

- `family`
- `family_id`
- datos de variante enriquecidos

### Filtros agregados

En listado de productos:

- `family_id`
- `barcode`

### Estado

Completado.

---

## Fase 2 - CRUD separado de familias

Se creo carpeta nueva de casos de uso:

- `back/central/services/modules/products/internal/app/usecasefamily`

Se expusieron rutas nuevas:

- `GET /products/families`
- `GET /products/families/:family_id`
- `POST /products/families`
- `PUT /products/families/:family_id`
- `DELETE /products/families/:family_id`

Capacidades implementadas:

- crear familia
- listar familias
- obtener familia por id
- actualizar familia
- soft delete de familia
- calcular `variant_count` por familia

### Estado

Completado.

---

## Fase 3 - Lectura de variantes por familia

Se agrego soporte para:

- `GET /products/families/:family_id`
  - ahora devuelve la familia y tambien `variants`
- `GET /products/families/:family_id/variants`
  - devuelve solo los SKUs hijos de esa familia

Capacidades implementadas:

- listar variantes por `family_id`
- preload de familia en productos
- respuestas con variantes incluidas

### Estado

Completado.

---

## Archivos clave tocados

### Backend `back/central`

- `back/central/services/modules/products/internal/domain/product.go`
- `back/central/services/modules/products/internal/domain/entities.go`
- `back/central/services/modules/products/internal/domain/ports.go`
- `back/central/services/modules/products/internal/domain/errors.go`
- `back/central/services/modules/products/internal/domain/variant_helpers.go`
- `back/central/services/modules/products/internal/app/usecaseproduct/usecases.go`
- `back/central/services/modules/products/internal/app/usecaseproduct/manage-integrations.go`
- `back/central/services/modules/products/internal/app/usecasefamily/constructor.go`
- `back/central/services/modules/products/internal/app/usecasefamily/usecases.go`
- `back/central/services/modules/products/internal/app/usecases/constructor.go`
- `back/central/services/modules/products/internal/infra/secondary/repository/repository.go`
- `back/central/services/modules/products/internal/infra/secondary/repository/mappers/mapper.go`
- `back/central/services/modules/products/internal/infra/secondary/repository/mappers/integration_mapper.go`
- `back/central/services/modules/products/internal/infra/primary/handlers/create-product.go`
- `back/central/services/modules/products/internal/infra/primary/handlers/update-product.go`
- `back/central/services/modules/products/internal/infra/primary/handlers/list-products.go`
- `back/central/services/modules/products/internal/infra/primary/handlers/family-handlers.go`
- `back/central/services/modules/products/internal/infra/primary/handlers/router.go`

### Backend `back/migration`

- `back/migration/shared/models/product.go`
- `back/migration/shared/models/product_family.go`
- `back/migration/shared/models/product_business_integration.go`
- `back/migration/internal/infra/repository/migrate_product_variants.go`
- `back/migration/internal/infra/repository/constructor.go`

---

## Verificacion ejecutada

Se verifico compilacion con:

```bash
GOCACHE=/tmp/go-build-cache go test ./services/modules/products/...
```

en:

- `back/central`

y con:

```bash
GOCACHE=/tmp/go-build-cache go test ./internal/infra/repository/...
```

en:

- `back/migration`

Resultado:

- compilacion OK
- sin tests unitarios existentes en esos paquetes

---

## Modelo funcional actual

Ejemplo:

- Familia: `Tenis Runner`
- Variantes:
  - `TR-BL-40`
  - `TR-BL-42`
  - `TR-NE-40`

La familia agrupa.
El SKU sigue operando inventario, picking, barcode y sincronizacion.

---

## Falta

## Fase 4 - Resolucion de variantes en integraciones y ordenes

Orden de matching implementado:

1. `external_variant_id`
2. `external_sku`
3. `external_barcode`
4. `sku` interno
5. `barcode` interno
6. `external_product_id`

- `orders` resuelve items por referencias externas antes de caer a SKU
- `GetOrCreateProduct` acepta `variant_id` o `product_id` sin SKU obligatorio
- upsert de mapping externo en `product_business_integrations` al crear producto
- DTO canonico de items soporta `external_barcode`
- VTEX popula `external_barcode` desde `EANID`
- logging de `unmapped variant` cuando no hay identificador resoluble

Pendiente menor:

- matching por `external_product_id + variant_attributes` (canales sin variant_id que tampoco tienen barcode)
- ampliar Shopify/MeLi para extraer barcode del payload si disponible

### Estado

Completado (principal). Items menores pendientes de bajo impacto.

---

## Fase 5 - Escritura explicita de mappings por variante

- `POST /products/:id/integrations` - registrar mapping
- `PUT /products/:id/integrations/:integration_id` - actualizar mapping externo
- `GET /products/:id/integrations` - listar mappings
- `DELETE /products/:id/integrations/:integration_id` - eliminar mapping
- `GET /products/lookup-by-external?integration_id=X&external_variant_id=Y` - resolver producto interno desde refs externas

### Estado

Completado.

---

## Fase 6 - Vista agrupada y operacion WMS

Pendiente mejorar UX operativa:

- listado agrupado por familia con expand de variantes
- filtros por ejes de variante (`color`, `talla`, etc.)
- semaforo por faltantes de variantes
- resumen agregado por familia:
  - total stock
  - variantes agotadas
  - variantes activas

### Estado

Pendiente.

---

## Fase 7 - Reglas operativas mas fuertes

Implementado:

- `DELETE /products/families/:family_id` devuelve 409 si la familia tiene variantes activas
- `HasFamilyActiveVariants` en repositorio y puerto

Pendiente:

- impedir merge accidental de variantes por cambios de atributos (revalidacion de firma)
- soportar cambio de familia de una variante con revalidacion
- auditoria de cambios de atributos

### Estado

Parcialmente completado.

---

## Riesgos abiertos

- hoy `orders` ya resuelve primero por `external_variant_id`, pero no todos los canales exponen aun `external_barcode` ni atributos de variante
- `product_business_integrations` sigue siendo la tabla usada para mapping; puede que en una fase posterior convenga separarla en una tabla mas explicita de variant mappings
- no se agregaron tests de dominio/repository para colisiones de `variant_signature`
- aun no existe endpoint dedicado para administrar mappings externos por variante desde `products`

---

## Recomendacion siguiente

Siguiente implementacion recomendada:

1. agregar endpoint de upsert y consulta de mapping externo por variante
2. soportar matching por `external_product_id + variant_attributes`
3. agregar pruebas para:
   - variantes duplicadas en misma familia
   - variante igual en familias distintas
   - detalle de familia con variantes
   - filtros por `family_id`
   - resolucion de orden por `external_barcode`
