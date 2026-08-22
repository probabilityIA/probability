# App movil: listados, paginado dinamico y memoria

Regla OBLIGATORIA para toda pantalla de la app Flutter (`mobile/mobile_central`)
que muestre una coleccion de registros. La app se disena para volumen alto
(cientos de miles de ordenes por negocio) y para telefonos de gama baja
(2 GB de RAM, Android Go). Ninguna pantalla puede asumir que el listado cabe
en memoria.

## 1. Siempre paginado dinamico con scroll, nunca paginado numerado

- **Prohibido** el paginado clasico con botones de pagina. En movil los
  objetivos de toque son malos y no ahorra nada del lado del servidor.
- **Prohibido** pedir una sola pagina y quedarse ahi. Un listado que llama
  `page=1&page_size=20` y no tiene `loadMore` esta roto: el usuario ve 20
  registros y cree que no hay mas.
- El patron unico es: `ListView.builder` + listener de scroll que dispara la
  pagina siguiente cuando faltan ~320 px para el final.
- Toda pantalla de listado muestra **"X de Y"** con el `total` que devuelve el
  API, para que el usuario no pierda la nocion del tamano real.

## 2. Nunca acumular sin limite

El scroll infinito ingenuo concatena paginas hasta morir. Prohibido.

- La coleccion vive en `PagedCollection` (`lib/shared/pagination/`), una lista
  **rala con ventana deslizante**: mantiene en memoria como maximo
  `maxPagesInMemory` paginas (por defecto 8 = 160 items). Las paginas que
  salen de la ventana se **expulsan por LRU** y sus posiciones quedan en
  `null`.
- Los indices **nunca se corren**: expulsar deja el hueco, no borra la
  posicion. Por eso el scroll no salta, que es el defecto clasico de recortar
  la cabeza de la lista.
- Un hueco se pinta como esqueleto de altura fija y **se vuelve a pedir solo**
  cuando entra al viewport.
- Costo de un registro expulsado: una referencia nula. 100.000 registros
  navegados = ~800 KB, no 100.000 objetos vivos.

## 3. El verdadero consumidor de memoria son las imagenes, no los modelos

Un modelo de orden pesa cientos de bytes; un PNG de 512x512 decodificado pesa
**1 MB** (ancho x alto x 4 bytes), sin importar que se pinte en 40 px.

- **Obligatorio** `cacheWidth` / `cacheHeight` en TODO `Image.network`,
  calculados como el tamano de pintado por el `devicePixelRatio`. Sin eso,
  Flutter decodifica a resolucion completa y lo guarda asi en el `ImageCache`.
- El tope global del `ImageCache` se fija en `main.dart`
  (`maximumSizeBytes`), dimensionado para gama baja. El default de Flutter
  (100 MB / 1000 imagenes) es demasiado para un telefono de 2 GB.
- Usar los widgets compartidos (`BrandLogo`, `NetworkAvatar`, la miniatura de
  producto), que ya aplican el downscale. Prohibido un `Image.network` suelto
  en una tarjeta de listado.

## 4. Como construir la lista

Usar `PaginatedListView` (`lib/shared/widgets/ui/paginated_list_view.dart`),
que ya resuelve el patron completo: estados de carga, error, vacio, refresh,
disparo de la pagina siguiente, esqueletos de los huecos y pie con "X de Y".

Si por algo hay que escribir un `ListView` a mano, cumplir igual:

- `ListView.builder` o `.separated`. **Prohibido** `ListView(children: [...])`
  o `Column` dentro de `SingleChildScrollView` con la coleccion completa: eso
  construye todos los widgets de una.
- `addAutomaticKeepAlives: false` - si no, cada tarjeta que sale de pantalla
  se queda viva.
- `addRepaintBoundaries: true` (default, no quitarlo).
- `cacheExtent` acotado (~600). El default se estira y construye de mas.
- Tarjetas de listado con constructor `const` donde se pueda, y sin animaciones
  ni sombras costosas por item.

## 4.1 Excepcion acotada: colecciones dentro de una pantalla de detalle

Dos pantallas construyen su coleccion completa de golpe porque el API la
devuelve entera y no existe endpoint paginado:

- `route_detail_screen.dart` - las paradas de la ruta.
- `order_detail_sections.dart` - los items de la orden.

Se aceptan mientras sean colecciones acotadas. **Umbral: si una ruta puede pasar
de ~80 paradas o una orden de ~80 items**, hay que pasar esa pantalla a
`CustomScrollView` + `SliverList.builder` para dejar de construir todo de una.
No agregar mas casos como estos.

## 5. Catalogos chicos

Los catalogos de menos de 50 registros (roles, permisos, recursos, acciones,
estados de orden / pago / fulfillment) quedan exentos del scroll infinito, igual
que en el backend (`.claude/rules/architecture.md`). Se piden de una con
`page_size` alto y se pintan igual con `ListView.builder`. **No** se pintan con
`Column`.

## 6. Del lado del backend

- Todo GET de listado pagina, y devuelve `total`. Ya es regla del repo.
- El `OFFSET` profundo degrada. Cuando un listado de la app pase a millones de
  filas, migrar ese endpoint a paginado por cursor (keyset) antes que subir el
  `page_size`. Anotarlo como alerta cuando se detecte.

## Violaciones criticas

- Listado sin `loadMore` (se queda en la primera pagina).
- Concatenar paginas sin ventana ni tope.
- Recortar la cabeza de la lista corriendo los indices (salta el scroll).
- `Image.network` en una tarjeta de listado sin `cacheWidth` / `cacheHeight`.
- `Column` / `ListView(children:)` con una coleccion que crece.
- Paginado numerado con botones de pagina.
