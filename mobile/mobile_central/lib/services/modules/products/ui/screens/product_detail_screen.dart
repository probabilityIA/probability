import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../../../../../shared/theme/app_colors.dart';
import '../../../../../shared/theme/app_tokens.dart';
import '../../../../../shared/utils/formatters.dart';
import '../../../../../shared/widgets/ui/ui.dart';
import '../../domain/entities.dart';
import '../providers/product_provider.dart';
import '../widgets/product_card.dart';

class ProductDetailScreen extends StatefulWidget {
  const ProductDetailScreen({super.key, required this.productId});

  final String productId;

  @override
  State<ProductDetailScreen> createState() => _ProductDetailScreenState();
}

class _ProductDetailScreenState extends State<ProductDetailScreen> {
  Product? _product;

  @override
  void initState() {
    super.initState();
    _product = context.read<ProductProvider>().productById(widget.productId);
    WidgetsBinding.instance.addPostFrameCallback((_) {
      context.read<ProductProvider>().fetchIntegrations(widget.productId);
    });
  }

  @override
  Widget build(BuildContext context) {
    final product = _product;

    return AppScaffold(
      title: product?.name ?? 'Producto',
      subtitle: product == null ? null : 'SKU ${product.sku}',
      onBack: () => Navigator.of(context).pop(),
      body: product == null
          ? const AppEmptyState(
              icon: Icons.sell_outlined,
              title: 'Producto no encontrado',
            )
          : ListView(
              padding: AppSpacing.page,
              children: [
                _Header(product: product),
                const SizedBox(height: 18),
                const AppSectionHeader(title: 'Precios'),
                _PricesCard(product: product),
                const SizedBox(height: 18),
                const AppSectionHeader(title: 'Inventario'),
                _StockCard(product: product),
                const SizedBox(height: 18),
                const AppSectionHeader(
                  title: 'Canales',
                  subtitle: 'Donde esta publicado este producto',
                ),
                const _ChannelsCard(),
                const SizedBox(height: 18),
                const AppSectionHeader(title: 'Dimensiones'),
                _DimensionsCard(product: product),
                const SizedBox(height: 26),
              ],
            ),
    );
  }
}

class _Header extends StatelessWidget {
  const _Header({required this.product});

  final Product product;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final badge = productStockBadge(product);

    return AppCard(
      child: Row(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          ProductThumb(product: product, size: 72),
          const SizedBox(width: 14),
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(product.name, style: theme.textTheme.titleMedium),
                const SizedBox(height: 4),
                Text('SKU ${product.sku}', style: theme.textTheme.bodySmall),
                const SizedBox(height: 10),
                Wrap(
                  spacing: 6,
                  runSpacing: 6,
                  children: [
                    AppStatusChip(dense: true, label: badge.label, tone: badge.tone),
                    AppStatusChip(
                      dense: true,
                      label: product.isActive ? 'Activo' : 'Inactivo',
                      tone: product.isActive ? AppStatusTone.brand : AppStatusTone.neutral,
                    ),
                  ],
                ),
              ],
            ),
          ),
        ],
      ),
    );
  }
}

class _PricesCard extends StatelessWidget {
  const _PricesCard({required this.product});

  final Product product;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final margin = product.costPrice == null
        ? null
        : product.price - product.costPrice!;
    final marginPercent = (margin == null || product.price == 0)
        ? null
        : (margin / product.price) * 100;

    return AppCard(
      child: Column(
        children: [
          AppKeyValueRow(
            label: 'Precio de venta',
            value: AppFormat.money(product.price),
            valueStyle: theme.textTheme.titleMedium?.copyWith(color: AppColors.primary),
          ),
          if (product.compareAtPrice != null) ...[
            const Divider(height: 18),
            AppKeyValueRow(
              label: 'Precio comparativo',
              value: AppFormat.money(product.compareAtPrice),
            ),
          ],
          if (product.costPrice != null) ...[
            const Divider(height: 18),
            AppKeyValueRow(label: 'Costo', value: AppFormat.money(product.costPrice)),
          ],
          if (margin != null) ...[
            const Divider(height: 18),
            AppKeyValueRow(
              label: 'Margen',
              value: marginPercent == null
                  ? AppFormat.money(margin)
                  : '${AppFormat.money(margin)}  (${marginPercent.toStringAsFixed(0)}%)',
            ),
          ],
        ],
      ),
    );
  }
}

class _StockCard extends StatelessWidget {
  const _StockCard({required this.product});

  final Product product;

  static const Map<String, String> _statusLabels = {
    'in_stock': 'Disponible',
    'low_stock': 'Stock bajo',
    'out_of_stock': 'Agotado',
    'on_backorder': 'Bajo pedido',
  };

  @override
  Widget build(BuildContext context) {
    return AppCard(
      child: Column(
        children: [
          AppKeyValueRow(
            label: 'Existencias',
            value: product.manageStock ? '${AppFormat.number(product.stock)} uds' : 'Sin control',
          ),
          const Divider(height: 18),
          AppKeyValueRow(
            label: 'Estado',
            value: _statusLabels[product.stockStatus] ?? product.stockStatus ?? '-',
          ),
          const Divider(height: 18),
          AppKeyValueRow(label: 'Moneda', value: product.currency),
        ],
      ),
    );
  }
}

class _ChannelsCard extends StatelessWidget {
  const _ChannelsCard();

  @override
  Widget build(BuildContext context) {
    return Consumer<ProductProvider>(
      builder: (context, provider, _) {
        final integrations = provider.integrations;
        if (integrations.isEmpty) {
          return AppCard(
            child: Text(
              'Este producto no esta publicado en ningun canal',
              style: Theme.of(context).textTheme.bodySmall,
            ),
          );
        }

        return AppCard(
          child: Column(
            children: [
              for (var i = 0; i < integrations.length; i++) ...[
                if (i > 0) const Divider(height: 20),
                Row(
                  children: [
                    BrandLogo(
                      name: integrations[i].integrationName ??
                          integrations[i].integrationType ??
                          '',
                      size: 36,
                      radius: 10,
                      padding: 5,
                    ),
                    const SizedBox(width: 12),
                    Expanded(
                      child: Column(
                        crossAxisAlignment: CrossAxisAlignment.start,
                        children: [
                          Text(
                            integrations[i].integrationName ??
                                integrations[i].integrationType ??
                                'Canal',
                            style: Theme.of(context).textTheme.titleSmall,
                          ),
                          const SizedBox(height: 2),
                          Text(
                            'ID externo ${integrations[i].externalProductId}',
                            style: Theme.of(context).textTheme.labelSmall,
                            maxLines: 1,
                            overflow: TextOverflow.ellipsis,
                          ),
                        ],
                      ),
                    ),
                  ],
                ),
              ],
            ],
          ),
        );
      },
    );
  }
}

class _DimensionsCard extends StatelessWidget {
  const _DimensionsCard({required this.product});

  final Product product;

  @override
  Widget build(BuildContext context) {
    final dimensions = [product.length, product.width, product.height]
        .map((e) => e == null ? '-' : AppFormat.number(e))
        .join(' x ');

    return AppCard(
      child: Column(
        children: [
          AppKeyValueRow(
            label: 'Peso',
            value: product.weight == null ? '-' : '${AppFormat.number(product.weight)} g',
          ),
          const Divider(height: 18),
          AppKeyValueRow(label: 'Largo x ancho x alto', value: '$dimensions cm'),
          const Divider(height: 18),
          AppKeyValueRow(
            label: 'Creado',
            value: AppFormat.date(AppFormat.parseDate(product.createdAt)),
          ),
        ],
      ),
    );
  }
}
