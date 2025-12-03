# Notify Email Service

Servicio de notificaciones por correo electrónico.

## Estructura del Proyecto

```
notify-email/
├── cmd/
│   └── main.go          # Punto de entrada de la aplicación
├── internal/
│   ├── app/             # Casos de uso y lógica de negocio
│   ├── domain/          # Entidades y contratos del dominio
│   └── infra/           # Implementaciones de infraestructura
├── shared/              # Código compartido y utilidades
└── go.mod              # Dependencias del proyecto
```

## Requisitos

- Go 1.23.0 o superior

## Compilación

```bash
go build ./cmd/main.go
```

## Ejecución

```bash
./main
```

## Desarrollo

El proyecto sigue una arquitectura limpia con separación de capas:
- **Domain**: Define las entidades y contratos (puertos/interfaces)
- **App**: Contiene la lógica de negocio y casos de uso
- **Infra**: Implementa las interfaces definidas en domain (repositorios, servicios externos)

