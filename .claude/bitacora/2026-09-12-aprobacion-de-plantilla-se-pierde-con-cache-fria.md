# Una plantilla aprobada por Meta se quedaba en "pending" para siempre

**Ticket:** sin ticket, pendiente de crear
**Negocio:** 26 (Demo), reproducido en local
**Canal:** WhatsApp (plantillas)

## Resumen

El webhook `message_template_status_update` llegaba bien, se parseaba bien y se
despachaba bien, pero si la clave de cache `whatsapp:waba:<id>` no existia en
Redis el handler cortaba antes de persistir. La plantilla quedaba `pending` en
nuestra base aunque Meta ya la hubiera aprobado.

## Sintoma

Al enviar 3 plantillas a revision (26_hola, 26_me_jodo, 26_dice_que_no), el mock
respondio `APPROVED` para las tres y el log del backend lo confirmaba, pero la
tabla seguia igual:

```
 id | name           | status
 31 | 26_hola        | pending
 32 | 26_me_jodo     | pending
 33 | 26_dice_que_no | pending
```

Log, dos lineas seguidas por plantilla:

```
INF Actualizacion de estado de plantilla recibida event=APPROVED template=26_hola
WRN no se pudo actualizar el estado de la plantilla en cache
    error="key not found: whatsapp:waba:1302830408357767"
```

## Diagnostico

**Hipotesis descartada (costo tiempo):** se penso en una carrera. El resultado
del submit es el que escribe `waba_id` en la fila, y la aprobacion llegaba unos
400 ms despues; parecia razonable que el `UPDATE ... WHERE waba_id = ? AND name = ?`
no encontrara la fila todavia. Se subio el retardo del mock de 400 ms a 3
segundos y **el problema se repitio identico**. La carrera no era.

La pista real estaba en el propio log: la linea de `WRN` es la rama de error de
la cache, y despues de ella no aparecia ninguna linea de "estado de plantilla
actualizado desde el webhook". El handler no seguia.

## Causa raiz

`usecasetemplates/status.go`, en `HandleStatusUpdate`:

```go
if err := u.templatesCache.UpdateStatusByWABA(...); err != nil {
    u.log.Warn(ctx).Err(err)...Msg("no se pudo actualizar el estado de la plantilla en cache")
    return nil            // <-- cortaba aca
}
...
if u.statusPublisher != nil {   // esto nunca se alcanzaba
    u.statusPublisher.PublishTemplateStatus(...)
}
```

El `return nil` estaba dentro de la rama de error. El publisher es el que manda
el estado a `whatsapp.templates.submit.results`, que es lo que el modulo de
notificaciones consume para escribir en la base. Sin cache, no habia escritura.

La cache `whatsapp:waba:<id>` se puebla cuando alguien consulta el estado de
plantillas. Si nadie lo consulto, o si Redis se reinicio, la clave no existe.

## Correccion

Se quita el `return nil`: la cache es un modelo de lectura, su fallo se registra
como advertencia pero no puede impedir la persistencia.

Commit `e1057247`.

## Verificacion

Contra el mock local: se reenviaron las 4 plantillas del flujo 1 y las 4 pasaron
a `approved` con su `waba_id` poblado. Despues la campana con flujo pudo
lanzarse, cosa que antes fallaba con "esa rama no contestaria".

## Impacto en produccion

No es solo del entorno de pruebas. En produccion la cache esta fria despues de
un reinicio de Redis o de un despliegue; cualquier aprobacion que Meta mande en
esa ventana se perdia y la plantilla quedaba inusable, sin error visible para el
negocio y sin forma de recuperarla salvo volver a consultar el estado a mano.

## Pendientes

- **No se reviso produccion.** Puede haber plantillas atascadas en `pending` por
  esta causa. La consulta para encontrarlas: plantillas con `meta_template_id`
  poblado y `status = 'pending'` con `last_synced_at` viejo.
- Crear el ticket correspondiente.
