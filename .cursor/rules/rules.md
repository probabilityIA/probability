# Reglas de Contexto - Connect Flow

## 📚 Contexto del Proyecto

**SIEMPRE** consulta el archivo `README.md` al inicio de cualquier tarea para entender:
- La arquitectura general del sistema
- Cómo funcionan las integraciones
- El SDK y sus componentes
- Flujos de sincronización (inventario, órdenes, estados)
- Guías de desarrollo y mejores prácticas
- Ejemplos y referencias de integraciones existentes

El `README.md` contiene el contexto completo del proyecto y se mantiene actualizado con la información más relevante.

**Connect Flow** es un **monolito de integraciones** que permite conectar Velocity con múltiples plataformas de ecommerce, facturación, mensajería y más. Cada integración se desarrolla de forma modular e independiente dentro del mismo repositorio.

### Características Principales

- **🔌 Sistema de Integración Modular**: Cada integración es independiente con su propia lógica backend y UI
- **🎯 SDK Unificado**: Core compartido que maneja autenticación, eventos, colas y más
- **🔄 Sincronización Bidireccional**: Productos, inventario, órdenes y estados
- **⚡ Sistema de Colas con NATS**: Procesamiento asíncrono y resiliente
- **🎨 UI Compartida**: Componentes reutilizables para todas las integraciones
- **📊 Rate Limiting**: Control de tasa de peticiones por integración
- **🔐 OAuth 2.0**: Flujos de autenticación seguros
- **📦 Webhooks**: Recepción de eventos en tiempo real

### Stack Tecnológico

- **Backend**: Go (Golang) 1.21+ con Echo framework
- **Frontend**: React 18+ con TypeScript + Vite
- **Base de Datos**: MySQL 8.0+ (via GORM)
- **Colas**: NATS JetStream
- **Scheduler**: gocron v2
- **UI Components**: shadcn/ui + Tailwind CSS

---

## 🏗️ Arquitectura

### Estructura General

```
connect_flow/
├── app/
│   ├── integrations/          # 🔌 Todas las integraciones
│   │   ├── tiendanube/        # Ejemplo de integración
│   │   ├── shopify/
│   │   ├── siigo/
│   │   ├── bsale/
│   │   ├── paris/
│   │   ├── whatsApp/
│   │   └── ...
│   │
│   └── shared/                # 📦 Código compartido
│       ├── sdk/               # SDK principal - Núcleo del sistema
│       ├── auth/              # Autenticación JWT
│       ├── models/            # Modelos de base de datos (GORM)
│       ├── lib/               # Librerías comunes
│       └── sharedRepository/  # Repository pattern
│
├── ui/                        # 🎨 Frontend React compartido
│   ├── components/            # Componentes UI reutilizables
│   ├── pages/                # Páginas principales
│   ├── lib/                  # Utilidades del frontend
│   └── integrations.ts       # Registro central de integraciones
│
├── docs/                      # 📚 Documentación
├── main.go                    # 🚪 Punto de entrada backend
└── .notes/                    # 📝 Tareas en curso y contexto
```

---

## 🎯 Arquitectura Hexagonal para Nuevas Integraciones

**IMPORTANTE**: Cada nueva integración en el backend **DEBE seguir arquitectura hexagonal**.

### Estructura de una Integración con Arquitectura Hexagonal

```
app/integrations/miintegracion/
│
├── internal/                  # Código backend privado
│   │
│   ├── application/          # Capa de Aplicación (Casos de Uso)
│   │   ├── usecase/          # Casos de uso específicos
│   │   │   ├── sync-inventory.go
│   │   │   ├── sync-orders.go
│   │   │   ├── sync-status.go
│   │   │   ├── webhook-handler.go
│   │   │   └── constructor.go
│   │   │
│   │   └── ports/            # Interfaces de casos de uso
│   │       └── IOrderIntegratorUseCase.go
│   │
│   ├── domain/               # Capa de Dominio
│   │   ├── entities/         # Entidades de negocio
│   │   │   ├── order.go
│   │   │   ├── product.go
│   │   │   └── status.go
│   │   │
│   │   ├── dtos/             # Data Transfer Objects
│   │   │   ├── dtos.go
│   │   │   ├── stock.go
│   │   │   └── order.go
│   │   │
│   │   ├── config/           # Configuración de dominio
│   │   │   └── config.go
│   │   │
│   │   ├── ports/            # Interfaces/Contratos (opcional)
│   │   │   └── repository.go
│   │   │
│   │   └── errors/           # Errores personalizados
│   │       └── errors.go
│   │
│   └── infrastructure/       # Capa de Infraestructura
│       │
│       ├── primary/          # Puertos Primarios (Entrada)
│       │   ├── handler/      # Handlers HTTP
│       │   │   ├── webhook.go
│       │   │   ├── sync-inventory.go
│       │   │   └── constructor.go
│       │   │
│       │   └── consumerNats/ # Consumers de NATS (opcional)
│       │       └── webhook-consumer.go
│       │
│       └── secondary/        # Puertos Secundarios (Salida)
│           ├── http/         # HTTP clients externos
│           │   ├── client.go
│           │   ├── get-order.go
│           │   ├── post-stock.go
│           │   ├── mappers/  # Mappers de datos
│           │   └── dtos/     # DTOs de request/response
│           │
│           └── repository/  # Repositorio de datos
│               └── repository.go
│
├── ui/                       # Componentes React específicos
│   ├── InstallView.tsx       # Vista de instalación/OAuth
│   ├── SettingsView.tsx      # Vista de configuración
│   ├── components/           # Componentes propios
│   ├── hooks/                # Hooks personalizados
│   ├── assets/               # Assets (logo, etc.)
│   └── index.ts              # Registro en el sistema UI
│
├── integrator.go             # Implementa sdk.Integrator interface
├── install.go                # Lógica de instalación
├── syncInventory.go          # Sincronización de inventario
├── syncOrders.go             # Sincronización de órdenes
├── syncStatus.go             # Sincronización de estados
├── routes.go                 # Rutas HTTP personalizadas (opcional)
└── README.md                 # Documentación específica
```

### Principios de Arquitectura Hexagonal

1. **Separación de Capas**:
   - **Domain**: Entidades, DTOs, interfaces - **NO depende de nada**
   - **Application**: Casos de uso - **Solo depende de Domain**
   - **Infrastructure**: Implementaciones - **Depende de Domain y Application**

2. **Dependencias**:
   - Las dependencias siempre apuntan **hacia adentro** (hacia Domain)
   - Domain **NO** debe importar de Application ni Infrastructure
   - Application **NO** debe importar de Infrastructure

3. **Puertos y Adaptadores**:
   - **Primary (Entrada)**: Handlers HTTP, consumers de NATS
   - **Secondary (Salida)**: HTTP clients, repositorios, bases de datos

4. **Casos de Uso**:
   - Contienen la lógica de negocio pura
   - Reciben interfaces (ports) como dependencias
   - No conocen detalles de implementación

### Principios SOLID

**SIEMPRE** respetar los principios SOLID en todo el código:

1. **Single Responsibility Principle (SRP)**:
   - Cada clase/función debe tener una única razón para cambiar
   - Separar responsabilidades claramente (validación, transformación, persistencia, etc.)
   - Casos de uso deben hacer UNA cosa y hacerla bien

2. **Open/Closed Principle (OCP)**:
   - Abierto para extensión, cerrado para modificación
   - Usar interfaces para permitir nuevas implementaciones sin modificar código existente
   - Los casos de uso deben ser extensibles mediante interfaces

3. **Liskov Substitution Principle (LSP)**:
   - Las implementaciones de interfaces deben ser intercambiables
   - Los mocks en tests deben comportarse como las implementaciones reales

4. **Interface Segregation Principle (ISP)**:
   - Interfaces pequeñas y específicas
   - No forzar implementaciones a depender de métodos que no usan
   - Separar interfaces de repositorio por entidad/operación si es necesario

5. **Dependency Inversion Principle (DIP)**:
   - Depender de abstracciones (interfaces), no de implementaciones concretas
   - Los casos de uso dependen de interfaces del dominio, no de infraestructura
   - Inyectar dependencias mediante constructores

---

## 📝 Notas de Tareas

Las notas específicas de tareas se encuentran en `.notes/`:
- Cada tarea puede tener su propio archivo `.md`
- Algunas notas se eliminarán cuando la tarea se complete
- Otras notas se crearán para nuevas tareas
- **SIEMPRE consulta `.notes/`** para contexto específico de tareas en progreso

### Plantilla Estándar de Tareas

**IMPORTANTE**: **TODAS** las tareas creadas en `.notes/` **DEBEN** seguir esta plantilla estándar para que la IA pueda escanearla rápidamente y entender el progreso.

#### Cómo Crear una Nueva Tarea

Cuando necesites crear una nueva tarea, usa el comando de Cursor:

**Comando**: `@.cursor/commands/crear-tarea.md [nombre-tarea] [descripción del objetivo]`

O simplemente di:
```
Usa @.cursor/commands/crear-tarea.md para crear una nueva tarea llamada [nombre-tarea]. 
El objetivo es [descripción breve de la tarea].
```

La IA usará el comando que contiene la plantilla completa y creará automáticamente el archivo en `.notes/[nombre-tarea].md` con la estructura estándar.

**Nota**: La plantilla de referencia está disponible en `.cursor/rules/plantilla_tarea.md` si necesitas consultarla.

#### Plantilla Estándar

```markdown
# 📝 [Nombre de la Tarea]

## 🎯 Objetivo General
Breve descripción de qué queremos lograr y por qué.

---

## 🛠 Contexto Técnico

**Ficheros Involucrados**: @archivo1.go, @archivo2.tsx (Usa @ para que Cursor los identifique)

**Dependencias/Herramientas**: 
- Ej: Docker, WhatsApp Cloud API, NATS, etc.

**Arquitectura**: 
- Ej: Arquitectura Hexagonal (Domain/Application/Infrastructure)

**Integración relacionada**: 
- Ej: WhatsApp, Shopify, etc.

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

**Fecha de inicio**: [fecha]
**Fecha estimada de finalización**: [fecha]

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

#### Uso de la Plantilla

1. **Al crear una tarea nueva**: 
   - Usa el comando: `@.cursor/commands/crear-tarea.md [nombre-tarea] [objetivo]`
   - O solicita: "Usa @.cursor/commands/crear-tarea.md para crear una nueva tarea llamada [nombre-tarea]. El objetivo es [descripción]."

2. **Durante el desarrollo**:
   - Actualiza el checklist marcando pasos completados: `[x] Paso completado`
   - Actualiza el "Estado de la Tarea" con el progreso actual
   - Documenta decisiones importantes en "Memoria de Decisiones"

3. **Al usar Composer (Ctrl+I)**:
   - Siempre referencia la nota: `@.notes/[nombre-tarea].md`
   - Esto ayuda a que Cursor se mantenga en los pasos definidos

4. **Sincronización**:
   - Después de cada cambio importante, actualiza: "Último cambio realizado"
   - Marca los pasos completados en el Plan de Ejecución

#### Características de la Plantilla

- ✅ **Estructura clara**: Fácil de escanear por la IA
- ✅ **Uso de @**: Permite a Cursor identificar archivos involucrados
- ✅ **Checklist**: IA sabe exactamente qué sigue
- ✅ **Estado visible**: Progreso y bloqueos claros
- ✅ **Memoria de decisiones**: Evita cambios innecesarios
- ✅ **Compatible con arquitectura hexagonal**: Incluye referencias a las capas

---

## 🔧 SDK Core

El SDK (`app/shared/sdk/`) es el corazón de Connect Flow. Proporciona:

- **Integrations Manager**: Registra y gestiona integraciones
- **Queue System (NATS JetStream)**: Sistema de colas para eventos
- **Rate Limiter**: Control de tasa de peticiones
- **Repository**: Acceso a base de datos via GORM
- **Inventory Manager**: Helpers de sincronización de inventario
- **Orders Manager**: Helpers de sincronización de órdenes
- **Auth & Security**: JWT y middleware de autenticación

### Interfaz Integrator

Todas las integraciones implementan `sdk.Integrator`:

```go
type Integrator interface {
    Settings() Settings
    OnStartup(ctx context.Context) error
    OnShutdown(ctx context.Context) error
    Install(ctx context.Context, integration models.ExternalIntegration) (models.ExternalIntegration, error)
    Health(ctx context.Context, integration models.ExternalIntegration) error
    SyncInventory(ctx context.Context, req SyncInventoryReq) error
    SyncOrders(ctx context.Context, req SyncOrdersReq) error
    SyncOrderStatus(ctx context.Context, req SyncOrderStatusReq) error
    CreateInvoice(ctx context.Context, req CreateInvoiceReq) error
    CancelInvoice(ctx context.Context, req CancelInvoiceReq) error
    CustomRoutes(route api.Route)
}
```

---

## 🚀 Crear Nueva Integración

### Usando el Scaffold

```bash
pnpm run new:integration
```

El scaffold genera la estructura base. Luego:

1. **Registrar en `main.go`**:
```go
sdk.RegisterIntegration(miintegracion.New)
```

2. **Registrar en `ui/integrations.ts`**:
```ts
import miintegracion from './@integrations/miintegracion/ui'
const integrators = new IntegratorsRegistry([
    // ...
    miintegracion,
])
```

### Requisitos para Nueva Integración

1. **Seguir arquitectura hexagonal** en el backend
2. **Implementar la interfaz `sdk.Integrator`**
3. **Crear componentes UI** (InstallView, SettingsView)
4. **Documentar** en README.md de la integración
5. **Registrar** en main.go y ui/integrations.ts

---

## 🎯 Principios de Trabajo

1. **Arquitectura Hexagonal**: **SIEMPRE** respetar arquitectura hexagonal en nuevas integraciones y código nuevo
   - Separar claramente Domain, Application e Infrastructure
   - Las dependencias deben apuntar hacia Domain
   - No mezclar capas

2. **Principios SOLID**: **SIEMPRE** respetar los principios SOLID
   - Cada función/clase con responsabilidad única
   - Depender de interfaces, no de implementaciones
   - Abierto para extensión, cerrado para modificación

3. **Testing**: **SIEMPRE** crear tests para código nuevo
   - Crear tests unitarios para casos de uso (application layer)
   - Crear tests de integración para handlers y repositorios cuando sea posible
   - Priorizar tests en Domain y Application layers (más testeable por arquitectura hexagonal)
   - Usar mocks/interfaces para aislar dependencias
   - Nombre de tests descriptivos: `Test_UseCase_Method_Scenario_ExpectedResult`

4. **Modularidad**: Cada integración es independiente y autocontenida
5. **Reutilización**: Usar el SDK y componentes compartidos cuando sea posible
6. **Documentación**: Mantener README.md actualizado en cada integración
7. **Consistencia**: Seguir la estructura y convenciones establecidas

---

## 🔍 Flujo de Trabajo Recomendado

1. **Leer `README.md`** del proyecto para entender el contexto general y la arquitectura
2. **Consultar `.notes/`** para contexto de tareas en curso
   - Si no existe una nota para la tarea, crear una usando la plantilla estándar
   - Referenciar siempre `@.notes/[nombre-tarea].md` en el Composer
3. **Revisar documentación** en `docs/guides/` para guías específicas
4. **Revisar integraciones existentes** como referencia (paris, whatsApp son buenos ejemplos de arquitectura hexagonal)
5. **Seguir los principios** de arquitectura hexagonal al crear nuevas integraciones
6. **Actualizar la nota de tarea** después de cada cambio importante:
   - Marcar pasos completados en el Plan de Ejecución
   - Actualizar el "Estado de la Tarea" con progreso y último cambio
   - Documentar decisiones en "Memoria de Decisiones"
7. **Actualizar documentación** si hay cambios en arquitectura o convenciones

---

## 📚 Documentación Adicional

- `README.md` - Documentación general del proyecto
- `docs/guides/` - Guías de desarrollo
- `docs/under-the-hood/` - Detalles técnicos internos
- Cada integración tiene su propio `README.md` con documentación específica

---

## 🧪 Testing

### Prioridad de Tests

**SIEMPRE** crear tests para código nuevo, priorizando en este orden:

1. **Domain Layer** (Más importante):
   - Entidades y validaciones
   - Lógica de negocio pura
   - DTOs y estructuras
   - Errores de dominio

2. **Application Layer** (Casos de Uso):
   - Tests unitarios de casos de uso
   - Mock de dependencias (repositorios, clients)
   - Validación de flujos de negocio
   - Manejo de errores

3. **Infrastructure Layer** (Cuando sea posible):
   - Tests de integración para repositorios
   - Tests de handlers HTTP (con mocks de casos de uso)
   - Tests de mappers

### Estructura de Tests

```
internal/
├── application/
│   └── usecase/
│       └── sync-orders.go
│       └── sync-orders_test.go  ← Test junto al código
└── domain/
    └── entities/
        └── order.go
        └── order_test.go  ← Test junto al código
```

### Convenciones de Testing

1. **Nombres descriptivos**: `Test_ProcessOrderStatus_WhenOrderNotFound_ReturnsError`
2. **Arrange-Act-Assert**: Estructurar tests en estas 3 secciones
3. **Mocks mediante interfaces**: Usar interfaces del dominio para crear mocks
4. **Table-driven tests**: Para múltiples casos similares en Go
5. **Cobertura mínima**: Aspirar a >70% en Domain y Application layers

### Ejemplo de Test de Caso de Uso

```go
func Test_NotifyUpdateStatus_WhenIntegrationNotFound_ReturnsError(t *testing.T) {
    // Arrange
    mockRepo := &MockRepository{}
    mockRepo.On("GetIntegrationByBusinessID", mock.Anything, "business123", mock.Anything).
        Return(nil, errors.New("not found"))
    
    usecase := NewSendMessageUsecase(mockRepo, mockClient, mockLogger)
    req := domain.NotifyWhatsAppRequest{
        BusinessID: "business123",
        OrderID:    "order123",
    }
    
    // Act
    result, err := usecase.NotifyUpdateStatus(context.Background(), req)
    
    // Assert
    assert.Error(t, err)
    assert.Empty(t, result.MessageID)
    mockRepo.AssertExpectations(t)
}
```

### Herramientas de Testing

- **testing**: Paquete estándar de Go
- **testify**: Para assertions y mocks (`assert`, `require`, `mock`)
- **gomock**: Alternativa para generar mocks desde interfaces
- **httptest**: Para tests de handlers HTTP

---

## ⚠️ Notas Importantes

- **Arquitectura Hexagonal**: **SIEMPRE** respetar la separación de capas
- **SOLID**: **SIEMPRE** aplicar principios SOLID en todo el código
- **Tests**: **SIEMPRE** crear tests para código nuevo, especialmente en Domain y Application layers
- **No mezclar capas**: Domain no debe conocer Application ni Infrastructure
- **Usar interfaces**: Los casos de uso deben recibir interfaces, no implementaciones concretas
- **Mappers**: Usar mappers para transformar entre DTOs de dominio y DTOs de infraestructura
- **Errores de dominio**: Crear errores específicos en `domain/errors/`
- **Configuración**: Centralizar configuración en `domain/config/`
- **Dependency Injection**: Inyectar dependencias mediante constructores, no crear instancias dentro de funciones
