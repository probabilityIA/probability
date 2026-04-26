# Lecciones Playwright - Pruebas Front Orders (produccion)

Aprendido durante prueba E2E de creacion de ordenes + cambio de estado en produccion (2026-04-24).

---

## 1. Dos botones "Siguiente" en la misma pagina

**El problema:** La pagina `/orders` tiene dos botones con texto "Siguiente":
- El de paginacion de la tabla (`class="btn btn-outline btn-sm"`)
- El del modal de cambio de estado (`class="px-5 py-2 ... bg-purple-600 ..."`)

`document.querySelector('button')` con `textContent === 'Siguiente'` siempre atrapa el primero (paginacion) y el modal nunca avanza.

**Solucion correcta:**
```javascript
const btn = Array.from(document.querySelectorAll('button')).find(
    b => b.textContent.trim() === 'Siguiente' && b.className.includes('bg-purple-600')
);
btn.click();
```

---

## 2. ProductSelector busca por nombre, no por SKU

**El problema:** Escribir el SKU ("PT01001") en el buscador del formulario de orden no devuelve resultados. El componente `ProductSelector.tsx` solo pasa `name: term` al endpoint, nunca `sku`.

**Solucion:** Buscar por una palabra del nombre del producto. Ejemplos:
- "PT01001" → sin resultados
- "Aislada" → encuentra "Proteina Aislada (ISO) - 2 Lb..."
- "Hidrolizado" → encuentra "Colageno Hidrolizado - 300g..."
- Usar palabras sin acento si el nombre tiene tilde (el input puede perderse con acentos)

---

## 3. Dropdown de productos aparece fuera del viewport

**El problema:** El dropdown de resultados del `ProductSelector` usa `position: absolute; z-10`. Cuando el campo de busqueda esta en la parte inferior del modal, el dropdown se renderiza hacia abajo y queda fuera de la pantalla visible. Los clicks normales fallan.

**Solucion:** Usar JS para clicar en el DOM directamente, sin importar visibilidad:
```javascript
const items = document.querySelectorAll('ul li button, [class*="z-10"] button, [class*="absolute"] button');
items[0].click(); // primer resultado
```

O mas especifico:
```javascript
const lis = document.querySelectorAll('[class*="absolute"] ul li button');
lis[0].click();
```

---

## 4. Inputs controlados por React no responden a value= directo

**El problema:** Escribir directamente `input.value = 'texto'` no dispara el re-render de React. El campo queda vacio o con valor previo.

**Solucion:** Usar el native value setter + evento `input`:
```javascript
const nativeSet = Object.getOwnPropertyDescriptor(window.HTMLInputElement.prototype, 'value').set;
nativeSet.call(input, 'nuevo valor');
input.dispatchEvent(new Event('input', { bubbles: true }));
input.dispatchEvent(new Event('change', { bubbles: true }));
```

---

## 5. pressSequentially duplica texto en inputs ya focuseados

**El problema:** `browser_type` o `pressSequentially` agrega texto al final del valor existente si el input ya tenia contenido (ej: el buscador de producto ya tenia "PT01001" y escribe "PT01001PT01001").

**Solucion:** Limpiar el campo antes con el native setter o seleccionar todo y reemplazar:
```javascript
nativeSet.call(input, '');
input.dispatchEvent(new Event('input', { bubbles: true }));
// luego escribir el nuevo valor
nativeSet.call(input, 'nuevo valor');
input.dispatchEvent(new Event('input', { bubbles: true }));
```

---

## 6. El business selector se resetea al confirmar cambio de estado

**El problema:** Cada vez que se confirma un cambio de estado en el modal, la pagina recarga la lista de ordenes y el selector de negocio vuelve al valor por defecto (vacio o el primero). Las ordenes del business filtrado desaparecen.

**Solucion:** Despues de cada `Confirmar Cambio`, siempre re-ejecutar:
```javascript
const sel = Array.from(document.querySelectorAll('select')).find(
    s => Array.from(s.options).some(o => o.text.includes('Demo'))
);
const opt = Array.from(sel.options).find(o => o.text.includes('Demo'));
sel.value = opt.value;
sel.dispatchEvent(new Event('change', { bubbles: true }));
```

---

## 7. Cards de estado se seleccionan con button.flex.items-start

**El problema:** Las opciones del modal "Cambiar Estado" son `<button class="flex items-start gap-3 p-3 rounded-lg border-2 ...">`. No tienen `role="option"` ni atributos especiales facilmente targeteable.

**Solucion correcta:**
```javascript
const cards = Array.from(document.querySelectorAll('button.flex.items-start'));
const target = cards.find(c => c.textContent.includes('Entregada'));
target.click();
```

---

## 8. El estado "Asignado a piloto" requiere campo adicional

**El problema:** Al seleccionar la opcion "Asignado a piloto", el modal muestra un campo extra requerido: "Nombre del piloto *". Clickear "Siguiente" sin llenarlo no avanza.

**Solucion:** Detectar si aparece un input adicional antes de clickear Siguiente:
```javascript
const pilotInput = document.querySelector('input[placeholder*="Carlos Ramirez"]');
if (pilotInput) {
    const nativeSet = Object.getOwnPropertyDescriptor(window.HTMLInputElement.prototype, 'value').set;
    nativeSet.call(pilotInput, 'Carlos Ramirez');
    pilotInput.dispatchEvent(new Event('input', { bubbles: true }));
}
```

---

## 9. Bodega (warehouse) es un select nativo con valores numericos

**El problema:** El campo Bodega en el formulario de orden es un `<select>` nativo. Su option "Principal" tiene `value="3"` (ID numerico), no el texto "principal". `select_option` con `values=["principal"]` falla.

**Solucion:** Consultar el valor real del option antes:
```javascript
const sel = document.querySelector('select[name="warehouse_id"]'); // o el que corresponda
const opt = Array.from(sel.options).find(o => o.text.toLowerCase().includes('principal'));
// opt.value === "3"
sel.value = opt.value;
sel.dispatchEvent(new Event('change', { bubbles: true }));
```

---

## 10. Flujo completo de estados para llegar a Entregada

Desde "Pendiente", la cadena de estados disponibles en produccion (Demo) es:

```
Pendiente
  -> Seleccionando productos
    -> Empacando
      -> Listo para despacho
        -> Asignado a piloto  (* requiere nombre del piloto)
          -> Recogido
            -> En camino
              -> En reparto final
                -> Entregada  <- confirm_sale disparado aqui
```

Cada paso = 3 acciones: click card + click Siguiente (bg-purple-600) + click Confirmar Cambio.

---

## 11. Verificar confirm_sale en logs de produccion

Despues de llegar a "Entregada", confirmar en logs:

```bash
ssh -i ".../probability.pem" ubuntu@ec2-3-224-189-33.compute-1.amazonaws.com \
  "cd /home/ubuntu/probability/infra/compose-prod && docker compose logs --tail 80 back-central 2>&1 | grep -E '(confirm_sale|Sale confirmed|inventory.confirmed)'"
```

Lineas esperadas:
```
Sale confirmed for order <uuid> ... success=true
Message published ... routing_key=inventory.confirmed
```

---

## 12. Super Admin debe seleccionar negocio antes de crear ordenes

**El problema:** Si no se selecciona un negocio en el combo "SUPER ADMIN", el formulario "Nueva Orden" no tiene `business_id` y la creacion falla o crea en el negocio incorrecto.

**Solucion:** Siempre seleccionar el negocio antes de abrir el modal:
```javascript
const sel = Array.from(document.querySelectorAll('select')).find(
    s => Array.from(s.options).some(o => o.text.includes('Demo'))
);
const opt = Array.from(sel.options).find(o => o.text.includes('Demo'));
sel.value = opt.value; // "26"
sel.dispatchEvent(new Event('change', { bubbles: true }));
```
