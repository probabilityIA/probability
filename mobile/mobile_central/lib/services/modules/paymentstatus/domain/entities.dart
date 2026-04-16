class PaymentStatusInfo {
  final int id;
  final String code;
  final String name;
  final String? description;
  final String? category;
  final String? color;
  final String? icon;
  final bool isActive;

  PaymentStatusInfo({required this.id, required this.code, required this.name, this.description, this.category, this.color, this.icon, required this.isActive});

  factory PaymentStatusInfo.fromJson(Map<String, dynamic> json) {
    return PaymentStatusInfo(
      id: json['id'] ?? 0, code: json['code'] ?? '', name: json['name'] ?? '',
      description: json['description'], category: json['category'], color: json['color'],
      icon: json['icon'], isActive: json['is_active'] ?? true,
    );
  }
}

class GetPaymentStatusesParams {
  final bool? isActive;
  GetPaymentStatusesParams({this.isActive});
  Map<String, dynamic> toQueryParams() {
    final params = <String, dynamic>{};
    if (isActive != null) params['is_active'] = isActive;
    return params;
  }
}
