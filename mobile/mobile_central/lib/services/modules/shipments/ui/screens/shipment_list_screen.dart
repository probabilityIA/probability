import 'dart:async';
import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../../../../../shared/theme/app_colors.dart';
import '../../../../../shared/theme/app_tokens.dart';
import '../../../../../shared/widgets/ui/ui.dart';
import '../providers/shipment_provider.dart';
import '../widgets/shipment_card.dart';
import 'shipment_detail_screen.dart';

class ShipmentListScreen extends StatefulWidget {
  const ShipmentListScreen({super.key, this.businessId});

  final int? businessId;

  @override
  State<ShipmentListScreen> createState() => _ShipmentListScreenState();
}

class _ShipmentListScreenState extends State<ShipmentListScreen> {
  final _searchController = TextEditingController();
  final _scrollController = ScrollController();
  Timer? _debounce;
  String _status = '';

  static const List<({String value, String label})> _statusOptions = [
    (value: '', label: 'Todas'),
    (value: 'created', label: 'Generadas'),
    (value: 'in_transit', label: 'En transito'),
    (value: 'delivered', label: 'Entregadas'),
    (value: 'returned', label: 'Devueltas'),
    (value: 'cancelled', label: 'Canceladas'),
  ];

  @override
  void initState() {
    super.initState();
    _scrollController.addListener(_onScroll);
    WidgetsBinding.instance.addPostFrameCallback((_) => _refresh());
  }

  @override
  void didUpdateWidget(ShipmentListScreen oldWidget) {
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
      context.read<ShipmentProvider>().loadMoreShipments(businessId: widget.businessId);
    }
  }

  void _refresh() {
    final provider = context.read<ShipmentProvider>();
    provider.setPage(1);
    provider.fetchShipments(businessId: widget.businessId);
  }

  void _onSearch(String value) {
    _debounce?.cancel();
    _debounce = Timer(const Duration(milliseconds: 400), () {
      if (!mounted) return;
      context.read<ShipmentProvider>().setShipmentFilters(search: value);
      _refresh();
    });
  }

  void _onStatus(String value) {
    setState(() => _status = value);
    context.read<ShipmentProvider>().setShipmentFilters(status: value);
    _refresh();
  }

  @override
  Widget build(BuildContext context) {
    return Column(
      children: [
        Padding(
          padding: const EdgeInsets.fromLTRB(16, 14, 16, 10),
          child: AppSearchField(
            controller: _searchController,
            hintText: 'Numero de guia',
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
          child: Consumer<ShipmentProvider>(
            builder: (context, provider, _) {
              if (provider.isLoading && provider.shipments.isEmpty) {
                return const AppListSkeleton();
              }
              if (provider.error != null && provider.shipments.isEmpty) {
                return AppErrorState(message: provider.error!, onRetry: _refresh);
              }
              if (provider.shipments.isEmpty) {
                return AppEmptyState(
                  icon: Icons.local_shipping_outlined,
                  title: 'Sin guias',
                  message: 'Cuando generes guias para tus ordenes las vas a ver aqui.',
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
                  itemCount: provider.shipments.length + 1,
                  separatorBuilder: (context, index) => const SizedBox(height: 10),
                  itemBuilder: (context, index) {
                    if (index == provider.shipments.length) {
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
                                      : '${provider.shipments.length} de ${provider.pagination?.total ?? provider.shipments.length} guias',
                                  style: Theme.of(context).textTheme.labelSmall,
                                ),
                        ),
                      );
                    }
                    final shipment = provider.shipments[index];
                    return ShipmentCard(
                      shipment: shipment,
                      onTap: () => Navigator.of(context).push(
                        MaterialPageRoute(
                          builder: (_) => ShipmentDetailScreen(shipmentId: shipment.id),
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
