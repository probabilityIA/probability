---
name: ortografia-front
description: >
  Revisa y corrige la ortografia del espanol en los textos de cara al usuario
  del front (labels, placeholders, toasts, mensajes de error, tours, correos):
  tildes faltantes, enies, signos de apertura y terminologia del producto.
  Usar cuando se escriba o modifique texto visible en front/central o
  front/website, cuando se pida "revisar ortografia", "faltan tildes",
  "corregir acentos", o antes de cerrar un PR que toque UI.
user-invocable: true
disable-model-invocation: false
allowed-tools: Bash, Read, Grep, Glob, Edit
argument-hint: "[ruta o modulo] [--fix]"
---

# Ortografia del espanol en el front

El front de Probability se escribio sin tildes en buena parte de sus textos:
`Codigo`, `Telefono`, `Integracion`, `envio`, `Anadir` (por `Código`,
`Teléfono`, `Integración`, `envío`, `Añadir`). Son ~1.700 casos en
~280 archivos. Este skill los detecta, distingue el texto de cara al usuario
del codigo, y los corrige respetando la regla de UTF-8 del repo.

**Regla de fondo: todo texto que ve un usuario esta bien escrito en espanol.**
Un identificador (`integrationId`, `codigo_postal`, una clave de objeto, una
ruta) no lleva tilde nunca; un `<label>`, un `placeholder`, un toast o un tour
la lleva siempre.

## Herramienta

```bash
python3 .claude/skills/ortografia-front/revisar.py [rutas...] [opciones]
```

Sin rutas revisa `front/central/src` y `front/website/src`. Ejecutar siempre
desde la raiz del repo.

| Opcion | Que hace |
|---|---|
| (nada) | reporta hallazgos seguros por archivo, ordenados por cantidad |
| `--fix` | aplica solo los hallazgos marcados `auto` |
| `--ambiguos` | ademas reporta tildes diacriticas (`esta`/`está`, `aun`/`aún`, `publico`) |
| `--sin-apertura` | omite la revision de `?` y `!` sin signo de apertura |
| `--limite N` | cuantos archivos imprime (por defecto 60) |

Codigo de salida: `0` sin hallazgos o tras `--fix`, `1` con hallazgos. Sirve
para un hook o un paso de CI.

Datos:

- `diccionario.tsv` - pares `incorrecto -> correcto` con modo `auto` o `manual`.
- `frases.tsv` - patrones de varias palabras (`no esta` -> `no está`,
  `Estas seguro` -> `¿Estás seguro`).
- `reglas.md` - las reglas de espanol y los casos que el script NO puede decidir.

### auditar.py: cuando el diccionario del skill no alcanza

`diccionario.tsv` son ~450 pares elegidos a mano. **No es el espanol.** Para
saber que se le escapa, `auditar.py` contrasta cada palabra del texto visible
contra un diccionario hunspell de verdad y propone las que se vuelven validas
al ponerles una tilde.

```bash
pip install spylls
cd .claude/skills/ortografia-front
curl -sSLO https://raw.githubusercontent.com/wooorm/dictionaries/main/dictionaries/es/index.dic
curl -sSLO https://raw.githubusercontent.com/wooorm/dictionaries/main/dictionaries/es/index.aff
cd - && python3 .claude/skills/ortografia-front/auditar.py
```

No corrige: imprime candidatos. Se revisan a mano y los reales se agregan a
`diccionario.tsv`. Buena parte son falsos positivos, porque el diccionario
espanol "arregla" palabras tecnicas o en ingles: `min` -> `mín`, `items` ->
`ítems`, `COD` -> `CÓD`, `ROI` -> `ROÍ`, `Record` -> `Récord`.

Correrlo despues de cada limpieza grande. En la primera pasada real encontro
122 formas (461 ocurrencias) que el diccionario a mano no veia.

## Como se escribe la tilde (regla de UTF-8 del repo)

`CLAUDE.md` prohibe caracteres non-ASCII en archivos de 500+ lineas por el bug
de highlight.js. Eso NO es una excusa para dejar el texto sin tilde: cambia
**como se escribe**, no **si se escribe**.

| Archivo | Como |
|---|---|
| menos de 500 lineas | caracter literal: `<label>Código</label>`, `'Teléfono'` |
| 500 lineas o mas | escape `\u00XX` |

En un archivo largo:

```tsx
showToast('Conexi\u00f3n exitosa', 'success');    // string: escape directo
<span>{'Facturaci\u00f3n'}</span>                 // texto JSX: envuelto en {'...'}
```

`revisar.py --fix` elige la forma sola segun el largo del archivo y marca los
archivos ASCII en el reporte con `[ASCII]`.

| letra | escape | letra | escape | letra | escape |
|---|---|---|---|---|---|
| á | `\u00e1` | ó | `\u00f3` | Á | `\u00c1` |
| é | `\u00e9` | ú | `\u00fa` | Ó | `\u00d3` |
| í | `\u00ed` | ñ | `\u00f1` | Ñ | `\u00d1` |
| ü | `\u00fc` | ¿ | `\u00bf` | ¡ | `\u00a1` |

**No usar** el patron `A{'ñ'}adir` que quedo en algunos archivos: mezcla
literal y escape, es ilegible y ademas deja non-ASCII en el archivo. Se
normaliza a `{'Añadir'}`.

## Flujo de trabajo

### Caso 1: escribiendo texto nuevo

No hace falta correr nada. Escribir el texto bien desde el principio, con la
forma que corresponda al largo del archivo, y consultar `reglas.md` ante la
duda. Al terminar, verificar el archivo tocado:

```bash
python3 .claude/skills/ortografia-front/revisar.py front/central/src/services/modules/orders
```

### Caso 2: limpiando ortografia existente

1. **Delimitar el alcance.** Un modulo por vez, nunca el front entero de una:
   un diff de 1.700 lineas no lo revisa nadie.

   ```bash
   python3 .claude/skills/ortografia-front/revisar.py front/central/src/services/modules/customers
   ```

2. **Leer el reporte** antes de corregir. Confirmar que los hallazgos son texto
   visible y no claves ni identificadores que el filtro dejo pasar.

3. **Aplicar** los seguros:

   ```bash
   python3 .claude/skills/ortografia-front/revisar.py front/central/src/services/modules/customers --fix
   ```

4. **Revisar el diff completo.** `git diff` linea por linea. Es correccion de
   texto: cualquier cambio fuera de un string o de un nodo de texto JSX es un
   bug del script y hay que revertirlo.

5. **Correr los tests.** El skill no toca los `*.test.*`, pero algunos
   **afirman sobre el texto de la UI**: si el codigo ahora dice
   `'No había productos por aplicar'` y el test espera `'No habia ...'`,
   el test falla y en este repo eso **bloquea el deploy**. Paso obligatorio,
   no opcional; ya rompio un build de `main` una vez.

   ```bash
   cd front/central && pnpm test
   ```

   Se corrige la **asercion** para que refleje el texto nuevo. Las
   descripciones de los `it(...)` y los datos mock (`'Juan Perez'`) se dejan
   como estan: no comparan contra nada y solo ensucian el diff.

6. **Verificar que compila y que no entro non-ASCII donde no debe:**

   ```bash
   cd front/central && npx tsc --noEmit
   for f in $(git diff --name-only); do
     n=$(wc -l < "$f"); if [ "$n" -ge 500 ] && grep -qP '[^\x00-\x7F]' "$f"; then
       echo "NON-ASCII en archivo largo: $f"; fi
   done
   ```

7. **Los `manual` a mano.** Los que el reporte marca `rev` (tilde diacritica,
   signos de apertura) se deciden leyendo la frase, uno por uno. Ver `reglas.md`.

8. Commit aparte, solo ortografia, sin mezclar con cambios de logica.

## Lo que el script NO decide

Estos casos se resuelven leyendo, y son la parte que de verdad importa:

- **Tilde diacritica.** `esta`/`está`, `el`/`él`, `si`/`sí`, `mas`/`más`,
  `aun`/`aún`, `solo` (hoy va SIN tilde siempre). Correr con `--ambiguos`.
- **Signos de apertura.** `?` y `!` van en pareja: `¿Eliminar este ticket?`.
  El script los detecta pero no los inserta, porque hay que ubicar donde
  empieza la pregunta.
- **Pasado vs presente.** `genero`/`generó`, `cambio`/`cambió`,
  `sincronizo`/`sincronizó`.
- **Terminologia y tono.** Mayusculas, ingles suelto, tuteo vs usted:
  `reglas.md`.

## Limitaciones conocidas del detector

Prefiere callarse antes que reportar codigo. Por eso NO revisa:

- Cadenas de una sola palabra en minuscula (`'telefono'`): son casi siempre
  claves, ids o valores enviados al backend. Si una de esas es un
  `placeholder` visible, se corrige a mano.
- Valores de atributos tecnicos (`className`, `href`, `fill`, `style`, ...).
- Interpolaciones de template (`${...}`): no revisa los signos de apertura ahi.
- Comentarios: no son texto de usuario y pueden quedarse sin tilde.
- Archivos `*.test.*`, `*.spec.*`, `*.stories.*`, `__tests__/`, `__mocks__/`:
  sus cadenas son fixtures y aserciones, no UI.
- **Literales protegidos.** Antes de revisar, el script recorre el front y
  junta toda cadena usada como identificador: argumentos de `hasPermission()`
  y valores de `resource:`, `permission:`, `code:`, `slug:`, `key:`, `event:`,
  `queue:`, `status:`, `provider:`. Una cadena igual a alguna de esas no se
  toca **en ningun lado**, aunque parezca un label.

  Esto no es teorico: los tours declaran `resource: 'Envios'` y
  `TourProvider` lo pasa a `hasPermission(resource, 'Read')`. Ponerle la tilde
  hubiera dejado el tour invisible para todo usuario no super admin. Lo mismo
  con `'Ordenes'`, `'Facturacion'`, `'Integraciones-Tipos-de-integracion'`:
  son filas de la base, no texto.

  Efecto de borde aceptado: si un label visible coincide exactamente con un
  identificador (`'Envios'` como `title` de un tour), tampoco se corrige.
  Se decide a mano.

Estado al crear el skill: `front/website` esta limpio; los ~1.760 hallazgos
estan todos en `front/central/src`.

## Prohibido

- Poner tilde en un identificador, clave de objeto, ruta, nombre de archivo,
  clase de Tailwind o valor enviado al backend. Solo texto visible.
- Correr `--fix` sobre el front entero en un solo commit.
- Dejar non-ASCII en un archivo de 500+ lineas (rompe highlight.js).
- "Corregir" `solo` a `sólo`: la RAE ya no lleva tilde ahi.
- Cambiar el texto de un mensaje mientras se le pone la tilde. Ortografia y
  redaccion son commits distintos.
