# WordPress + WooCommerce de Pruebas (local)

Entorno local en Docker para probar la integracion de WooCommerce de Probability
(test de conexion, sync de ordenes por API REST y webhooks `order.created`).

No toca produccion. Usa volumenes nombrados de Docker, asi que no escribe archivos
dentro del repo.

## Requisitos

- Docker + Docker Compose (v1 o v2; `setup.sh` detecta cual)
- `curl` (para el script de setup)

> El backend de Probability NO hace falta para montar WordPress/WooCommerce.
> Solo se necesita despues, para probar webhooks e integracion:
> `./scripts/dev-services.sh start backend` (`:3050`).

## Levantar todo (primera vez)

```bash
cd wordpress
./setup.sh
```

El script:
1. Levanta MariaDB + WordPress (puerto host `8088`).
2. Instala WordPress core (admin/admin).
3. Activa permalinks "pretty" (requisito de la REST API de WooCommerce).
4. Instala y activa WooCommerce.
5. Crea una API key REST (read_write) y la imprime.
6. Crea el webhook `order.created` apuntando al backend local.

Al final imprime **Store URL + Consumer Key + Consumer Secret** para crear la
integracion en Probability. La primera vez descarga imagenes (~400 MB).

> **Re-ejecutar `setup.sh`:** la instalacion de WP/WooCommerce es idempotente
> (se omite si ya existe), pero **crea una API key y un webhook nuevos cada vez**.
> Si solo querias re-levantar, usa `docker compose up -d` (ver abajo).

### Recuperar las credenciales

Si perdiste las claves que imprimio el script:

```bash
docker exec woo_test_db mariadb -uwordpress -pwordpress wordpress \
  -e "SELECT consumer_key, consumer_secret, truncated_key FROM wp_woocommerce_api_keys;"
```

`consumer_key` esta hasheada; usa `truncated_key` para identificar cual es y el
`consumer_secret` (en claro). Si no sabes cual `ck_` corresponde, borra esa fila
y vuelve a correr `setup.sh` para generar un par nuevo.

## Uso diario

> Si tu maquina solo tiene Compose v1, reemplaza `docker compose` por `docker-compose`.
> (`setup.sh` detecta cual tienes automaticamente.)

```bash
cd wordpress
docker compose up -d db wordpress    # arrancar
docker compose stop                  # detener (conserva datos)
docker compose logs -f wordpress     # ver logs
docker compose down                  # detener y borrar contenedores (conserva volumenes)
docker compose down -v               # BORRA TODO (datos incluidos) -> empezar de cero
```

## Accesos

| Que | Donde |
|-----|-------|
| Tienda | http://localhost:8088 |
| wp-admin | http://localhost:8088/wp-admin (admin / admin) |
| DB | MariaDB en contenedor `woo_test_db` (wordpress/wordpress) |

## Conectar con Probability

Al crear la integracion WooCommerce en el frontend:

- **Store URL:** `http://localhost:8088`
- **Consumer Key / Secret:** los que imprime `setup.sh`

El backend (proceso host en `:3050`) alcanza WooCommerce via `localhost:8088`.
WooCommerce (contenedor) alcanza el backend via
`http://host.docker.internal:3050/api/v1/woocommerce/webhook`.

### Basic Auth sobre HTTP

El backend autentica contra la REST API de WooCommerce con **Basic Auth**, que
WooCommerce solo acepta cuando `is_ssl()` es true. En prod eso pasa por HTTPS;
en local (HTTP plano) `setup.sh` instala un mu-plugin
(`wp-content/mu-plugins/probability-rest-ssl.php`) que marca solo las requests
`/wp-json/` como SSL. Sin el, el back recibiria 401 al conectar.

### Webhooks / HMAC

El backend valida HMAC solo si `WOOCOMMERCE_WEBHOOK_SECRET` esta seteado en su `.env`.
Para local, dejarlo vacio = se omite la validacion (mas simple).

Los webhooks de WooCommerce se entregan via wp-cron (async). El servicio `wpcron`
del compose (contenedor `woo_test_cron`) ejecuta `wp cron event run --due-now` y
`wp action-scheduler run` cada 20s, simulando el trafico de una tienda real:
las ordenes llegan solas a Probability en menos de un minuto.

Si necesitas forzar la entrega inmediata:

```bash
docker compose run --rm wpcli wp action-scheduler run
```

## Probar la integracion contra PRODUCCION (tienda local -> prod)

Flujo completo: el WooCommerce local se expone a internet con un tunel, se crea
la integracion desde el panel de produccion, el backend de prod registra los
webhooks automaticamente en la tienda, y las ordenes creadas en local llegan
solas al modulo de Ordenes de prod.

### 1. Levantar la tienda local

```bash
cd wordpress
./setup.sh            # primera vez
docker compose up -d  # siguientes veces (incluye el contenedor wpcron)
```

### 2. Generar credenciales API REST

```bash
docker compose run --rm wpcli wp eval '$ck="ck_".wc_rand_hash(); $cs="cs_".wc_rand_hash(); global $wpdb; $wpdb->insert($wpdb->prefix."woocommerce_api_keys", array("user_id"=>1,"description"=>"prueba-prod","permissions"=>"read_write","consumer_key"=>wc_api_hash($ck),"consumer_secret"=>$cs,"truncated_key"=>substr($ck,-7))); echo $ck." ".$cs;'
```

Guarda el `ck_...` y el `cs_...` que imprime (no se pueden recuperar despues).

### 3. Exponer la tienda con un tunel publico

```bash
# una sola vez: descargar cloudflared
curl -sSL -o ~/bin/cloudflared https://github.com/cloudflare/cloudflared/releases/latest/download/cloudflared-linux-amd64
chmod +x ~/bin/cloudflared

# levantar el tunel (dejar corriendo en una terminal)
~/bin/cloudflared tunnel --url http://localhost:8088
```

Copia la URL `https://xxxx.trycloudflare.com` que imprime. OJO: la URL cambia
en cada arranque del tunel (ver "Si el tunel cambia" abajo).

### 4. Crear la integracion en produccion

En `https://www.probabilityia.com.co` -> Integraciones -> Crear Integracion ->
E-commerce -> WooCommerce:

- **URL de la tienda:** la URL del tunel (sin barra final)
- **Consumer Key / Secret:** los del paso 2
- **Negocio:** Demo (para pruebas)

"Probar Conexion" debe pasar (prod alcanza la tienda via el tunel). Al crear,
el backend registra solo los webhooks `order.created` y `order.updated` en la
tienda con `?integration_id=<id>` en el delivery URL (verificalos en wp-admin ->
WooCommerce -> Ajustes -> Avanzado -> Webhooks).

### 5. Probar

Crea una orden en `http://localhost:8088/wp-admin` (admin/admin) con cualquier
metodo de pago. En menos de un minuto aparece en Ordenes del negocio Demo en
prod. Con metodo `cod` (contra entrega) ademas llega con `cod_total` y entra
al modulo de Recaudo COD.

### Si el tunel cambia de URL

El quick tunnel de cloudflare es efimero. Al reiniciarlo:

1. Editar la integracion en prod y actualizar la URL de la tienda.
2. Borrar los webhooks viejos en wp-admin y recrearlos: lo mas facil es borrar
   la integracion en prod y crearla de nuevo (los webhooks se registran solos).

### Limpieza despues de probar

- Borrar la integracion en prod (Integraciones -> eliminar).
- Borrar los webhooks que apunten a probabilityia.com.co en wp-admin.
- Las ordenes de prueba en prod se eliminan desde el modulo de Ordenes.

## Reproducir el bug del `image.id` (HISTORICO — ya corregido)

El bug original: WooCommerce manda `line_items[].image.id` como string y el
backend rechazaba la orden con `cannot unmarshal string into ...
WooLineItemImage.line_items.image.id of type int64`. Corregido con un tipo
flexible que acepta numero, string o null (con tests).

Ver logs del back: `./scripts/dev-services.sh logs backend 200`.
