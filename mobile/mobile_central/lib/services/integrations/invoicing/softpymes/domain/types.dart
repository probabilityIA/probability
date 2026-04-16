class SoftpymesConfig {
  final String? companyNit;
  final String? companyName;
  final String? referer;
  final String? defaultCustomerNit;
  final int? resolutionId;
  final String? branchCode;
  final String? customerBranchCode;
  final String? sellerNit;

  SoftpymesConfig({
    this.companyNit,
    this.companyName,
    this.referer,
    this.defaultCustomerNit,
    this.resolutionId,
    this.branchCode,
    this.customerBranchCode,
    this.sellerNit,
  });

  factory SoftpymesConfig.fromJson(Map<String, dynamic> json) {
    return SoftpymesConfig(
      companyNit: json['company_nit'],
      companyName: json['company_name'],
      referer: json['referer'],
      defaultCustomerNit: json['default_customer_nit'],
      resolutionId: json['resolution_id'],
      branchCode: json['branch_code'],
      customerBranchCode: json['customer_branch_code'],
      sellerNit: json['seller_nit'],
    );
  }

  Map<String, dynamic> toJson() {
    final json = <String, dynamic>{};
    if (companyNit != null) json['company_nit'] = companyNit;
    if (companyName != null) json['company_name'] = companyName;
    if (referer != null) json['referer'] = referer;
    if (defaultCustomerNit != null) {
      json['default_customer_nit'] = defaultCustomerNit;
    }
    if (resolutionId != null) json['resolution_id'] = resolutionId;
    if (branchCode != null) json['branch_code'] = branchCode;
    if (customerBranchCode != null) {
      json['customer_branch_code'] = customerBranchCode;
    }
    if (sellerNit != null) json['seller_nit'] = sellerNit;
    return json;
  }
}

class SoftpymesCredentials {
  final String? apiKey;
  final String? apiSecret;

  SoftpymesCredentials({
    this.apiKey,
    this.apiSecret,
  });

  factory SoftpymesCredentials.fromJson(Map<String, dynamic> json) {
    return SoftpymesCredentials(
      apiKey: json['api_key'],
      apiSecret: json['api_secret'],
    );
  }

  Map<String, dynamic> toJson() {
    final json = <String, dynamic>{};
    if (apiKey != null) json['api_key'] = apiKey;
    if (apiSecret != null) json['api_secret'] = apiSecret;
    return json;
  }
}
