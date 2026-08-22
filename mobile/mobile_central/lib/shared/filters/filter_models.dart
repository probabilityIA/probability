import 'package:flutter/material.dart';

class SearchField {
  const SearchField({
    required this.key,
    required this.label,
    required this.hint,
  });

  final String key;
  final String label;
  final String hint;
}

class FilterOption {
  const FilterOption({required this.value, required this.label});

  final String value;
  final String label;
}

class FilterDimension {
  const FilterDimension({
    required this.key,
    required this.label,
    required this.icon,
    required this.options,
    this.imageUrl,
    this.accent,
  });

  final String key;
  final String label;
  final IconData icon;
  final List<FilterOption> options;
  final String? imageUrl;
  final Color? accent;

  FilterOption? optionFor(String value) {
    for (final option in options) {
      if (option.value == value) return option;
    }
    return null;
  }
}

class ActiveFilter {
  const ActiveFilter({
    required this.dimensionKey,
    required this.dimensionLabel,
    required this.value,
    required this.valueLabel,
    this.imageUrl,
    this.accent,
  });

  final String dimensionKey;
  final String dimensionLabel;
  final String value;
  final String valueLabel;
  final String? imageUrl;
  final Color? accent;
}

class FilterSelection {
  const FilterSelection([this.values = const {}]);

  final Map<String, String> values;

  bool get isEmpty => values.isEmpty;
  bool get isNotEmpty => values.isNotEmpty;
  int get count => values.length;

  String? operator [](String key) => values[key];

  FilterSelection withValue(String key, String? value) {
    final next = Map<String, String>.from(values);
    if (value == null || value.isEmpty) {
      next.remove(key);
    } else {
      next[key] = value;
    }
    return FilterSelection(next);
  }

  FilterSelection without(String key) => withValue(key, null);

  FilterSelection get cleared => const FilterSelection();

  bool? boolFor(String key) {
    final value = values[key];
    if (value == null) return null;
    return value == 'true';
  }

  List<ActiveFilter> describe(List<FilterDimension> dimensions) {
    final out = <ActiveFilter>[];
    for (final dimension in dimensions) {
      final value = values[dimension.key];
      if (value == null) continue;
      final option = dimension.optionFor(value);
      if (option == null) continue;
      out.add(ActiveFilter(
        dimensionKey: dimension.key,
        dimensionLabel: dimension.label,
        value: value,
        valueLabel: option.label,
        imageUrl: dimension.imageUrl,
        accent: dimension.accent,
      ));
    }
    return out;
  }
}
