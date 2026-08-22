import 'dart:async';
import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
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
  String _status = '';

  static const List<({String value, String label})> _statusOptions = [
    (value: '', label: 'Todas'),
    (value: 'pending', label: 'Pendientes'),
    (value: 'confirmed', label: 'Confirmadas'),
    (value: 'processing', label: 'En proceso'),
    (value: 'shipped', label: 'Enviadas'),
    (value: 'delivered', label: 'Entregadas'),
    (value: 'cancelled', label: 'Canceladas'),
  ];

  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addPostFrameCallback((_) => _refresh());
  }

  @override
  void didUpdateWidget(OrderListScreen oldWidget) {
    super.didUpdateWidget(oldWidget);
    if (oldWidget.businessId != widget.businessId) {
      context.read<OrderProvider>().resetFilters();
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

  void _onSearch(String value) {
    _debounce?.cancel();
    _debounce = Timer(const Duration(milliseconds: 400), () {
      if (!mounted) return;
      context.read<OrderProvider>().setFilters(orderNumber: value);
      _refresh();
    });
  }

  void _onStatus(String value) {
    setState(() => _status = value);
    context.read<OrderProvider>().setFilters(status: value);
    _refresh();
  }

  void _openDetail(Order order) {
    Navigator.of(context).push(
      MaterialPageRoute(builder: (_) => OrderDetailScreen(orderId: order.id)),
    );
  }

  @override
  Widget build(BuildContext context) {
    return Column(
      children: [
        Padding(
          padding: const EdgeInsets.fromLTRB(16, 14, 16, 10),
          child: AppSearchField(
            controller: _searchController,
            hintText: 'Numero, cliente o guia',
            onChanged: _onSearch,
          ),
        ),
        AppFilterChips(
          options: _statusOptions,
          selected: _status,
          onSelected: _onStatus,
        ),
        const SizedBox(height: 4),
        Expanded(
          child: Consumer<OrderProvider>(
            builder: (context, provider, _) {
              return PaginatedListView<Order>(
                controller: provider.list,
                unitLabel: 'ordenes',
                placeholderHeight: 158,
                emptyIcon: Icons.receipt_long_outlined,
                emptyTitle: 'Sin ordenes',
                emptyMessage:
                    _searchController.text.isNotEmpty || _status.isNotEmpty
                        ? 'Ninguna orden coincide con el filtro aplicado.'
                        : 'Cuando entren pedidos desde tus canales los vas a ver aqui.',
                itemBuilder: (context, order, index) =>
                    OrderCard(order: order, onTap: () => _openDetail(order)),
              );
            },
          ),
        ),
      ],
    );
  }
}
