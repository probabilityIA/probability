---
name: goal
description: >
  Ejecuta un objetivo de forma autonoma hasta cumplir un criterio de listo
  verificable. Usar cuando el usuario pide "goal: X", "termina X sin preguntar",
  "avanza hasta que X funcione" o cuando se combina con /loop para iterar
  sobre un mismo objetivo. No usar para preguntas ni para cambios de una linea.
user-invocable: true
disable-model-invocation: true
argument-hint: "<objetivo> [-- criterio de listo]"
---

# Goal - ejecucion autonoma con criterio de listo

Recibes un objetivo y trabajas sin pedir confirmaciones intermedias hasta que
el criterio de listo se cumpla, se agote el presupuesto o choques con un limite
duro. Las reglas del repo (`CLAUDE.md`, `.claude/rules/`) aplican completas;
este skill solo define COMO se ejecuta un objetivo, no relaja nada.

## Input

`$ARGUMENTS` = objetivo en lenguaje natural. Opcionalmente ` -- ` seguido del
criterio de listo. Si no viene criterio, lo defines tu en la fase 1 y lo dejas
escrito antes de tocar codigo.

Ejemplos:

- `/goal integrar Falabella Seller Center -- test_connection real y ordenes en cola canonical`
- `/goal que pase el E2E CU-03 de orders`
- `/goal cerrar la alerta wallet-cobro-guias-no-atomico`

## Fases

### 1. Contexto (sin escribir nada)

1. `ls .claude/alerts/` y leer las alertas que toquen el modulo.
2. Buscar en `.claude/bitacora/` un caso parecido.
3. `./scripts/dev-db-switch.sh status`: si apunta a prod y el objetivo escribe
   datos, cambiar a local ANTES de seguir (`dev-db-switch.sh local`).
4. `git status` y rama actual. Si estas en `main`, crear rama
   `feat/<tema>` o `fix/<tema>`.
5. Escribir el criterio de listo en
   `.claude/plan/goal-<tema>.md` con: objetivo, criterio, alcance, fuera de
   alcance, limites. Ese archivo es el contrato: no se amplia el alcance sin
   anotarlo ahi.

### 2. Plan

Lista de pasos pequenos y verificables. Cada paso dice como se comprueba
(comando, test, consulta SELECT, respuesta HTTP). Sin paso de verificacion no
hay paso.

### 3. Ejecucion

Iterar: hacer un paso -> verificar -> actualizar `goal-<tema>.md` (hecho /
bloqueado / pendiente) -> siguiente. Reglas:

- Backend: siempre compilar (`cd back/central && go build ./...`) y correr
  `go vet ./...` sobre los paquetes tocados antes de dar un paso por hecho.
- Frontend: `npm run lint` y `npx tsc --noEmit` en `front/central`.
- Si un paso falla 3 veces con el mismo enfoque, cambiar de enfoque o marcarlo
  bloqueado; nunca repetir lo mismo una cuarta vez.
- Commits chicos por paso cumplido, mensaje en espanol, sin push.
- Sin comentarios en Go/TS; archivos de 500+ lineas en ASCII puro.

### 4. Verificacion final

Ejecutar el criterio de listo tal como quedo escrito. Si no se cumple, el
objetivo NO esta listo aunque todos los pasos digan hecho.

### 5. Reporte

Terminar con un resumen que se sostenga solo:

- criterio de listo: cumplido / no cumplido y la evidencia (salida real de
  comando o test, no descripcion)
- commits creados (hash + mensaje)
- que quedo fuera y por que
- si quedo trabajo critico sin cerrar: crear o actualizar la alerta en
  `.claude/alerts/`
- si se diagnostico algo no trivial: entrada en `.claude/bitacora/`

## Limites duros (paran el goal, no se negocian)

Nunca, aunque el objetivo lo pida de forma implicita:

- `git push`, `gh pr merge`, tags, borrar ramas remotas.
- Iniciar, reiniciar o detener backend/frontend. Si el objetivo necesita el
  servicio arriba y no lo esta, marcar bloqueado y pedirlo al final.
- INSERT/UPDATE/DELETE contra el tunel de produccion (`5433`) ni contra
  ningun proveedor externo real (guias, cobros, plantillas de WhatsApp).
- Cambiar security groups, IAM, RDS, iptables o `.env`.
- Correr migraciones contra produccion. En local si, una a la vez.
- Ampliar el alcance: "ya que estoy" no existe. Si aparece algo, va a
  `goal-<tema>.md` como pendiente.
- Editar un caso de prueba o desactivar una validacion para que "pase".

Al chocar con un limite: parar ese paso, dejarlo bloqueado con el motivo, seguir
con los pasos que no dependan de el y reportarlo al final.

## Presupuesto

Por defecto 25 iteraciones (paso + verificacion) o 2 horas de reloj, lo que
llegue primero. Se puede pasar `--max N`. Al agotarse: reporte con lo hecho y lo
pendiente, no seguir.

## Combinacion con /loop

`/loop` puede reinvocar `/goal <mismo objetivo>`; en cada vuelta se lee
`goal-<tema>.md` para retomar donde quedo y no repetir pasos hechos. Ver
`.claude/rules/loop.md` para lo que un loop puede y no puede hacer.
