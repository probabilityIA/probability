import 'dart:async';
import 'dart:convert';

import 'package:dio/dio.dart';

import '../config/environment.dart';

class SseEvent {
  const SseEvent({required this.type, required this.data});

  final String type;
  final Map<String, dynamic> data;
}

abstract class SseSource {
  Stream<SseEvent> get events;

  Future<void> connect({required int businessId, required List<String> eventTypes});

  void disconnect();

  void dispose();
}

class SseClient implements SseSource {
  SseClient({Dio? dio, String? baseUrl})
      : _dio = dio ?? Dio(),
        _baseUrl = baseUrl ?? Environment.apiBaseUrl;

  final Dio _dio;
  final String _baseUrl;

  final StreamController<SseEvent> _events = StreamController<SseEvent>.broadcast();
  CancelToken? _cancel;
  Timer? _retry;
  bool _closed = false;
  int _attempt = 0;

  @override
  Stream<SseEvent> get events => _events.stream;
  bool get isConnected => _cancel != null && !(_cancel?.isCancelled ?? true);

  @override
  Future<void> connect({
    required int businessId,
    required List<String> eventTypes,
  }) async {
    if (_closed) return;
    disconnect();
    _cancel = CancelToken();
    await _open(businessId: businessId, eventTypes: eventTypes);
  }

  Future<void> _open({
    required int businessId,
    required List<String> eventTypes,
  }) async {
    final token = _cancel;
    if (token == null || _closed) return;

    try {
      final response = await _dio.get<ResponseBody>(
        '$_baseUrl/notify/sse/order-notify',
        queryParameters: <String, dynamic>{
          if (businessId > 0) 'business_id': businessId,
          if (eventTypes.isNotEmpty) 'event_types': eventTypes.join(','),
        },
        options: Options(
          responseType: ResponseType.stream,
          headers: const {'Accept': 'text/event-stream'},
          receiveTimeout: Duration.zero,
        ),
        cancelToken: token,
      );

      _attempt = 0;
      final body = response.data;
      if (body == null) {
        _scheduleRetry(businessId: businessId, eventTypes: eventTypes);
        return;
      }

      var buffer = '';
      await for (final chunk in body.stream) {
        if (_closed || token.isCancelled) return;
        buffer += utf8.decode(chunk, allowMalformed: true);
        var cut = buffer.indexOf('\n\n');
        while (cut >= 0) {
          _emit(buffer.substring(0, cut));
          buffer = buffer.substring(cut + 2);
          cut = buffer.indexOf('\n\n');
        }
        if (buffer.length > 65536) buffer = '';
      }
    } catch (_) {
      if (_closed || token.isCancelled) return;
    }

    _scheduleRetry(businessId: businessId, eventTypes: eventTypes);
  }

  void _emit(String block) {
    var type = '';
    final payload = StringBuffer();

    for (final raw in block.split('\n')) {
      final line = raw.trimRight();
      if (line.isEmpty || line.startsWith(':')) continue;
      if (line.startsWith('event:')) {
        type = line.substring(6).trim();
      } else if (line.startsWith('data:')) {
        payload.write(line.substring(5).trim());
      }
    }

    if (payload.isEmpty) return;

    try {
      final decoded = jsonDecode(payload.toString());
      if (decoded is! Map<String, dynamic>) return;
      final resolved = _typeOf(decoded, type);
      if (resolved.isEmpty) return;
      _events.add(SseEvent(type: resolved, data: decoded));
    } catch (_) {
      return;
    }
  }

  String _typeOf(Map<String, dynamic> decoded, String fallback) {
    final inner = decoded['type'];
    if (inner is String && inner.isNotEmpty) return inner;
    final metadata = decoded['metadata'];
    if (metadata is Map && metadata['event_type'] is String) {
      return metadata['event_type'] as String;
    }
    return fallback;
  }

  void _scheduleRetry({
    required int businessId,
    required List<String> eventTypes,
  }) {
    if (_closed) return;
    final token = _cancel;
    if (token == null || token.isCancelled) return;

    _attempt = _attempt >= 5 ? 5 : _attempt + 1;
    final seconds = <int>[1, 2, 4, 8, 15, 30][_attempt - 1];
    _retry?.cancel();
    _retry = Timer(Duration(seconds: seconds), () {
      _open(businessId: businessId, eventTypes: eventTypes);
    });
  }

  @override
  void disconnect() {
    _retry?.cancel();
    _retry = null;
    _cancel?.cancel();
    _cancel = null;
  }

  @override
  void dispose() {
    _closed = true;
    disconnect();
    _events.close();
  }
}
