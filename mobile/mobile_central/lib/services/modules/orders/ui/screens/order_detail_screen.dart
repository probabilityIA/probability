import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../providers/order_provider.dart';
import '../../domain/entities.dart';

class OrderDetailScreen extends StatefulWidget {
  final String orderId;

  const OrderDetailScreen({super.key, required this.orderId});

  @override
  State<OrderDetailScreen> createState() => _OrderDetailScreenState();
}

class _OrderDetailScreenState extends State<OrderDetailScreen> {
  Order? _order;
  bool _isLoading = true;
  String? _error;

  @override
  void initState() {
    super.initState();
    _loadOrder();
  }

  Future<void> _loadOrder() async {
    setState(() {
      _isLoading = true;
      _error = null;
    });

    final provider = context.read<OrderProvider>();
    final order = await provider.getOrderById(widget.orderId);

    if (!mounted) return;

    setState(() {
      _order = order;
      _isLoading = false;
      if (order == null) {
        _error = provider.error ?? 'No se pudo cargar la orden';
      }
    });
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: Text(
          _order != null
              ? '#${_order!.orderNumber.isNotEmpty ? _order!.orderNumber : _order!.id}'
              : 'Detalle de Orden',
        ),
      ),
      body: _buildBody(),
    );
  }

  Widget _buildBody() {
    final colorScheme = Theme.of(context).colorScheme;

    if (_isLoading) {
      return const Center(child: CircularProgressIndicator());
    }

    if (_error != null) {
      return Center(
        child: Padding(
          padding: const EdgeInsets.all(24),
          child: Column(
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              Icon(Icons.error_outline, size: 48, color: colorScheme.error),
              const SizedBox(height: 16),
              Text(
                _error!,
                textAlign: TextAlign.center,
                style: TextStyle(color: colorScheme.error),
              ),
              const SizedBox(height: 16),
              FilledButton.icon(
                onPressed: _loadOrder,
                icon: const Icon(Icons.refresh),
                label: const Text('Reintentar'),
              ),
            ],
          ),
        ),
      );
    }

    final order = _order!;

    return RefreshIndicator(
      onRefresh: _loadOrder,
      child: SingleChildScrollView(
        physics: const AlwaysScrollableScrollPhysics(),
        padding: const EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            _buildStatusRow(order),
            const SizedBox(height: 16),
            _buildGeneralInfoCard(order),
            const SizedBox(height: 12),
            _buildCustomerCard(order),
            const SizedBox(height: 12),
            _buildAddressCard(order),
            const SizedBox(height: 12),
            _buildFinancialCard(order),
            const SizedBox(height: 12),
            _buildItemsCard(order),
            const SizedBox(height: 12),
            _buildDatesCard(order),
            const SizedBox(height: 24),
          ],
        ),
      ),
    );
  }

  Widget _buildStatusRow(Order order) {
    return Row(
      children: [
        _buildStatusChip(
          label: order.orderStatus?.name ?? order.status,
          color: _parseHexColor(order.orderStatus?.color),
          icon: Icons.local_shipping_outlined,
          title: 'Orden',
        ),
        const SizedBox(width: 8),
        _buildStatusChip(
          label: order.paymentStatus?.name ??
              (order.isPaid ? 'Pagado' : 'No pagado'),
          color: _parseHexColor(order.paymentStatus?.color) ??
              (order.isPaid
                  ? const Color(0xFF16A34A)
                  : const Color(0xFFDC2626)),
          icon: Icons.payment_outlined,
          title: 'Pago',
        ),
        const SizedBox(width: 8),
        if (order.fulfillmentStatus != null)
          _buildStatusChip(
            label: order.fulfillmentStatus!.name,
            color: _parseHexColor(order.fulfillmentStatus!.color),
            icon: Icons.inventory_2_outlined,
            title: 'Fulfillment',
          ),
      ],
    );
  }

  Widget _buildStatusChip({
    required String label,
    Color? color,
    required IconData icon,
    required String title,
  }) {
    final bgColor = color ?? Theme.of(context).colorScheme.surfaceContainerHighest;
    final textColor = _textColorForBg(bgColor);

    return Expanded(
      child: Container(
        padding: const EdgeInsets.symmetric(vertical: 10, horizontal: 8),
        decoration: BoxDecoration(
          color: bgColor,
          borderRadius: BorderRadius.circular(12),
        ),
        child: Column(
          children: [
            Icon(icon, size: 18, color: textColor),
            const SizedBox(height: 4),
            Text(
              title,
              style: TextStyle(
                fontSize: 9,
                fontWeight: FontWeight.w600,
                color: textColor.withAlpha(179),
              ),
            ),
            const SizedBox(height: 2),
            Text(
              label,
              textAlign: TextAlign.center,
              style: TextStyle(
                fontSize: 11,
                fontWeight: FontWeight.bold,
                color: textColor,
              ),
              maxLines: 2,
              overflow: TextOverflow.ellipsis,
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildGeneralInfoCard(Order order) {
    final colorScheme = Theme.of(context).colorScheme;

    return Card(
      child: Padding(
        padding: const EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Row(
              children: [
                Icon(Icons.info_outline, size: 18, color: colorScheme.primary),
                const SizedBox(width: 8),
                Text(
                  'Informacion General',
                  style: Theme.of(context).textTheme.titleSmall?.copyWith(
                        fontWeight: FontWeight.bold,
                      ),
                ),
              ],
            ),
            const Divider(height: 20),
            _buildInfoRow('Nro. Orden', order.orderNumber),
            if (order.internalNumber.isNotEmpty)
              _buildInfoRow('Numero Interno', order.internalNumber),
            if (order.externalId.isNotEmpty)
              _buildInfoRow('ID Externo', order.externalId),
            _buildInfoRow('Plataforma', order.platform),
            if (order.integrationName != null &&
                order.integrationName!.isNotEmpty)
              _buildInfoRow('Integracion', order.integrationName!),
            if (order.orderTypeName.isNotEmpty)
              _buildInfoRow('Tipo', order.orderTypeName),
            if (order.warehouseName.isNotEmpty)
              _buildInfoRow('Bodega', order.warehouseName),
            if (order.driverName.isNotEmpty)
              _buildInfoRow('Conductor', order.driverName),
            if (order.notes != null && order.notes!.isNotEmpty)
              _buildInfoRow('Notas', order.notes!),
          ],
        ),
      ),
    );
  }

  Widget _buildCustomerCard(Order order) {
    final colorScheme = Theme.of(context).colorScheme;

    return Card(
      child: Padding(
        padding: const EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Row(
              children: [
                Icon(Icons.person_outline,
                    size: 18, color: colorScheme.primary),
                const SizedBox(width: 8),
                Text(
                  'Cliente',
                  style: Theme.of(context).textTheme.titleSmall?.copyWith(
                        fontWeight: FontWeight.bold,
                      ),
                ),
              ],
            ),
            const Divider(height: 20),
            _buildInfoRow('Nombre', order.customerName),
            if (order.customerEmail.isNotEmpty)
              _buildInfoRow('Email', order.customerEmail),
            if (order.customerPhone.isNotEmpty)
              _buildInfoRow('Telefono', order.customerPhone),
            if (order.customerDni.isNotEmpty)
              _buildInfoRow('DNI', order.customerDni),
          ],
        ),
      ),
    );
  }

  Widget _buildAddressCard(Order order) {
    final colorScheme = Theme.of(context).colorScheme;
    final hasAddress = order.shippingStreet.isNotEmpty ||
        order.shippingCity.isNotEmpty ||
        order.shippingState.isNotEmpty;

    if (!hasAddress) return const SizedBox.shrink();

    return Card(
      child: Padding(
        padding: const EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Row(
              children: [
                Icon(Icons.location_on_outlined,
                    size: 18, color: colorScheme.primary),
                const SizedBox(width: 8),
                Text(
                  'Direccion de Envio',
                  style: Theme.of(context).textTheme.titleSmall?.copyWith(
                        fontWeight: FontWeight.bold,
                      ),
                ),
              ],
            ),
            const Divider(height: 20),
            if (order.shippingStreet.isNotEmpty)
              _buildInfoRow('Direccion', order.shippingStreet),
            if (order.shippingBarrio != null &&
                order.shippingBarrio!.isNotEmpty)
              _buildInfoRow('Barrio', order.shippingBarrio!),
            if (order.shippingCity.isNotEmpty)
              _buildInfoRow('Ciudad', order.shippingCity),
            if (order.shippingState.isNotEmpty)
              _buildInfoRow('Departamento', order.shippingState),
            if (order.shippingCountry.isNotEmpty)
              _buildInfoRow('Pais', order.shippingCountry),
            if (order.shippingPostalCode.isNotEmpty)
              _buildInfoRow('Codigo Postal', order.shippingPostalCode),
          ],
        ),
      ),
    );
  }

  Widget _buildFinancialCard(Order order) {
    final colorScheme = Theme.of(context).colorScheme;

    return Card(
      child: Padding(
        padding: const EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Row(
              children: [
                Icon(Icons.attach_money, size: 18, color: colorScheme.primary),
                const SizedBox(width: 8),
                Text(
                  'Resumen Financiero',
                  style: Theme.of(context).textTheme.titleSmall?.copyWith(
                        fontWeight: FontWeight.bold,
                      ),
                ),
              ],
            ),
            const Divider(height: 20),
            _buildAmountRow('Subtotal', order.subtotal, order.currency),
            if (order.tax > 0)
              _buildAmountRow('Impuesto', order.tax, order.currency),
            if (order.discount > 0)
              _buildAmountRow(
                  'Descuento', -order.discount, order.currency,
                  isDiscount: true),
            _buildAmountRow('Envio', order.shippingCost, order.currency),
            if (order.shippingDiscount != null && order.shippingDiscount! > 0)
              _buildAmountRow(
                  'Desc. Envio', -order.shippingDiscount!, order.currency,
                  isDiscount: true),
            const Divider(height: 16),
            _buildAmountRow('Total', order.totalAmount, order.currency,
                isTotal: true),
          ],
        ),
      ),
    );
  }

  Widget _buildItemsCard(Order order) {
    final colorScheme = Theme.of(context).colorScheme;

    final items = _parseItems(order);
    if (items.isEmpty) {
      return Card(
        child: Padding(
          padding: const EdgeInsets.all(16),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Row(
                children: [
                  Icon(Icons.shopping_bag_outlined,
                      size: 18, color: colorScheme.primary),
                  const SizedBox(width: 8),
                  Text(
                    'Productos',
                    style: Theme.of(context).textTheme.titleSmall?.copyWith(
                          fontWeight: FontWeight.bold,
                        ),
                  ),
                ],
              ),
              const Divider(height: 20),
              Center(
                child: Text(
                  'No hay informacion de productos',
                  style: Theme.of(context).textTheme.bodySmall?.copyWith(
                        color: colorScheme.outline,
                      ),
                ),
              ),
            ],
          ),
        ),
      );
    }

    return Card(
      child: Padding(
        padding: const EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Row(
              children: [
                Icon(Icons.shopping_bag_outlined,
                    size: 18, color: colorScheme.primary),
                const SizedBox(width: 8),
                Text(
                  'Productos (${items.length})',
                  style: Theme.of(context).textTheme.titleSmall?.copyWith(
                        fontWeight: FontWeight.bold,
                      ),
                ),
              ],
            ),
            const Divider(height: 20),
            ...items.map((item) => _buildItemRow(item, order.currency)),
          ],
        ),
      ),
    );
  }

  Widget _buildItemRow(Map<String, dynamic> item, String currency) {
    final name = (item['product_name'] ?? item['name'] ?? item['title'] ?? '-')
        .toString();
    final sku = (item['product_sku'] ?? item['sku'] ?? '').toString();
    final qty = (item['quantity'] ?? 0);
    final price = _toDouble(item['unit_price'] ?? item['price'] ?? 0);
    final totalPrice = _toDouble(item['total_price'] ?? (price * qty));

    return Padding(
      padding: const EdgeInsets.only(bottom: 12),
      child: Row(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  name,
                  style: Theme.of(context).textTheme.bodyMedium?.copyWith(
                        fontWeight: FontWeight.w500,
                      ),
                ),
                if (sku.isNotEmpty)
                  Text(
                    'SKU: $sku',
                    style: Theme.of(context).textTheme.bodySmall?.copyWith(
                          color: Theme.of(context).colorScheme.outline,
                        ),
                  ),
              ],
            ),
          ),
          const SizedBox(width: 8),
          Column(
            crossAxisAlignment: CrossAxisAlignment.end,
            children: [
              Text(
                '$currency ${_formatAmount(totalPrice)}',
                style: Theme.of(context).textTheme.bodyMedium?.copyWith(
                      fontWeight: FontWeight.bold,
                    ),
              ),
              Text(
                '${qty}x ${_formatAmount(price)}',
                style: Theme.of(context).textTheme.bodySmall?.copyWith(
                      color: Theme.of(context).colorScheme.outline,
                    ),
              ),
            ],
          ),
        ],
      ),
    );
  }

  Widget _buildDatesCard(Order order) {
    final colorScheme = Theme.of(context).colorScheme;

    return Card(
      child: Padding(
        padding: const EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Row(
              children: [
                Icon(Icons.schedule, size: 18, color: colorScheme.primary),
                const SizedBox(width: 8),
                Text(
                  'Cronologia',
                  style: Theme.of(context).textTheme.titleSmall?.copyWith(
                        fontWeight: FontWeight.bold,
                      ),
                ),
              ],
            ),
            const Divider(height: 20),
            _buildInfoRow('Creado', _formatDateTime(order.createdAt)),
            if (order.occurredAt.isNotEmpty)
              _buildInfoRow('Ocurrio', _formatDateTime(order.occurredAt)),
            if (order.importedAt.isNotEmpty)
              _buildInfoRow('Importado', _formatDateTime(order.importedAt)),
            _buildInfoRow('Actualizado', _formatDateTime(order.updatedAt)),
            if (order.paidAt != null && order.paidAt!.isNotEmpty)
              _buildInfoRow('Pagado', _formatDateTime(order.paidAt!)),
            if (order.deliveredAt != null && order.deliveredAt!.isNotEmpty)
              _buildInfoRow('Entregado', _formatDateTime(order.deliveredAt!)),
          ],
        ),
      ),
    );
  }

  Widget _buildInfoRow(String label, String value) {
    if (value.isEmpty) return const SizedBox.shrink();

    return Padding(
      padding: const EdgeInsets.only(bottom: 8),
      child: Row(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          SizedBox(
            width: 120,
            child: Text(
              label,
              style: Theme.of(context).textTheme.bodySmall?.copyWith(
                    color: Theme.of(context).colorScheme.outline,
                    fontWeight: FontWeight.w500,
                  ),
            ),
          ),
          Expanded(
            child: Text(
              value,
              style: Theme.of(context).textTheme.bodyMedium,
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildAmountRow(String label, double amount, String currency,
      {bool isTotal = false, bool isDiscount = false}) {
    final textTheme = Theme.of(context).textTheme;
    final colorScheme = Theme.of(context).colorScheme;

    return Padding(
      padding: const EdgeInsets.only(bottom: 4),
      child: Row(
        mainAxisAlignment: MainAxisAlignment.spaceBetween,
        children: [
          Text(
            label,
            style: isTotal
                ? textTheme.titleSmall?.copyWith(fontWeight: FontWeight.bold)
                : textTheme.bodyMedium,
          ),
          Text(
            '$currency ${_formatAmount(amount)}',
            style: isTotal
                ? textTheme.titleSmall?.copyWith(
                    fontWeight: FontWeight.bold,
                    color: colorScheme.primary,
                  )
                : isDiscount
                    ? textTheme.bodyMedium?.copyWith(
                        color: Colors.green.shade700,
                      )
                    : textTheme.bodyMedium?.copyWith(
                        fontWeight: FontWeight.w500,
                      ),
          ),
        ],
      ),
    );
  }

  // Helpers

  List<Map<String, dynamic>> _parseItems(Order order) {
    final raw = order.orderItems ?? order.items;
    if (raw == null) return [];
    if (raw is List) {
      return raw
          .whereType<Map<String, dynamic>>()
          .toList();
    }
    return [];
  }

  static double _toDouble(dynamic value) {
    if (value is double) return value;
    if (value is int) return value.toDouble();
    if (value is String) return double.tryParse(value) ?? 0;
    return 0;
  }

  static Color? _parseHexColor(String? hex) {
    if (hex == null || hex.isEmpty) return null;
    final clean = hex.replaceFirst('#', '');
    if (clean.length != 6) return null;
    final value = int.tryParse(clean, radix: 16);
    if (value == null) return null;
    return Color(0xFF000000 | value);
  }

  static Color _textColorForBg(Color bg) {
    final luminance =
        (0.299 * bg.r + 0.587 * bg.g + 0.114 * bg.b) / 255;
    return luminance > 0.5 ? Colors.black87 : Colors.white;
  }

  static String _formatAmount(double amount) {
    final absAmount = amount.abs();
    final prefix = amount < 0 ? '-' : '';
    if (absAmount == absAmount.roundToDouble()) {
      return '$prefix${absAmount.toStringAsFixed(0).replaceAllMapped(RegExp(r'(\d)(?=(\d{3})+(?!\d))'), (m) => '${m[1]},')}';
    }
    final parts = absAmount.toStringAsFixed(2).split('.');
    final integer = parts[0].replaceAllMapped(
        RegExp(r'(\d)(?=(\d{3})+(?!\d))'), (m) => '${m[1]},');
    return '$prefix$integer.${parts[1]}';
  }

  static String _formatDateTime(String dateStr) {
    if (dateStr.isEmpty) return '-';
    try {
      final date = DateTime.parse(dateStr);
      final months = [
        'ene', 'feb', 'mar', 'abr', 'may', 'jun',
        'jul', 'ago', 'sep', 'oct', 'nov', 'dic',
      ];
      final hour = date.hour.toString().padLeft(2, '0');
      final minute = date.minute.toString().padLeft(2, '0');
      return '${date.day} ${months[date.month - 1]} ${date.year}, $hour:$minute';
    } catch (_) {
      return dateStr;
    }
  }
}
