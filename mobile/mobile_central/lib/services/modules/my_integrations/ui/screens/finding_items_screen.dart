import 'package:flutter/material.dart';
import 'package:provider/provider.dart';

import '../../../../../shared/theme/app_colors.dart';
import '../../../../../shared/theme/app_tokens.dart';
import '../../../../../shared/utils/formatters.dart';
import '../../../../../shared/widgets/ui/ui.dart';
import '../../domain/saved_comparison_entities.dart';
import '../providers/comparison_lists_provider.dart';

class FindingItemsScreen extends StatefulWidget {
  const FindingItemsScreen({
    super.key,
    required this.finding,
    this.businessId,
  });

  final Finding finding;
  final int? businessId;

  @override
  State<FindingItemsScreen> createState() => _FindingItemsScreenState();
}

class _FindingItemsScreenState extends State<FindingItemsScreen> {
  final TextEditingController _search = TextEditingController();

  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addPostFrameCallback((_) {
      if (!mounted) return;
      context.read<FindingItemsProvider>().load(
            code: widget.finding.code,
            businessId: widget.businessId,
          );
    });
  }

  @override
  void dispose() {
    _search.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return Consumer<FindingItemsProvider>(
      builder: (context, provider, _) {
        return AppScaffold(
          title: widget.finding.title,
          subtitle: '${AppFormat.number(widget.finding.count)} productos',
          onBack: () => Navigator.of(context).maybePop(),
          body: Column(
            children: [
              Padding(
                padding: const EdgeInsets.fromLTRB(16, 10, 16, 8),
                child: Column(
                  children: [
                    AppCard(
                      padding: const EdgeInsets.all(12),
                      child: Text(
                        widget.finding.detail,
                        style: Theme.of(context).textTheme.bodySmall,
                      ),
                    ),
                    const SizedBox(height: 8),
                    TextField(
                      controller: _search,
                      textInputAction: TextInputAction.search,
                      decoration: const InputDecoration(
                        isDense: true,
                        hintText: 'Buscar SKU o producto',
                        prefixIcon: Icon(Icons.search_rounded, size: 19),
                      ),
                      onSubmitted: provider.setSearch,
                    ),
                  ],
                ),
              ),
              Expanded(
                child: PaginatedListView<FindingItem>(
                  controller: provider.list,
                  unitLabel: 'productos',
                  emptyIcon: Icons.check_circle_outline_rounded,
                  emptyTitle: 'Nada por revisar',
                  emptyMessage: 'Este hallazgo ya no tiene productos.',
                  placeholderHeight: 96,
                  onRefresh: () => provider.load(
                    code: widget.finding.code,
                    businessId: widget.businessId,
                  ),
                  itemBuilder: (context, item, index) => _ItemCard(item: item),
                ),
              ),
            ],
          ),
        );
      },
    );
  }
}

class _ItemCard extends StatelessWidget {
  const _ItemCard({required this.item});

  final FindingItem item;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    return AppCard(
      padding: const EdgeInsets.all(12),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text(
            item.name == null || item.name!.isEmpty ? item.sku : item.name!,
            style: theme.textTheme.titleSmall,
            maxLines: 2,
            overflow: TextOverflow.ellipsis,
          ),
          const SizedBox(height: 2),
          Text('SKU ${item.sku.isEmpty ? "-" : item.sku}',
              style: theme.textTheme.labelSmall),
          if (item.isPaired) ...[
            const SizedBox(height: 10),
            _Pair(item: item),
          ],
          if (item.detail != null && item.detail!.isNotEmpty) ...[
            const SizedBox(height: 8),
            Text(item.detail!, style: theme.textTheme.bodySmall),
          ],
          if (item.isCrossChannel) ...[
            const SizedBox(height: 10),
            if (item.presentIn.isNotEmpty)
              _ChannelLine(
                label: 'esta en',
                names: item.presentIn,
                color: AppColors.success,
              ),
            if (item.missingIn.isNotEmpty) ...[
              const SizedBox(height: 6),
              _ChannelLine(
                label: 'falta en',
                names: item.missingIn,
                color: AppColors.error,
              ),
            ],
          ] else if (item.channels.isNotEmpty) ...[
            const SizedBox(height: 10),
            _ChannelLine(
              label: 'en',
              names: item.channels,
              color: AppColors.textMuted,
            ),
          ],
        ],
      ),
    );
  }
}

class _Pair extends StatelessWidget {
  const _Pair({required this.item});

  final FindingItem item;

  @override
  Widget build(BuildContext context) {
    return Row(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Expanded(
          child: _Side(
            label: 'En el canal',
            sku: item.counterpartSku ?? '-',
            name: item.counterpartName,
            qty: item.channelQty,
            flagged: item.fixSide == 'channel',
          ),
        ),
        const Padding(
          padding: EdgeInsets.symmetric(horizontal: 8, vertical: 14),
          child: Icon(Icons.compare_arrows_rounded,
              size: 15, color: AppColors.textDisabled),
        ),
        Expanded(
          child: _Side(
            label: 'En Probability',
            sku: item.sku,
            name: item.name,
            qty: item.ownQty,
            flagged: item.fixSide == 'probability',
          ),
        ),
      ],
    );
  }
}

class _Side extends StatelessWidget {
  const _Side({
    required this.label,
    required this.sku,
    required this.name,
    required this.qty,
    required this.flagged,
  });

  final String label;
  final String sku;
  final String? name;
  final int? qty;
  final bool flagged;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final color = flagged ? AppColors.error : AppColors.textSecondary;

    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 8),
      decoration: BoxDecoration(
        color: flagged
            ? AppColors.error.withValues(alpha: 0.06)
            : AppColors.surfaceMuted,
        borderRadius: AppRadius.smAll,
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text(
            label.toUpperCase(),
            style: theme.textTheme.labelSmall?.copyWith(fontSize: 9),
          ),
          const SizedBox(height: 3),
          Text(
            sku,
            style: theme.textTheme.bodySmall?.copyWith(
              fontWeight: FontWeight.w700,
              color: color,
            ),
            maxLines: 1,
            overflow: TextOverflow.ellipsis,
          ),
          if (name != null && name!.isNotEmpty) ...[
            const SizedBox(height: 2),
            Text(
              name!,
              style: theme.textTheme.labelSmall,
              maxLines: 1,
              overflow: TextOverflow.ellipsis,
            ),
          ],
          if (qty != null) ...[
            const SizedBox(height: 3),
            Text('${AppFormat.number(qty)} u.',
                style: theme.textTheme.labelSmall),
          ],
        ],
      ),
    );
  }
}

class _ChannelLine extends StatelessWidget {
  const _ChannelLine({
    required this.label,
    required this.names,
    required this.color,
  });

  final String label;
  final List<String> names;
  final Color color;

  @override
  Widget build(BuildContext context) {
    return Row(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(
          label,
          style: TextStyle(
            fontFamily: 'Inter',
            fontSize: 11,
            fontWeight: FontWeight.w600,
            color: color,
          ),
        ),
        const SizedBox(width: 8),
        Expanded(
          child: Wrap(
            spacing: 6,
            runSpacing: 6,
            children: names
                .map((name) => Container(
                      padding:
                          const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
                      decoration: BoxDecoration(
                        color: color.withValues(alpha: 0.10),
                        borderRadius: AppRadius.pillAll,
                      ),
                      child: Text(
                        name,
                        style: TextStyle(
                          fontFamily: 'Inter',
                          fontSize: 11,
                          fontWeight: FontWeight.w600,
                          color: color,
                        ),
                      ),
                    ))
                .toList(),
          ),
        ),
      ],
    );
  }
}
