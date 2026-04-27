# 🏗️ Configuración de Claude Code - Estructura

## Flujo de Configuración

```
Tu Máquina (~/ = tu home directory)
│
├── ~/.claude/
│   ├── settings.json          ← Preferencias globales de Claude
│   ├── mcp.json               ← MCPs globales (si los usas en todos los proyectos)
│   └── mcp-postgres.js        ← Script Node.js para MCP de PostgreSQL
│
└── ~/.claude.json             ← ⭐ TUS MCPs (AQUÍ ES DONDE ESTÁN AHORA)
                                 chromium-devtools
                                 playwright
                                 postgres-probability


Carpeta del Proyecto (probability/)
│
├── .claude/                   ← PROYECTOS Y REFERENCIAS LOCALES
│   ├── settings.json          ← Configuración del proyecto (NO COMMITEADA)
│   ├── mcp.json               ← Referencia de MCPs (NO COMMITEADA)
│   ├── mcp-config-template.json ← Template para otros devs (SÍ COMMITEADA)
│   ├── MCP_SETUP_GUIDE.md     ← Guía de configuración (SÍ COMMITEADA)
│   └── STRUCTURE.md           ← Este archivo (SÍ COMMITEADA)
│
├── .nvmrc                     ← Node.js 20.19.0 (SÍ COMMITEADA)
├── .gitignore                 ← Contiene: .claude.json, .claude/mcp.json, .claude/settings.json
├── CLAUDE.md                  ← Documentación principal del proyecto (SÍ COMMITEADA)
└── ... resto del proyecto
```

---

## En .gitignore

```
# Local Claude Code configuration (each developer has their own)
.claude/mcp.json
.claude/settings.json
.claude.json
```

**¿Por qué?**
- Cada developer tiene sus propios MCPs basados en su máquina local
- Credenciales de base de datos nunca se suben (DATABASE_URL)
- Configuración personal no necesita ser compartida

---

## Archivos SÍNCOMMITEADOS (Para Todo el Equipo)

| Archivo | Propósito |
|---------|-----------|
| `.nvmrc` | Especifica Node.js 20.19.0 para todos |
| `.claude/MCP_SETUP_GUIDE.md` | Guía para nuevos MCPs |
| `.claude/mcp-config-template.json` | Template de referencia |
| `CLAUDE.md` | Documentación del proyecto |

---

## Configuración Por Developer

### Daniel Camacho (TÚ)
**Ubicación**: `~/.claude.json` (NO en el repo)

```json
{
  "mcpServers": {
    "chrome-devtools": { ... },
    "playwright": { ... },
    "postgres-probability": { ... }
  }
}
```

### Otro Developer (Juan, etc.)
1. Clona el proyecto
2. Ejecuta: `nvm use` (carga Node.js 20.19.0 del .nvmrc)
3. Lee `.claude/MCP_SETUP_GUIDE.md`
4. Copia `mcp-config-template.json` a su `~/.claude.json`
5. Ajusta DATABASE_URL con sus credenciales locales
6. Reinicia Claude Code
7. ¡Listo!

---

## Cómo Cambiará Para Otros Developers

```bash
# Developer #2 (Juan)
cd probability
nvm use                                    # → Node 20.19.0 ✓

# Ve la guía
cat .claude/MCP_SETUP_GUIDE.md

# Copia el template
cat .claude/mcp-config-template.json       # → Referencia

# Agrega a su ~/.claude.json personal
{
  "mcpServers": {
    "chrome-devtools": { ... },            # Mismo
    "playwright": { ... },                 # Mismo
    "postgres-probability": {              # AJUSTA CREDENCIALES
      "command": "node",
      "args": ["/Users/juan/.claude/..."], # Juan's path
      "env": { "DATABASE_URL": "..." }     # Juan's DB creds
    }
  }
}

# Reinicia Claude Code
# ¡Listo!
```

---

## Ventajas de Esta Estructura

✅ **Seguridad**: Credenciales nunca se suben al repo  
✅ **Flexibilidad**: Cada dev puede agregar más MCPs locales sin afectar a otros  
✅ **Documentación**: Guía clara para nuevos developers  
✅ **Consistencia**: Node.js 20.19.0 forzado para todos con `.nvmrc`  
✅ **Mantenible**: Solo archivos esenciales en el repo  

---

## Quick Reference

| Pregunta | Respuesta |
|----------|-----------|
| ¿Dónde están mis MCPs? | `~/.claude.json` (global, NO en el repo) |
| ¿Qué versión de Node.js? | 20.19.0 (especificado en `.nvmrc`) |
| ¿Cómo agrego un MCP? | Lee `.claude/MCP_SETUP_GUIDE.md` |
| ¿Qué le digo a otros devs? | "Clona y sigue `.claude/MCP_SETUP_GUIDE.md`" |
| ¿Se sube al repo? | NUNCA MCPs / SÍ plantillas y guías |
