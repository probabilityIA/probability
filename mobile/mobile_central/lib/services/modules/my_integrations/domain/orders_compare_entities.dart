enum OrderCompareAction { create, inSync, onlyInProbability }

OrderCompareAction orderCompareActionFrom(String? code) {
  switch (code) {
    case 'create':
      return OrderCompareAction.create;
    case 'only_in_probability':
      return OrderCompareAction.onlyInProbability;
    default:
      return OrderCompareAction.inSync;
  }
}

class OrderCompareRow {
  const OrderCompareRow({
    required this.externalId,
    required this.number,
    required this.customerName,
    required this.channelStatus,
    required this.rawStatus,
    required this.total,
    required this.currency,
    required this.items,
    required this.createdAt,
    required this.action,
    required this.movesInventory,
    required this.statusMismatch,
    required this.totalMismatch,
    this.localStatus,
    this.orderId,
    this.orderNumber,
    this.localTotal,
    this.url,
    this.inventoryNote,
  });

  final String externalId;
  final String number;
  final String customerName;
  final String channelStatus;
  final String rawStatus;
  final String? localStatus;
  final String? orderId;
  final String? orderNumber;
  final double total;
  final double? localTotal;
  final String currency;
  final int items;
  final String createdAt;
  final String? url;
  final OrderCompareAction action;
  final bool movesInventory;
  final String? inventoryNote;
  final bool statusMismatch;
  final bool totalMismatch;

  bool get canCreate => action == OrderCompareAction.create;

  bool get hasMismatch => statusMismatch || totalMismatch;

  String get label {
    if (number.isNotEmpty) return number;
    if (orderNumber != null && orderNumber!.isNotEmpty) return orderNumber!;
    return externalId;
  }

  double get amount => total != 0 ? total : (localTotal ?? 0);

  DateTime? get createdOn {
    final parsed = DateTime.tryParse(createdAt);
    if (parsed == null || parsed.year < 2000) return null;
    return parsed.toLocal();
  }

  static double _double(dynamic value) {
    if (value is double) return value;
    if (value is num) return value.toDouble();
    if (value is String) return double.tryParse(value) ?? 0;
    return 0;
  }

  static int _int(dynamic value) {
    if (value is int) return value;
    if (value is num) return value.round();
    if (value is String) return int.tryParse(value) ?? 0;
    return 0;
  }

  static String? _text(dynamic value) {
    if (value == null) return null;
    final text = value.toString();
    return text.isEmpty ? null : text;
  }

  factory OrderCompareRow.fromJson(Map<String, dynamic> json) => OrderCompareRow(
        externalId: json['external_id']?.toString() ?? '',
        number: json['number']?.toString() ?? '',
        customerName: json['customer_name']?.toString() ?? '',
        channelStatus: json['channel_status']?.toString() ?? '',
        rawStatus: json['raw_status']?.toString() ?? '',
        localStatus: _text(json['local_status']),
        orderId: _text(json['order_id']),
        orderNumber: _text(json['order_number']),
        total: _double(json['total']),
        localTotal: json['local_total'] == null ? null : _double(json['local_total']),
        currency: json['currency']?.toString() ?? 'COP',
        items: _int(json['items']),
        createdAt: json['created_at']?.toString() ?? '',
        url: _text(json['url']),
        action: orderCompareActionFrom(json['action']?.toString()),
        movesInventory: json['moves_inventory'] == true,
        inventoryNote: _text(json['inventory_note']),
        statusMismatch: json['status_mismatch'] == true,
        totalMismatch: json['total_mismatch'] == true,
      );
}

class OrderCompareTotals {
  const OrderCompareTotals({
    this.total = 0,
    this.toCreate = 0,
    this.inSync = 0,
    this.onlyInProbability = 0,
    this.withoutInventory = 0,
    this.withStatusMismatch = 0,
  });

  final int total;
  final int toCreate;
  final int inSync;
  final int onlyInProbability;
  final int withoutInventory;
  final int withStatusMismatch;

  int get inChannel => inSync + toCreate;

  static int _int(dynamic value) {
    if (value is int) return value;
    if (value is num) return value.round();
    if (value is String) return int.tryParse(value) ?? 0;
    return 0;
  }

  factory OrderCompareTotals.fromJson(Map<String, dynamic>? json) {
    if (json == null) return const OrderCompareTotals();
    return OrderCompareTotals(
      total: _int(json['total']),
      toCreate: _int(json['to_create']),
      inSync: _int(json['in_sync']),
      onlyInProbability: _int(json['only_in_probability']),
      withoutInventory: _int(json['without_inventory']),
      withStatusMismatch: _int(json['with_status_mismatch']),
    );
  }
}

class OrdersComparePage {
  const OrdersComparePage({
    required this.rows,
    required this.totals,
    required this.page,
    required this.pageSize,
    required this.total,
    required this.totalPages,
    this.checkedAt,
  });

  final List<OrderCompareRow> rows;
  final OrderCompareTotals totals;
  final int page;
  final int pageSize;
  final int total;
  final int totalPages;
  final String? checkedAt;

  static int _int(dynamic value) {
    if (value is int) return value;
    if (value is num) return value.round();
    if (value is String) return int.tryParse(value) ?? 0;
    return 0;
  }

  factory OrdersComparePage.fromJson(Map<String, dynamic> json) {
    final rows = (json['rows'] as List<dynamic>?) ?? const <dynamic>[];
    return OrdersComparePage(
      rows: rows
          .whereType<Map<String, dynamic>>()
          .map(OrderCompareRow.fromJson)
          .toList(),
      totals: OrderCompareTotals.fromJson(
        json['totals'] is Map<String, dynamic>
            ? json['totals'] as Map<String, dynamic>
            : null,
      ),
      page: _int(json['page']) == 0 ? 1 : _int(json['page']),
      pageSize: _int(json['page_size']) == 0 ? 20 : _int(json['page_size']),
      total: _int(json['total']),
      totalPages: _int(json['total_pages']),
      checkedAt: json['checked_at']?.toString(),
    );
  }
}

class OrdersCompareQuery {
  const OrdersCompareQuery({
    required this.integrationId,
    this.businessId,
    this.from,
    this.to,
    this.page = 1,
    this.pageSize = 20,
    this.onlyDiff = false,
    this.search,
  });

  final int integrationId;
  final int? businessId;
  final String? from;
  final String? to;
  final int page;
  final int pageSize;
  final bool onlyDiff;
  final String? search;

  Map<String, dynamic> toQueryParams() => <String, dynamic>{
        'integration_id': integrationId,
        'page': page,
        'page_size': pageSize,
        if (businessId != null) 'business_id': businessId,
        if (from != null && from!.isNotEmpty) 'from': from,
        if (to != null && to!.isNotEmpty) 'to': to,
        if (onlyDiff) 'only_diff': 'true',
        if (search != null && search!.isNotEmpty) 'q': search,
      };
}

class OrdersApplyResult {
  const OrdersApplyResult({
    this.queued = const <String>[],
    this.skipped = const <String>[],
    this.failed = const <String, String>{},
    this.withoutInventory = const <String>[],
    this.note,
  });

  final List<String> queued;
  final List<String> skipped;
  final Map<String, String> failed;
  final List<String> withoutInventory;
  final String? note;

  static List<String> _list(dynamic value) => value is List
      ? value.map((e) => e.toString()).toList()
      : const <String>[];

  factory OrdersApplyResult.fromJson(Map<String, dynamic>? json) {
    if (json == null) return const OrdersApplyResult();
    final failed = json['failed'];
    return OrdersApplyResult(
      queued: _list(json['queued']),
      skipped: _list(json['skipped']),
      withoutInventory: _list(json['without_inventory']),
      failed: failed is Map
          ? failed.map((key, value) => MapEntry(key.toString(), value.toString()))
          : const <String, String>{},
      note: json['note']?.toString(),
    );
  }

  String get summary {
    final parts = <String>[];
    if (queued.isNotEmpty) {
      parts.add('${queued.length} orden${queued.length == 1 ? '' : 'es'} '
          'enviada${queued.length == 1 ? '' : 's'} a crear');
    }
    if (skipped.isNotEmpty) parts.add('${skipped.length} ya existian');
    if (failed.isNotEmpty) parts.add('${failed.length} fallaron');
    return parts.isEmpty ? 'No se creo ninguna orden' : parts.join(', ');
  }
}
