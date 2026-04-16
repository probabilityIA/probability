import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../providers/order_provider.dart';
import '../../domain/entities.dart';
import 'order_detail_screen.dart';

class OrderListScreen extends StatefulWidget {
  final int? businessId;

  const OrderListScreen({super.key, this.businessId});

  @override
  State<OrderListScreen> createState() => _OrderListScreenState();
}

class _OrderListScreenState extends State<OrderListScreen> {
  final _searchController = TextEditingController();
  String? _selectedStatus;

  static const _statusFilters = <String, String>{
    '': 'Todos',
    'pending': 'Pendiente',
    'processing': 'Procesando',
    'shipped': 'Enviado',
    'delivered': 'Entregado',
    'cancelled': 'Cancelado',
  };

  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addPostFrameCallback((_) {
      context.read<OrderProvider>().fetchOrders(businessId: widget.businessId);
    });
  }

  @override
  void didUpdateWidget(OrderListScreen oldWidget) {
    super.didUpdateWidget(oldWidget);
    if (oldWidget.businessId != widget.businessId) {
      context.read<OrderProvider>().resetFilters();
      context.read<OrderProvider>().fetchOrders(businessId: widget.businessId);
    }
  }

  @override
  void dispose() {
    _searchController.dispose();
    super.dispose();
  }

  void _onSearch(String value) {
    final provider = context.read<OrderProvider>();
    provider.setFilters(orderNumber: value);
    provider.fetchOrders(businessId: widget.businessId);
  }

  void _onStatusFilter(String? status) {
    setState(() {
      _selectedStatus = status;
    });
    final provider = context.read<OrderProvider>();
    provider.setFilters(
      status: (status != null && status.isNotEmpty) ? status : '',
    );
    provider.fetchOrders(businessId: widget.businessId);
  }

  void _goToPage(int page) {
    final provider = context.read<OrderProvider>();
    provider.setPage(page);
    provider.fetchOrders(businessId: widget.businessId);
  }

  void _navigateToDetail(Order order) {
    Navigator.of(context).push(
      MaterialPageRoute(
        builder: (_) => OrderDetailScreen(orderId: order.id),
      ),
    );
  }

  @override
  Widget build(BuildContext context) {
    final colorScheme = Theme.of(context).colorScheme;

    return Scaffold(
      appBar: AppBar(
        title: const Text('Ordenes'),
      ),
      body: Consumer<OrderProvider>(
        builder: (context, provider, child) {
          return Column(
            children: [
              // Search bar
              Padding(
                padding:
                    const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
                child: TextField(
                  controller: _searchController,
                  decoration: InputDecoration(
                    hintText: 'Buscar por numero de orden...',
                    prefixIcon: const Icon(Icons.search),
                    suffixIcon: _searchController.text.isNotEmpty
                        ? IconButton(
                            icon: const Icon(Icons.clear),
                            onPressed: () {
                              _searchController.clear();
                              _onSearch('');
                            },
                          )
                        : null,
                    isDense: true,
                  ),
                  onSubmitted: _onSearch,
                  textInputAction: TextInputAction.search,
                ),
              ),

              // Status filter chips
              SizedBox(
                height: 48,
                child: ListView(
                  scrollDirection: Axis.horizontal,
                  padding: const EdgeInsets.symmetric(horizontal: 12),
                  children: _statusFilters.entries.map((entry) {
                    final isSelected = (_selectedStatus ?? '') == entry.key;
                    return Padding(
                      padding: const EdgeInsets.symmetric(horizontal: 4),
                      child: FilterChip(
                        label: Text(entry.value),
                        selected: isSelected,
                        onSelected: (_) => _onStatusFilter(entry.key),
                        selectedColor: colorScheme.primaryContainer,
                        showCheckmark: false,
                      ),
                    );
                  }).toList(),
                ),
              ),

              const SizedBox(height: 4),

              // Content
              Expanded(
                child: _buildContent(provider, colorScheme),
              ),
            ],
          );
        },
      ),
    );
  }

  Widget _buildContent(OrderProvider provider, ColorScheme colorScheme) {
    if (provider.isLoading && provider.orders.isEmpty) {
      return const Center(child: CircularProgressIndicator());
    }

    if (provider.error != null && provider.orders.isEmpty) {
      return Center(
        child: Padding(
          padding: const EdgeInsets.all(24),
          child: Column(
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              Icon(Icons.error_outline, size: 48, color: colorScheme.error),
              const SizedBox(height: 16),
              Text(
                provider.error!,
                textAlign: TextAlign.center,
                style: TextStyle(color: colorScheme.error),
              ),
              const SizedBox(height: 16),
              FilledButton.icon(
                onPressed: () =>
                    provider.fetchOrders(businessId: widget.businessId),
                icon: const Icon(Icons.refresh),
                label: const Text('Reintentar'),
              ),
            ],
          ),
        ),
      );
    }

    if (provider.orders.isEmpty) {
      return Center(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Icon(Icons.receipt_long_outlined,
                size: 64, color: colorScheme.outline),
            const SizedBox(height: 16),
            Text(
              'No hay ordenes',
              style: Theme.of(context).textTheme.titleMedium?.copyWith(
                    color: colorScheme.outline,
                  ),
            ),
          ],
        ),
      );
    }

    return Column(
      children: [
        Expanded(
          child: RefreshIndicator(
            onRefresh: () =>
                provider.fetchOrders(businessId: widget.businessId),
            child: ListView.builder(
              padding: const EdgeInsets.symmetric(horizontal: 16),
              itemCount: provider.orders.length,
              itemBuilder: (context, index) {
                final order = provider.orders[index];
                return _OrderCard(
                  order: order,
                  onTap: () => _navigateToDetail(order),
                );
              },
            ),
          ),
        ),

        // Pagination
        if (provider.pagination != null) _buildPagination(provider),
      ],
    );
  }

  Widget _buildPagination(OrderProvider provider) {
    final pagination = provider.pagination!;
    if (pagination.lastPage <= 1) return const SizedBox.shrink();

    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
      decoration: BoxDecoration(
        color: Theme.of(context).colorScheme.surface,
        border: Border(
          top: BorderSide(
            color: Theme.of(context).colorScheme.outlineVariant,
          ),
        ),
      ),
      child: Row(
        mainAxisAlignment: MainAxisAlignment.spaceBetween,
        children: [
          Text(
            '${pagination.total} resultados',
            style: Theme.of(context).textTheme.bodySmall,
          ),
          Row(
            children: [
              IconButton(
                icon: const Icon(Icons.chevron_left),
                onPressed: pagination.hasPrev
                    ? () => _goToPage(pagination.currentPage - 1)
                    : null,
                iconSize: 20,
              ),
              Container(
                padding:
                    const EdgeInsets.symmetric(horizontal: 12, vertical: 4),
                decoration: BoxDecoration(
                  color: Theme.of(context).colorScheme.primaryContainer,
                  borderRadius: BorderRadius.circular(16),
                ),
                child: Text(
                  '${pagination.currentPage} / ${pagination.lastPage}',
                  style: Theme.of(context).textTheme.bodySmall?.copyWith(
                        fontWeight: FontWeight.bold,
                      ),
                ),
              ),
              IconButton(
                icon: const Icon(Icons.chevron_right),
                onPressed: pagination.hasNext
                    ? () => _goToPage(pagination.currentPage + 1)
                    : null,
                iconSize: 20,
              ),
            ],
          ),
        ],
      ),
    );
  }
}

class _OrderCard extends StatelessWidget {
  final Order order;
  final VoidCallback onTap;

  const _OrderCard({required this.order, required this.onTap});

  @override
  Widget build(BuildContext context) {
    final colorScheme = Theme.of(context).colorScheme;
    final textTheme = Theme.of(context).textTheme;

    return Card(
      margin: const EdgeInsets.only(bottom: 8),
      child: InkWell(
        onTap: onTap,
        borderRadius: BorderRadius.circular(12),
        child: Padding(
          padding: const EdgeInsets.all(12),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              // Top row: order number + status badge
              Row(
                children: [
                  // Platform icon
                  if (order.integrationLogoUrl != null &&
                      order.integrationLogoUrl!.isNotEmpty)
                    Padding(
                      padding: const EdgeInsets.only(right: 8),
                      child: ClipOval(
                        child: Image.network(
                          order.integrationLogoUrl!,
                          width: 28,
                          height: 28,
                          fit: BoxFit.contain,
                          errorBuilder: (context, error, stackTrace) => CircleAvatar(
                            radius: 14,
                            backgroundColor: colorScheme.surfaceContainerHighest,
                            child: Text(
                              order.platform.isNotEmpty
                                  ? order.platform[0].toUpperCase()
                                  : '?',
                              style: textTheme.labelSmall,
                            ),
                          ),
                        ),
                      ),
                    ),
                  Expanded(
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        Text(
                          '#${order.orderNumber.isNotEmpty ? order.orderNumber : order.externalId.isNotEmpty ? order.externalId : order.id}',
                          style: textTheme.titleSmall?.copyWith(
                            fontWeight: FontWeight.bold,
                          ),
                        ),
                        if (order.internalNumber.isNotEmpty)
                          Text(
                            'Interno: ${order.internalNumber}',
                            style: textTheme.bodySmall?.copyWith(
                              color: colorScheme.outline,
                            ),
                          ),
                      ],
                    ),
                  ),
                  _StatusBadge(
                    label: order.orderStatus?.name ?? order.status,
                    color: _parseHexColor(order.orderStatus?.color),
                  ),
                ],
              ),

              const SizedBox(height: 8),

              // Customer name
              Row(
                children: [
                  Icon(Icons.person_outline,
                      size: 16, color: colorScheme.outline),
                  const SizedBox(width: 4),
                  Expanded(
                    child: Text(
                      order.customerName.isNotEmpty
                          ? order.customerName
                          : 'Sin cliente',
                      style: textTheme.bodyMedium,
                      overflow: TextOverflow.ellipsis,
                    ),
                  ),
                ],
              ),

              const SizedBox(height: 8),

              // Bottom row: total + payment status + date
              Row(
                children: [
                  // Total amount
                  Container(
                    padding:
                        const EdgeInsets.symmetric(horizontal: 8, vertical: 2),
                    decoration: BoxDecoration(
                      color: colorScheme.primaryContainer,
                      borderRadius: BorderRadius.circular(8),
                    ),
                    child: Text(
                      '${order.currency} ${_formatAmount(order.totalAmount)}',
                      style: textTheme.labelMedium?.copyWith(
                        fontWeight: FontWeight.bold,
                        color: colorScheme.onPrimaryContainer,
                      ),
                    ),
                  ),

                  const SizedBox(width: 8),

                  // Payment status
                  _PaymentBadge(
                    label: order.paymentStatus?.name,
                    isPaid: order.isPaid,
                  ),

                  const Spacer(),

                  // Date
                  Text(
                    _formatDate(order.createdAt),
                    style: textTheme.bodySmall?.copyWith(
                      color: colorScheme.outline,
                    ),
                  ),
                ],
              ),
            ],
          ),
        ),
      ),
    );
  }

  static Color? _parseHexColor(String? hex) {
    if (hex == null || hex.isEmpty) return null;
    final clean = hex.replaceFirst('#', '');
    if (clean.length != 6) return null;
    final value = int.tryParse(clean, radix: 16);
    if (value == null) return null;
    return Color(0xFF000000 | value);
  }

  static String _formatAmount(double amount) {
    if (amount == amount.roundToDouble()) {
      return amount.toStringAsFixed(0).replaceAllMapped(
          RegExp(r'(\d)(?=(\d{3})+(?!\d))'), (m) => '${m[1]},');
    }
    final parts = amount.toStringAsFixed(2).split('.');
    final integer = parts[0].replaceAllMapped(
        RegExp(r'(\d)(?=(\d{3})+(?!\d))'), (m) => '${m[1]},');
    return '$integer.${parts[1]}';
  }

  static String _formatDate(String dateStr) {
    if (dateStr.isEmpty) return '';
    try {
      final date = DateTime.parse(dateStr);
      final months = [
        'ene', 'feb', 'mar', 'abr', 'may', 'jun',
        'jul', 'ago', 'sep', 'oct', 'nov', 'dic',
      ];
      return '${date.day} ${months[date.month - 1]} ${date.year}';
    } catch (_) {
      return dateStr;
    }
  }
}

class _StatusBadge extends StatelessWidget {
  final String label;
  final Color? color;

  const _StatusBadge({required this.label, this.color});

  @override
  Widget build(BuildContext context) {
    final bgColor = color ?? _fallbackColor(label);
    final textColor = _textColorForBg(bgColor);

    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 4),
      decoration: BoxDecoration(
        color: bgColor,
        borderRadius: BorderRadius.circular(16),
      ),
      child: Text(
        label,
        style: TextStyle(
          fontSize: 11,
          fontWeight: FontWeight.w600,
          color: textColor,
        ),
      ),
    );
  }

  static Color _fallbackColor(String status) {
    switch (status.toLowerCase()) {
      case 'pending':
      case 'pendiente':
        return const Color(0xFFFEF3C7);
      case 'processing':
      case 'procesando':
        return const Color(0xFFDBEAFE);
      case 'shipped':
      case 'enviado':
        return const Color(0xFFEDE9FE);
      case 'delivered':
      case 'entregado':
        return const Color(0xFFDCFCE7);
      case 'cancelled':
      case 'cancelado':
        return const Color(0xFFFEE2E2);
      default:
        return const Color(0xFFF3F4F6);
    }
  }

  static Color _textColorForBg(Color bg) {
    final luminance =
        (0.299 * bg.r + 0.587 * bg.g + 0.114 * bg.b) / 255;
    return luminance > 0.5 ? Colors.black87 : Colors.white;
  }
}

class _PaymentBadge extends StatelessWidget {
  final String? label;
  final bool isPaid;

  const _PaymentBadge({this.label, required this.isPaid});

  @override
  Widget build(BuildContext context) {
    final displayLabel = label ?? (isPaid ? 'Pagado' : 'No pagado');
    final bgColor = isPaid ? const Color(0xFFDCFCE7) : const Color(0xFFFEE2E2);
    final textColor =
        isPaid ? const Color(0xFF166534) : const Color(0xFF991B1B);

    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 2),
      decoration: BoxDecoration(
        color: bgColor,
        borderRadius: BorderRadius.circular(8),
      ),
      child: Text(
        displayLabel,
        style: TextStyle(
          fontSize: 10,
          fontWeight: FontWeight.w600,
          color: textColor,
        ),
      ),
    );
  }
}
