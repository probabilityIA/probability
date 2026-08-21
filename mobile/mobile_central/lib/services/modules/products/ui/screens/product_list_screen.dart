import 'dart:async';
import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../../../../../shared/theme/app_colors.dart';
import '../../../../../shared/theme/app_tokens.dart';
import '../../../../../shared/widgets/ui/ui.dart';
import '../providers/product_provider.dart';
import '../widgets/product_card.dart';
import 'product_detail_screen.dart';

class ProductListScreen extends StatefulWidget {
  const ProductListScreen({super.key, this.businessId});

  final int? businessId;

  @override
  State<ProductListScreen> createState() => _ProductListScreenState();
}

class _ProductListScreenState extends State<ProductListScreen> {
  final _searchController = TextEditingController();
  final _scrollController = ScrollController();
  Timer? _debounce;
  String _filter = '';

  static const List<({String value, String label})> _filters = [
    (value: '', label: 'Todos'),
    (value: 'in_stock', label: 'Con stock'),
    (value: 'low_stock', label: 'Stock bajo'),
    (value: 'out_of_stock', label: 'Agotados'),
    (value: 'inactive', label: 'Inactivos'),
  ];

  @override
  void initState() {
    super.initState();
    _scrollController.addListener(_onScroll);
    WidgetsBinding.instance.addPostFrameCallback((_) => _refresh());
  }

  @override
  void didUpdateWidget(ProductListScreen oldWidget) {
    super.didUpdateWidget(oldWidget);
    if (oldWidget.businessId != widget.businessId) _refresh();
  }

  @override
  void dispose() {
    _debounce?.cancel();
    _scrollController.removeListener(_onScroll);
    _scrollController.dispose();
    _searchController.dispose();
    super.dispose();
  }

  void _onScroll() {
    if (!_scrollController.hasClients) return;
    final position = _scrollController.position;
    if (position.pixels >= position.maxScrollExtent - 320) {
      context.read<ProductProvider>().loadMore(businessId: widget.businessId);
    }
  }

  void _refresh() {
    final provider = context.read<ProductProvider>();
    provider.setPage(1);
    provider.fetchProducts(businessId: widget.businessId);
  }

  void _onSearch(String value) {
    _debounce?.cancel();
    _debounce = Timer(const Duration(milliseconds: 400), () {
      if (!mounted) return;
      context.read<ProductProvider>().setFilters(name: value);
      _refresh();
    });
  }

  @override
  Widget build(BuildContext context) {
    return Column(
      children: [
        Padding(
          padding: const EdgeInsets.fromLTRB(16, 14, 16, 10),
          child: AppSearchField(
            controller: _searchController,
            hintText: 'Nombre o SKU',
            onChanged: _onSearch,
          ),
        ),
        AppFilterChips(
          options: _filters,
          selected: _filter,
          onSelected: (value) => setState(() => _filter = value),
        ),
        const SizedBox(height: 4),
        Expanded(
          child: Consumer<ProductProvider>(
            builder: (context, provider, _) {
              if (provider.isLoading && provider.products.isEmpty) {
                return const AppListSkeleton();
              }
              if (provider.error != null && provider.products.isEmpty) {
                return AppErrorState(message: provider.error!, onRetry: _refresh);
              }

              final rows = provider.products.where((product) {
                switch (_filter) {
                  case 'in_stock':
                    return product.stock >= 20;
                  case 'low_stock':
                    return product.stock > 0 && product.stock < 20;
                  case 'out_of_stock':
                    return product.stock <= 0;
                  case 'inactive':
                    return !product.isActive;
                  default:
                    return true;
                }
              }).toList();

              if (rows.isEmpty) {
                return AppEmptyState(
                  icon: Icons.sell_outlined,
                  title: 'Sin productos',
                  message: 'Ninguno coincide con el filtro aplicado.',
                  actionLabel: 'Actualizar',
                  onAction: _refresh,
                );
              }

              return RefreshIndicator(
                onRefresh: () async => _refresh(),
                color: AppColors.primary,
                child: ListView.separated(
                  controller: _scrollController,
                  physics: const AlwaysScrollableScrollPhysics(),
                  padding: AppSpacing.page,
                  itemCount: rows.length + 1,
                  separatorBuilder: (context, index) => const SizedBox(height: 10),
                  itemBuilder: (context, index) {
                    if (index == rows.length) {
                      return Padding(
                        padding: const EdgeInsets.symmetric(vertical: 18),
                        child: Center(
                          child: provider.isLoadingMore
                              ? const SizedBox(
                                  height: 22,
                                  width: 22,
                                  child: CircularProgressIndicator(strokeWidth: 2.2),
                                )
                              : Text(
                                  provider.hasMore
                                      ? 'Desliza para cargar mas'
                                      : '${rows.length} de ${provider.pagination?.total ?? rows.length} productos',
                                  style: Theme.of(context).textTheme.labelSmall,
                                ),
                        ),
                      );
                    }
                    final product = rows[index];
                    return ProductCard(
                      product: product,
                      onTap: () => Navigator.of(context).push(
                        MaterialPageRoute(
                          builder: (_) => ProductDetailScreen(productId: product.id),
                        ),
                      ),
                    );
                  },
                ),
              );
            },
          ),
        ),
      ],
    );
  }
}
