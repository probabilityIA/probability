# probabilityIA - Landing Page

Landing page para probabilityIA, una plataforma de inteligencia artificial que predice y reduce devoluciones en tiempo real para eCommerce.

## 🚀 Tecnologías

- **Astro** - Framework web moderno
- **Tailwind CSS** - Framework de CSS utility-first
- **Preact** - Biblioteca JavaScript ligera para componentes interactivos
- **TypeScript** - Tipado estático

## 📁 Estructura del Proyecto

```
/
├── public/          # Archivos estáticos (imágenes, favicon, etc.)
├── src/
│   ├── components/ # Componentes reutilizables
│   │   ├── Header.astro
│   │   ├── HeroSection.astro
│   │   ├── StatsSection.astro
│   │   ├── ROICalculator.astro
│   │   ├── ROICalculator.tsx (componente interactivo)
│   │   ├── IntegrationsSection.astro
│   │   └── ContactSection.astro
│   ├── layouts/    # Layouts base
│   │   └── Layout.astro
│   ├── pages/      # Páginas/rutas
│   │   └── index.astro
│   └── styles/     # Estilos globales
│       └── global.css
└── package.json
```

## 🧞 Comandos

Todos los comandos se ejecutan desde la raíz del proyecto:

| Comando                | Acción                                           |
| :--------------------- | :----------------------------------------------- |
| `npm install`          | Instala las dependencias                         |
| `npm run dev`          | Inicia el servidor de desarrollo en `localhost:4321` |
| `npm run build`        | Construye el sitio para producción en `./dist/`  |
| `npm run preview`      | Previsualiza la build localmente                 |

## 🎨 Características

- **Diseño Responsive** - Optimizado para móviles, tablets y desktop
- **Calculadora ROI Interactiva** - Simulador de ahorro potencial con sliders
- **Secciones Incluidas**:
  - Hero Section con CTA
  - Estadísticas de impacto
  - Calculadora de ROI
  - Integraciones con plataformas eCommerce
  - Sección de contacto

## 🌐 Desarrollo

Para iniciar el servidor de desarrollo:

```bash
npm run dev
```

El sitio estará disponible en `http://localhost:4321`

## 📦 Build para Producción

```bash
npm run build
```

Los archivos estáticos se generarán en la carpeta `dist/`.
