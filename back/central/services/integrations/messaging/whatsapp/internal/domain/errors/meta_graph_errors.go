package errors

import "fmt"

// MetaGraphError representa un error parseado de la Graph API de Meta
type MetaGraphError struct {
	Code         int
	Subcode      int
	Message      string
	StatusCode   int
	PhoneNumberID uint
}

func (e *MetaGraphError) Error() string {
	return e.FriendlyMessage()
}

// FriendlyMessage retorna un mensaje claro en español basado en el código de error de Meta
func (e *MetaGraphError) FriendlyMessage() string {
	code := e.Code
	subcode := e.Subcode

	switch {
	// ── Autenticación ──
	case code == 190:
		switch subcode {
		case 463:
			return "El Access Token ha expirado. Genera uno nuevo en Meta Business Manager"
		case 460:
			return "El Access Token ha sido revocado. Genera uno nuevo en Meta Business Manager"
		default:
			return "Access Token invalido o expirado. Verifica el token en Meta Business Manager"
		}

	// ── Objeto no existe / permisos insuficientes ──
	case code == 100 && subcode == 33:
		return fmt.Sprintf("El Phone Number ID '%d' no existe o no tiene permisos. Verifica el ID en Meta Business Manager → WhatsApp → Numeros de telefono", e.PhoneNumberID)
	case code == 100 && subcode == 2018109:
		return "Parametros invalidos en el mensaje. Verifica que el numero de destino tenga formato internacional (ej: +573001234567)"
	case code == 100:
		return fmt.Sprintf("Solicitud invalida (code 100). Verifica que el Phone Number ID '%d' sea correcto y tenga permisos asignados", e.PhoneNumberID)

	// ── Límites de envío ──
	case code == 130429:
		return "Se excedio el limite de mensajes. WhatsApp tiene restricciones de envio por minuto. Intenta de nuevo en unos minutos"
	case code == 131048:
		return "No se permite enviar mensajes a este numero. El usuario no ha iniciado conversacion o la ventana de 24 horas ha expirado"
	case code == 131026:
		return "El numero de destino no esta registrado en WhatsApp"
	case code == 131030:
		return "El numero de destino no esta en la lista de permitidos. En cuentas de prueba, agrega el numero en Meta Business Manager → WhatsApp → API Setup → Allowed Numbers"

	// ── Plantillas ──
	case code == 132000:
		return "La plantilla de mensaje no existe o no esta aprobada. Verifica el nombre de la plantilla en Meta Business Manager"
	case code == 132001:
		return "Los parametros de la plantilla no coinciden. Verifica que se envian todas las variables requeridas"
	case code == 132005:
		return "La plantilla fue pausada o deshabilitada por Meta. Revisa su estado en Meta Business Manager"
	case code == 132012:
		return "La plantilla tiene demasiados parametros. Simplifica la plantilla en Meta Business Manager"

	// ── Cuenta / Negocio ──
	case code == 131031:
		return "La cuenta de WhatsApp Business no esta verificada. Completa la verificacion en Meta Business Manager"
	case code == 131009:
		return "El numero de telefono no esta registrado en la API de WhatsApp Cloud. Registra el numero en Meta Business Manager"

	// ── Permisos de API ──
	case code == 10 || code == 200:
		return "Permisos insuficientes en el Access Token. Asegurate de que el token tenga el permiso 'whatsapp_business_messaging'"

	// ── Errores de servidor de Meta ──
	case code == 1 || code == 2:
		return "Error temporal en los servidores de Meta. Intenta de nuevo en unos minutos"
	case code == 4:
		return "Se excedio el limite de llamadas a la API de Meta. Intenta de nuevo en unos minutos"

	// ── HTTP genéricos (cuando Meta no retorna código de error en JSON) ──
	case e.StatusCode == 401:
		return "No autorizado. El Access Token es invalido o no tiene permisos"
	case e.StatusCode == 403:
		return "Acceso denegado. Verifica los permisos de la app en Meta Business Manager"
	case e.StatusCode == 404:
		return fmt.Sprintf("Recurso no encontrado. Verifica la URL de la API y que el Phone Number ID '%d' exista", e.PhoneNumberID)
	case e.StatusCode == 429:
		return "Demasiadas solicitudes. Espera unos minutos antes de reintentar"
	case e.StatusCode >= 500:
		return "Error interno en los servidores de Meta. Intenta de nuevo mas tarde"

	default:
		if e.Message != "" {
			return fmt.Sprintf("Error de WhatsApp API (code: %d): %s", code, e.Message)
		}
		return fmt.Sprintf("Error %d de WhatsApp API", e.StatusCode)
	}
}

// NewMetaGraphError crea un MetaGraphError a partir de los datos parseados
func NewMetaGraphError(code, subcode int, message string, statusCode int, phoneNumberID uint) *MetaGraphError {
	return &MetaGraphError{
		Code:          code,
		Subcode:       subcode,
		Message:       message,
		StatusCode:    statusCode,
		PhoneNumberID: phoneNumberID,
	}
}

// NewMetaGraphErrorUnparseable crea un error cuando no se pudo parsear el JSON de Meta
func NewMetaGraphErrorUnparseable(statusCode int) *MetaGraphError {
	return &MetaGraphError{
		StatusCode: statusCode,
	}
}
