import 'dart:async';
import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../../../../../shared/widgets/ui/ui.dart';
import '../../domain/entities.dart';
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
  Timer? _debounce;

  @override
  void initState() {
    super.initState();
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
    _searchController.dispose();
    super.dispose();
  }

  void _refresh() {
    context.read<ProductProvider>().fetchProducts(businessId: widget.businessId);
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
        const SizedBox(height: 4),
        Expanded(
          child: Consumer<ProductProvider>(
            builder: (context, provider, _) {
              return PaginatedListView<Product>(
                controller: provider.list,
                unitLabel: 'productos',
                placeholderHeight: 96,
                emptyIcon: Icons.sell_outlined,
                emptyTitle: 'Sin productos',
                emptyMessage: _searchController.text.isEmpty
                    ? 'Cuando sincronices tu catalogo lo vas a ver aqui.'
                    : 'Ningun producto coincide con la busqueda.',
                itemBuilder: (context, product, index) => ProductCard(
                  product: product,
                  onTap: () => Navigator.of(context).push(
                    MaterialPageRoute(
                      builder: (_) => ProductDetailScreen(productId: product.id),
                    ),
                  ),
                ),
              );
            },
          ),
        ),
      ],
    );
  }
}
