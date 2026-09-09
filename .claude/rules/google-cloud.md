# Google Cloud - acceso por CLI

Desde 2026-09-09 hay `gcloud` instalado y autenticado. Antes de mandar al usuario
a la consola web, revisar si el comando existe: casi todo lo de GCP se hace por
CLI, igual que AWS.

## Instalacion y perfil

SDK en `~/google-cloud-sdk` (tarball, no apt). El PATH ya quedo en `~/.zshrc`.

Las **configuraciones** de gcloud son el equivalente de los perfiles de AWS:

| AWS | gcloud |
|---|---|
| `--profile probability` | `--configuration=probability` |
| `AWS_PROFILE` | `CLOUDSDK_ACTIVE_CONFIG_NAME` |
| `aws configure --profile X` | `gcloud config configurations create X` |

La configuracion activa es `probability`: cuenta `probabilitysas@gmail.com`,
proyecto `probabilityia`. La configuracion `default` quedo vacia a proposito.

```bash
export PATH="$HOME/google-cloud-sdk/bin:$PATH"
export CLOUDSDK_ACTIVE_CONFIG_NAME=probability
gcloud auth list                 # cuentas logueadas (varias a la vez, como AWS)
gcloud config list               # cuenta y proyecto de la configuracion activa
```

Varias cuentas pueden estar logueadas simultaneamente y se alternan con
`gcloud config configurations activate <nombre>` sin volver a autenticarse.

## Que hay en la cuenta

Organizacion `probabilitysas-org`, ID **432676620088**.

| Proyecto | Nombre | Numero |
|---|---|---|
| `probabilityia` | ProbabilityIA | 901639821821 |
| `probability-app-6de7c` | Probability App | 176535439905 |
| `project-38264ba0-30a2-41a9-a4e` | My First Project | 911443558270 |

**El proyecto del producto es `probabilityia`.** Los otros dos no se tocan sin
preguntar.

`probabilitysas@gmail.com` tiene `roles/owner` en `probabilityia` y
`roles/resourcemanager.organizationAdmin` en la organizacion, o sea que puede
crear proyectos y dar accesos a nivel de org.

## ADC es global, no por configuracion

`gcloud auth application-default login` escribe UN solo archivo
`~/.config/gcloud/application_default_credentials.json`, compartido por todas
las configuraciones. Las librerias (Python, Go, Node) leen ese archivo, asi que
un script corre con la ultima cuenta que hizo ADC login, sin importar que
configuracion este activa.

Para separarlo de verdad: cuenta de servicio por proyecto y
`GOOGLE_APPLICATION_CREDENTIALS` apuntando a su JSON.

## Analitica del sitio publico

`front/website/src/layouts/Layout.astro` tiene GA4 (`G-PR13PJFXXF`) y el meta de
verificacion de Search Console. Son dos fuentes distintas:

- **Search Console**: clics HACIA la web desde resultados de busqueda (consultas,
  impresiones, CTR, posicion).
- **GA4**: clics DENTRO de la web (paginas vistas, origen del trafico, eventos).

Hoy GA4 solo mide automaticamente enlaces salientes y descargas. **Los clics en
botones internos (demo, WhatsApp, planes) no se estan registrando**: hace falta
`gtag('event', ...)` en cada CTA.

Para consultarlas por CLI/script hace falta habilitar
`analyticsdata.googleapis.com` y `searchconsole.googleapis.com`, y ademas
invitar a la identidad (usuario o cuenta de servicio) DENTRO de cada producto:
GA4 en Administrar > Gestion de accesos, Search Console en Configuracion >
Usuarios y permisos. Sin esa invitacion la API responde 403 aunque el permiso de
GCP este bien.

## Google Maps - geocodificacion

Desde 2026-09-09 la geocodificacion corre con keys del proyecto `probabilityia`,
no con la cuenta personal de nadie. Facturacion vinculada:
`01D184-410BBB-DA9172`.

APIs habilitadas: `geocoding-backend`, `places-backend`, `maps-backend`.

Hay **dos keys**, ambas restringidas por IP y limitadas a Geocoding y Places:

| Key | Restringida a | Para |
|---|---|---|
| `Probability Backend - Geocoding` | `3.224.189.33` (EC2 de produccion) | produccion |
| `Probability Dev - Geocoding` | IP publica del equipo local | desarrollo |

**Las keys NO van en el repo.** Estan en `~/.config/probability/maps-keys.txt`
(chmod 600) y en el `.env` de cada entorno, que esta gitignored.

La key **solo la usa el backend**. El navegador nunca la ve: el front llama a
`/api/v1/geocode` y el backend consulta a Google
(`cmd/internal/routes/geocode.go`). Si algun dia una pantalla necesita el mapa
JavaScript, esa key es OTRA, restringida por dominio HTTP, nunca la del backend.

La key de desarrollo esta atada a la IP de casa: **cuando cambie el proveedor de
internet, deja de funcionar** y responde `REQUEST_DENIED` diciendo desde que IP
llego el request. Se actualiza con:

```bash
gcloud services api-keys update <KEY_ID> --allowed-ips="<IP_NUEVA>"
gcloud services api-keys list --format="table(displayName,uid)"
```

**Cuota diaria: 500 llamadas** en `geocoding-backend` y otras 500 en
`places-backend` (el endpoint de busqueda de direcciones usa Places Text
Search). Existe para que un bucle en el codigo no genere una factura sorpresa,
no porque el volumen real este cerca. Si un dia la operacion legitima la choca,
se sube; el sintoma es `OVER_QUERY_LIMIT`.

```bash
gcloud alpha services quota update --service=geocoding-backend.googleapis.com \
  --consumer=projects/901639821821 \
  --metric=geocoding-backend.googleapis.com/billable_default \
  --unit="1/d/{project}" --value=<N> --force
```

**Prohibido repartir el trafico entre la cuenta de la empresa y una personal
para duplicar el tramo gratuito.** Lo prohiben los terminos de Maps Platform, es
trivial de detectar (misma IP, mismo dominio, mismo perfil de pagos) y el
castigo es la suspension de las dos cuentas, lo que deja produccion sin
geocodificar. El tramo gratis es de 10.000 llamadas al mes; el excedente son
~5 USD por cada 1.000. El ahorro real esta en cachear por direccion normalizada,
no en abrir cuentas.

## Firebase / notificaciones push

Proyecto `probability-app-6de7c` (el de Firebase, distinto de `probabilityia`).
Apps registradas: Android `com.Probability`, iOS `com.probabilityia.mobileCentral`.

FCM es gratis y sin limite; no consume el tramo gratuito de nada.

Cuenta de servicio que envia: `probability-fcm-sender@probability-app-6de7c.iam.gserviceaccount.com`,
con `roles/firebasecloudmessaging.admin` y nada mas. Su llave JSON esta en
`~/.config/probability/fcm-sender.json` (chmod 600) y **nunca en el repo**. Aca
la llave si se justifica: el EC2 no tiene gcloud y el envio es desatendido.

Ojo con el rol: `roles/firebasemessaging.admin` **no existe** y la API lo rechaza
con `INVALID_ARGUMENT`. El bueno es `roles/firebasecloudmessaging.admin`.

Detalle del modulo: `back/central/services/modules/push/README.md`.

## Lo que NO tiene CLI

- Agregar o quitar usuarios en Search Console: solo consola web.
- La cuenta Google en si (`probabilitysas@gmail.com`): es Gmail, no Workspace,
  asi que no hay Admin SDK. Contrasena, 2FA y sesiones solo por web.
- Roles de GA4: se gestionan por la Admin API o la consola, no por `gcloud`.

## Reglas

1. **Dar accesos con el rol mas bajo que sirva.** `roles/owner` solo para socios
   administradores; para mirar, `roles/viewer`. Un `owner` de mas es un dueno de
   la organizacion, no un invitado.
2. **Nunca crear llaves JSON de cuenta de servicio "por comodidad".** Para
   operar a mano esta el login interactivo; la llave es para lo desatendido
   (cron, CI) y es una credencial larga que no expira. Si se crea, va fuera del
   repo o en un path gitignored, nunca en un commit.
3. **Prohibido borrar proyectos** (`gcloud projects delete`) y prohibido tocar
   la cuenta de facturacion sin autorizacion explicita del usuario.
4. Al dar de alta a alguien, dejarlo escrito aca con el rol y la fecha.

## `roles/owner` no se puede otorgar por CLI a un correo externo

`gcloud projects add-iam-policy-binding ... --role=roles/owner` falla con
`ORG_MUST_INVITE_EXTERNAL_OWNERS` cuando el correo no pertenece al dominio de la
organizacion (cualquier Gmail suelto lo es). Google exige el flujo de invitacion
por correo, que solo existe en la consola web:

> IAM y administracion > Otorgar acceso > correo > rol Propietario > Guardar.
> Le llega una invitacion y el rol no aplica hasta que la acepta.

El equivalente que SI se otorga por CLI, y que alcanza para operar el proyecto
entero, es `roles/editor` + `roles/resourcemanager.projectIamAdmin`. Lo unico
que le falta frente a un owner de verdad es poder eliminar el proyecto y
gestionar su vinculo de facturacion.

## Altas

| Fecha | Correo | Rol | Alcance |
|---|---|---|---|
| 2026-09-09 | `secamc93@gmail.com` | `roles/editor` + `roles/resourcemanager.projectIamAdmin` | **revocado el mismo dia**, fue una prueba |

El control total se ejerce con `probabilitysas@gmail.com`, que es el unico
miembro del proyecto. No agregar owners nuevos sin acordarlo.

## Dar acceso solo a la analitica de la web

La analitica NO pasa por IAM de GCP. Un usuario invitado a GA4 y Search Console
ve el trafico del sitio publico y **nada** de la aplicacion: ni el proyecto de
nube, ni el RDS de AWS, ni datos de clientes.

- **GA4**: Administrar > Gestion de accesos a la propiedad > `+` > correo >
  `Analista` (lee y arma informes) o `Lector` (solo lee). Sin marcar
  "costo/ingresos" si no hace falta.
- **Search Console**: Configuracion > Usuarios y permisos > Agregar usuario >
  `Restringido` (solo lectura).

Search Console no tiene API de usuarios en absoluto: ese siempre es a mano. GA4
si tiene (`accessBindings`), ver mas abajo.

## Consultar la analitica desde la terminal

`./scripts/web-analytics.sh` (`props`, `sites`, `pages`, `sources`, `events`,
`queries`, `clicks`, todos con `[dias]` opcional, por defecto 28).

| Que | Valor |
|---|---|
| Propiedad GA4 de la web | `553391903` ("Probability Web"), measurement ID `G-84XRXZKPSW` |
| Sitio en Search Console | `sc-domain:probabilityia.com.co` |
| Identidad que lee | `probability-analytics@probabilityia.iam.gserviceaccount.com` |

La propiedad `491299983` (`probability-app-6de7c`) la creo Firebase sola, no
tiene ningun flujo de datos y **siempre devuelve vacio**. No confundirlas.

### Autenticacion: impersonacion, no llaves

**El login de usuario NO sirve para Analytics.** Google bloquea el scope
`analytics.readonly` con el client ID por defecto de gcloud: el navegador
responde "Esta aplicacion esta bloqueada". No perder tiempo ahi.

El camino que funciona es acunar un token de la cuenta de servicio con los
scopes de analitica, usando el login de usuario como firmante:

```bash
UT=$(gcloud auth print-access-token)
curl -s -X POST "https://iamcredentials.googleapis.com/v1/projects/-/serviceAccounts/${SA}:generateAccessToken" \
  -H "Authorization: Bearer $UT" -H "Content-Type: application/json" \
  -d '{"scope":["https://www.googleapis.com/auth/analytics.readonly","https://www.googleapis.com/auth/webmasters.readonly"],"lifetime":"3600s"}'
```

`gcloud auth print-access-token --impersonate-service-account=... --scopes=...`
**ignora el flag `--scopes`** y devuelve un token de `cloud-platform`, que las
APIs de Analytics rechazan. Hay que pegarle a `iamcredentials` directo.

Requisito: `roles/iam.serviceAccountTokenCreator` sobre la SA (ya lo tiene
`probabilitysas@gmail.com`) y `iamcredentials.googleapis.com` habilitada.

### Gotchas de la Admin API de GA4

- `accessBindings` (dar acceso a usuarios) vive en **`v1alpha`**, no en
  `v1beta`: con `v1beta` responde un 404 en HTML.
- Crear una propiedad exige ser **Editor de la CUENTA**, no de la propiedad.
- Cambiar roles de usuarios exige ser **Administrador de la cuenta**; con
  Editor da 403 aunque la llamada sea correcta.
