class Resource {
  final int id;
  final String name;
  final String? description;
  final int? businessTypeId;
  final String? businessTypeName;

  Resource({
    required this.id,
    required this.name,
    this.description,
    this.businessTypeId,
    this.businessTypeName,
  });

  factory Resource.fromJson(Map<String, dynamic> json) {
    return Resource(
      id: json['id'] ?? 0,
      name: json['name'] ?? '',
      description: json['description'],
      businessTypeId: json['business_type_id'],
      businessTypeName: json['business_type_name'],
    );
  }
}

class GetResourcesParams {
  final int? page;
  final int? pageSize;
  final String? name;
  final String? description;
  final int? businessTypeId;

  GetResourcesParams({
    this.page,
    this.pageSize,
    this.name,
    this.description,
    this.businessTypeId,
  });

  Map<String, dynamic> toQueryParams() {
    final params = <String, dynamic>{};
    if (page != null) params['page'] = page;
    if (pageSize != null) params['page_size'] = pageSize;
    if (name != null && name!.isNotEmpty) params['name'] = name;
    if (description != null && description!.isNotEmpty) {
      params['description'] = description;
    }
    if (businessTypeId != null) params['business_type_id'] = businessTypeId;
    return params;
  }
}

class CreateResourceDTO {
  final String name;
  final String? description;
  final int? businessTypeId;

  CreateResourceDTO({required this.name, this.description, this.businessTypeId});

  Map<String, dynamic> toJson() {
    final json = <String, dynamic>{'name': name};
    if (description != null) json['description'] = description;
    if (businessTypeId != null) json['business_type_id'] = businessTypeId;
    return json;
  }
}

class UpdateResourceDTO {
  final String? name;
  final String? description;
  final int? businessTypeId;

  UpdateResourceDTO({this.name, this.description, this.businessTypeId});

  Map<String, dynamic> toJson() {
    final json = <String, dynamic>{};
    if (name != null) json['name'] = name;
    if (description != null) json['description'] = description;
    if (businessTypeId != null) json['business_type_id'] = businessTypeId;
    return json;
  }
}
