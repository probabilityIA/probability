# Las alarmas de CloudWatch no le avisaban a nadie

Fecha: 2026-08-21

## Resumen

El tema SNS `probability-alertas` tenia **cero suscriptores**, asi que las 7
alarmas de CloudWatch sobre EC2 y RDS llevaban tiempo disparando al vacio. Se
conecto `secamc93@gmail.com`, se verifico con un envio de prueba y de paso se
activo CloudTrail, que tampoco existia.

## Sintoma

No hubo sintoma. Ese es justamente el punto: nada fallaba de forma visible.
Aparecio al auditar el monitoreo de la cuenta, no por un incidente.

```bash
aws sns list-subscriptions-by-topic --profile probability --region us-east-1 \
  --topic-arn arn:aws:sns:us-east-1:476702565908:probability-alertas
# -> Subscriptions: []
```

Las 7 alarmas estaban bien construidas, con `AlarmActions` apuntando al tema
correcto y todas en estado `OK`:

| Alarma | Umbral |
|---|---|
| `ec2-cpu-alta` | CPUUtilization > 85% |
| `ec2-creditos-cpu-bajos` | CPUCreditBalance < 30 |
| `ec2-status-check-falla` | StatusCheckFailed >= 1 |
| `rds-almacenamiento-bajo` | FreeStorageSpace < 5 GB |
| `rds-conexiones-altas` | DatabaseConnections > 60 |
| `rds-creditos-cpu-bajos` | CPUCreditBalance < 50 |
| `rds-memoria-libre-baja` | FreeableMemory < 100 MB |

O sea: la deteccion funcionaba, la notificacion no. Si el RDS se hubiera quedado
sin espacio o el EC2 se hubiera caido, el aviso se quedaba dentro de AWS.

## Causa raiz

Se creo el tema SNS y se cablearon las alarmas contra el, pero nunca se suscribio
un destinatario. Un tema SNS sin suscriptores acepta el `Publish` y responde con
un `MessageId` valido: **no hay error, no hay warning, no hay forma de notarlo
salvo mirandolo a proposito.**

Confunde ademas que Cost Anomaly Detection SI notificaba (a
`probabilitysas@gmail.com`), porque usa su propio mecanismo de suscripcion y no
pasa por este tema. Ver alertas de costos llegando da la falsa impresion de que
el monitoreo entero esta conectado.

## Correccion

```bash
aws sns subscribe --profile probability --region us-east-1 \
  --topic-arn arn:aws:sns:us-east-1:476702565908:probability-alertas \
  --protocol email --notification-endpoint secamc93@gmail.com
```

La suscripcion queda en `PendingConfirmation` hasta que la persona confirma desde
el correo de AWS. **El paso de confirmacion no se puede automatizar** y el enlace
expira en 3 dias. Sin ese clic no sirve de nada.

## Verificacion

`SubscriptionArn` paso de `PendingConfirmation` al ARN
`...probability-alertas:23dc3fee-c805-4a55-87ce-64ef88843e39`, y se mando un
`sns publish` manual de prueba (MessageId `df4aedd3-e8cd-56ce-952c-8ed0e10e4dca`)
que llego al correo.

## Lo otro que se hizo en la misma sesion

- **CloudTrail**: no existia ningun trail. Se creo `probability-trail`
  (multi-region, validacion de integridad) sobre el bucket
  `probability-cloudtrail-476702565908`, con lifecycle de 90 dias y **sin data
  events**, que son los que cuestan. Detalle en `.claude/rules/infra-ops.md`.
- **Usuario `santiago.camacho`**: se le adjunto `ProbabilityOpsSantiago` para S3,
  lectura de EC2/RDS y shell por SSM. Se decidio ampliar el usuario existente en
  vez de crear uno nuevo con guion, para no duplicar identidades.

## Hipotesis y decisiones descartadas

- **Mandar los logs de los contenedores a CloudWatch Logs.** Descartado por
  costo: la ingesta son $0.50/GB y el backend, durante el incidente del consumidor
  en bucle (`.claude/rules/colas-errores-permanentes.md`), llego a 4,2 MB/min, que
  proyecta ~180 GB/mes. Habria que medir el volumen real antes de encenderlo.
- **Levantar Prometheus + Grafana** (`monitoring-api`, `monitoring-web`, carpetas
  `grafana/` y `alloy/` del compose de produccion). Descartado: el EC2 es un
  `t4g.small` con 8 contenedores encima y no le sobra RAM.
- **Route53 health checks / CloudWatch Synthetics** para uptime externo.
  Descartado: Route53 cuesta $0.50/mes por chequeo y el tier gratis de Synthetics
  (100 corridas/mes) da un chequeo cada 7 horas, que no sirve.
- **Detailed monitoring del EC2**: se dejo apagado a proposito. Las metricas de 1
  minuto se cobran; las de 5 minutos son gratis y alcanzan.
- **Retencion de CloudTrail a 1 ano**: se propuso y el usuario prefirio dejar 90
  dias. Vale saber que `lookup-events` ya da 90 dias gratis sin trail, asi que el
  valor del trail hoy es la copia con validacion de integridad, no el histrico.

## Pendientes

Todos dentro de capa gratuita, ninguno hecho:

1. **Metricas de RAM y disco del EC2** via CloudWatch agent. Es el hueco mas
    grande: hoy nadie vigila la memoria ni el disco del `t4g.small`, y es lo que
    con mas probabilidad lo tumbe (disco lleno de imagenes Docker, o RAM agotada).
    Son ~3 metricas custom de las 10 gratis.
2. **Budget de costos** (2 gratis). Avisa antes de la factura, no despues.
3. **Dashboard de CloudWatch** (3 gratis).
4. **Revisar periodicamente que la suscripcion SNS siga activa.** Cada correo de
    AWS trae un enlace de "unsubscribe" al pie: un clic accidental devuelve el
    sistema al estado de esta entrada, en silencio.
