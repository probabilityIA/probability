# 🔧 MCP (Model Context Protocol) Setup Guide

Esta guía explica cómo configurar MCPs nuevos para el proyecto Probability.

---

## 📋 Principios

- **Local Only**: Los MCPs se configuran localmente en `~/.claude.json` de cada desarrollador
- **No al Repo**: Los archivos `.claude/mcp.json` y `.claude.json` están en `.gitignore`
- **Por Desarrollador**: Cada developer mantiene sus propios MCPs (no compartidos)
- **Node.js**: El proyecto usa Node 20.19.0 (especificado en `.nvmrc`)

---

## ✅ MCPs Actualmente Configurados

| MCP | Propósito | Comando |
|-----|-----------|---------|
| `chrome-devtools` | Automatización y debugging del navegador | `npx chrome-devtools-mcp@latest` |
| `playwright` | Testing automatizado y web scraping | `npx @playwright/mcp@latest` |
| `postgres-probability` | Consultas directas a la BD | Node.js script local |

---

## 🚀 Cómo Agregar un Nuevo MCP

### Paso 1: Verifica que funciona localmente
```bash
# Asegúrate de estar en el proyecto
cd /Users/danielcamacho/development/probability
nvm use  # Carga Node 20.19.0 (del .nvmrc)

# Intenta ejecutar el MCP
npx nombre-del-mcp-package@latest --version
```

### Paso 2: Agrega a tu `~/.claude.json` global

```json
{
  "mcpServers": {
    "nombre-nuevo-mcp": {
      "command": "npx",
      "args": ["nombre-del-mcp-package@latest"]
    }
  }
}
```

**Para MCPs con configuración especial** (como postgres-probability con variables de entorno):

```json
{
  "mcpServers": {
    "mi-mcp-custom": {
      "command": "node",
      "args": ["/ruta/absoluta/al/script.js"],
      "env": {
        "VAR_IMPORTANTE": "valor"
      }
    }
  }
}
```

### Paso 3: Reinicia Claude Code
- Cierra Claude Code completamente
- Reabre Claude Code
- Ejecuta `/mcp` para verificar que aparece el nuevo MCP

### Paso 4: Documenta (ESTE ARCHIVO)
- Agrega el MCP a la tabla de MCPs Configurados arriba
- Describe su propósito en una línea

---

## 📁 Archivos Relacionados

```
probability/
├── .nvmrc                    # ← Node.js 20.19.0 (requerido)
├── .gitignore               # ← Contiene .claude.json, .claude/mcp.json
├── .claude/
│   ├── mcp.json            # ← NO COMMITEAR (local reference)
│   ├── settings.json       # ← NO COMMITEAR (local settings)
│   └── MCP_SETUP_GUIDE.md  # ← ESTE ARCHIVO
└── CLAUDE.md               # ← Documentación principal
```

**Tu archivo global** (NO en el repo):
```
~/.claude.json              # ← Donde están TUS MCPs
```

---

## 🔍 Troubleshooting

### El MCP no aparece en `/mcp`
1. ¿Reiniciaste Claude Code? (cierra completamente, no solo la sesión)
2. ¿Está en `~/.claude.json` (global), no en `.claude.json` (proyecto)?
3. ¿El JSON está bien formado? Valida con `jq .mcpServers ~/.claude.json`

### El MCP falla al conectar
1. ¿Funciona ejecutarlo manualmente? `npx nombre-del-mcp@latest --version`
2. ¿Tienes Node 20.19.0? Ejecuta `nvm use` en el proyecto
3. ¿Las rutas absolutas son correctas? (para scripts Node.js)

### Conflicto de versiones de Node.js
- Siempre ejecuta `nvm use` en el proyecto (carga `.nvmrc`)
- Si algo falló antes, reinicia la terminal/shell

---

## 📝 Ejemplo: Agregar un MCP para Redis

```bash
# 1. Prueba local
nvm use
npx redis-mcp@latest --version

# 2. Agrega a ~/.claude.json
{
  "mcpServers": {
    "redis": {
      "command": "npx",
      "args": ["redis-mcp@latest"],
      "env": {
        "REDIS_URL": "redis://localhost:6379"
      }
    }
  }
}

# 3. Reinicia Claude Code
# 4. Verifica: /mcp → deberías ver "redis · ✔ connected"
```

---

## 🔗 Referencias

- **Claude Code MCP Docs**: https://code.claude.com/docs/en/mcp
- **Node.js Version Management**: Con `nvm`, automáticamente carga la versión en `.nvmrc`
- **Este Proyecto**: CLAUDE.md contiene más detalles del stack
