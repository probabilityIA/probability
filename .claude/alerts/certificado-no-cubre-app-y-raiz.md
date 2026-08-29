# Certificado TLS no cubre app.probabilityia.com.co ni el dominio raiz

Fecha: 2026-08-28
Detectado de paso mientras se probaba el endpoint `/woo-store/start` contra
produccion.

## Que pasa

El certificado que sirve nginx en produccion tiene un solo nombre:

```
subject=CN = www.probabilityia.com.co
X509v3 Subject Alternative Name:
    DNS:www.probabilityia.com.co
notAfter=Oct 16 06:26:37 2026 GMT
```

Route53 tiene tres registros A apuntando al mismo servidor `3.224.189.33`:

| Nombre | Responde |
|---|---|
| `www.probabilityia.com.co` | OK, HTTP 200 |
| `app.probabilityia.com.co` | **falla TLS** |
| `probabilityia.com.co` (raiz) | **falla TLS** |

Reproducir:

```bash
curl -sS https://app.probabilityia.com.co/health
# curl: (60) SSL: no alternative certificate subject name matches
#       target host name 'app.probabilityia.com.co'
```

## Por que importa

`app.` es el nombre que se ve como el del panel, y el raiz es lo que la gente
teclea. Cualquiera que entre por ahi ve la pantalla roja de "conexion no
privada" del navegador. Solo funciona si se escribe `www.` explicitamente.

No es un problema de DNS: los tres registros resuelven bien. Es el certificado.

## Urgente

- Emitir el certificado con los tres nombres (`probabilityia.com.co`,
  `www.`, `app.`) y recargar nginx. Revisar como se emite hoy en
  `infra/nginx/` para no romper la renovacion automatica.

## Importante

- Decidir cual es el nombre canonico y redirigir los otros dos con 301, en vez
  de servir el mismo sitio en tres nombres distintos.
- Al agregar un subdominio nuevo, agregarlo tambien al certificado. Este caso
  sugiere que `app.` se creo en DNS y nadie toco el cert.

## Criterio para cerrar

`curl -sS https://app.probabilityia.com.co/health` y el mismo contra el raiz
devuelven 200 sin error de TLS, y la renovacion automatica sigue funcionando
(verificar que el siguiente ciclo renueve los tres nombres).
