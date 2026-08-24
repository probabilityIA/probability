import 'dart:async';
import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../../../../../shared/filters/filter_models.dart';
import '../../../../../shared/widgets/ui/ui.dart';
import '../../domain/entities.dart';
import '../providers/order_provider.dart';
import '../widgets/order_card.dart';
import 'order_detail_screen.dart';

class OrderListScreen extends StatefulWidget {
  const OrderListScreen({super.key, this.businessId});

  final int? businessId;

  @override
  State<OrderListScreen> createState() => _OrderListScreenState();
}

class _OrderListScreenState extends State<OrderListScreen> {
  final _searchController = TextEditingController();
  Timer? _debounce;
  String _field = 'order_number';
  FilterSelection _selection = const FilterSelection();

  static const List<SearchField> _searchFields = [
    SearchField(
      key: 'order_number',
      label: 'N\u00ba orden',
      hint: 'Numero de la orden',
    ),
    SearchField(
      key: 'internal_number',
      label: 'N\u00ba interno',
      hint: 'Numero interno',
    ),
    SearchField(
      key: 'customer_email',
      label: 'Correo',
      hint: 'Correo del cliente',
    ),
    SearchField(
      key: 'customer_phone',
      label: 'Telefono',
      hint: 'Telefono del cliente',
    ),
  ];

  static const List<FilterOption> _channelOptions = [
    FilterOption(value: 'shopify', label: 'Shopify'),
    FilterOption(value: 'woocommerce', label: 'WooCommerce'),
    FilterOption(value: 'mercadolibre', label: 'MercadoLibre'),
    FilterOption(value: 'amazon', label: 'Amazon'),
    FilterOption(value: 'whatsapp', label: 'WhatsApp'),
    FilterOption(value: 'manual', label: 'Manual'),
  ];

  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addPostFrameCallback((_) {
      _refresh();
      context.read<OrderProvider>().loadStatusOptions(
            businessId: widget.businessId,
          );
    });
  }

  @override
  void didUpdateWidget(OrderListScreen oldWidget) {
    super.didUpdateWidget(oldWidget);
    if (oldWidget.businessId != widget.businessId) {
      context.read<OrderProvider>().resetFilters();
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
    context.read<OrderProvider>().fetchOrders(businessId: widget.businessId);
  }

  List<FilterDimension> _dimensions(OrderProvider provider) {
    return [
      FilterDimension(
        key: 'status',
        label: 'Estado',
        icon: Icons.flag_outlined,
        options: [
          for (final status in provider.statusOptions)
            FilterOption(value: status.code, label: status.name),
        ],
      ),
      const FilterDimension(
        key: 'platform',
        label: 'Canal',
        icon: Icons.hub_outlined,
        options: _channelOptions,
      ),
      const FilterDimension(
        key: 'paid',
        label: 'Pago',
        icon: Icons.payments_outlined,
        options: [
          FilterOption(value: 'true', label: 'Pagada'),
          FilterOption(value: 'false', label: 'Sin pago'),
        ],
      ),
      const FilterDimension(
        key: 'cod',
        label: 'Contra entrega',
        icon: Icons.local_atm_outlined,
        options: [
          FilterOption(value: 'true', label: 'Si'),
          FilterOption(value: 'false', label: 'No'),
        ],
      ),
    ];
  }

  void _applySelection(FilterSelection selection) {
    setState(() => _selection = selection);
    context.read<OrderProvider>().applyFilters(
          status: selection['status'],
          platform: selection['platform'],
          isPaid: selection.boolFor('paid'),
          isCod: selection.boolFor('cod'),
        );
    _refresh();
  }

  void _onField(String field) {
    setState(() => _field = field);
    context.read<OrderProvider>().setSearch(field: field);
    if (_searchController.text.trim().isNotEmpty) _refresh();
  }

  void _onSearch(String value) {
    _debounce?.cancel();
    _debounce = Timer(const Duration(milliseconds: 400), () {
      if (!mounted) return;
      context.read<OrderProvider>().setSearch(field: _field, term: value);
      _refresh();
    });
  }

  String? _summary(OrderProvider provider) {
    if (_selection.isEmpty && _searchController.text.trim().isEmpty) return null;
    if (provider.list.isLoading) return null;
    return '${provider.list.total} de ${provider.unfilteredTotal} ordenes';
  }

  void _openDetail(Order order) {
    Navigator.of(context).push(
      MaterialPageRoute(builder: (_) => OrderDetailScreen(orderId: order.id)),
    );
  }

  @override
  Widget build(BuildContext context) {
    return Consumer<OrderProvider>(
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
              dimensions: _dimensions(provider),
              selection: _selection,
              onSelectionChanged: _applySelection,
              summary: _summary(provider),
            ),
            Expanded(
              child: PaginatedListView<Order>(
                controller: provider.list,
                unitLabel: 'ordenes',
                placeholderHeight: 158,
                emptyIcon: Icons.receipt_long_outlined,
                emptyTitle: 'Sin ordenes',
                emptyMessage: filtering
                    ? 'Ninguna orden coincide con los filtros aplicados.'
                    : 'Cuando entren pedidos desde tus canales los vas a ver aqui.',
                itemBuilder: (context, order, index) =>
                    OrderCard(order: order, onTap: () => _openDetail(order)),
              ),
            ),
          ],
        );
      },
    );
  }
}

