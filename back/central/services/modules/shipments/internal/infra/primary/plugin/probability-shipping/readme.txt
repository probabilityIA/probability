=== Probability Shipping ===
Contributors: probability
Requires at least: 6.0
Tested up to: 6.7
Requires PHP: 7.4
Stable tag: 1.4.0
License: GPLv2 or later

Cotiza tarifas de transportadoras (EnvioClick y otras) en el checkout de
WooCommerce consultando la API de Probability. El cliente final ve las opciones
de envio con su precio antes de pagar.

== Instalacion ==

1. En Probability, entra a Integraciones -> WooCommerce y copia tu "Clave de conexion".
2. En tu WordPress: Plugins -> Anadir nuevo -> Subir plugin -> selecciona el .zip
   de Probability Shipping -> Instalar -> Activar.
3. Ve a WooCommerce -> Ajustes -> Envio -> tu zona -> Anadir metodo de envio ->
   "Probability (Transportadoras)".
4. Abre el metodo, pega la "Clave de conexion" y guarda.
5. Listo: en el checkout apareceran las tarifas reales de las transportadoras.

== Changelog ==

= 1.4.0 =
* Ciudad como desplegable con buscador (select2) restringido a los municipios
  del departamento, igual que el campo Departamento (solo checkout clasico).

= 1.3.0 =
* Soporte para el checkout de bloques (WooCommerce Blocks): logos de
  transportadoras en las opciones de envio y validacion de la direccion de
  destino. El checkout clasico sigue funcionando igual.

= 1.2.0 =
* Logos de transportadoras en las opciones de envio del checkout.
* Sugerencias de municipios (DANE) segun el departamento y validacion de la
  direccion de destino contra la API de Probability.

= 1.1.0 =
* Clave de conexion unica (URL + integracion + token) para configurar en un paso.
* Autenticacion por token.

= 1.0.0 =
* Version inicial.
