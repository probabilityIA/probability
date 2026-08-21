import 'dart:async';
import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../../../../../shared/theme/app_colors.dart';
import '../../../../../shared/theme/app_tokens.dart';
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
  final _scrollController = ScrollController();
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
    _scrollController.addListener(_onScroll);
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
    _scrollController.removeListener(_onScroll);
    _scrollController.dispose();
    _searchController.dispose();
    super.dispose();
  }

  void _onScroll() {
    if (!_scrollController.hasClients) return;
    final position = _scrollController.position;
    if (position.pixels >= position.maxScrollExtent - 320) {
      context.read<OrderProvider>().loadMore(businessId: widget.businessId);
    }
  }

  void _refresh() {
    final provider = context.read<OrderProvider>();
    provider.setPage(1);
    provider.fetchOrders(businessId: widget.businessId);
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
              if (provider.isLoading && provider.orders.isEmpty) {
                return const AppListSkeleton();
              }

              if (provider.error != null && provider.orders.isEmpty) {
                return AppErrorState(message: provider.error!, onRetry: _refresh);
              }

              if (provider.orders.isEmpty) {
                return AppEmptyState(
                  icon: Icons.receipt_long_outlined,
                  title: 'Sin ordenes',
                  message: _searchController.text.isNotEmpty || _status.isNotEmpty
                      ? 'Ninguna orden coincide con el filtro aplicado.'
                      : 'Cuando entren pedidos desde tus canales los vas a ver aqui.',
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
                  itemCount: provider.orders.length + 1,
                  separatorBuilder: (context, index) => const SizedBox(height: 10),
                  itemBuilder: (context, index) {
                    if (index == provider.orders.length) {
                      return _ListFooter(
                        loading: provider.isLoadingMore,
                        hasMore: provider.hasMore,
                        total: provider.pagination?.total ?? provider.orders.length,
                        shown: provider.orders.length,
                      );
                    }
                    final order = provider.orders[index];
                    return OrderCard(order: order, onTap: () => _openDetail(order));
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

class _ListFooter extends StatelessWidget {
  const _ListFooter({
    required this.loading,
    required this.hasMore,
    required this.total,
    required this.shown,
  });

  final bool loading;
  final bool hasMore;
  final int total;
  final int shown;

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 18),
      child: Center(
        child: loading
            ? const SizedBox(
                height: 22,
                width: 22,
                child: CircularProgressIndicator(strokeWidth: 2.2),
              )
            : Text(
                hasMore ? 'Desliza para cargar mas' : '$shown de $total ordenes',
                style: Theme.of(context).textTheme.labelSmall,
              ),
      ),
    );
  }
}
