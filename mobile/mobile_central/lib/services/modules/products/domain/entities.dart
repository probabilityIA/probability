class Product {
  final String id;
  final String createdAt;
  final String updatedAt;
  final String? deletedAt;
  final int businessId;
  final int? integrationId;
  final String? integrationType;
  final String? externalId;
  final String sku;
  final String name;
  final String? description;
  final double price;
  final double? compareAtPrice;
  final double? costPrice;
  final String currency;
  final int stock;
  final String? stockStatus;
  final bool manageStock;
  final double? weight;
  final double? height;
  final double? width;
  final double? length;
  final String? imageUrl;
  final List<String>? images;
  final String? thumbnail;
  final String status;
  final bool isActive;
  final dynamic metadata;

  Product({
    required this.id,
    required this.createdAt,
    required this.updatedAt,
    this.deletedAt,
    required this.businessId,
    this.integrationId,
    this.integrationType,
    this.externalId,
    required this.sku,
    required this.name,
    this.description,
    required this.price,
    this.compareAtPrice,
    this.costPrice,
    required this.currency,
    required this.stock,
    this.stockStatus,
    required this.manageStock,
    this.weight,
    this.height,
    this.width,
    this.length,
    this.imageUrl,
    this.images,
    this.thumbnail,
    required this.status,
    required this.isActive,
    this.metadata,
  });

  factory Product.fromJson(Map<String, dynamic> json) {
    return Product(
      id: json['id']?.toString() ?? '',
      createdAt: json['created_at'] ?? '',
      updatedAt: json['updated_at'] ?? '',
      deletedAt: json['deleted_at'],
      businessId: json['business_id'] ?? 0,
      integrationId: json['integration_id'],
      integrationType: json['integration_type'],
      externalId: json['external_id'],
      sku: json['sku'] ?? '',
      name: json['name'] ?? '',
      description: json['description'],
      price: (json['price'] ?? 0).toDouble(),
      compareAtPrice: json['compare_at_price']?.toDouble(),
      costPrice: json['cost_price']?.toDouble(),
      currency: json['currency'] ?? 'COP',
      stock: json['stock'] ?? 0,
      stockStatus: json['stock_status'],
      manageStock: json['manage_stock'] ?? false,
      weight: json['weight']?.toDouble(),
      height: json['height']?.toDouble(),
      width: json['width']?.toDouble(),
      length: json['length']?.toDouble(),
      imageUrl: json['image_url'],
      images: (json['images'] as List<dynamic>?)?.map((e) => e.toString()).toList(),
      thumbnail: json['thumbnail'],
      status: json['status'] ?? '',
      isActive: json['is_active'] ?? true,
      metadata: json['metadata'],
    );
  }
}

class ProductIntegration {
  final int id;
  final String productId;
  final int integrationId;
  final String? integrationType;
  final String? integrationName;
  final String externalProductId;
  final String createdAt;
  final String updatedAt;

  ProductIntegration({
    required this.id,
    required this.productId,
    required this.integrationId,
    this.integrationType,
    this.integrationName,
    required this.externalProductId,
    required this.createdAt,
    required this.updatedAt,
  });

  factory ProductIntegration.fromJson(Map<String, dynamic> json) {
    return ProductIntegration(
      id: json['id'] ?? 0,
      productId: json['product_id']?.toString() ?? '',
      integrationId: json['integration_id'] ?? 0,
      integrationType: json['integration_type'],
      integrationName: json['integration_name'],
      externalProductId: json['external_product_id'] ?? '',
      createdAt: json['created_at'] ?? '',
      updatedAt: json['updated_at'] ?? '',
    );
  }
}

class GetProductsParams {
  final int? page;
  final int? pageSize;
  final int? businessId;
  final int? integrationId;
  final String? integrationType;
  final String? sku;
  final String? skus;
  final String? name;
  final String? externalId;
  final String? sortBy;
  final String? sortOrder;
  final String? startDate;
  final String? endDate;

  GetProductsParams({
    this.page,
    this.pageSize,
    this.businessId,
    this.integrationId,
    this.integrationType,
    this.sku,
    this.skus,
    this.name,
    this.externalId,
    this.sortBy,
    this.sortOrder,
    this.startDate,
    this.endDate,
  });

  Map<String, dynamic> toQueryParams() {
    final params = <String, dynamic>{};
    if (page != null) params['page'] = page;
    if (pageSize != null) params['page_size'] = pageSize;
    if (businessId != null) params['business_id'] = businessId;
    if (integrationId != null) params['integration_id'] = integrationId;
    if (integrationType != null) params['integration_type'] = integrationType;
    if (sku != null) params['sku'] = sku;
    if (skus != null) params['skus'] = skus;
    if (name != null) params['name'] = name;
    if (externalId != null) params['external_id'] = externalId;
    if (sortBy != null) params['sort_by'] = sortBy;
    if (sortOrder != null) params['sort_order'] = sortOrder;
    if (startDate != null) params['start_date'] = startDate;
    if (endDate != null) params['end_date'] = endDate;
    return params;
  }
}

class CreateProductDTO {
  final int businessId;
  final String sku;
  final String name;
  final double price;
  final int stock;
  final String? description;
  final double? compareAtPrice;
  final double? costPrice;
  final String? currency;
  final String? stockStatus;
  final bool? manageStock;
  final double? weight;
  final double? height;
  final double? width;
  final double? length;
  final String? status;
  final bool? isActive;

  CreateProductDTO({
    required this.businessId,
    required this.sku,
    required this.name,
    required this.price,
    required this.stock,
    this.description,
    this.compareAtPrice,
    this.costPrice,
    this.currency,
    this.stockStatus,
    this.manageStock,
    this.weight,
    this.height,
    this.width,
    this.length,
    this.status,
    this.isActive,
  });

  Map<String, dynamic> toJson() {
    final json = <String, dynamic>{
      'business_id': businessId,
      'sku': sku,
      'name': name,
      'price': price,
      'stock': stock,
    };
    if (description != null) json['description'] = description;
    if (compareAtPrice != null) json['compare_at_price'] = compareAtPrice;
    if (costPrice != null) json['cost_price'] = costPrice;
    if (currency != null) json['currency'] = currency;
    if (stockStatus != null) json['stock_status'] = stockStatus;
    if (manageStock != null) json['manage_stock'] = manageStock;
    if (weight != null) json['weight'] = weight;
    if (height != null) json['height'] = height;
    if (width != null) json['width'] = width;
    if (length != null) json['length'] = length;
    if (status != null) json['status'] = status;
    if (isActive != null) json['is_active'] = isActive;
    return json;
  }
}

class UpdateProductDTO {
  final String? sku;
  final String? name;
  final String? description;
  final double? price;
  final double? compareAtPrice;
  final double? costPrice;
  final String? currency;
  final int? stock;
  final String? stockStatus;
  final bool? manageStock;
  final double? weight;
  final double? height;
  final double? width;
  final double? length;
  final String? status;
  final bool? isActive;

  UpdateProductDTO({
    this.sku,
    this.name,
    this.description,
    this.price,
    this.compareAtPrice,
    this.costPrice,
    this.currency,
    this.stock,
    this.stockStatus,
    this.manageStock,
    this.weight,
    this.height,
    this.width,
    this.length,
    this.status,
    this.isActive,
  });

  Map<String, dynamic> toJson() {
    final json = <String, dynamic>{};
    if (sku != null) json['sku'] = sku;
    if (name != null) json['name'] = name;
    if (description != null) json['description'] = description;
    if (price != null) json['price'] = price;
    if (compareAtPrice != null) json['compare_at_price'] = compareAtPrice;
    if (costPrice != null) json['cost_price'] = costPrice;
    if (currency != null) json['currency'] = currency;
    if (stock != null) json['stock'] = stock;
    if (stockStatus != null) json['stock_status'] = stockStatus;
    if (manageStock != null) json['manage_stock'] = manageStock;
    if (weight != null) json['weight'] = weight;
    if (height != null) json['height'] = height;
    if (width != null) json['width'] = width;
    if (length != null) json['length'] = length;
    if (status != null) json['status'] = status;
    if (isActive != null) json['is_active'] = isActive;
    return json;
  }
}

class AddProductIntegrationDTO {
  final int integrationId;
  final String externalProductId;

  AddProductIntegrationDTO({
    required this.integrationId,
    required this.externalProductId,
  });

  Map<String, dynamic> toJson() => {
        'integration_id': integrationId,
        'external_product_id': externalProductId,
      };
}
