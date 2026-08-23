enum SyncNodeState { idle, queued, active, scan, done, error }

enum SyncMode { idle, inventory, products }

enum SyncEnvironment { overview, products, data, inventory, ordersCompare, invoicing }

enum SyncRunKind { inventory, products }

enum ProductActionKey {
  associate,
  createInChannel,
  createInProbability,
  updateInProbability,
  createBothSides,
}

extension SyncRunKindX on SyncRunKind {
  String get code => this == SyncRunKind.inventory ? 'inventory' : 'products';

  static SyncRunKind fromCode(String? code) =>
      code == 'products' ? SyncRunKind.products : SyncRunKind.inventory;
}

abstract class SyncResult {
  const SyncResult();
}

class InventorySyncResult extends SyncResult {
  const InventorySyncResult({
    this.total = 0,
    this.updated = 0,
    this.unchanged = 0,
    this.skipped = 0,
    this.failed = 0,
  });

  final int total;
  final int updated;
  final int unchanged;
  final int skipped;
  final int failed;
}

class ProductsSyncResult extends SyncResult {
  const ProductsSyncResult({
    this.matched = 0,
    this.notAssociated = 0,
    this.onlyInProbability = 0,
    this.onlyInChannel = 0,
    this.channelNoSku = 0,
    this.skuChanged = 0,
    this.skuTypo = 0,
  });

  final int matched;
  final int notAssociated;
  final int onlyInProbability;
  final int onlyInChannel;
  final int channelNoSku;
  final int skuChanged;
  final int skuTypo;

  int get pending => notAssociated + onlyInChannel + channelNoSku + skuChanged + skuTypo;
}

class SyncErrorResult extends SyncResult {
  const SyncErrorResult(this.message);

  final String message;
}

class SyncProgress {
  const SyncProgress({this.processed = 0, this.total = 0});

  final int processed;
  final int total;

  double? get ratio => total > 0 ? processed / total : null;
}

class ProductActionResult {
  const ProductActionResult({
    required this.ok,
    required this.message,
    this.pending = false,
  });

  final bool ok;
  final String message;
  final bool pending;
}

class SyncRunDetail {
  const SyncRunDetail({
    required this.sku,
    required this.label,
    required this.tone,
    this.group,
    this.matchedBy,
    this.matchedValue,
    this.parentRef,
    this.parentLabel,
    this.variantLabel,
  });

  final String sku;
  final String label;
  final String tone;
  final String? group;
  final String? matchedBy;
  final String? matchedValue;
  final String? parentRef;
  final String? parentLabel;
  final String? variantLabel;

  factory SyncRunDetail.fromJson(Map<String, dynamic> json) => SyncRunDetail(
        sku: json['sku']?.toString() ?? '',
        label: json['label']?.toString() ?? '',
        tone: json['tone']?.toString() ?? 'ok',
        group: json['group']?.toString(),
        matchedBy: json['matched_by']?.toString(),
        matchedValue: json['matched_value']?.toString(),
        parentRef: json['parent_ref']?.toString(),
        parentLabel: json['parent_label']?.toString(),
        variantLabel: json['variant_label']?.toString(),
      );
}

class SyncRunRecord {
  const SyncRunRecord({
    required this.integrationId,
    required this.kind,
    this.status = '',
    this.message,
    this.finishedAt,
    this.total = 0,
    this.updated = 0,
    this.unchanged = 0,
    this.skipped = 0,
    this.failed = 0,
    this.matched = 0,
    this.notAssociated = 0,
    this.onlyInProbability = 0,
    this.onlyInChannel = 0,
    this.channelNoSku = 0,
    this.skuChanged = 0,
    this.skuTypo = 0,
  });

  final int integrationId;
  final SyncRunKind kind;
  final String status;
  final String? message;
  final String? finishedAt;
  final int total;
  final int updated;
  final int unchanged;
  final int skipped;
  final int failed;
  final int matched;
  final int notAssociated;
  final int onlyInProbability;
  final int onlyInChannel;
  final int channelNoSku;
  final int skuChanged;
  final int skuTypo;

  DateTime? get finishedOn =>
      finishedAt == null ? null : DateTime.tryParse(finishedAt!)?.toLocal();

  static int _int(dynamic value) {
    if (value is int) return value;
    if (value is num) return value.round();
    if (value is String) return int.tryParse(value) ?? 0;
    return 0;
  }

  factory SyncRunRecord.fromJson(Map<String, dynamic> json) => SyncRunRecord(
        integrationId: _int(json['integration_id']),
        kind: SyncRunKindX.fromCode(json['kind']?.toString()),
        status: json['status']?.toString() ?? '',
        message: json['message']?.toString(),
        finishedAt: json['finished_at']?.toString(),
        total: _int(json['total']),
        updated: _int(json['updated']),
        unchanged: _int(json['unchanged']),
        skipped: _int(json['skipped']),
        failed: _int(json['failed']),
        matched: _int(json['matched']),
        notAssociated: _int(json['not_associated']),
        onlyInProbability: _int(json['only_in_probability']),
        onlyInChannel: _int(json['only_in_channel']),
        channelNoSku: _int(json['channel_no_sku']),
        skuChanged: _int(json['sku_changed']),
        skuTypo: _int(json['sku_typo']),
      );

  SyncResult toResult() => kind == SyncRunKind.products
      ? ProductsSyncResult(
          matched: matched,
          notAssociated: notAssociated,
          onlyInProbability: onlyInProbability,
          onlyInChannel: onlyInChannel,
          channelNoSku: channelNoSku,
          skuChanged: skuChanged,
          skuTypo: skuTypo,
        )
      : InventorySyncResult(
          total: total,
          updated: updated,
          unchanged: unchanged,
          skipped: skipped,
          failed: failed,
        );
}

class SyncRunItemsQuery {
  const SyncRunItemsQuery({
    required this.integrationId,
    required this.kind,
    this.group,
    this.search,
    this.page = 1,
    this.pageSize = 50,
    this.businessId,
  });

  final int integrationId;
  final SyncRunKind kind;
  final String? group;
  final String? search;
  final int page;
  final int pageSize;
  final int? businessId;

  Map<String, dynamic> toQueryParams() => <String, dynamic>{
        'integration_id': integrationId,
        'kind': kind.code,
        'page': page,
        'page_size': pageSize,
        if (group != null && group!.isNotEmpty) 'group': group,
        if (search != null && search!.isNotEmpty) 'q': search,
        if (businessId != null) 'business_id': businessId,
      };
}

class SyncStartResult {
  const SyncStartResult({
    required this.success,
    this.correlationId,
    this.message,
  });

  final bool success;
  final String? correlationId;
  final String? message;

  bool get tracked => success && (correlationId?.isNotEmpty ?? false);

  factory SyncStartResult.fromJson(Map<String, dynamic>? json) {
    if (json == null) return const SyncStartResult(success: false);
    final envelope = json['data'] is Map<String, dynamic>
        ? json['data'] as Map<String, dynamic>
        : json;
    return SyncStartResult(
      success: json['success'] != false,
      correlationId: (envelope['correlation_id'] ?? json['correlation_id'])?.toString(),
      message: (json['message'] ?? envelope['message'])?.toString(),
    );
  }
}

String buildApplyMessage(
  int total,
  int created,
  int updated,
  int failed,
  List<dynamic>? failedItems,
) {
  final parts = <String>[];
  if (created > 0) parts.add('$created creados');
  if (updated > 0) parts.add('$updated actualizados');
  if (failed > 0) parts.add('$failed con error');
  if (parts.isEmpty) {
    return total == 0 ? 'No habia productos por aplicar' : 'Sin cambios';
  }

  var message = '${parts.join(', ')} de $total';
  final items = failedItems ?? const <dynamic>[];
  if (items.isEmpty) return message;

  String? firstError;
  String? firstSku;
  for (final item in items) {
    if (item is! Map) continue;
    final error = item['error']?.toString();
    if (error != null && error.isNotEmpty) {
      firstError = error;
      break;
    }
  }
  final head = items.first;
  if (head is Map) firstSku = head['sku']?.toString();

  if (firstError == null) return message;
  message += ' - ${firstSku != null && firstSku.isNotEmpty ? '$firstSku: ' : ''}$firstError';
  if (items.length > 1) message += ' (y ${items.length - 1} mas)';
  return message;
}
