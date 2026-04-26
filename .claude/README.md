# Configuración de Claude Code para este proyecto

## Configuración MCP

Este proyecto usa **configuración MCP local** (`.claude/mcp.json`), no la global.

### Reglas importantes:

1. **Siempre usar `.claude/mcp.json`** para configurar servidores MCP
2. **NO usar** `~/.config/claude/mcp.json` (archivo global vacío)
3. Cada proyecto tiene su propia configuración y tokens

### Servidores MCP configurados en este proyecto:

- **postgres-reserve**: Servidor PostgreSQL de producción
- **github**: Integración con GitHub (token específico del proyecto)
- **puppeteer**: Automatización de navegador

### Cómo actualizar la configuración MCP:

```bash
# Editar configuración local
vim .claude/mcp.json

# Reiniciar Claude Code para aplicar cambios
# (Ctrl+C y luego reiniciar)
```

### Verificar configuración activa:

```bash
claude mcp list
```

## Permisos

Los permisos específicos del proyecto se configuran en `.claude/settings.local.json`.
