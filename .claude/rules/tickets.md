# Tickets - el registro de todo bug, arreglo y desarrollo

Desde 2026-09-03 **todo trabajo de bug, correccion o desarrollo se registra en un
ticket** del modulo de tickets. El ticket es la fuente de verdad de que se hizo y
por que; la bitacora y las alertas siguen existiendo, pero cuelgan de el.

## Reglas

1. **Antes de empezar, busca el ticket.** Si el usuario reporta algo, revisa si
   ya existe uno abierto sobre el tema (`GET /tickets` con filtros, o el numero
   que el usuario pase). Trabajar sin ticket solo se acepta si el usuario lo dice
   explicitamente.
2. **Si no existe, crealo** antes de tocar codigo, con el tipo correcto (`bug`,
   `improvement`, `feature`, `data`, `integration`, `support`).
3. **Al terminar, comenta el ticket** con el diagnostico y la correccion. Ese
   comentario es el entregable, no un resumen en el chat.
4. **Mueve el estado** cuando cambie de verdad: `in_development` al empezar,
   `testing` al desplegar, `resolved` cuando este verificado en produccion.
   Nunca marcar `resolved` sin evidencia.
5. **Enlaza los commits** en el comentario (hash corto) y la entrada de bitacora
   si la hay. Al reves tambien: la bitacora referencia el ticket.

## Que lleva el comentario de cierre

En este orden, y con numeros reales, no descripciones vagas:

- **Que reporto el cliente**, en sus terminos.
- **Causa raiz**, una sola vez y sin rodeos.
- **Puntos de entrada afectados** si el mismo error vive en varios lugares.
- **Impacto medido**: cuantos registros, cuanta plata, que negocios. Con la
  consulta que lo sustenta si aplica.
- **Como se corrigio**, archivo por archivo y commit por commit.
- **Verificacion**: que se probo, contra que entorno, y el resultado exacto.
- **Hipotesis descartadas**, sobre todo las que costaron tiempo.
- **Pendiente**: correcciones de datos, notas credito, decisiones del usuario.

Si algo no se verifico, se dice. Un comentario que afirma mas de lo que se probo
es peor que no comentar.

## API

Todas las rutas exigen JWT y **super admin** (`middleware.RequireSuperAdmin`).

```
GET    /api/v1/tickets                 lista, con filtros y paginacion
POST   /api/v1/tickets                 crear
GET    /api/v1/tickets/:id             detalle
PUT    /api/v1/tickets/:id             actualizar
PATCH  /api/v1/tickets/:id/status      cambiar estado
PATCH  /api/v1/tickets/:id/assign      asignar
PATCH  /api/v1/tickets/:id/sprint      mover de sprint
PATCH  /api/v1/tickets/:id/area        cambiar area
PATCH  /api/v1/tickets/:id/escalate    escalar
GET    /api/v1/tickets/:id/comments    listar comentarios
POST   /api/v1/tickets/:id/comments    comentar
GET    /api/v1/tickets/:id/history     historial
POST   /api/v1/tickets/:id/attachments adjuntar
```

Comentar: `{"body": "...", "is_internal": false}`. `is_internal` solo lo respeta
el backend si quien comenta es super admin; sirve para notas que el negocio no
debe ver.

En el front el ticket se abre en `/tickets?ticket=<id>`.

## Valores validos

`back/central/services/modules/tickets/internal/app/constants.go` es la fuente:

| Campo | Valores |
|---|---|
| status | `open`, `in_review`, `in_development`, `testing`, `blocked`, `resolved`, `closed`, `wont_fix` |
| priority | `low`, `medium`, `high`, `critical` |
| type | `bug`, `improvement`, `feature`, `data`, `integration`, `support`, `complaint`, `claim`, `question` |
| severity | `low`, `medium`, `high` |
| source | `internal`, `business` |
| area | `comercial`, `soporte`, `desarrollo` |

Cuentan como cerrados: `resolved`, `closed`, `wont_fix`.

## Relacion con bitacora y alertas

| | ticket | bitacora | alerta |
|---|---|---|---|
| Que es | la unidad de trabajo | el historico del diagnostico | lo urgente sin terminar |
| Vive en | base de datos | `.claude/bitacora/` | `.claude/alerts/` |
| Se cierra | al verificar en produccion | nunca, es historico | al resolver lo urgente |

Un caso grande deja las tres cosas: el ticket con el cierre, la entrada de
bitacora con el detalle tecnico, y una alerta si quedo trabajo critico. Un bug
menor deja solo el ticket.

Reglas relacionadas: `.claude/rules/bitacora.md`, `.claude/rules/alerts.md`.
