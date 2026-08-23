import 'package:flutter/material.dart';

import '../../filters/filter_models.dart';
import '../../theme/app_colors.dart';
import '../../theme/app_tokens.dart';
import 'brand_logo.dart';

class AppFilterBar extends StatelessWidget {
  const AppFilterBar({
    super.key,
    required this.controller,
    required this.searchFields,
    required this.selectedField,
    required this.onFieldChanged,
    required this.onSearchChanged,
    this.dimensions = const [],
    this.selection = const FilterSelection(),
    this.onSelectionChanged,
    this.summary,
    this.trailing,
  });

  final TextEditingController controller;
  final List<SearchField> searchFields;
  final String selectedField;
  final ValueChanged<String> onFieldChanged;
  final ValueChanged<String> onSearchChanged;
  final List<FilterDimension> dimensions;
  final FilterSelection selection;
  final ValueChanged<FilterSelection>? onSelectionChanged;
  final String? summary;
  final Widget? trailing;

  SearchField get _field => searchFields.firstWhere(
        (f) => f.key == selectedField,
        orElse: () => searchFields.first,
      );

  Future<void> _openFilters(BuildContext context) async {
    final result = await showModalBottomSheet<FilterSelection>(
      context: context,
      isScrollControlled: true,
      backgroundColor: AppColors.surface,
      shape: const RoundedRectangleBorder(
        borderRadius: BorderRadius.vertical(top: Radius.circular(AppRadius.xl)),
      ),
      builder: (_) => _FilterSheet(dimensions: dimensions, selection: selection),
    );
    if (result != null) onSelectionChanged?.call(result);
  }

  @override
  Widget build(BuildContext context) {
    final active = selection.describe(dimensions);
    final hasFilters = dimensions.isNotEmpty && onSelectionChanged != null;

    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Padding(
          padding: const EdgeInsets.fromLTRB(16, 12, 16, 0),
          child: Row(
            children: [
              if (hasFilters) ...[
                _FilterButton(
                  count: selection.count,
                  onTap: () => _openFilters(context),
                ),
                const SizedBox(width: 10),
              ],
              Expanded(
                child: _SearchInput(
                  controller: controller,
                  field: _field,
                  fields: searchFields,
                  onFieldChanged: onFieldChanged,
                  onChanged: onSearchChanged,
                ),
              ),
              if (trailing != null) ...[
                const SizedBox(width: 6),
                trailing!,
              ],
            ],
          ),
        ),
        if (active.isNotEmpty)
          Padding(
            padding: const EdgeInsets.fromLTRB(16, 10, 16, 0),
            child: Wrap(
              spacing: 7,
              runSpacing: 7,
              crossAxisAlignment: WrapCrossAlignment.center,
              children: [
                for (final filter in active)
                  _ActiveChip(
                    filter: filter,
                    onRemove: () => onSelectionChanged
                        ?.call(selection.without(filter.dimensionKey)),
                  ),
                _ClearAllButton(
                  onTap: () => onSelectionChanged?.call(selection.cleared),
                ),
              ],
            ),
          ),
        if (summary != null)
          Padding(
            padding: const EdgeInsets.fromLTRB(16, 10, 16, 0),
            child: Text(summary!, style: Theme.of(context).textTheme.labelSmall),
          ),
        const SizedBox(height: 10),
      ],
    );
  }
}

class _FilterButton extends StatelessWidget {
  const _FilterButton({required this.count, required this.onTap});

  final int count;
  final VoidCallback onTap;

  @override
  Widget build(BuildContext context) {
    final on = count > 0;
    final scheme = Theme.of(context).colorScheme;
    return Material(
      color: on ? scheme.primaryContainer : AppColors.surface,
      borderRadius: AppRadius.mdAll,
      child: InkWell(
        borderRadius: AppRadius.mdAll,
        onTap: onTap,
        child: Container(
          height: 48,
          width: 52,
          alignment: Alignment.center,
          decoration: BoxDecoration(
            borderRadius: AppRadius.mdAll,
            border: Border.all(
              color: on ? scheme.primary : AppColors.border,
            ),
          ),
          child: Stack(
            clipBehavior: Clip.none,
            children: [
              Icon(
                Icons.tune_rounded,
                size: 21,
                color: on ? scheme.primary : AppColors.textMuted,
              ),
              if (on)
                Positioned(
                  right: -8,
                  top: -6,
                  child: Container(
                    padding: const EdgeInsets.symmetric(horizontal: 5, vertical: 1),
                    decoration: BoxDecoration(
                      color: scheme.primary,
                      borderRadius: BorderRadius.circular(999),
                    ),
                    child: Text(
                      '$count',
                      style: TextStyle(
                        fontFamily: 'Inter',
                        fontSize: 10,
                        fontWeight: FontWeight.w700,
                        color: scheme.onPrimary,
                        height: 1.3,
                      ),
                    ),
                  ),
                ),
            ],
          ),
        ),
      ),
    );
  }
}

class _SearchInput extends StatelessWidget {
  const _SearchInput({
    required this.controller,
    required this.field,
    required this.fields,
    required this.onFieldChanged,
    required this.onChanged,
  });

  final TextEditingController controller;
  final SearchField field;
  final List<SearchField> fields;
  final ValueChanged<String> onFieldChanged;
  final ValueChanged<String> onChanged;

  @override
  Widget build(BuildContext context) {
    return ValueListenableBuilder<TextEditingValue>(
      valueListenable: controller,
      builder: (context, value, _) {
        return TextField(
          controller: controller,
          onChanged: onChanged,
          textInputAction: TextInputAction.search,
          decoration: InputDecoration(
            hintText: field.hint,
            filled: true,
            fillColor: AppColors.surface,
            contentPadding: const EdgeInsets.symmetric(vertical: 13),
            prefixIcon: fields.length < 2
                ? const Icon(Icons.search, size: 20)
                : _FieldPicker(
                    field: field,
                    fields: fields,
                    onChanged: onFieldChanged,
                  ),
            prefixIconConstraints: const BoxConstraints(minWidth: 0),
            suffixIcon: value.text.isEmpty
                ? null
                : IconButton(
                    icon: const Icon(Icons.close, size: 18),
                    onPressed: () {
                      controller.clear();
                      onChanged('');
                    },
                  ),
          ),
        );
      },
    );
  }
}

class _FieldPicker extends StatelessWidget {
  const _FieldPicker({
    required this.field,
    required this.fields,
    required this.onChanged,
  });

  final SearchField field;
  final List<SearchField> fields;
  final ValueChanged<String> onChanged;

  @override
  Widget build(BuildContext context) {
    return PopupMenuButton<String>(
      tooltip: 'Buscar por',
      initialValue: field.key,
      position: PopupMenuPosition.under,
      onSelected: onChanged,
      itemBuilder: (context) => [
        for (final f in fields)
          PopupMenuItem<String>(
            value: f.key,
            child: Row(
              children: [
                Icon(
                  f.key == field.key
                      ? Icons.radio_button_checked
                      : Icons.radio_button_unchecked,
                  size: 17,
                  color: f.key == field.key
                      ? AppColors.primary
                      : AppColors.textDisabled,
                ),
                const SizedBox(width: 10),
                Text(f.label),
              ],
            ),
          ),
      ],
      child: Padding(
        padding: const EdgeInsets.only(left: 12, right: 8),
        child: Row(
          mainAxisSize: MainAxisSize.min,
          children: [
            const Icon(Icons.search, size: 19, color: AppColors.textMuted),
            const SizedBox(width: 7),
            Text(
              field.label,
              style: const TextStyle(
                fontFamily: 'Inter',
                fontSize: 12.5,
                fontWeight: FontWeight.w600,
                color: AppColors.textSecondary,
              ),
            ),
            const Icon(Icons.arrow_drop_down_rounded,
                size: 20, color: AppColors.textMuted),
            Container(
              width: 1,
              height: 22,
              margin: const EdgeInsets.only(left: 4, right: 2),
              color: AppColors.border,
            ),
          ],
        ),
      ),
    );
  }
}

class _ActiveChip extends StatelessWidget {
  const _ActiveChip({required this.filter, required this.onRemove});

  final ActiveFilter filter;
  final VoidCallback onRemove;

  @override
  Widget build(BuildContext context) {
    final scheme = Theme.of(context).colorScheme;
    return Material(
      color: scheme.primaryContainer,
      borderRadius: AppRadius.pillAll,
      child: InkWell(
        borderRadius: AppRadius.pillAll,
        onTap: onRemove,
        child: Padding(
          padding: const EdgeInsets.fromLTRB(9, 5, 8, 5),
          child: Row(
            mainAxisSize: MainAxisSize.min,
            children: [
              if (filter.imageUrl != null || filter.accent != null) ...[
                BrandLogo(
                  name: filter.dimensionLabel,
                  imageUrl: filter.imageUrl,
                  size: 18,
                  radius: 5,
                  padding: 2,
                ),
                const SizedBox(width: 6),
              ],
              Text(
                filter.imageUrl != null || filter.accent != null
                    ? filter.valueLabel
                    : '${filter.dimensionLabel}: ${filter.valueLabel}',
                style: TextStyle(
                  fontFamily: 'Inter',
                  fontSize: 12,
                  fontWeight: FontWeight.w600,
                  color: scheme.onPrimaryContainer,
                ),
              ),
              const SizedBox(width: 5),
              Icon(Icons.close_rounded, size: 15, color: scheme.primary),
            ],
          ),
        ),
      ),
    );
  }
}

class _ClearAllButton extends StatelessWidget {
  const _ClearAllButton({required this.onTap});

  final VoidCallback onTap;

  @override
  Widget build(BuildContext context) {
    return TextButton(
      onPressed: onTap,
      style: TextButton.styleFrom(
        minimumSize: Size.zero,
        padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
        tapTargetSize: MaterialTapTargetSize.shrinkWrap,
      ),
      child: const Text(
        'Limpiar todo',
        style: TextStyle(
          fontFamily: 'Inter',
          fontSize: 12,
          fontWeight: FontWeight.w600,
          color: AppColors.textMuted,
        ),
      ),
    );
  }
}

class _FilterSheet extends StatefulWidget {
  const _FilterSheet({required this.dimensions, required this.selection});

  final List<FilterDimension> dimensions;
  final FilterSelection selection;

  @override
  State<_FilterSheet> createState() => _FilterSheetState();
}

class _FilterSheetState extends State<_FilterSheet> {
  late FilterSelection _draft = widget.selection;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final maxHeight = MediaQuery.sizeOf(context).height * 0.82;

    return SafeArea(
      top: false,
      child: ConstrainedBox(
        constraints: BoxConstraints(maxHeight: maxHeight),
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            Padding(
              padding: const EdgeInsets.fromLTRB(20, 14, 12, 6),
              child: Row(
                children: [
                  Expanded(
                    child: Text('Filtros', style: theme.textTheme.titleMedium),
                  ),
                  if (_draft.isNotEmpty)
                    TextButton(
                      onPressed: () => setState(() => _draft = _draft.cleared),
                      child: const Text('Limpiar'),
                    ),
                ],
              ),
            ),
            Flexible(
              child: ListView.builder(
                shrinkWrap: true,
                padding: const EdgeInsets.fromLTRB(20, 4, 20, 12),
                itemCount: widget.dimensions.length,
                itemBuilder: (context, index) {
                  final dimension = widget.dimensions[index];
                  final current = _draft[dimension.key];
                  return Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Row(
                        children: [
                          if (dimension.imageUrl != null ||
                              dimension.accent != null)
                            BrandLogo(
                              name: dimension.label,
                              imageUrl: dimension.imageUrl,
                              size: 22,
                              radius: 6,
                              padding: 3,
                            )
                          else
                            Icon(dimension.icon,
                                size: 16, color: AppColors.textMuted),
                          const SizedBox(width: 8),
                          Expanded(
                            child: Text(
                              dimension.label,
                              style: theme.textTheme.titleSmall,
                              maxLines: 1,
                              overflow: TextOverflow.ellipsis,
                            ),
                          ),
                        ],
                      ),
                      const SizedBox(height: 9),
                      Wrap(
                        spacing: 7,
                        runSpacing: 7,
                        children: [
                          for (final option in dimension.options)
                            ChoiceChip(
                              label: Text(option.label),
                              selected: current == option.value,
                              onSelected: (on) => setState(() {
                                _draft = _draft.withValue(
                                  dimension.key,
                                  on ? option.value : null,
                                );
                              }),
                              labelStyle: TextStyle(
                                fontFamily: 'Inter',
                                fontSize: 12.5,
                                fontWeight: FontWeight.w600,
                                color: current == option.value
                                    ? theme.colorScheme.onPrimaryContainer
                                    : AppColors.textSecondary,
                              ),
                            ),
                        ],
                      ),
                      SizedBox(
                        height: index == widget.dimensions.length - 1 ? 6 : 20,
                      ),
                    ],
                  );
                },
              ),
            ),
            Padding(
              padding: const EdgeInsets.fromLTRB(20, 6, 20, 14),
              child: SizedBox(
                width: double.infinity,
                child: FilledButton(
                  onPressed: () => Navigator.pop(context, _draft),
                  child: Text(
                    _draft.isEmpty
                        ? 'Ver todo'
                        : 'Aplicar ${_draft.count} filtro${_draft.count == 1 ? '' : 's'}',
                  ),
                ),
              ),
            ),
          ],
        ),
      ),
    );
  }
}
