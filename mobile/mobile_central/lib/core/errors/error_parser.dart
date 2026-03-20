import 'package:dio/dio.dart';

/// Traduce excepciones de Dio a mensajes claros en español para el usuario.
String parseError(dynamic e) {
  if (e is DioException) {
    switch (e.type) {
      case DioExceptionType.connectionError:
        return 'No se pudo conectar al servidor. Verifica tu conexión a internet o que el servidor esté activo.';
      case DioExceptionType.connectionTimeout:
        return 'La conexión con el servidor tardó demasiado. Intenta de nuevo.';
      case DioExceptionType.sendTimeout:
        return 'No se pudo enviar la solicitud. Verifica tu conexión a internet.';
      case DioExceptionType.receiveTimeout:
        return 'El servidor tardó demasiado en responder. Intenta de nuevo.';
      case DioExceptionType.badResponse:
        final statusCode = e.response?.statusCode;
        final data = e.response?.data;
        if (data is Map) {
          if (data.containsKey('error')) return data['error'].toString();
          if (data.containsKey('message')) return data['message'].toString();
        }
        switch (statusCode) {
          case 400:
            return 'Datos inválidos. Verifica la información ingresada.';
          case 401:
            return 'Sesión expirada o credenciales incorrectas. Inicia sesión de nuevo.';
          case 403:
            return 'No tienes permisos para realizar esta acción.';
          case 404:
            return 'Recurso no encontrado.';
          case 409:
            return 'Conflicto: el recurso ya existe o fue modificado.';
          case 422:
            return 'Los datos enviados no son válidos. Revisa el formulario.';
          case 500:
            return 'Error interno del servidor. Intenta más tarde.';
          case 502:
          case 503:
            return 'Servidor no disponible. Intenta más tarde.';
          default:
            return 'Error del servidor (código $statusCode). Intenta más tarde.';
        }
      case DioExceptionType.cancel:
        return 'La solicitud fue cancelada.';
      case DioExceptionType.badCertificate:
        return 'Error de seguridad en la conexión. Contacta al administrador.';
      case DioExceptionType.unknown:
        return 'Error de conexión. Verifica tu conexión a internet.';
    }
  }
  return 'Ocurrió un error inesperado. Intenta de nuevo.';
}
