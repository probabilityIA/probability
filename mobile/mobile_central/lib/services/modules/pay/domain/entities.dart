class PaymentGatewayType {
  final int id;
  final String name;
  final String code;
  final String? imageUrl;
  final bool isActive;
  final bool inDevelopment;

  PaymentGatewayType({
    required this.id,
    required this.name,
    required this.code,
    this.imageUrl,
    required this.isActive,
    required this.inDevelopment,
  });

  factory PaymentGatewayType.fromJson(Map<String, dynamic> json) {
    return PaymentGatewayType(
      id: json['id'] ?? 0,
      name: json['name'] ?? '',
      code: json['code'] ?? '',
      imageUrl: json['image_url'],
      isActive: json['is_active'] ?? true,
      inDevelopment: json['in_development'] ?? false,
    );
  }
}

class PaymentGatewayTypesResponse {
  final bool success;
  final List<PaymentGatewayType> data;
  final String? message;

  PaymentGatewayTypesResponse({
    required this.success,
    required this.data,
    this.message,
  });

  factory PaymentGatewayTypesResponse.fromJson(Map<String, dynamic> json) {
    return PaymentGatewayTypesResponse(
      success: json['success'] ?? false,
      data: (json['data'] as List<dynamic>?)
              ?.map((e) => PaymentGatewayType.fromJson(e))
              .toList() ??
          [],
      message: json['message'],
    );
  }
}
