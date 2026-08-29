# woo-store: registro DNS automatico

La tienda WooCommerce de pruebas (`i-0e41ea3a2f1747cd3`) ya **no tiene IP
elastica**. Se solto el 2026-08-28 porque AWS cobra `PublicIPv4:IdleAddress`
(~USD 3.47/mes) mientras la instancia esta apagada, y esa instancia vive
apagada: se auto-apaga tras 1 hora sin trafico.

En su lugar, la instancia registra su propia IP publica al arrancar.

## Como funciona

1. Al bootear, `woo-dns.service` (oneshot, `Before=docker.service`) ejecuta
   `/usr/local/bin/woo-dns-update.sh`.
2. El script lee la IP publica de IMDSv2 (reintenta 10 veces cada 5 s, porque
   la IP tarda unos segundos en existir).
3. Compara contra el valor actual del registro y hace `UPSERT` de
   `woo.probabilityia.com.co` (TTL 60) solo si cambio.
4. Corre antes que Docker, asi Caddy levanta con el DNS ya correcto y puede
   renovar el certificado Let's Encrypt por reto HTTP-01.

Log: `/var/log/woo-dns.log`.

## Permisos

Rol `woo-store-dns-role` (instance profile del mismo nombre):

- `ProbabilityWooStoreDNS` (inline): `route53:ChangeResourceRecordSets` sobre la
  zona `Z06743392M4CV4HQ9DTIT`, **acotado por condicion** a que el registro sea
  exactamente `woo.probabilityia.com.co`, tipo `A`, accion `UPSERT`.
  Lleva `Null: false` sobre la misma clave: sin eso, `ForAllValues` con la clave
  ausente evalua a verdadero y el permiso se abriria a toda la zona, incluidos
  los registros de produccion.
- `AmazonSSMManagedInstanceCore` (gestionada): acceso por SSM, ya que esa
  instancia no tenia agente registrado y no hay `.pem`.

Verificado el 2026-08-28: el rol recibe `AccessDenied` al intentar mover
`app.probabilityia.com.co` y al intentar `DELETE` de su propio registro.

## Reinstalar

```bash
aws ssm start-session --profile probability --region us-east-1 \
  --target i-0e41ea3a2f1747cd3
sudo install -m 755 woo-dns-update.sh /usr/local/bin/woo-dns-update.sh
sudo install -m 644 woo-dns.service /etc/systemd/system/woo-dns.service
sudo systemctl daemon-reload && sudo systemctl enable --now woo-dns.service
```

## Cuidado

- El TTL del registro debe quedarse en **60**. Subirlo hace que la tienda quede
  apuntando a una IP muerta despues de cada encendido.
- La subred `subnet-0fce07b0057acaa4a` debe conservar
  `MapPublicIpOnLaunch=true`. Sin eso la instancia arranca sin IP publica y
  queda incomunicada.
