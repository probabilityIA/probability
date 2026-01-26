Crea una nueva nota de tarea en `.notes/[nombre-tarea].md` usando la plantilla estándar definida en `.cursor/rules/rules.md` (sección "Plantilla Estándar de Tareas").

El archivo debe seguir exactamente esta estructura:

```markdown
# 📝 [Nombre de la Tarea]

## 🎯 Objetivo General
[Descripción breve de qué queremos lograr y por qué]

---

## 🛠 Contexto Técnico

**Ficheros Involucrados**: @archivo1.go, @archivo2.tsx (Usa @ para que Cursor los identifique)

**Dependencias/Herramientas**: 
- [Lista de dependencias, APIs, herramientas necesarias]

**Arquitectura**: 
- Arquitectura Hexagonal (Domain/Application/Infrastructure)

**Integración relacionada**: 
- [Nombre de la integración si aplica: WhatsApp, Shopify, etc.]

---

## 📋 Plan de Ejecución (Paso a Paso)

Usa este checklist para que la IA sepa dónde estamos.

### [ ] Paso 1: Análisis y Preparación
- [ ] Revisar lógica actual en @archivo_relevante.go
- [ ] Definir contratos/interfaces en domain
- [ ] Identificar dependencias necesarias

### [ ] Paso 2: Implementación
- [ ] Escribir lógica de negocio en el dominio (domain/)
- [ ] Crear casos de uso en application/
- [ ] Crear adaptador/repositorio en infrastructure/
- [ ] Implementar handlers HTTP si aplica

### [ ] Paso 3: Testing
- [ ] Test unitarios de domain
- [ ] Test unitarios de casos de uso (con mocks)
- [ ] Test de handlers/infraestructura
- [ ] Test de integración si aplica

### [ ] Paso 4: Validación y Documentación
- [ ] Verificación visual o de API
- [ ] Actualizar documentación si es necesario
- [ ] Revisar que sigue arquitectura hexagonal
- [ ] Verificar principios SOLID

---

## 🚦 Estado de la Tarea

**Progreso actual**: 0%

**Bloqueos**: Ninguno

**Último cambio realizado**: Ninguno

**Fecha de inicio**: [fecha actual]
**Fecha estimada de finalización**: [fecha estimada]

---

## 🧠 Memoria de Decisiones

Anota aquí por qué decidiste hacer algo de una forma específica para que la IA no te proponga cambiarlo después.

### Decisión 1: [Título]
- **Qué**: Descripción de la decisión
- **Por qué**: Razón técnica o de negocio
- **Alternativas consideradas**: Qué otras opciones se evaluaron

---

## 📌 Notas de Cierre / Próximos Pasos

- [ ] Tarea pendiente relacionada 1
- [ ] Tarea pendiente relacionada 2

**Observaciones finales**: [Notas adicionales cuando se complete la tarea]
```

**Instrucciones**:
1. Reemplaza `[Nombre de la Tarea]` con el nombre descriptivo de la tarea
2. Completa el "Objetivo General" con la descripción proporcionada por el usuario
3. Identifica y lista los archivos involucrados usando @ para que Cursor los identifique
4. Completa las dependencias y herramientas necesarias
5. Si es una tarea relacionada con una integración específica, indícalo
6. Deja los checklists sin marcar (serán marcados durante el desarrollo)
7. Establece la fecha de inicio como la fecha actual
8. Deja "Memoria de Decisiones" vacío (se llenará durante el desarrollo)

El archivo debe crearse en `.notes/[nombre-tarea].md` donde `[nombre-tarea]` es un nombre descriptivo en minúsculas con guiones (ej: `log-3823-whatsapp-integration.md`).
