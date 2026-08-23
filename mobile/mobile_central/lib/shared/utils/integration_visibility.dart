class IntegrationVisibility {
  const IntegrationVisibility._();

  static const Set<String> hiddenCategories = {'internal', 'platform'};

  static const Set<String> hiddenTypes = {'envioclick', 'envio_click'};

  static bool isHiddenCategory(String? category) {
    if (category == null) return false;
    return hiddenCategories.contains(category.trim().toLowerCase());
  }

  static bool isHiddenType(String? type) {
    if (type == null) return false;
    return hiddenTypes.contains(type.trim().toLowerCase().replaceAll(' ', ''));
  }

  static bool isVisible({String? category, String? type, String? name}) {
    if (isHiddenCategory(category)) return false;
    if (isHiddenType(type)) return false;
    if (isHiddenType(name)) return false;
    return true;
  }
}
