class ActionEntity {
  final int id;
  final String name;
  final String? description;

  ActionEntity({
    required this.id,
    required this.name,
    this.description,
  });

  factory ActionEntity.fromJson(Map<String, dynamic> json) {
    return ActionEntity(
      id: json['id'] ?? 0,
      name: json['name'] ?? '',
      description: json['description'],
    );
  }
}

class GetActionsParams {
  final int? page;
  final int? pageSize;
  final String? name;

  GetActionsParams({this.page, this.pageSize, this.name});

  Map<String, dynamic> toQueryParams() {
    final params = <String, dynamic>{};
    if (page != null) params['page'] = page;
    if (pageSize != null) params['page_size'] = pageSize;
    if (name != null && name!.isNotEmpty) params['name'] = name;
    return params;
  }
}

class CreateActionDTO {
  final String name;
  final String? description;

  CreateActionDTO({required this.name, this.description});

  Map<String, dynamic> toJson() {
    final json = <String, dynamic>{'name': name};
    if (description != null) json['description'] = description;
    return json;
  }
}

class UpdateActionDTO {
  final String? name;
  final String? description;

  UpdateActionDTO({this.name, this.description});

  Map<String, dynamic> toJson() {
    final json = <String, dynamic>{};
    if (name != null) json['name'] = name;
    if (description != null) json['description'] = description;
    return json;
  }
}
