# Rediseno del Hero y nueva seccion de transportadoras (website publica)

**Fecha:** 2026-09-10
**Modulo:** `front/website` (sitio publico Astro)
**Tipo:** Cambio de diseno/UI, no incidente. Se documenta por el volumen de
iteraciones y las decisiones de diseno tomadas, para que quien retome esto
no tenga que reconstruir el contexto.

## Que se hizo

### Hero (`src/components/HeroSection.astro`)

- Se probo primero un hero oscuro con una caja de producto animada (video
  con chroma-key en canvas, rutas SVG tipo "mapa"). Se descarto: el usuario
  pidio fondo blanco, caja estatica y mas grande.
- Version final: hero claro (`#FBFAFE`), texto centrado, sin la caja (se
  movio a la nueva seccion de transportadoras). Badge, H1 con gradiente
  violeta, CTAs, stat strip con count-up. Sin elemento visual a la derecha.

### Nueva seccion `CarriersSection.astro`

Seccion nueva, insertada en `index.astro` justo despues del Hero:

- **Texto izquierda:** badge "Cobertura nacional", H2 "Una sola plataforma,
  todas las transportadoras del pais" (ancho de columna 720px, tipografia
  grande `lg:text-[52px]`), grid de 8 logos de transportadoras
  (`grid-cols-4 sm:grid-cols-6`), 3 checks de valor (texto `18px`).
- **Logos:** Coordinadora, Servientrega, Interrapidisimo, Envia, Deprisa,
  Pibox, Mensajeros Urbanos, TCC. Todos en
  `s3://probability-media-assets/public/carriers/imagen_*`.
  **Rappi quedo pendiente**: no hay logo en S3 ni en el repo. Si el usuario
  lo pasa, agregar como un elemento mas al array `carriers` en el frontmatter.
- **Caja de producto:** imagen `hero-box-poster.png` en
  `s3://probability-media-assets/public/website/`, generada a partir de una
  foto de producto (`Gemini_Generated_Image_...jpeg` en Downloads del
  usuario) con el fondo blanco removido via ImageMagick (`floodfill` desde
  las esquinas + trim). Se probo tambien un video real girando con
  chroma-key en canvas (`hero-box-rotating.mp4`), pero el usuario pidio
  volver a la caja estatica.
- **Posicionamiento de la caja:** `position: absolute` (fuera del flujo)
  dentro de la seccion, con tamano fluido via `clamp()` y CSS custom
  property `--box-size` (NO breakpoints fijos: un tamano fijo por
  breakpoint se veia distinto en cada monitor — a veces se veia un sticker
  del logo, a veces dos, dependiendo del ancho real). Se posiciona con
  `top: 50%` + `margin-top: calc(var(--box-size) * -0.5)` para centrar
  verticalmente, y `right: calc(var(--box-size) * -0.42)` para que sangre
  fuera del borde derecho. `z-index: 30` y estar fuera del flujo (no
  margenes negativos en flow) es lo que permite que la caja se superponga
  visualmente sobre el Hero (arriba) y sobre `DiagnosticoSection` (abajo)
  en vez de quedar recortada por la altura de su propia seccion.
- **Contenedor con ancho maximo:** el wrapper de texto+logos vive dentro de
  `max-w-[1440px] mx-auto`. Sin este limite, en monitores grandes (2000px+)
  el texto quedaba fijo a la izquierda y quedaba un hueco vacio enorme en
  el medio antes de llegar a la caja.
- **Fondo animado:** grid de "mapa" tipo Google Maps (CSS
  `repeating-linear-gradient` + `perspective`/`rotateX` para el efecto de
  inclinacion), una ruta SVG principal punto A -> punto B con un punto que
  viaja por el trayecto (`animateMotion`), y 8 "rutas laser" geometricas
  adicionales (path en angulo recto con `stroke-dasharray` animado
  simulando luz corriendo por el cable) en tonos violeta. Explicitamente
  se descarto generar esto como un video con IA (Gemini) por: (1) mp4 no
  soporta alpha, se habria necesitado chroma-key igual que con el video de
  la caja; (2) los modelos de video de IA no son confiables para geometria
  precisa (rutas rectas, texto de calles).

## Pendiente / conocido

- **Sombra residual en la base de la caja**: la foto original tenia una
  sombra de piso que el floodfill no quito del todo. Se intento limpiar con
  mas fuzz (rompia la textura de la caja), con un mask-image CSS en
  degradado (dejaba un remanente en un lado por la perspectiva de la foto)
  y con un polygon manual en ImageMagick erasing la zona (mejoro bastante
  pero no se termino de verificar 100%). El usuario dijo que lo iba a
  revisar el mismo. Si hace falta retomarlo: el archivo de trabajo
  (`hero-box-clean2-trim.png`) quedo en el scratchpad de la sesion, no se
  subio a S3 todavia porque no se confirmo que estuviera limpio.
- **Logo de Rappi**: falta el archivo. Usuario mencionado en la lista de
  transportadoras del pedido original pero no esta ni en S3 ni en el repo.
- El navbar del sitio (`Header.astro`) sigue con fondo blanco fijo; no se
  toco en esta sesion (quedo mencionado como posible ajuste, sin decision).

## Assets subidos a S3

Bucket `probability-media-assets` (publico, con CORS GET agregado para
permitir que un `<canvas>` lea pixeles de imagenes/videos del bucket desde
el navegador — necesario para el experimento de chroma-key del video):

- `public/website/hero-package.png` (ilustracion IA descartada, primera
  iteracion del hero oscuro)
- `public/website/hero-box-rotating.mp4` + `hero-box-poster.jpg` (video real
  girando, descartado — se volvio a caja estatica)
- `public/website/hero-box-poster.png` (imagen final en uso, fondo
  removido)

## Commit

`b32749e8` — `feat(website): rediseno del hero y nueva seccion de
transportadoras`, pusheado a `main`.
