enum FindingSeverity { error, warn, info }

FindingSeverity findingSeverityFrom(String? code) {
  switch (code) {
    case 'error':
      return FindingSeverity.error;
    case 'warn':
      return FindingSeverity.warn;
    default:
      return FindingSeverity.info;
  }
}

int _int(dynamic value) {
  if (value is int) return value;
  if (value is num) return value.round();
  if (value is String) return int.tryParse(value) ?? 0;
  return 0;
}

DateTime? _when(dynamic value) {
  if (value == null) return null;
  final parsed = DateTime.tryParse(value.toString());
  if (parsed == null || parsed.year < 2000) return null;
  return parsed.toLocal();
}

class Finding {
  const Finding({
    required this.code,
    required this.severity,
    required this.title,
    required this.detail,
    required this.count,
    this.channels = const <String>[],
  });

  final String code;
  final FindingSeverity severity;
  final String title;
  final String detail;
  final int count;
  final List<String> channels;

  factory Finding.fromJson(Map<String, dynamic> json) => Finding(
        code: json['code']?.toString() ?? '',
        severity: findingSeverityFrom(json['severity']?.toString()),
        title: json['title']?.toString() ?? '',
        detail: json['detail']?.toString() ?? '',
        count: _int(json['count']),
        channels: (json['channels'] as List<dynamic>?)
                ?.map((e) => e.toString())
                .toList() ??
            const <String>[],
      );
}

class FindingChannelSummary {
  const FindingChannelSummary({
    required this.integrationId,
    required this.integrationName,
    required this.channelCode,
    this.matched = 0,
    this.notAssociated = 0,
    this.onlyInChannel = 0,
    this.channelNoSku = 0,
    this.skuChanged = 0,
    this.skuTypo = 0,
    this.comparedAt,
  });

  final int integrationId;
  final String integrationName;
  final String channelCode;
  final int matched;
  final int notAssociated;
  final int onlyInChannel;
  final int channelNoSku;
  final int skuChanged;
  final int skuTypo;
  final DateTime? comparedAt;

  int get pending =>
      notAssociated + onlyInChannel + channelNoSku + skuChanged + skuTypo;

  bool get isClean => pending == 0;

  factory FindingChannelSummary.fromJson(Map<String, dynamic> json) =>
      FindingChannelSummary(
        integrationId: _int(json['integration_id']),
        integrationName: json['integration_name']?.toString() ?? '',
        channelCode: json['channel_code']?.toString() ?? '',
        matched: _int(json['matched']),
        notAssociated: _int(json['not_associated']),
        onlyInChannel: _int(json['only_in_channel']),
        channelNoSku: _int(json['channel_no_sku']),
        skuChanged: _int(json['sku_changed']),
        skuTypo: _int(json['sku_typo']),
        comparedAt: _when(json['compared_at']),
      );
}

class FindingsReport {
  const FindingsReport({
    this.total = 0,
    this.findings = const <Finding>[],
    this.channels = const <FindingChannelSummary>[],
  });

  final int total;
  final List<Finding> findings;
  final List<FindingChannelSummary> channels;

  bool get isEmpty => findings.isEmpty && channels.isEmpty;

  DateTime? get lastComparedAt {
    DateTime? latest;
    for (final channel in channels) {
      final when = channel.comparedAt;
      if (when == null) continue;
      if (latest == null || when.isAfter(latest)) latest = when;
    }
    return latest;
  }

  factory FindingsReport.fromJson(Map<String, dynamic>? json) {
    if (json == null) return const FindingsReport();
    return FindingsReport(
      total: _int(json['total']),
      findings: (json['findings'] as List<dynamic>?)
              ?.whereType<Map<String, dynamic>>()
              .map(Finding.fromJson)
              .toList() ??
          const <Finding>[],
      channels: (json['channels'] as List<dynamic>?)
              ?.whereType<Map<String, dynamic>>()
              .map(FindingChannelSummary.fromJson)
              .toList() ??
          const <FindingChannelSummary>[],
    );
  }
}

class DataSummaryCell {
  const DataSummaryCell({
    required this.integrationId,
    required this.integrationName,
    required this.channelCode,
    this.canFill = 0,
    this.canOverwrite = 0,
  });

  final int integrationId;
  final String integrationName;
  final String channelCode;
  final int canFill;
  final int canOverwrite;

  bool get hasSomething => canFill > 0 || canOverwrite > 0;

  factory DataSummaryCell.fromJson(Map<String, dynamic> json) => DataSummaryCell(
        integrationId: _int(json['integration_id']),
        integrationName: json['integration_name']?.toString() ?? '',
        channelCode: json['channel_code']?.toString() ?? '',
        canFill: _int(json['can_fill']),
        canOverwrite: _int(json['can_overwrite']),
      );
}

class DataSummaryRow {
  const DataSummaryRow({
    required this.field,
    required this.label,
    this.note,
    this.cells = const <DataSummaryCell>[],
  });

  final String field;
  final String label;
  final String? note;
  final List<DataSummaryCell> cells;

  int get totalFill =>
      cells.fold(0, (sum, cell) => sum + cell.canFill);

  int get totalOverwrite =>
      cells.fold(0, (sum, cell) => sum + cell.canOverwrite);

  bool get hasSomething => totalFill > 0 || totalOverwrite > 0;

  factory DataSummaryRow.fromJson(Map<String, dynamic> json) => DataSummaryRow(
        field: json['field']?.toString() ?? '',
        label: json['label']?.toString() ?? '',
        note: json['note']?.toString(),
        cells: (json['cells'] as List<dynamic>?)
                ?.whereType<Map<String, dynamic>>()
                .map(DataSummaryCell.fromJson)
                .toList() ??
            const <DataSummaryCell>[],
      );
}

class DataSummary {
  const DataSummary({this.rows = const <DataSummaryRow>[], this.snapshotAt});

  final List<DataSummaryRow> rows;
  final DateTime? snapshotAt;

  bool get isEmpty => rows.isEmpty;

  List<DataSummaryRow> get actionable =>
      rows.where((row) => row.hasSomething).toList();

  factory DataSummary.fromJson(Map<String, dynamic>? json) {
    if (json == null) return const DataSummary();
    return DataSummary(
      rows: (json['data'] as List<dynamic>?)
              ?.whereType<Map<String, dynamic>>()
              .map(DataSummaryRow.fromJson)
              .toList() ??
          const <DataSummaryRow>[],
      snapshotAt: _when(json['snapshot_at']),
    );
  }
}

enum InventoryCompareAction { update, unchanged, skip }

InventoryCompareAction inventoryActionFrom(String? code) {
  switch (code) {
    case 'update':
      return InventoryCompareAction.update;
    case 'skip':
      return InventoryCompareAction.skip;
    default:
      return InventoryCompareAction.unchanged;
  }
}

class InventoryCompareRow {
  const InventoryCompareRow({
    required this.productId,
    required this.sku,
    required this.name,
    required this.action,
    this.imageUrl,
    this.probabilityQty,
    this.channelQty,
    this.delta,
    this.reason,
  });

  final String productId;
  final String sku;
  final String name;
  final InventoryCompareAction action;
  final String? imageUrl;
  final int? probabilityQty;
  final int? channelQty;
  final int? delta;
  final String? reason;

  bool get needsUpdate => action == InventoryCompareAction.update;

  static int? _nullableInt(dynamic value) {
    if (value == null) return null;
    return _int(value);
  }

  factory InventoryCompareRow.fromJson(Map<String, dynamic> json) =>
      InventoryCompareRow(
        productId: json['product_id']?.toString() ?? '',
        sku: json['sku']?.toString() ?? '',
        name: json['name']?.toString() ?? '',
        action: inventoryActionFrom(json['action']?.toString()),
        imageUrl: json['image_url']?.toString(),
        probabilityQty: _nullableInt(json['probability_qty']),
        channelQty: _nullableInt(json['channel_qty']),
        delta: _nullableInt(json['delta']),
        reason: json['reason']?.toString(),
      );
}

class InventoryCompareTotals {
  const InventoryCompareTotals({
    this.total = 0,
    this.toUpdate = 0,
    this.unchanged = 0,
    this.skipped = 0,
  });

  final int total;
  final int toUpdate;
  final int unchanged;
  final int skipped;

  factory InventoryCompareTotals.fromJson(Map<String, dynamic>? json) {
    if (json == null) return const InventoryCompareTotals();
    return InventoryCompareTotals(
      total: _int(json['total']),
      toUpdate: _int(json['to_update']),
      unchanged: _int(json['unchanged']),
      skipped: _int(json['skipped']),
    );
  }
}

class InventoryComparePage {
  const InventoryComparePage({
    required this.rows,
    required this.totals,
    required this.page,
    required this.pageSize,
    required this.total,
    required this.totalPages,
    this.checkedAt,
    this.fromCache = false,
  });

  final List<InventoryCompareRow> rows;
  final InventoryCompareTotals totals;
  final int page;
  final int pageSize;
  final int total;
  final int totalPages;
  final DateTime? checkedAt;
  final bool fromCache;

  factory InventoryComparePage.fromJson(Map<String, dynamic> json) =>
      InventoryComparePage(
        rows: (json['rows'] as List<dynamic>?)
                ?.whereType<Map<String, dynamic>>()
                .map(InventoryCompareRow.fromJson)
                .toList() ??
            const <InventoryCompareRow>[],
        totals: InventoryCompareTotals.fromJson(
          json['totals'] is Map<String, dynamic>
              ? json['totals'] as Map<String, dynamic>
              : null,
        ),
        page: _int(json['page']) == 0 ? 1 : _int(json['page']),
        pageSize: _int(json['page_size']) == 0 ? 100 : _int(json['page_size']),
        total: _int(json['total']),
        totalPages: _int(json['total_pages']),
        checkedAt: _when(json['checked_at']),
        fromCache: json['from_cache'] == true,
      );
}

class InventoryCompareQuery {
  const InventoryCompareQuery({
    required this.integrationId,
    this.businessId,
    this.page = 1,
    this.pageSize = 100,
    this.snapshot = true,
    this.onlyDiff = false,
    this.search,
    this.skus,
  });

  final int integrationId;
  final int? businessId;
  final int page;
  final int pageSize;
  final bool snapshot;
  final bool onlyDiff;
  final String? search;
  final List<String>? skus;

  Map<String, dynamic> toBody() => <String, dynamic>{
        'integration_id': integrationId,
        'page': page,
        'page_size': pageSize,
        if (businessId != null) 'business_id': businessId,
        if (snapshot) 'source': 'snapshot',
        if (onlyDiff) 'only_diff': true,
        if (search != null && search!.isNotEmpty) 'q': search,
        if (skus != null && skus!.isNotEmpty) 'skus': skus,
      };
}
