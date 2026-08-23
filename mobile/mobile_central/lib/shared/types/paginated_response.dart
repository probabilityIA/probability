class PaginatedResponse<T> {
  final List<T> data;
  final Pagination pagination;

  PaginatedResponse({required this.data, required this.pagination});
}

class Pagination {
  final int currentPage;
  final int perPage;
  final int total;
  final int lastPage;
  final bool hasNext;
  final bool hasPrev;

  Pagination({
    required this.currentPage,
    required this.perPage,
    required this.total,
    required this.lastPage,
    required this.hasNext,
    required this.hasPrev,
  });

  factory Pagination.fromJson(Map<String, dynamic> json) {
    final currentPage = _int(json['current_page'] ?? json['page'], 1);
    final perPage = _int(json['per_page'] ?? json['page_size'], 10);
    final total = _int(json['total'], 0);
    final lastPage = _int(
      json['last_page'] ?? json['total_pages'],
      perPage > 0 ? ((total + perPage - 1) ~/ perPage) : 1,
    );

    return Pagination(
      currentPage: currentPage,
      perPage: perPage,
      total: total,
      lastPage: lastPage,
      hasNext: json['has_next'] ?? (currentPage < lastPage),
      hasPrev: json['has_prev'] ?? (currentPage > 1),
    );
  }

  factory Pagination.fromEnvelope(Map<String, dynamic> envelope) {
    final nested = envelope['pagination'];
    if (nested is Map<String, dynamic> && nested.isNotEmpty) {
      return Pagination.fromJson(nested);
    }
    return Pagination.fromJson(envelope);
  }

  static int _int(dynamic value, int fallback) {
    if (value is int) return value;
    if (value is num) return value.toInt();
    if (value is String) return int.tryParse(value) ?? fallback;
    return fallback;
  }
}

class SingleResponse<T> {
  final bool success;
  final T data;
  final String? message;

  SingleResponse({
    required this.success,
    required this.data,
    this.message,
  });
}
