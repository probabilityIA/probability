class SyncProviderSpec {
  const SyncProviderSpec({
    required this.typeId,
    required this.key,
    required this.label,
    required this.prefix,
    required this.inventoryEventPrefix,
    this.productEventPrefix,
    this.reconcilePath = '/products/reconcile',
    this.supportsCompareInventory = false,
    this.supportsAssociate = false,
    this.applyModes = const <String>[],
    this.applyBodies = const <String, ProductApplyBody>{},
  });

  final int typeId;
  final String key;
  final String label;
  final String prefix;
  final String inventoryEventPrefix;
  final String? productEventPrefix;
  final String reconcilePath;
  final bool supportsCompareInventory;
  final bool supportsAssociate;
  final List<String> applyModes;
  final Map<String, ProductApplyBody> applyBodies;

  String get productPrefix => productEventPrefix ?? inventoryEventPrefix;

  String get syncInventoryPath => '$prefix/inventory/sync';
  String get compareInventoryPath => '$prefix/inventory/compare';
  String get reconcileProductsPath => '$prefix$reconcilePath';
  String get applyProductsPath => '$prefix/products/apply';
  String get associateProductsPath => '$prefix/products/associate';

  String get onlyInChannelField => 'only_in_$key';
  String get channelNoSkuField => '${key}_no_sku';

  bool supportsApply(String action) => applyModes.contains(action);

  ProductApplyBody applyBodyFor(String action) =>
      applyBodies[action] ?? const ProductApplyBody();
}

class ProductApplyBody {
  const ProductApplyBody({this.direction, this.mode});

  final String? direction;
  final String? mode;

  Map<String, dynamic> toJson() => <String, dynamic>{
        if (direction != null) 'direction': direction,
        if (mode != null) 'mode': mode,
      };
}

const Map<int, SyncProviderSpec> syncProviders = <int, SyncProviderSpec>{
  1: SyncProviderSpec(
    typeId: 1,
    key: 'shopify',
    label: 'Shopify',
    prefix: '/integrations/shopify',
    inventoryEventPrefix: 'shopify',
    supportsAssociate: true,
    applyModes: ['createInChannel', 'createInProbability'],
    applyBodies: {
      'createInChannel': ProductApplyBody(direction: 'to_shopify'),
      'createInProbability': ProductApplyBody(direction: 'to_probability'),
    },
  ),
  3: SyncProviderSpec(
    typeId: 3,
    key: 'meli',
    label: 'Mercado Libre',
    prefix: '/integrations/meli',
    inventoryEventPrefix: 'meli',
    supportsCompareInventory: true,
    supportsAssociate: true,
    applyModes: ['createInChannel', 'createInProbability'],
    applyBodies: {
      'createInChannel': ProductApplyBody(direction: 'to_meli'),
      'createInProbability': ProductApplyBody(direction: 'to_probability'),
    },
  ),
  4: SyncProviderSpec(
    typeId: 4,
    key: 'woocommerce',
    label: 'WooCommerce',
    prefix: '/woocommerce',
    inventoryEventPrefix: 'woo',
    productEventPrefix: 'woocommerce',
    supportsCompareInventory: true,
    supportsAssociate: true,
    applyModes: ['createInChannel', 'createInProbability'],
    applyBodies: {
      'createInChannel': ProductApplyBody(direction: 'to_woo'),
      'createInProbability': ProductApplyBody(direction: 'to_probability'),
    },
  ),
  8: SyncProviderSpec(
    typeId: 8,
    key: 'siigo',
    label: 'Siigo',
    prefix: '/siigo',
    inventoryEventPrefix: 'siigo',
    reconcilePath: '/products/reconcile/start',
    supportsCompareInventory: true,
    applyModes: ['createInProbability'],
  ),
  16: SyncProviderSpec(
    typeId: 16,
    key: 'vtex',
    label: 'VTEX',
    prefix: '/vtex',
    inventoryEventPrefix: 'vtex',
    supportsAssociate: true,
    applyModes: ['createInProbability', 'updateInProbability'],
    applyBodies: {
      'createInProbability': ProductApplyBody(direction: 'to_probability', mode: 'create'),
      'updateInProbability': ProductApplyBody(direction: 'to_probability', mode: 'update'),
    },
  ),
  17: SyncProviderSpec(
    typeId: 17,
    key: 'tiendanube',
    label: 'Tiendanube',
    prefix: '/tiendanube',
    inventoryEventPrefix: 'tiendanube',
    supportsCompareInventory: true,
    supportsAssociate: true,
    applyModes: ['createInChannel', 'createInProbability', 'updateInProbability'],
    applyBodies: {
      'createInChannel': ProductApplyBody(direction: 'to_tiendanube', mode: 'create'),
      'createInProbability': ProductApplyBody(direction: 'to_probability', mode: 'create'),
      'updateInProbability': ProductApplyBody(direction: 'to_probability', mode: 'update'),
    },
  ),
  33: SyncProviderSpec(
    typeId: 33,
    key: 'jumpseller',
    label: 'Jumpseller',
    prefix: '/jumpseller',
    inventoryEventPrefix: 'jumpseller',
    supportsAssociate: true,
    applyModes: ['createInChannel', 'createInProbability', 'updateInProbability'],
    applyBodies: {
      'createInChannel': ProductApplyBody(direction: 'to_jumpseller', mode: 'create'),
      'createInProbability': ProductApplyBody(direction: 'to_probability', mode: 'create'),
      'updateInProbability': ProductApplyBody(direction: 'to_probability', mode: 'update'),
    },
  ),
};

const List<int> ordersCompareTypeIds = <int>[1, 3, 4, 16, 17, 33];

SyncProviderSpec? syncProviderFor(dynamic integrationTypeId) {
  if (integrationTypeId == null) return null;
  final id = integrationTypeId is int
      ? integrationTypeId
      : int.tryParse(integrationTypeId.toString());
  if (id == null) return null;
  return syncProviders[id];
}

const List<String> _inventorySuffixes = ['started', 'item', 'progress', 'completed'];

const List<String> _productSyncSuffixes = ['started', 'progress', 'completed'];

List<String> get globalSyncEventTypes {
  final types = <String>[];
  for (final spec in syncProviders.values) {
    for (final suffix in _inventorySuffixes) {
      types.add('${spec.inventoryEventPrefix}.inventory.sync.$suffix');
    }
    types.add('${spec.inventoryEventPrefix}.product.reconcile.started');
    types.add('${spec.inventoryEventPrefix}.product.reconcile.completed');
    for (final suffix in _productSyncSuffixes) {
      types.add('${spec.productPrefix}.product.sync.$suffix');
    }
  }
  return types;
}
