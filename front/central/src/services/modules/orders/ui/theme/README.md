# OrderDetails Theme

Carpeta reservada para futuros temas dinámicos de OrderDetails.

## Nota

El componente OrderDetails actualmente usa estilos estándar de Tailwind CSS.

Para implementar colores personalizados por business en el futuro:
1. Crear un hook que lea colores de `business.primary_color`, etc.
2. Usar CSS variables en el elemento root del modal
3. Aplicar las variables en lugar de clases estáticas
