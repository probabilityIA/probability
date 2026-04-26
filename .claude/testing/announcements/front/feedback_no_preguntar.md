---
name: no-sobrecomplicar-pruebas
description: Para pruebas E2E con Playwright MCP, usar las herramientas directamente sin delegar a agentes
type: feedback
---

Cuando el usuario pide pruebas E2E con Playwright, usar los tools MCP Playwright directamente. No delegar a agentes debug-tester ni crear flujos complejos de inicio de servicios.

**Why:** El usuario espera que simplemente se abra el browser y se ejecuten los tests. Los servicios ya pueden estar corriendo o el usuario los maneja por su cuenta.

**How to apply:** Usar mcp__playwright__browser_navigate, mcp__playwright__browser_snapshot, etc. directamente.
