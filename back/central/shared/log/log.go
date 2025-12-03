package log

import (
	"context"
	"os"
	"runtime"
	"strings"

	"github.com/rs/zerolog"
)

type ILogger interface {
	Info(ctx ...context.Context) *zerolog.Event
	Error(ctx ...context.Context) *zerolog.Event
	Warn(ctx ...context.Context) *zerolog.Event
	Debug(ctx ...context.Context) *zerolog.Event
	Fatal(ctx ...context.Context) *zerolog.Event
	Panic(ctx ...context.Context) *zerolog.Event
	With() zerolog.Context
	// Nuevos métodos para logger contextual
	WithService(service string) ILogger
	WithModule(module string) ILogger
	WithBusinessID(businessID uint) ILogger
}

type logger struct {
	log        zerolog.Logger
	service    string
	module     string
	businessID uint
}

var defaultLogger *logger

func New() ILogger {
	if defaultLogger == nil {
		// Configurar el logger con formato de consola bonito
		consoleWriter := zerolog.ConsoleWriter{
			Out:        os.Stdout,        // Cambiar a stdout para que se vea en consola
			TimeFormat: "01-02 15:04:05", // Fecha corta (mes-día) y hora
		}

		defaultLogger = &logger{
			log: zerolog.New(consoleWriter).
				With().
				Timestamp().
				Logger().
				Hook(&tracingHook{}),
		}

		// Configurar el nivel de log desde variable de entorno
		level := os.Getenv("LOG_LEVEL")
		switch level {
		case "debug":
			defaultLogger.log = defaultLogger.log.Level(zerolog.DebugLevel)
		case "warn":
			defaultLogger.log = defaultLogger.log.Level(zerolog.WarnLevel)
		case "error":
			defaultLogger.log = defaultLogger.log.Level(zerolog.ErrorLevel)
		default:
			defaultLogger.log = defaultLogger.log.Level(zerolog.InfoLevel)
		}

		// Forzar el logger por defecto del contexto
		zerolog.DefaultContextLogger = &defaultLogger.log
	}
	return defaultLogger
}

// NewWithContext crea un logger con contexto automático de servicio y módulo
func NewWithContext() ILogger {
	service, module := extractServiceAndModule()
	return &logger{
		log:        defaultLogger.log,
		service:    service,
		module:     module,
		businessID: 0,
	}
}

func (l *logger) Info(ctx ...context.Context) *zerolog.Event {
	if len(ctx) > 0 {
		event := zerolog.Ctx(ctx[0]).Info().Ctx(ctx[0])
		return l.addContextualFields(event, ctx[0])
	}
	return l.addContextualFields(l.log.Info(), nil)
}

func (l *logger) Error(ctx ...context.Context) *zerolog.Event {
	if len(ctx) > 0 {
		event := zerolog.Ctx(ctx[0]).Error().Ctx(ctx[0])
		return l.addContextualFields(event, ctx[0])
	}
	return l.addContextualFields(l.log.Error(), nil)
}

func (l *logger) Warn(ctx ...context.Context) *zerolog.Event {
	if len(ctx) > 0 {
		event := zerolog.Ctx(ctx[0]).Warn().Ctx(ctx[0])
		return l.addContextualFields(event, ctx[0])
	}
	return l.addContextualFields(l.log.Warn(), nil)
}

func (l *logger) Debug(ctx ...context.Context) *zerolog.Event {
	if len(ctx) > 0 {
		event := zerolog.Ctx(ctx[0]).Debug().Ctx(ctx[0])
		return l.addContextualFields(event, ctx[0])
	}
	return l.addContextualFields(l.log.Debug(), nil)
}

func (l *logger) Fatal(ctx ...context.Context) *zerolog.Event {
	if len(ctx) > 0 {
		event := zerolog.Ctx(ctx[0]).Fatal()
		return l.addContextualFields(event, ctx[0])
	}
	return l.addContextualFields(l.log.WithLevel(zerolog.FatalLevel), nil)
}

func (l *logger) Panic(ctx ...context.Context) *zerolog.Event {
	if len(ctx) > 0 {
		event := zerolog.Ctx(ctx[0]).Panic().Ctx(ctx[0])
		return l.addContextualFields(event, ctx[0])
	}
	return l.addContextualFields(l.log.WithLevel(zerolog.PanicLevel), nil)
}

// addContextualFields agrega los campos contextuales al evento de log
func (l *logger) addContextualFields(event *zerolog.Event, ctx context.Context) *zerolog.Event {
	// Agregar campos del logger si están configurados
	if l.service != "" {
		event = event.Str("service", l.service)
	}
	if l.module != "" {
		event = event.Str("module", l.module)
	}
	if l.businessID != 0 {
		event = event.Uint("business_id", l.businessID)
	}

	// Agregar nombre de función automáticamente
	skipFunc := false
	if ctx != nil {
		if skip, ok := SkipFunctionFromCtx(ctx); ok && skip {
			skipFunc = true
		}
	}

	if !skipFunc {
		funcName := getFunctionName()
		if funcName != "" {
			event = event.Str("function", funcName)
		}
	}

	// Agregar campos del contexto si están disponibles
	if ctx != nil {
		if service, ok := ServiceFromCtx(ctx); ok {
			event = event.Str("service", service)
		}
		if module, ok := ModuleFromCtx(ctx); ok {
			event = event.Str("module", module)
		}
		if businessID, ok := BusinessIDFromCtx(ctx); ok {
			event = event.Uint("business_id", businessID)
		}
		if userID, ok := UserIDFromCtx(ctx); ok {
			event = event.Uint("user_id", userID)
		}
		if duration, ok := DurationFromCtx(ctx); ok {
			event = event.Str("duration", duration)
		}
		if statusCode, ok := StatusCodeFromCtx(ctx); ok {
			event = event.Int("status_code", statusCode)
		}
		if function, ok := FunctionFromCtx(ctx); ok {
			event = event.Str("function", function)
		}
	}

	return event
}

func (l *logger) With() zerolog.Context {
	return l.log.With()
}

// WithService crea un nuevo logger con el servicio especificado
func (l *logger) WithService(service string) ILogger {
	return &logger{
		log:        l.log,
		service:    service,
		module:     l.module,
		businessID: l.businessID,
	}
}

// WithModule crea un nuevo logger con el módulo especificado
func (l *logger) WithModule(module string) ILogger {
	return &logger{
		log:        l.log,
		service:    l.service,
		module:     module,
		businessID: l.businessID,
	}
}

// WithBusinessID crea un nuevo logger con el business_id especificado
func (l *logger) WithBusinessID(businessID uint) ILogger {
	return &logger{
		log:        l.log,
		service:    l.service,
		module:     l.module,
		businessID: businessID,
	}
}

func Init() {
	New()
}

// extractServiceAndModule extrae automáticamente el servicio y módulo del stack trace
func extractServiceAndModule() (service, module string) {
	// Obtener el stack trace
	pc := make([]uintptr, 15)
	n := runtime.Callers(2, pc) // Saltar esta función y NewWithContext
	if n == 0 {
		return "unknown", "unknown"
	}

	// Buscar la función que no sea del paquete log o runtime
	for i := 0; i < n; i++ {
		fn := runtime.FuncForPC(pc[i])
		if fn == nil {
			continue
		}

		funcName := fn.Name()
		// Filtrar funciones del logger y runtime
		if !strings.Contains(funcName, "log.") &&
			!strings.Contains(funcName, "runtime.") &&
			!strings.Contains(funcName, "zerolog.") {

			// Extraer servicio y módulo del path completo
			// Ejemplo: central/services/horizontalproperty/internal/infra/primary/handlers/handlers.CreatePropertyUnit
			parts := strings.Split(funcName, "/")
			if len(parts) >= 3 {
				// Buscar el servicio en el path
				for j, part := range parts {
					if part == "services" && j+1 < len(parts) {
						service = parts[j+1]
						break
					}
				}

				// Buscar el módulo (handler, usecase, repository, etc.)
				for _, part := range parts {
					if strings.Contains(part, "handler") ||
						strings.Contains(part, "usecase") ||
						strings.Contains(part, "repository") ||
						strings.Contains(part, "controller") {
						module = part
						break
					}
				}

				if service != "" && module != "" {
					return service, module
				}
			}
		}
	}

	return "unknown", "unknown"
}

// getFunctionName obtiene el nombre de la función que está ejecutando el log
func getFunctionName() string {
	// Obtener el stack trace
	pc := make([]uintptr, 10)
	n := runtime.Callers(0, pc)
	if n == 0 {
		return "unknown"
	}

	// Buscar la función que no sea del paquete log o runtime
	for i := 0; i < n; i++ {
		fn := runtime.FuncForPC(pc[i])
		if fn == nil {
			continue
		}

		funcName := fn.Name()
		// Filtrar funciones del logger y runtime
		if !strings.Contains(funcName, "log.") &&
			!strings.Contains(funcName, "runtime.") &&
			!strings.Contains(funcName, "zerolog.") {
			// Extraer solo el nombre de la función sin el paquete completo
			parts := strings.Split(funcName, ".")
			if len(parts) > 0 {
				return parts[len(parts)-1]
			}
			return funcName
		}
	}
	return "unknown"
}

type tracingHook struct{}

func (h *tracingHook) Run(e *zerolog.Event, level zerolog.Level, msg string) {
	// Si el mensaje es solo espacios o vacío, no agregar metadatos
	if strings.TrimSpace(msg) == "" {
		return
	}

	ctx := e.GetCtx()
	reqId, ok := ReqId(ctx)
	if ok {
		e.Str("req_id", reqId)
	}

	// Agregar metadatos estructurados si están disponibles en el contexto
	if service, ok := ServiceFromCtx(ctx); ok {
		e.Str("service", service)
	}
	if module, ok := ModuleFromCtx(ctx); ok {
		e.Str("module", module)
	}
	if businessID, ok := BusinessIDFromCtx(ctx); ok {
		e.Uint("business_id", businessID)
	}
}

type reqIdKey struct{}

var reqId reqIdKey

type serviceKey struct{}

var service serviceKey

type moduleKey struct{}

var module moduleKey

type businessIDKey struct{}

var businessID businessIDKey

type userIDKey struct{}

var userID userIDKey

type durationKey struct{}

var duration durationKey

type statusCodeKey struct{}

var statusCode statusCodeKey

type functionKey struct{}

var function functionKey

func ReqId(ctx context.Context) (string, bool) {
	id, ok := ctx.Value(reqId).(string)
	return id, ok
}

func ServiceFromCtx(ctx context.Context) (string, bool) {
	service, ok := ctx.Value(service).(string)
	return service, ok
}

func ModuleFromCtx(ctx context.Context) (string, bool) {
	module, ok := ctx.Value(module).(string)
	return module, ok
}

func BusinessIDFromCtx(ctx context.Context) (uint, bool) {
	businessID, ok := ctx.Value(businessID).(uint)
	return businessID, ok
}

func UserIDFromCtx(ctx context.Context) (uint, bool) {
	userID, ok := ctx.Value(userID).(uint)
	return userID, ok
}

func DurationFromCtx(ctx context.Context) (string, bool) {
	duration, ok := ctx.Value(duration).(string)
	return duration, ok
}

func StatusCodeFromCtx(ctx context.Context) (int, bool) {
	statusCode, ok := ctx.Value(statusCode).(int)
	return statusCode, ok
}

func FunctionFromCtx(ctx context.Context) (string, bool) {
	function, ok := ctx.Value(function).(string)
	return function, ok
}

// WithServiceCtx agrega el servicio al contexto
func WithServiceCtx(ctx context.Context, service string) context.Context {
	return context.WithValue(ctx, service, service)
}

// WithModuleCtx agrega el módulo al contexto
func WithModuleCtx(ctx context.Context, module string) context.Context {
	return context.WithValue(ctx, module, module)
}

// WithBusinessIDCtx agrega el business_id al contexto
func WithBusinessIDCtx(ctx context.Context, businessID uint) context.Context {
	return context.WithValue(ctx, businessID, businessID)
}

// WithUserIDCtx agrega el user_id al contexto
func WithUserIDCtx(ctx context.Context, userID uint) context.Context {
	return context.WithValue(ctx, userID, userID)
}

// WithDurationCtx agrega la duración al contexto
func WithDurationCtx(ctx context.Context, duration string) context.Context {
	return context.WithValue(ctx, duration, duration)
}

// WithStatusCodeCtx agrega el status code al contexto
func WithStatusCodeCtx(ctx context.Context, statusCode int) context.Context {
	return context.WithValue(ctx, statusCode, statusCode)
}

// WithFunctionCtx agrega el nombre de función al contexto
func WithFunctionCtx(ctx context.Context, function string) context.Context {
	return context.WithValue(ctx, function, function)
}

type skipFunctionKey struct{}

var skipFunction skipFunctionKey

func SkipFunctionFromCtx(ctx context.Context) (bool, bool) {
	skip, ok := ctx.Value(skipFunction).(bool)
	return skip, ok
}

// WithSkipFunctionCtx agrega la instrucción de saltar el nombre de función al contexto
func WithSkipFunctionCtx(ctx context.Context) context.Context {
	return context.WithValue(ctx, skipFunction, true)
}
