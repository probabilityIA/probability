# API key de OpenRouter hardcodeada en el codigo

Fecha: 2026-08-05
Detectado al escribir los tests de `modules/ai`.

## Contexto

`back/central/services/modules/ai/internal/infra/secondary/openrouter/client.go:18`
declara la credencial como constante publica del paquete:

```go
const (
    APIKey = "sk-or-v1-d4353...a506"
    APIURL = "https://openrouter.ai/api/v1/chat/completions"
    Model  = "xiaomi/mimo-v2-flash:free"
)
```

La key esta commiteada en el historial de git, o sea que ya se considera
comprometida aunque el repo sea privado: cualquiera con acceso al clon (o a un
fork, o a un backup del repo) puede gastar el saldo de la cuenta de OpenRouter.

## Items

### Urgente

- [ ] Revocar la key en https://openrouter.ai/keys. Es lo unico que corta el
      riesgo; reemplazarla en el codigo sin revocar no sirve porque la vieja
      sigue viva en el historial.
- [ ] Emitir una key nueva y moverla a variable de entorno
      (`OPENROUTER_API_KEY`), leida via `env.IConfig` como el resto de los
      secretos del backend. La constante `APIKey` debe desaparecer.
- [ ] Cargar la variable en `infra/compose-prod` y en los secrets del workflow
      de deploy.

### Importante

- [ ] `APIURL` tambien es constante, asi que el cliente no es inyectable y no se
      puede testear la respuesta del modelo con un `httptest.Server`. Pasar la
      URL base por constructor deja el parseo (limpieza de markdown, choices
      vacios, JSON invalido) cubierto por tests.
- [ ] Barrer el resto del repo por credenciales en codigo antes de cerrar la
      alerta. Hoy esta es la unica coincidencia de `sk-or-v1|AKIA|sk-ant-` en
      archivos `.go` fuera de tests.

### Deseable

- [ ] Considerar limpiar el historial (`git filter-repo`) o al menos dejar
      registrado que ese commit contiene un secreto revocado.

## Criterio para cerrar

La key vieja esta revocada, la nueva vive solo en variables de entorno, y
`grep -rn "sk-or-v1" back/` no devuelve nada.
