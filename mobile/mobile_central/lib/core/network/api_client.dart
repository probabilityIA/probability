import 'package:dio/dio.dart';
import '../config/environment.dart';

class ApiClient {
  final Dio _dio;
  String? _token;

  ApiClient({String? token})
      : _token = token,
        _dio = Dio(BaseOptions(
          baseUrl: Environment.apiBaseUrl,
          connectTimeout: const Duration(seconds: 30),
          receiveTimeout: const Duration(seconds: 30),
          headers: {'Content-Type': 'application/json'},
        ));

  void setToken(String? token) {
    _token = token;
  }

  Map<String, dynamic> _authHeaders() {
    if (_token != null) {
      return {'Authorization': 'Bearer $_token'};
    }
    return {};
  }

  Future<Response> get(
    String path, {
    Map<String, dynamic>? queryParameters,
  }) async {
    return _dio.get(
      path,
      queryParameters: queryParameters,
      options: Options(headers: _authHeaders()),
    );
  }

  Future<Response> post(
    String path, {
    dynamic data,
    Map<String, dynamic>? queryParameters,
  }) async {
    return _dio.post(
      path,
      data: data,
      queryParameters: queryParameters,
      options: Options(headers: _authHeaders()),
    );
  }

  Future<Response> put(
    String path, {
    dynamic data,
    Map<String, dynamic>? queryParameters,
  }) async {
    return _dio.put(
      path,
      data: data,
      queryParameters: queryParameters,
      options: Options(headers: _authHeaders()),
    );
  }

  Future<Response> delete(
    String path, {
    Map<String, dynamic>? queryParameters,
  }) async {
    return _dio.delete(
      path,
      queryParameters: queryParameters,
      options: Options(headers: _authHeaders()),
    );
  }

  Future<Response> postFormData(
    String path, {
    required FormData data,
    Map<String, dynamic>? queryParameters,
  }) async {
    final headers = _authHeaders();
    headers.remove('Content-Type');
    return _dio.post(
      path,
      data: data,
      queryParameters: queryParameters,
      options: Options(headers: headers),
    );
  }

  Future<Response> putFormData(
    String path, {
    required FormData data,
    Map<String, dynamic>? queryParameters,
  }) async {
    final headers = _authHeaders();
    headers.remove('Content-Type');
    return _dio.put(
      path,
      data: data,
      queryParameters: queryParameters,
      options: Options(headers: headers),
    );
  }
}
