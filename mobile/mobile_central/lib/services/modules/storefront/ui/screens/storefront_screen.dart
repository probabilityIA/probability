import 'dart:async';
import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../../../../../shared/theme/app_colors.dart';
import '../../../../../shared/theme/app_tokens.dart';
import '../../../../../shared/utils/formatters.dart';
import '../../../../../shared/widgets/ui/ui.dart';
import '../../domain/entities.dart';
import '../providers/storefront_provider.dart';

class StorefrontScreen extends StatefulWidget {
  const StorefrontScreen({super.key, this.businessId});

  final int? businessId;

  @override
  State<StorefrontScreen> createState() => _StorefrontScreenState();
}

class _StorefrontScreenState extends State<StorefrontScreen> {
  final _searchController = TextEditingController();
  Timer? _debounce;
  bool _onlyFeatured = false;

  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addPostFrameCallback((_) => _refresh());
  }

  @override
  void dispose() {
    _debounce?.cancel();
    _searchController.dispose();
    super.dispose();
  }

  void _refresh() {
    context.read<StorefrontProvider>().fetchCatalog(businessId: widget.businessId);
  }

  void _onSearch(String value) {
    _debounce?.cancel();
    _debounce = Timer(const Duration(milliseconds: 400), () {
      if (mounted) setState(() {});
    });
  }

  @override
  Widget build(BuildContext context) {
    return Consumer<StorefrontProvider>(
      builder: (context, provider, _) {
        if (provider.isLoading && provider.catalogProducts.isEmpty) {
          return const AppListSkeleton();
        }
        if (provider.error != null && provider.catalogProducts.isEmpty) {
          return AppErrorState(message: provider.error!, onRetry: _refresh);
        }
        if (provider.catalogProducts.isEmpty) {
          return AppEmptyState(
            icon: Icons.storefront_outlined,
            title: 'Catalogo vacio',
            message: 'Publica productos para que aparezcan en tu tienda online.',
            actionLabel: 'Actualizar',
            onAction: _refresh,
          );
        }

        final query = _searchController.text.trim().toLowerCase();
        final rows = provider.catalogProducts.where((product) {
          final matchesQuery = query.isEmpty ||
              product.name.toLowerCase().contains(query) ||
              product.sku.toLowerCase().contains(query);
          final matchesFeatured = !_onlyFeatured || product.isFeatured;
          return matchesQuery && matchesFeatured;
        }).toList();

        return RefreshIndicator(
          onRefresh: () async => _refresh(),
          color: AppColors.primary,
          child: Column(
            children: [
              Padding(
                padding: const EdgeInsets.fromLTRB(16, 14, 16, 10),
                child: AppSearchField(
                  controller: _searchController,
                  hintText: 'Buscar en el catalogo',
                  onChanged: _onSearch,
                ),
              ),
              AppFilterChips(
                options: const [
                  (value: 'all', label: 'Todo el catalogo'),
                  (value: 'featured', label: 'Destacados'),
                ],
                selected: _onlyFeatured ? 'featured' : 'all',
                onSelected: (value) => setState(() => _onlyFeatured = value == 'featured'),
              ),
              const SizedBox(height: 4),
              Expanded(
                child: rows.isEmpty
                    ? const AppEmptyState(
                        icon: Icons.search_off_rounded,
                        title: 'Sin resultados',
                        message: 'Ningun producto coincide con el filtro.',
                      )
                    : GridView.builder(
                        physics: const AlwaysScrollableScrollPhysics(),
                        padding: AppSpacing.page,
                        gridDelegate: const SliverGridDelegateWithFixedCrossAxisCount(
                          crossAxisCount: 2,
                          mainAxisSpacing: 11,
                          crossAxisSpacing: 11,
                          mainAxisExtent: 196,
                        ),
                        itemCount: rows.length,
                        itemBuilder: (context, index) => _CatalogCard(product: rows[index]),
                      ),
              ),
            ],
          ),
        );
      },
    );
  }
}

class _CatalogCard extends StatelessWidget {
  const _CatalogCard({required this.product});

  final StorefrontProduct product;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final soldOut = product.stockQuantity <= 0;

    return AppCard(
      padding: const EdgeInsets.all(12),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Expanded(
            child: Container(
              width: double.infinity,
              decoration: BoxDecoration(
                color: AppColors.surfaceMuted,
                borderRadius: AppRadius.mdAll,
              ),
              child: Stack(
                children: [
                  const Center(
                    child: Icon(Icons.image_outlined, size: 28, color: AppColors.textDisabled),
                  ),
                  if (product.isFeatured)
                    Positioned(
                      top: 6,
                      left: 6,
                      child: Container(
                        padding: const EdgeInsets.symmetric(horizontal: 7, vertical: 3),
                        decoration: BoxDecoration(
                          color: AppColors.primary,
                          borderRadius: AppRadius.pillAll,
                        ),
                        child: const Text(
                          'Destacado',
                          style: TextStyle(
                            fontFamily: 'Inter',
                            fontSize: 9.5,
                            fontWeight: FontWeight.w700,
                            color: Colors.white,
                          ),
                        ),
                      ),
                    ),
                ],
              ),
            ),
          ),
          const SizedBox(height: 10),
          Text(
            product.name,
            style: theme.textTheme.titleSmall?.copyWith(fontSize: 13),
            maxLines: 2,
            overflow: TextOverflow.ellipsis,
          ),
          const SizedBox(height: 5),
          Row(
            children: [
              Text(
                AppFormat.money(product.price),
                style: theme.textTheme.titleSmall?.copyWith(color: AppColors.primary),
              ),
              const Spacer(),
              if (soldOut)
                const AppStatusChip(dense: true, label: 'Agotado', tone: AppStatusTone.error)
              else
                Text('${product.stockQuantity} uds', style: theme.textTheme.labelSmall),
            ],
          ),
        ],
      ),
    );
  }
}
