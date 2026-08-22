import 'dart:async';
import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../../../../../shared/widgets/ui/ui.dart';
import '../../domain/entities.dart';
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
    _searchController.dispose();
    super.dispose();
  }

  void _refresh() {
    final provider = context.read<ShipmentProvider>();
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
              return PaginatedListView<Shipment>(
                controller: provider.list,
                unitLabel: 'guias',
                placeholderHeight: 150,
                emptyIcon: Icons.local_shipping_outlined,
                emptyTitle: 'Sin guias',
                emptyMessage:
                    'Cuando generes guias para tus ordenes las vas a ver aqui.',
                itemBuilder: (context, shipment, index) => ShipmentCard(
                  shipment: shipment,
                  onTap: () => Navigator.of(context).push(
                    MaterialPageRoute(
                      builder: (_) => ShipmentDetailScreen(shipmentId: shipment.id),
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
