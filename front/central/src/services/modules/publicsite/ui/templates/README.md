# Sistema de Templates — Sitio Web Publico

## Como funciona

Cada negocio tiene un campo `template` en su `BusinessWebsiteConfig` que determina que set de componentes React se usa para renderizar su sitio web publico (`/tienda/{slug}`).

### Flujo

```
1. Usuario visita /tienda/mi-negocio
2. layout.tsx carga el business desde la API
3. Lee config.template (ej: "default")
4. getTemplate("default") retorna el set de componentes
5. Layout usa template.Nav, template.Footer, template.WhatsAppButton
6. page.tsx usa template.HeroSection, template.AboutSection, etc.
7. productos/page.tsx usa template.ProductCard
8. contacto/page.tsx usa template.ContactSection
```

### Archivos clave

```
publicsite/ui/templates/
├── types.ts          ← Contrato: interface TemplateComponents
├── registry.ts       ← Registro: mapea nombres a templates
├── default/
│   └── index.ts      ← Re-exporta los componentes existentes
└── README.md         ← Este archivo
```

---

## Crear un template nuevo

### Paso 1: Crear la carpeta

```
templates/mi-template/
├── index.ts
├── Nav.tsx
├── Footer.tsx
├── HeroSection.tsx
├── AboutSection.tsx
├── FeaturedProducts.tsx
├── TestimonialsSection.tsx
├── LocationSection.tsx
├── ContactSection.tsx
├── SocialMediaLinks.tsx
├── WhatsAppButton.tsx
└── ProductCard.tsx
```

### Paso 2: Implementar los componentes

Cada componente recibe los mismos props que los del template default. Los tipos estan definidos en `publicsite/domain/types.ts`.

```tsx
// templates/mi-template/HeroSection.tsx
import { PublicBusiness, HeroContent } from '../../../domain/types';

interface Props {
    content: HeroContent | null;
    business: PublicBusiness;
    slug: string;
}

export function HeroSection({ content, business, slug }: Props) {
    // Tu diseno completamente personalizado
    return (
        <section className="relative h-screen flex items-center">
            {/* ... */}
        </section>
    );
}
```

**Variables CSS disponibles:**

Los colores de marca del negocio se inyectan como CSS variables en el layout:

```css
var(--brand-primary)     /* Color principal (header, footer) */
var(--brand-secondary)   /* Botones, CTAs, acentos */
var(--brand-tertiary)    /* Color terciario */
var(--brand-quaternary)  /* Color cuaternario */
```

### Paso 3: Crear el barrel (index.ts)

```typescript
// templates/mi-template/index.ts
import { TemplateComponents } from '../types';
import { Nav } from './Nav';
import { Footer } from './Footer';
// ... importar todos los componentes

export const miTemplate: TemplateComponents = {
    Nav,
    Footer,
    WhatsAppButton,
    HeroSection,
    AboutSection,
    FeaturedProducts,
    TestimonialsSection,
    LocationSection,
    ContactSection,
    SocialMediaLinks,
    ProductCard,
};
```

### Paso 4: Registrar en registry.ts

```typescript
// templates/registry.ts
import { miTemplate } from './mi-template';

const templates: Record<string, TemplateComponents> = {
    default: defaultTemplate,
    'mi-template': miTemplate,  // <-- agregar aqui
};

export function getAvailableTemplates() {
    return [
        { id: 'default', name: 'Clasico', description: 'Diseno limpio y profesional' },
        { id: 'mi-template', name: 'Mi Template', description: 'Descripcion del template' },
    ];
}
```

### Paso 5: Probar

1. Ir a **Mi Sitio Web** en el panel admin
2. Seleccionar "Mi Template" en el selector de plantilla
3. Guardar
4. Abrir **/tienda/{slug}** y verificar que usa los componentes nuevos

---

## Reutilizar componentes del template default

No es necesario crear todos los componentes desde cero. Se pueden reutilizar componentes del default:

```typescript
// templates/mi-template/index.ts
import { TemplateComponents } from '../types';
import { defaultTemplate } from '../default';
import { HeroSection } from './HeroSection';  // solo este es custom

export const miTemplate: TemplateComponents = {
    ...defaultTemplate,          // reutilizar todo el default
    HeroSection,                 // sobreescribir solo el hero
};
```

---

## HomePage override (control total)

Si un template necesita control total del homepage (sin renderizado seccion por seccion), puede proveer un componente `HomePage`:

```typescript
export const miTemplate: TemplateComponents = {
    ...defaultTemplate,
    HomePage: MiHomePage,  // Toma control total del homepage
};
```

```tsx
// templates/mi-template/HomePage.tsx
import { PublicBusiness, WebsiteConfig } from '../../../domain/types';

interface Props {
    business: PublicBusiness;
    slug: string;
    config: WebsiteConfig;
}

export function MiHomePage({ business, slug, config }: Props) {
    // Diseno completamente libre — no usa el renderizado por secciones
    return (
        <div>
            {/* Lo que quieras */}
        </div>
    );
}
```

Cuando `HomePage` esta presente, `page.tsx` lo usa en lugar del renderizado seccion por seccion.

---

## Contrato de componentes (TemplateComponents)

```typescript
interface TemplateComponents {
    Nav: (props: { business: PublicBusiness }) => JSX.Element;
    Footer: (props: { business: PublicBusiness }) => JSX.Element;
    WhatsAppButton: (props: { content: WhatsAppContent }) => JSX.Element;

    HeroSection: (props: { content: HeroContent | null; business: PublicBusiness; slug: string }) => JSX.Element;
    AboutSection: (props: { content: AboutContent }) => JSX.Element;
    FeaturedProducts: (props: { products: PublicProduct[]; slug: string }) => JSX.Element;
    TestimonialsSection: (props: { content: Testimonial[] }) => JSX.Element;
    LocationSection: (props: { content: LocationContent }) => JSX.Element;
    ContactSection: (props: { slug: string; content: ContactContent | null }) => JSX.Element;
    SocialMediaLinks: (props: { content: SocialMediaContent }) => JSX.Element;
    ProductCard: (props: { product: PublicProduct; slug: string }) => JSX.Element;

    HomePage?: (props: { business: PublicBusiness; slug: string; config: WebsiteConfig }) => JSX.Element;
}
```

---

## Notas importantes

- **Fallback seguro**: Si un negocio tiene un template que ya no existe en el codigo, se usa "default" automaticamente.
- **Server Components**: Los componentes de layout (Nav, Footer) y secciones son Server Components por defecto. Si necesitas interactividad, agregar `'use client'` al inicio del archivo (como ContactSection y WhatsAppButton).
- **Componentes compartidos**: CatalogSearch y CatalogPagination NO son parte del template — son compartidos entre todos los templates. Solo los componentes listados en `TemplateComponents` son intercambiables.
- **Datos del negocio**: Todos los templates reciben los mismos datos. La diferencia es solo visual/de presentacion.
