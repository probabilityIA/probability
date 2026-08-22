import 'dart:async';
import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../../../../../shared/filters/filter_models.dart';
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
  String _field = 'tracking_number';
  FilterSelection _selection = const FilterSelection();

  static const List<SearchField> _searchFields = [
    SearchField(
      key: 'tracking_number',
      label: 'Guia',
      hint: 'Numero de guia',
    ),
    SearchField(
      key: 'customer_name',
      label: 'Cliente',
      hint: 'Nombre del cliente',
    ),
    SearchField(
      key: 'order_id',
      label: 'Orden',
      hint: 'Id de la orden',
    ),
  ];

  static const List<FilterDimension> _dimensions = [
    FilterDimension(
      key: 'status',
      label: 'Estado',
      icon: Icons.flag_outlined,
      options: [
        FilterOption(value: 'pending', label: 'Pendiente'),
        FilterOption(value: 'created', label: 'Generada'),
        FilterOption(value: 'in_transit', label: 'En transito'),
        FilterOption(value: 'delivered', label: 'Entregada'),
        FilterOption(value: 'returned', label: 'Devuelta'),
        FilterOption(value: 'cancelled', label: 'Cancelada'),
      ],
    ),
    FilterDimension(
      key: 'carrier',
      label: 'Transportadora',
      icon: Icons.local_shipping_outlined,
      options: [
        FilterOption(value: 'Coordinadora', label: 'Coordinadora'),
        FilterOption(value: 'Interrapidisimo', label: 'Interrapidisimo'),
        FilterOption(value: 'Servientrega', label: 'Servientrega'),
        FilterOption(value: 'Envia', label: 'Envia'),
        FilterOption(value: 'TCC', label: 'TCC'),
      ],
    ),
  ];

  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addPostFrameCallback((_) => _refresh());
  }

  @override
  void didUpdateWidget(ShipmentListScreen oldWidget) {
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
    context.read<ShipmentProvider>().fetchShipments(
          businessId: widget.businessId,
        );
  }

  void _applySelection(FilterSelection selection) {
    setState(() => _selection = selection);
    context.read<ShipmentProvider>().applyFilters(
          status: selection['status'],
          carrier: selection['carrier'],
        );
    _refresh();
  }

  void _onField(String field) {
    setState(() => _field = field);
    context.read<ShipmentProvider>().setSearch(field: field);
    if (_searchController.text.trim().isNotEmpty) _refresh();
  }

  void _onSearch(String value) {
    _debounce?.cancel();
    _debounce = Timer(const Duration(milliseconds: 400), () {
      if (!mounted) return;
      context.read<ShipmentProvider>().setSearch(field: _field, term: value);
      _refresh();
    });
  }

  @override
  Widget build(BuildContext context) {
    return Consumer<ShipmentProvider>(
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
                  ? '${provider.list.total} de ${provider.unfilteredTotal} guias'
                  : null,
            ),
            Expanded(
              child: PaginatedListView<Shipment>(
                controller: provider.list,
                unitLabel: 'guias',
                placeholderHeight: 150,
                emptyIcon: Icons.local_shipping_outlined,
                emptyTitle: 'Sin guias',
                emptyMessage: filtering
                    ? 'Ninguna guia coincide con los filtros aplicados.'
                    : 'Cuando generes guias para tus ordenes las vas a ver aqui.',
                itemBuilder: (context, shipment, index) => ShipmentCard(
                  shipment: shipment,
                  onTap: () => Navigator.of(context).push(
                    MaterialPageRoute(
                      builder: (_) =>
                          ShipmentDetailScreen(shipmentId: shipment.id),
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

