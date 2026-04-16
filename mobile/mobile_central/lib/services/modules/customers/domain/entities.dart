class CustomerInfo {
  final int id;
  final int businessId;
  final String name;
  final String? email;
  final String phone;
  final String? dni;
  final String createdAt;
  final String updatedAt;

  CustomerInfo({
    required this.id,
    required this.businessId,
    required this.name,
    this.email,
    required this.phone,
    this.dni,
    required this.createdAt,
    required this.updatedAt,
  });

  factory CustomerInfo.fromJson(Map<String, dynamic> json) {
    return CustomerInfo(
      id: json['id'] ?? 0,
      businessId: json['business_id'] ?? 0,
      name: json['name'] ?? '',
      email: json['email'],
      phone: json['phone'] ?? '',
      dni: json['dni'],
      createdAt: json['created_at'] ?? '',
      updatedAt: json['updated_at'] ?? '',
    );
  }
}

class CustomerDetail extends CustomerInfo {
  final int orderCount;
  final double totalSpent;
  final String? lastOrderAt;

  CustomerDetail({
    required super.id,
    required super.businessId,
    required super.name,
    super.email,
    required super.phone,
    super.dni,
    required super.createdAt,
    required super.updatedAt,
    required this.orderCount,
    required this.totalSpent,
    this.lastOrderAt,
  });

  factory CustomerDetail.fromJson(Map<String, dynamic> json) {
    return CustomerDetail(
      id: json['id'] ?? 0,
      businessId: json['business_id'] ?? 0,
      name: json['name'] ?? '',
      email: json['email'],
      phone: json['phone'] ?? '',
      dni: json['dni'],
      createdAt: json['created_at'] ?? '',
      updatedAt: json['updated_at'] ?? '',
      orderCount: json['order_count'] ?? 0,
      totalSpent: (json['total_spent'] ?? 0).toDouble(),
      lastOrderAt: json['last_order_at'],
    );
  }
}

class GetCustomersParams {
  final int? page;
  final int? pageSize;
  final String? search;
  final int? businessId;

  GetCustomersParams({this.page, this.pageSize, this.search, this.businessId});

  Map<String, dynamic> toQueryParams() {
    final params = <String, dynamic>{};
    if (page != null) params['page'] = page;
    if (pageSize != null) params['page_size'] = pageSize;
    if (search != null && search!.isNotEmpty) params['search'] = search;
    if (businessId != null) params['business_id'] = businessId;
    return params;
  }
}

class CreateCustomerDTO {
  final String name;
  final String? email;
  final String? phone;
  final String? dni;

  CreateCustomerDTO({required this.name, this.email, this.phone, this.dni});

  Map<String, dynamic> toJson() {
    final json = <String, dynamic>{'name': name};
    if (email != null) json['email'] = email;
    if (phone != null) json['phone'] = phone;
    if (dni != null) json['dni'] = dni;
    return json;
  }
}

class UpdateCustomerDTO {
  final String name;
  final String? email;
  final String? phone;
  final String? dni;

  UpdateCustomerDTO({required this.name, this.email, this.phone, this.dni});

  Map<String, dynamic> toJson() {
    final json = <String, dynamic>{'name': name};
    if (email != null) json['email'] = email;
    if (phone != null) json['phone'] = phone;
    if (dni != null) json['dni'] = dni;
    return json;
  }
}
