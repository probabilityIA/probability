class StorefrontProduct {
  final String id;
  final String name;
  final String description;
  final String shortDescription;
  final double price;
  final double? compareAtPrice;
  final String currency;
  final String imageUrl;
  final List<String>? images;
  final String sku;
  final int stockQuantity;
  final String category;
  final String brand;
  final bool isFeatured;
  final String createdAt;

  StorefrontProduct({
    required this.id,
    required this.name,
    required this.description,
    required this.shortDescription,
    required this.price,
    this.compareAtPrice,
    required this.currency,
    required this.imageUrl,
    this.images,
    required this.sku,
    required this.stockQuantity,
    required this.category,
    required this.brand,
    required this.isFeatured,
    required this.createdAt,
  });

  factory StorefrontProduct.fromJson(Map<String, dynamic> json) {
    return StorefrontProduct(
      id: json['id']?.toString() ?? '',
      name: json['name'] ?? '',
      description: json['description'] ?? '',
      shortDescription: json['short_description'] ?? '',
      price: (json['price'] ?? 0).toDouble(),
      compareAtPrice: json['compare_at_price']?.toDouble(),
      currency: json['currency'] ?? 'COP',
      imageUrl: json['image_url'] ?? '',
      images: (json['images'] as List<dynamic>?)?.map((e) => e.toString()).toList(),
      sku: json['sku'] ?? '',
      stockQuantity: json['stock_quantity'] ?? 0,
      category: json['category'] ?? '',
      brand: json['brand'] ?? '',
      isFeatured: json['is_featured'] ?? false,
      createdAt: json['created_at'] ?? '',
    );
  }
}

class StorefrontOrder {
  final String id;
  final String orderNumber;
  final String status;
  final double totalAmount;
  final String currency;
  final String createdAt;
  final List<StorefrontOrderItem> items;

  StorefrontOrder({
    required this.id,
    required this.orderNumber,
    required this.status,
    required this.totalAmount,
    required this.currency,
    required this.createdAt,
    required this.items,
  });

  factory StorefrontOrder.fromJson(Map<String, dynamic> json) {
    return StorefrontOrder(
      id: json['id']?.toString() ?? '',
      orderNumber: json['order_number'] ?? '',
      status: json['status'] ?? '',
      totalAmount: (json['total_amount'] ?? 0).toDouble(),
      currency: json['currency'] ?? 'COP',
      createdAt: json['created_at'] ?? '',
      items: (json['items'] as List<dynamic>?)
              ?.map((e) => StorefrontOrderItem.fromJson(e))
              .toList() ??
          [],
    );
  }
}

class StorefrontOrderItem {
  final String productName;
  final int quantity;
  final double unitPrice;
  final double totalPrice;
  final String? imageUrl;

  StorefrontOrderItem({
    required this.productName,
    required this.quantity,
    required this.unitPrice,
    required this.totalPrice,
    this.imageUrl,
  });

  factory StorefrontOrderItem.fromJson(Map<String, dynamic> json) {
    return StorefrontOrderItem(
      productName: json['product_name'] ?? '',
      quantity: json['quantity'] ?? 0,
      unitPrice: (json['unit_price'] ?? 0).toDouble(),
      totalPrice: (json['total_price'] ?? 0).toDouble(),
      imageUrl: json['image_url'],
    );
  }
}

class CreateStorefrontOrderDTO {
  final List<CreateStorefrontOrderItemDTO> items;
  final String? notes;
  final StorefrontAddress? address;

  CreateStorefrontOrderDTO({
    required this.items,
    this.notes,
    this.address,
  });

  Map<String, dynamic> toJson() {
    final json = <String, dynamic>{
      'items': items.map((e) => e.toJson()).toList(),
    };
    if (notes != null) json['notes'] = notes;
    if (address != null) json['address'] = address!.toJson();
    return json;
  }
}

class CreateStorefrontOrderItemDTO {
  final String productId;
  final int quantity;

  CreateStorefrontOrderItemDTO({
    required this.productId,
    required this.quantity,
  });

  Map<String, dynamic> toJson() => {
        'product_id': productId,
        'quantity': quantity,
      };
}

class StorefrontAddress {
  final String firstName;
  final String? lastName;
  final String? phone;
  final String street;
  final String? street2;
  final String city;
  final String? state;
  final String? country;
  final String? postalCode;
  final String? instructions;

  StorefrontAddress({
    required this.firstName,
    this.lastName,
    this.phone,
    required this.street,
    this.street2,
    required this.city,
    this.state,
    this.country,
    this.postalCode,
    this.instructions,
  });

  Map<String, dynamic> toJson() {
    final json = <String, dynamic>{
      'first_name': firstName,
      'street': street,
      'city': city,
    };
    if (lastName != null) json['last_name'] = lastName;
    if (phone != null) json['phone'] = phone;
    if (street2 != null) json['street2'] = street2;
    if (state != null) json['state'] = state;
    if (country != null) json['country'] = country;
    if (postalCode != null) json['postal_code'] = postalCode;
    if (instructions != null) json['instructions'] = instructions;
    return json;
  }
}

class RegisterDTO {
  final String name;
  final String email;
  final String password;
  final String? phone;
  final String? dni;
  final String businessCode;

  RegisterDTO({
    required this.name,
    required this.email,
    required this.password,
    this.phone,
    this.dni,
    required this.businessCode,
  });

  Map<String, dynamic> toJson() {
    final json = <String, dynamic>{
      'name': name,
      'email': email,
      'password': password,
      'business_code': businessCode,
    };
    if (phone != null) json['phone'] = phone;
    if (dni != null) json['dni'] = dni;
    return json;
  }
}

class GetCatalogParams {
  final int? page;
  final int? pageSize;
  final String? search;
  final String? category;
  final int? businessId;

  GetCatalogParams({
    this.page,
    this.pageSize,
    this.search,
    this.category,
    this.businessId,
  });

  Map<String, dynamic> toQueryParams() {
    final params = <String, dynamic>{};
    if (page != null) params['page'] = page;
    if (pageSize != null) params['page_size'] = pageSize;
    if (search != null) params['search'] = search;
    if (category != null) params['category'] = category;
    if (businessId != null) params['business_id'] = businessId;
    return params;
  }
}

class GetOrdersParams {
  final int? page;
  final int? pageSize;
  final int? businessId;

  GetOrdersParams({
    this.page,
    this.pageSize,
    this.businessId,
  });

  Map<String, dynamic> toQueryParams() {
    final params = <String, dynamic>{};
    if (page != null) params['page'] = page;
    if (pageSize != null) params['page_size'] = pageSize;
    if (businessId != null) params['business_id'] = businessId;
    return params;
  }
}
