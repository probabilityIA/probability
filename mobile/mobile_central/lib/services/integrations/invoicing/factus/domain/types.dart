class FactusConfig {
  final int? numberingRangeId;
  final String? defaultTaxRate;
  final String? paymentForm;
  final String? paymentMethodCode;
  final String? legalOrganizationId;
  final String? tributeId;
  final String? identificationDocumentId;
  final String? municipalityId;

  FactusConfig({
    this.numberingRangeId,
    this.defaultTaxRate,
    this.paymentForm,
    this.paymentMethodCode,
    this.legalOrganizationId,
    this.tributeId,
    this.identificationDocumentId,
    this.municipalityId,
  });

  factory FactusConfig.fromJson(Map<String, dynamic> json) {
    return FactusConfig(
      numberingRangeId: json['numbering_range_id'],
      defaultTaxRate: json['default_tax_rate'],
      paymentForm: json['payment_form'],
      paymentMethodCode: json['payment_method_code'],
      legalOrganizationId: json['legal_organization_id'],
      tributeId: json['tribute_id'],
      identificationDocumentId: json['identification_document_id'],
      municipalityId: json['municipality_id'],
    );
  }

  Map<String, dynamic> toJson() {
    final json = <String, dynamic>{};
    if (numberingRangeId != null) {
      json['numbering_range_id'] = numberingRangeId;
    }
    if (defaultTaxRate != null) json['default_tax_rate'] = defaultTaxRate;
    if (paymentForm != null) json['payment_form'] = paymentForm;
    if (paymentMethodCode != null) {
      json['payment_method_code'] = paymentMethodCode;
    }
    if (legalOrganizationId != null) {
      json['legal_organization_id'] = legalOrganizationId;
    }
    if (tributeId != null) json['tribute_id'] = tributeId;
    if (identificationDocumentId != null) {
      json['identification_document_id'] = identificationDocumentId;
    }
    if (municipalityId != null) json['municipality_id'] = municipalityId;
    return json;
  }
}

class FactusCredentials {
  final String? clientId;
  final String? clientSecret;
  final String? username;
  final String? password;
  final String? apiUrl;

  FactusCredentials({
    this.clientId,
    this.clientSecret,
    this.username,
    this.password,
    this.apiUrl,
  });

  factory FactusCredentials.fromJson(Map<String, dynamic> json) {
    return FactusCredentials(
      clientId: json['client_id'],
      clientSecret: json['client_secret'],
      username: json['username'],
      password: json['password'],
      apiUrl: json['api_url'],
    );
  }

  Map<String, dynamic> toJson() {
    final json = <String, dynamic>{};
    if (clientId != null) json['client_id'] = clientId;
    if (clientSecret != null) json['client_secret'] = clientSecret;
    if (username != null) json['username'] = username;
    if (password != null) json['password'] = password;
    if (apiUrl != null) json['api_url'] = apiUrl;
    return json;
  }
}
