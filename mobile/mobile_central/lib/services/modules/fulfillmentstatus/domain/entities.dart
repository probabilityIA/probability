class FulfillmentStatusInfo {
  final int id;
  final String code;
  final String name;
  final String? description;
  final String? category;
  final String? color;

  FulfillmentStatusInfo({required this.id, required this.code, required this.name, this.description, this.category, this.color});

  factory FulfillmentStatusInfo.fromJson(Map<String, dynamic> json) {
    return FulfillmentStatusInfo(
      id: json['id'] ?? 0, code: json['code'] ?? '', name: json['name'] ?? '',
      description: json['description'], category: json['category'], color: json['color'],
    );
  }
}
