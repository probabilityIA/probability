import 'package:flutter/material.dart';
import '../../../../../shared/theme/app_colors.dart';
import '../../../../../shared/theme/app_tokens.dart';
import '../../../../../shared/utils/formatters.dart';
import '../../../../../shared/widgets/ui/ui.dart';
import '../../domain/entities.dart';

({String label, AppStatusTone tone}) productStockBadge(Product product) {
  if (!product.manageStock) return (label: 'Sin control', tone: AppStatusTone.neutral);
  if (product.stock <= 0) return (label: 'Agotado', tone: AppStatusTone.error);
  if (product.stock < 20) return (label: 'Stock bajo', tone: AppStatusTone.warning);
  return (label: 'Disponible', tone: AppStatusTone.success);
}

class ProductThumb extends StatelessWidget {
  const ProductThumb({super.key, required this.product, this.size = 54});

  final Product product;
  final double size;

  @override
  Widget build(BuildContext context) {
    final url = product.thumbnail ?? product.imageUrl;
    if (url != null && url.trim().isNotEmpty) {
      return ClipRRect(
        borderRadius: AppRadius.mdAll,
        child: Image.network(
          url,
          width: size,
          height: size,
          fit: BoxFit.cover,
          webHtmlElementStrategy: WebHtmlElementStrategy.fallback,
          errorBuilder: (context, error, stack) => _Placeholder(size: size),
        ),
      );
    }
    return _Placeholder(size: size);
  }
}

class _Placeholder extends StatelessWidget {
  const _Placeholder({required this.size});

  final double size;

  @override
  Widget build(BuildContext context) {
    return Container(
      width: size,
      height: size,
      decoration: BoxDecoration(
        color: AppColors.surfaceMuted,
        borderRadius: AppRadius.mdAll,
      ),
      child: Icon(Icons.image_outlined, size: size * 0.36, color: AppColors.textDisabled),
    );
  }
}

class ProductCard extends StatelessWidget {
  const ProductCard({super.key, required this.product, this.onTap});

  final Product product;
  final VoidCallback? onTap;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final badge = productStockBadge(product);

    return AppCard(
      padding: const EdgeInsets.all(13),
      onTap: onTap,
      child: Row(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          ProductThumb(product: product),
          const SizedBox(width: 12),
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  product.name,
                  style: theme.textTheme.titleSmall,
                  maxLines: 2,
                  overflow: TextOverflow.ellipsis,
                ),
                const SizedBox(height: 3),
                Text('SKU ${product.sku}', style: theme.textTheme.labelSmall),
                const SizedBox(height: 9),
                Row(
                  children: [
                    AppStatusChip(dense: true, label: badge.label, tone: badge.tone),
                    const SizedBox(width: 6),
                    if (product.manageStock)
                      Text(
                        '${AppFormat.number(product.stock)} uds',
                        style: theme.textTheme.labelSmall,
                      ),
                    if (!product.isActive) ...[
                      const SizedBox(width: 6),
                      const AppStatusChip(dense: true, label: 'Inactivo', tone: AppStatusTone.neutral),
                    ],
                  ],
                ),
              ],
            ),
          ),
          const SizedBox(width: 10),
          Column(
            crossAxisAlignment: CrossAxisAlignment.end,
            children: [
              Text(
                AppFormat.money(product.price),
                style: theme.textTheme.titleSmall?.copyWith(fontSize: 15),
              ),
              if (product.compareAtPrice != null && product.compareAtPrice! > product.price) ...[
                const SizedBox(height: 2),
                Text(
                  AppFormat.money(product.compareAtPrice),
                  style: theme.textTheme.labelSmall?.copyWith(
                    decoration: TextDecoration.lineThrough,
                  ),
                ),
              ],
            ],
          ),
        ],
      ),
    );
  }
}
