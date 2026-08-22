import 'dart:async';
import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../../../../../shared/widgets/ui/ui.dart';
import '../../../../../shared/filters/filter_models.dart';
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
  String _field = 'name';
  FilterSelection _selection = const FilterSelection();

  static const List<SearchField> _searchFields = [
    SearchField(key: 'name', label: 'Nombre', hint: 'Nombre del producto'),
    SearchField(key: 'sku', label: 'SKU', hint: 'SKU exacto'),
  ];

  static const List<FilterDimension> _dimensions = [
    FilterDimension(
      key: 'channel',
      label: 'Canal',
      icon: Icons.hub_outlined,
      options: [
        FilterOption(value: 'shopify', label: 'Shopify'),
        FilterOption(value: 'woocommerce', label: 'WooCommerce'),
        FilterOption(value: 'mercadolibre', label: 'MercadoLibre'),
        FilterOption(value: 'amazon', label: 'Amazon'),
        FilterOption(value: 'siigo', label: 'Siigo'),
      ],
    ),
  ];

  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addPostFrameCallback((_) => _refresh());
  }

  @override
  void didUpdateWidget(ProductListScreen oldWidget) {
    super.didUpdateWidget(oldWidget);
    if (oldWidget.businessId != widget.businessId) {
      setState(() => _selection = const FilterSelection());
      _searchController.clear();
      _refresh();
    }
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

  void _applySelection(FilterSelection selection) {
    setState(() => _selection = selection);
    context.read<ProductProvider>().applyFilters(channel: selection['channel']);
    _refresh();
  }

  void _onField(String field) {
    setState(() => _field = field);
    context.read<ProductProvider>().setSearch(field: field);
    if (_searchController.text.trim().isNotEmpty) _refresh();
  }

  void _onSearch(String value) {
    _debounce?.cancel();
    _debounce = Timer(const Duration(milliseconds: 400), () {
      if (!mounted) return;
      context.read<ProductProvider>().setSearch(field: _field, term: value);
      _refresh();
    });
  }

  @override
  Widget build(BuildContext context) {
    return Consumer<ProductProvider>(
      builder: (context, provider, _) {
        final filtering =
            _selection.isNotEmpty || _searchController.text.trim().isNotEmpty;

        return Column(
          children: [
            AppFilterBar(
              controller: _searchController,
              searchFields: _searchFields,
              selectedField: _field,
              onFieldChanged: _onField,
              onSearchChanged: _onSearch,
              dimensions: _dimensions,
              selection: _selection,
              onSelectionChanged: _applySelection,
              summary: filtering && !provider.list.isLoading
                  ? '${provider.list.total} de ${provider.unfilteredTotal} productos'
                  : null,
            ),
            Expanded(
              child: PaginatedListView<Product>(
                controller: provider.list,
                unitLabel: 'productos',
                placeholderHeight: 96,
                emptyIcon: Icons.sell_outlined,
                emptyTitle: 'Sin productos',
                emptyMessage: filtering
                    ? 'Ningun producto coincide con los filtros aplicados.'
                    : 'Cuando sincronices tu catalogo lo vas a ver aqui.',
                itemBuilder: (context, product, index) => ProductCard(
                  product: product,
                  onTap: () => Navigator.of(context).push(
                    MaterialPageRoute(
                      builder: (_) => ProductDetailScreen(productId: product.id),
                    ),
                  ),
                ),
              ),
            ),
          ],
        );
      },
    );
  }
}
