import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:provider/provider.dart';
import '../../../../../shared/theme/app_colors.dart';
import '../../../../../shared/theme/app_tokens.dart';
import '../../../../../shared/widgets/ui/ui.dart';
import '../../domain/entities.dart';
import '../providers/order_provider.dart';
import '../../../mobile/ui/providers/order_full_provider.dart';
import '../widgets/order_detail_sections.dart';
import '../widgets/order_related_cards.dart';

class OrderDetailScreen extends StatefulWidget {
  const OrderDetailScreen({super.key, required this.orderId});

  final String orderId;

  @override
  State<OrderDetailScreen> createState() => _OrderDetailScreenState();
}

class _OrderDetailScreenState extends State<OrderDetailScreen> {
  Order? _order;
  bool _loading = true;
  String? _error;

  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addPostFrameCallback((_) => _load());
  }

  Future<void> _load() async {
    setState(() {
      _loading = true;
      _error = null;
    });

    final provider = context.read<OrderProvider>();
    final fullProvider = context.read<OrderFullProvider>();

    final results = await Future.wait([
      provider.getOrderById(widget.orderId),
      fullProvider.load(widget.orderId),
    ]);

    if (!mounted) return;
    final order = results.first as Order?;
    setState(() {
      _order = order;
      _error = order == null ? (provider.error ?? 'No se encontro la orden') : null;
      _loading = false;
    });
  }

  void _copyNumber() {
    final number = _order?.orderNumber;
    if (number == null) return;
    Clipboard.setData(ClipboardData(text: number));
    ScaffoldMessenger.of(context).showSnackBar(
      SnackBar(content: Text('Copiado: $number')),
    );
  }

  Future<void> _confirmCancel() async {
    final confirmed = await showDialog<bool>(
      context: context,
      builder: (dialogContext) => AlertDialog(
        title: const Text('Cancelar orden'),
        content: Text(
          'Vas a cancelar la orden ${_order?.orderNumber ?? ''}. Esta accion no se puede deshacer.',
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(dialogContext, false),
            child: const Text('Volver'),
          ),
          FilledButton(
            style: FilledButton.styleFrom(backgroundColor: AppColors.error),
            onPressed: () => Navigator.pop(dialogContext, true),
            child: const Text('Cancelar orden'),
          ),
        ],
      ),
    );

    if (confirmed != true || !mounted) return;

    final provider = context.read<OrderProvider>();
    final ok = await provider.updateOrder(
      widget.orderId,
      UpdateOrderDTO(status: 'cancelled'),
    );

    if (!mounted) return;
    ScaffoldMessenger.of(context).showSnackBar(
      SnackBar(
        content: Text(ok ? 'Orden cancelada' : (provider.error ?? 'No se pudo cancelar')),
      ),
    );
    if (ok) _load();
  }

  @override
  Widget build(BuildContext context) {
    final order = _order;

    return AppScaffold(
      title: order?.orderNumber ?? 'Orden',
      subtitle: order == null ? null : 'Detalle del pedido',
      onBack: () => Navigator.of(context).pop(),
      actions: [
        if (order != null)
          IconButton(
            icon: const Icon(Icons.copy_rounded, size: 20),
            tooltip: 'Copiar numero',
            onPressed: _copyNumber,
          ),
        if (order != null)
          PopupMenuButton<String>(
            icon: const Icon(Icons.more_vert_rounded, size: 20),
            onSelected: (value) {
              if (value == 'cancel') _confirmCancel();
              if (value == 'refresh') _load();
            },
            itemBuilder: (context) => [
              const PopupMenuItem(
                value: 'refresh',
                child: Row(
                  children: [
                    Icon(Icons.refresh, size: 18),
                    SizedBox(width: 10),
                    Text('Actualizar'),
                  ],
                ),
              ),
              if (order.status != 'cancelled')
                const PopupMenuItem(
                  value: 'cancel',
                  child: Row(
                    children: [
                      Icon(Icons.cancel_outlined, size: 18, color: AppColors.error),
                      SizedBox(width: 10),
                      Text('Cancelar orden', style: TextStyle(color: AppColors.error)),
                    ],
                  ),
                ),
            ],
          ),
      ],
      body: _buildBody(order),
    );
  }

  Widget _buildBody(Order? order) {
    if (_loading) return const AppLoading(label: 'Cargando orden');
    if (_error != null) return AppErrorState(message: _error!, onRetry: _load);
    if (order == null) {
      return const AppEmptyState(
        icon: Icons.receipt_long_outlined,
        title: 'Orden no encontrada',
      );
    }

    return RefreshIndicator(
      onRefresh: _load,
      color: AppColors.primary,
      child: ListView(
        physics: const AlwaysScrollableScrollPhysics(),
        padding: AppSpacing.page,
        children: [
          OrderHeaderCard(order: order),
          const SizedBox(height: 18),
          AppSectionHeader(title: 'Productos', subtitle: '${order.lineItems.length} items'),
          OrderItemsCard(items: order.lineItems),
          const SizedBox(height: 18),
          const AppSectionHeader(title: 'Totales'),
          OrderTotalsCard(order: order),
          const SizedBox(height: 18),
          const AppSectionHeader(title: 'Cliente y entrega'),
          OrderCustomerCard(order: order),
          Consumer<OrderFullProvider>(
            builder: (context, provider, _) {
              final shipment = provider.orderFull?.shipment;
              final invoice = provider.orderFull?.invoice;
              return Column(
                children: [
                  const SizedBox(height: 18),
                  const AppSectionHeader(
                    title: 'Envio',
                    subtitle: 'Guia asociada a la orden',
                  ),
                  if (shipment != null)
                    OrderShipmentCard(shipment: shipment)
                  else
                    OrderShippingCard(order: order),
                  if (invoice != null) ...[
                    const SizedBox(height: 18),
                    const AppSectionHeader(title: 'Facturacion'),
                    OrderInvoiceCard(invoice: invoice),
                  ],
                ],
              );
            },
          ),
          const SizedBox(height: 18),
          const AppSectionHeader(title: 'Trazabilidad'),
          AppCard(
            child: Column(
              children: [
                AppKeyValueRow(label: 'Canal', value: order.integrationName ?? order.integrationType),
                const Divider(height: 18),
                AppKeyValueRow(label: 'Id externo', value: order.externalId.isEmpty ? '-' : order.externalId),
                const Divider(height: 18),
                AppKeyValueRow(label: 'Numero interno', value: order.internalNumber.isEmpty ? '-' : order.internalNumber),
                const Divider(height: 18),
                AppKeyValueRow(label: 'Creada por', value: order.userName.isEmpty ? '-' : order.userName),
              ],
            ),
          ),
          const SizedBox(height: 26),
        ],
      ),
    );
  }
}
