# /loop - ejecucion repetida

`/loop` reejecuta un prompt o skill cada cierto intervalo (o con ritmo propio si
no se da intervalo). Corre sin nadie mirando, asi que cada vuelta tiene que ser
segura de repetir. Estas reglas aplican a toda iteracion.

## Para que se usa

- Vigilar algo externo que el harness no notifica: un `gh run`, un deploy, una
  cola de RabbitMQ, un tunel SSM.
- Reintentar un `/goal` hasta cumplir su criterio de listo
  (`.claude/skills/goal/SKILL.md`).
- Sondear un endpoint o log local durante una prueba.

No se usa para tareas de una sola vez ni como cron de produccion: eso va en
`/schedule` o en un cron real del EC2.

## Cada iteracion debe ser idempotente

Repetir la vuelta 10 veces tiene que dejar el sistema igual que repetirla 1 vez.

- Antes de crear algo, comprobar si ya existe (rama, commit, archivo, registro).
- Nunca acumular: no reencolar, no duplicar filas, no crear otra alerta igual.
- Escribir el estado en un archivo (`.claude/plan/goal-<tema>.md` o el que el
  prompt indique) y leerlo al inicio de la vuelta siguiente.

## Prohibido dentro de un loop

Todo lo prohibido fuera de un loop sigue prohibido, y ademas:

- `git push`, merges, tags. Un loop puede commitear en rama; publicar es del
  usuario.
- Iniciar, reiniciar o detener servicios (`dev-services.sh start|restart|stop`,
  `docker restart`, `docker compose up`).
- Escribir en produccion: base por el tunel `5433`, RDS, S3 de deploy, Redis o
  RabbitMQ de prod, Graph API de Meta, cotizaciones o guias reales.
- Llamadas a AWS que cambien estado (`authorize-security-group-ingress`,
  `modify-db-instance`, `ssm send-command` con comandos que escriban). Solo
  lectura: `describe-*`, `get-*`, `list-*`, `lookup-events`.
- Mandar mensajes a canales externos (WhatsApp, correo, Slack) o abrir PRs.
- Consultar el MCP de postgres de produccion en cada vuelta para "ver si
  cambio algo": si hace falta vigilar datos, hacerlo contra local (`5434`) o
  con una consulta puntual acordada con el usuario.

## Cadencia

- **Intervalo por defecto: 60 segundos.** Si el usuario no da intervalo, el
  loop corre cada 60 s (`/loop 60s ...`). No alargarlo por cuenta propia.
- Solo se usa otro intervalo si el usuario lo indica explicitamente.
- Con 60 s la vuelta tiene que ser barata: leer estado, comprobar, actuar solo
  si cambio algo. Nada de recompilar todo o releer el repo en cada vuelta.

## Parar

El loop se detiene solo cuando:

- se cumplio el objetivo o la condicion vigilada,
- tres vueltas seguidas fallan con el mismo error,
- se choco con un limite duro del `/goal`,
- pasaron 3 horas sin avance (ninguna vuelta marco `noop: false`),
- el presupuesto del `/goal` se agoto.

Al parar, dejar un reporte final que se sostenga solo: que se logro, evidencia,
que quedo pendiente. Si quedo algo critico, alerta en `.claude/alerts/`.

## Ruido

- Vuelta sin cambios: `noop: true`, sin escribir archivos ni commits vacios.
- Vuelta con cambios: `noop: false` y una linea de que cambio.
- No repetir en cada vuelta el contexto completo; solo el delta.
