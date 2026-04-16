import 'dart:math' as math;
import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../../domain/entities.dart';
import '../providers/dashboard_provider.dart';

class DashboardScreen extends StatefulWidget {
  final int? businessId;

  const DashboardScreen({super.key, this.businessId});

  @override
  State<DashboardScreen> createState() => _DashboardScreenState();
}

class _DashboardScreenState extends State<DashboardScreen> {
  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addPostFrameCallback((_) {
      _loadStats();
    });
  }

  void _loadStats() {
    context.read<DashboardProvider>().fetchStats(businessId: widget.businessId);
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Dashboard'),
        centerTitle: false,
      ),
      body: Consumer<DashboardProvider>(
        builder: (context, provider, _) {
          if (provider.isLoading && provider.stats == null) {
            return const _LoadingState();
          }
          if (provider.error != null && provider.stats == null) {
            return _ErrorState(
              error: provider.error!,
              onRetry: _loadStats,
            );
          }
          final stats = provider.stats;
          if (stats == null) {
            return const _EmptyState();
          }
          return RefreshIndicator(
            onRefresh: () async => _loadStats(),
            child: SingleChildScrollView(
              physics: const AlwaysScrollableScrollPhysics(),
              padding: const EdgeInsets.fromLTRB(16, 8, 16, 32),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  _SummaryCardsRow(stats: stats),
                  const SizedBox(height: 24),
                  if (stats.ordersByIntegrationType.isNotEmpty) ...[
                    _OrdersByChannelSection(
                        items: stats.ordersByIntegrationType),
                    const SizedBox(height: 24),
                  ],
                  if (stats.topCustomers.isNotEmpty) ...[
                    _TopCustomersSection(customers: stats.topCustomers),
                    const SizedBox(height: 24),
                  ],
                  if (stats.topProducts.isNotEmpty) ...[
                    _TopProductsSection(products: stats.topProducts),
                    const SizedBox(height: 24),
                  ],
                  if (stats.shipmentsByStatus.isNotEmpty) ...[
                    _ShipmentsByStatusSection(items: stats.shipmentsByStatus),
                    const SizedBox(height: 24),
                  ],
                  if (stats.shipmentsByCarrier.isNotEmpty) ...[
                    _ShipmentsByCarrierSection(items: stats.shipmentsByCarrier),
                    const SizedBox(height: 24),
                  ],
                  if (stats.ordersByLocation.isNotEmpty) ...[
                    _OrdersByDepartmentSection(items: stats.ordersByLocation),
                    const SizedBox(height: 24),
                  ],
                ],
              ),
            ),
          );
        },
      ),
    );
  }
}

// ---------------------------------------------------------------------------
// Loading state
// ---------------------------------------------------------------------------

class _LoadingState extends StatelessWidget {
  const _LoadingState();

  @override
  Widget build(BuildContext context) {
    final colorScheme = Theme.of(context).colorScheme;
    return Center(
      child: Column(
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          SizedBox(
            width: 48,
            height: 48,
            child: CircularProgressIndicator(
              strokeWidth: 3,
              color: colorScheme.primary,
            ),
          ),
          const SizedBox(height: 16),
          Text(
            'Cargando estadisticas...',
            style: TextStyle(
              fontSize: 14,
              color: colorScheme.onSurfaceVariant,
            ),
          ),
        ],
      ),
    );
  }
}

// ---------------------------------------------------------------------------
// Error state
// ---------------------------------------------------------------------------

class _ErrorState extends StatelessWidget {
  final String error;
  final VoidCallback onRetry;

  const _ErrorState({required this.error, required this.onRetry});

  @override
  Widget build(BuildContext context) {
    final colorScheme = Theme.of(context).colorScheme;
    return Center(
      child: Padding(
        padding: const EdgeInsets.all(32),
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Icon(Icons.cloud_off_rounded, size: 56, color: colorScheme.error),
            const SizedBox(height: 16),
            Text(
              'Error al cargar datos',
              style: TextStyle(
                fontSize: 18,
                fontWeight: FontWeight.w600,
                color: colorScheme.onSurface,
              ),
            ),
            const SizedBox(height: 8),
            Text(
              error,
              textAlign: TextAlign.center,
              style: TextStyle(
                fontSize: 14,
                color: colorScheme.onSurfaceVariant,
              ),
            ),
            const SizedBox(height: 24),
            FilledButton.icon(
              onPressed: onRetry,
              icon: const Icon(Icons.refresh_rounded),
              label: const Text('Reintentar'),
            ),
          ],
        ),
      ),
    );
  }
}

// ---------------------------------------------------------------------------
// Empty state
// ---------------------------------------------------------------------------

class _EmptyState extends StatelessWidget {
  const _EmptyState();

  @override
  Widget build(BuildContext context) {
    final colorScheme = Theme.of(context).colorScheme;
    return Center(
      child: Column(
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          Icon(
            Icons.dashboard_outlined,
            size: 64,
            color: colorScheme.onSurfaceVariant.withValues(alpha: 0.5),
          ),
          const SizedBox(height: 16),
          Text(
            'Sin datos disponibles',
            style: TextStyle(
              fontSize: 16,
              color: colorScheme.onSurfaceVariant,
            ),
          ),
        ],
      ),
    );
  }
}

// ---------------------------------------------------------------------------
// Section header
// ---------------------------------------------------------------------------

class _SectionHeader extends StatelessWidget {
  final IconData icon;
  final String title;

  const _SectionHeader({required this.icon, required this.title});

  @override
  Widget build(BuildContext context) {
    final colorScheme = Theme.of(context).colorScheme;
    return Padding(
      padding: const EdgeInsets.only(bottom: 12),
      child: Row(
        children: [
          Icon(icon, size: 20, color: colorScheme.primary),
          const SizedBox(width: 8),
          Text(
            title,
            style: TextStyle(
              fontSize: 16,
              fontWeight: FontWeight.w700,
              color: colorScheme.onSurface,
            ),
          ),
        ],
      ),
    );
  }
}

// ---------------------------------------------------------------------------
// 1) Summary cards row
// ---------------------------------------------------------------------------

class _SummaryCardsRow extends StatelessWidget {
  final DashboardStats stats;

  const _SummaryCardsRow({required this.stats});

  @override
  Widget build(BuildContext context) {
    final totalChannels = stats.ordersByIntegrationType.length;
    final totalProducts = stats.topProducts.fold<int>(
      0,
      (sum, p) => sum + p.totalSold,
    );
    final totalShipments = stats.shipmentsByStatus.fold<int>(
      0,
      (sum, s) => sum + s.count,
    );

    return SingleChildScrollView(
      scrollDirection: Axis.horizontal,
      child: Row(
        children: [
          _StatCard(
            icon: Icons.shopping_bag_rounded,
            label: 'Total Ordenes',
            value: _formatNumber(stats.totalOrders),
            color: Colors.deepPurple,
          ),
          const SizedBox(width: 12),
          _StatCard(
            icon: Icons.storefront_rounded,
            label: 'Canales',
            value: '$totalChannels',
            color: Colors.blue,
          ),
          const SizedBox(width: 12),
          _StatCard(
            icon: Icons.inventory_2_rounded,
            label: 'Uds. Vendidas',
            value: _formatNumber(totalProducts),
            color: Colors.teal,
          ),
          const SizedBox(width: 12),
          _StatCard(
            icon: Icons.local_shipping_rounded,
            label: 'Envios',
            value: _formatNumber(totalShipments),
            color: Colors.orange,
          ),
        ],
      ),
    );
  }
}

class _StatCard extends StatelessWidget {
  final IconData icon;
  final String label;
  final String value;
  final Color color;

  const _StatCard({
    required this.icon,
    required this.label,
    required this.value,
    required this.color,
  });

  @override
  Widget build(BuildContext context) {
    final colorScheme = Theme.of(context).colorScheme;
    return SizedBox(
      width: 150,
      child: Card(
        elevation: 0,
        color: color.withValues(alpha: 0.08),
        shape: RoundedRectangleBorder(
          borderRadius: BorderRadius.circular(16),
          side: BorderSide(color: color.withValues(alpha: 0.2)),
        ),
        child: Padding(
          padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 20),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Container(
                padding: const EdgeInsets.all(8),
                decoration: BoxDecoration(
                  color: color.withValues(alpha: 0.15),
                  borderRadius: BorderRadius.circular(10),
                ),
                child: Icon(icon, size: 22, color: color),
              ),
              const SizedBox(height: 14),
              Text(
                value,
                style: TextStyle(
                  fontSize: 26,
                  fontWeight: FontWeight.w800,
                  color: colorScheme.onSurface,
                ),
              ),
              const SizedBox(height: 2),
              Text(
                label,
                style: TextStyle(
                  fontSize: 12,
                  fontWeight: FontWeight.w500,
                  color: colorScheme.onSurfaceVariant,
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }
}

// ---------------------------------------------------------------------------
// 2) Orders by channel
// ---------------------------------------------------------------------------

class _OrdersByChannelSection extends StatelessWidget {
  final List<OrderCountByIntegrationType> items;

  const _OrdersByChannelSection({required this.items});

  static IconData _iconForChannel(String type) {
    final lower = type.toLowerCase();
    if (lower.contains('shopify')) return Icons.shopping_cart_rounded;
    if (lower.contains('amazon')) return Icons.store_rounded;
    if (lower.contains('meli') || lower.contains('mercadolibre')) {
      return Icons.handshake_rounded;
    }
    if (lower.contains('whatsapp')) return Icons.chat_rounded;
    if (lower.contains('manual')) return Icons.edit_note_rounded;
    if (lower.contains('web') || lower.contains('tienda')) {
      return Icons.language_rounded;
    }
    return Icons.storefront_rounded;
  }

  static Color _colorForChannel(String type) {
    final lower = type.toLowerCase();
    if (lower.contains('shopify')) return Colors.green;
    if (lower.contains('amazon')) return Colors.orange;
    if (lower.contains('meli') || lower.contains('mercadolibre')) {
      return Colors.yellow.shade800;
    }
    if (lower.contains('whatsapp')) return Colors.teal;
    if (lower.contains('manual')) return Colors.blueGrey;
    if (lower.contains('web') || lower.contains('tienda')) return Colors.blue;
    return Colors.deepPurple;
  }

  @override
  Widget build(BuildContext context) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        const _SectionHeader(
          icon: Icons.hub_rounded,
          title: 'Ordenes por Canal',
        ),
        SingleChildScrollView(
          scrollDirection: Axis.horizontal,
          child: Row(
            children: items.map((item) {
              final color = _colorForChannel(item.integrationType);
              final icon = _iconForChannel(item.integrationType);
              return Padding(
                padding: const EdgeInsets.only(right: 10),
                child: _ChannelChip(
                  icon: icon,
                  label: item.integrationType,
                  count: item.count,
                  color: color,
                ),
              );
            }).toList(),
          ),
        ),
      ],
    );
  }
}

class _ChannelChip extends StatelessWidget {
  final IconData icon;
  final String label;
  final int count;
  final Color color;

  const _ChannelChip({
    required this.icon,
    required this.label,
    required this.count,
    required this.color,
  });

  @override
  Widget build(BuildContext context) {
    final colorScheme = Theme.of(context).colorScheme;
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 14, vertical: 10),
      decoration: BoxDecoration(
        color: color.withValues(alpha: 0.08),
        borderRadius: BorderRadius.circular(12),
        border: Border.all(color: color.withValues(alpha: 0.25)),
      ),
      child: Row(
        mainAxisSize: MainAxisSize.min,
        children: [
          Icon(icon, size: 18, color: color),
          const SizedBox(width: 8),
          Text(
            label,
            style: TextStyle(
              fontSize: 13,
              fontWeight: FontWeight.w600,
              color: colorScheme.onSurface,
            ),
          ),
          const SizedBox(width: 8),
          Container(
            padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 2),
            decoration: BoxDecoration(
              color: color.withValues(alpha: 0.18),
              borderRadius: BorderRadius.circular(10),
            ),
            child: Text(
              _formatNumber(count),
              style: TextStyle(
                fontSize: 12,
                fontWeight: FontWeight.w700,
                color: color,
              ),
            ),
          ),
        ],
      ),
    );
  }
}

// ---------------------------------------------------------------------------
// 3) Top customers
// ---------------------------------------------------------------------------

class _TopCustomersSection extends StatelessWidget {
  final List<TopCustomer> customers;

  const _TopCustomersSection({required this.customers});

  @override
  Widget build(BuildContext context) {
    final colorScheme = Theme.of(context).colorScheme;
    final top = customers.take(5).toList();

    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        const _SectionHeader(
          icon: Icons.people_rounded,
          title: 'Top Clientes',
        ),
        Card(
          elevation: 0,
          color: colorScheme.surfaceContainerLow,
          shape: RoundedRectangleBorder(
            borderRadius: BorderRadius.circular(16),
          ),
          child: ListView.separated(
            shrinkWrap: true,
            physics: const NeverScrollableScrollPhysics(),
            itemCount: top.length,
            separatorBuilder: (context, idx) => Divider(
              height: 1,
              indent: 68,
              color: colorScheme.outlineVariant.withValues(alpha: 0.5),
            ),
            itemBuilder: (context, index) {
              final c = top[index];
              final initial = c.customerName.isNotEmpty
                  ? c.customerName[0].toUpperCase()
                  : '?';
              final avatarColors = [
                Colors.deepPurple,
                Colors.blue,
                Colors.teal,
                Colors.orange,
                Colors.pink,
              ];
              final avatarColor = avatarColors[index % avatarColors.length];

              return ListTile(
                contentPadding: const EdgeInsets.symmetric(
                  horizontal: 16,
                  vertical: 4,
                ),
                leading: CircleAvatar(
                  radius: 20,
                  backgroundColor: avatarColor.withValues(alpha: 0.15),
                  child: Text(
                    initial,
                    style: TextStyle(
                      fontWeight: FontWeight.w700,
                      color: avatarColor,
                      fontSize: 16,
                    ),
                  ),
                ),
                title: Text(
                  c.customerName,
                  style: const TextStyle(
                    fontWeight: FontWeight.w600,
                    fontSize: 14,
                  ),
                  maxLines: 1,
                  overflow: TextOverflow.ellipsis,
                ),
                subtitle: Text(
                  c.customerEmail,
                  style: TextStyle(
                    fontSize: 12,
                    color: colorScheme.onSurfaceVariant,
                  ),
                  maxLines: 1,
                  overflow: TextOverflow.ellipsis,
                ),
                trailing: Container(
                  padding: const EdgeInsets.symmetric(
                    horizontal: 10,
                    vertical: 4,
                  ),
                  decoration: BoxDecoration(
                    color: colorScheme.primaryContainer,
                    borderRadius: BorderRadius.circular(12),
                  ),
                  child: Text(
                    '${c.orderCount}',
                    style: TextStyle(
                      fontWeight: FontWeight.w700,
                      fontSize: 13,
                      color: colorScheme.onPrimaryContainer,
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

// ---------------------------------------------------------------------------
// 4) Top products
// ---------------------------------------------------------------------------

class _TopProductsSection extends StatelessWidget {
  final List<TopProduct> products;

  const _TopProductsSection({required this.products});

  @override
  Widget build(BuildContext context) {
    final colorScheme = Theme.of(context).colorScheme;
    final top = products.take(5).toList();

    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        const _SectionHeader(
          icon: Icons.inventory_2_rounded,
          title: 'Top Productos',
        ),
        Card(
          elevation: 0,
          color: colorScheme.surfaceContainerLow,
          shape: RoundedRectangleBorder(
            borderRadius: BorderRadius.circular(16),
          ),
          child: ListView.separated(
            shrinkWrap: true,
            physics: const NeverScrollableScrollPhysics(),
            itemCount: top.length,
            separatorBuilder: (context, idx) => Divider(
              height: 1,
              indent: 68,
              color: colorScheme.outlineVariant.withValues(alpha: 0.5),
            ),
            itemBuilder: (context, index) {
              final p = top[index];
              return ListTile(
                contentPadding: const EdgeInsets.symmetric(
                  horizontal: 16,
                  vertical: 4,
                ),
                leading: Container(
                  width: 40,
                  height: 40,
                  decoration: BoxDecoration(
                    color: Colors.deepPurple.withValues(alpha: 0.1),
                    borderRadius: BorderRadius.circular(10),
                  ),
                  child: Center(
                    child: Text(
                      '#${index + 1}',
                      style: TextStyle(
                        fontWeight: FontWeight.w800,
                        fontSize: 14,
                        color: Colors.deepPurple.shade700,
                      ),
                    ),
                  ),
                ),
                title: Text(
                  p.productName,
                  style: const TextStyle(
                    fontWeight: FontWeight.w600,
                    fontSize: 14,
                  ),
                  maxLines: 1,
                  overflow: TextOverflow.ellipsis,
                ),
                subtitle: Text(
                  'SKU: ${p.sku}',
                  style: TextStyle(
                    fontSize: 12,
                    color: colorScheme.onSurfaceVariant,
                  ),
                ),
                trailing: Column(
                  mainAxisAlignment: MainAxisAlignment.center,
                  crossAxisAlignment: CrossAxisAlignment.end,
                  children: [
                    Text(
                      '${_formatNumber(p.totalSold)} uds.',
                      style: TextStyle(
                        fontWeight: FontWeight.w700,
                        fontSize: 13,
                        color: colorScheme.onSurface,
                      ),
                    ),
                    Text(
                      '${p.orderCount} ordenes',
                      style: TextStyle(
                        fontSize: 11,
                        color: colorScheme.onSurfaceVariant,
                      ),
                    ),
                  ],
                ),
              );
            },
          ),
        ),
      ],
    );
  }
}

// ---------------------------------------------------------------------------
// 5) Shipments by status
// ---------------------------------------------------------------------------

class _ShipmentsByStatusSection extends StatelessWidget {
  final List<ShipmentsByStatus> items;

  const _ShipmentsByStatusSection({required this.items});

  static Color _colorForStatus(String status) {
    final lower = status.toLowerCase();
    if (lower.contains('entreg') || lower.contains('deliver')) {
      return Colors.green;
    }
    if (lower.contains('cancel') || lower.contains('devol') ||
        lower.contains('return') || lower.contains('fail')) {
      return Colors.red;
    }
    if (lower.contains('pendi') || lower.contains('process') ||
        lower.contains('transit') || lower.contains('enviad')) {
      return Colors.orange;
    }
    if (lower.contains('recog') || lower.contains('pickup')) {
      return Colors.blue;
    }
    return Colors.blueGrey;
  }

  @override
  Widget build(BuildContext context) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        const _SectionHeader(
          icon: Icons.local_shipping_rounded,
          title: 'Envios por Estado',
        ),
        Wrap(
          spacing: 8,
          runSpacing: 8,
          children: items.map((item) {
            final color = _colorForStatus(item.status);
            return Container(
              padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 8),
              decoration: BoxDecoration(
                color: color.withValues(alpha: 0.1),
                borderRadius: BorderRadius.circular(20),
                border: Border.all(color: color.withValues(alpha: 0.3)),
              ),
              child: Row(
                mainAxisSize: MainAxisSize.min,
                children: [
                  Container(
                    width: 8,
                    height: 8,
                    decoration: BoxDecoration(
                      color: color,
                      shape: BoxShape.circle,
                    ),
                  ),
                  const SizedBox(width: 8),
                  Text(
                    item.status,
                    style: TextStyle(
                      fontSize: 12,
                      fontWeight: FontWeight.w600,
                      color: color.withValues(alpha: 0.9),
                    ),
                  ),
                  const SizedBox(width: 6),
                  Text(
                    '${item.count}',
                    style: TextStyle(
                      fontSize: 12,
                      fontWeight: FontWeight.w800,
                      color: color,
                    ),
                  ),
                ],
              ),
            );
          }).toList(),
        ),
      ],
    );
  }
}

// ---------------------------------------------------------------------------
// 6) Shipments by carrier
// ---------------------------------------------------------------------------

class _ShipmentsByCarrierSection extends StatelessWidget {
  final List<ShipmentsByCarrier> items;

  const _ShipmentsByCarrierSection({required this.items});

  @override
  Widget build(BuildContext context) {
    final colorScheme = Theme.of(context).colorScheme;
    final maxCount = items.fold<int>(0, (m, i) => math.max(m, i.count));

    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        const _SectionHeader(
          icon: Icons.flight_rounded,
          title: 'Envios por Transportadora',
        ),
        Card(
          elevation: 0,
          color: colorScheme.surfaceContainerLow,
          shape: RoundedRectangleBorder(
            borderRadius: BorderRadius.circular(16),
          ),
          child: Padding(
            padding: const EdgeInsets.all(16),
            child: Column(
              children: items.map((item) {
                final ratio = maxCount > 0 ? item.count / maxCount : 0.0;
                return Padding(
                  padding: const EdgeInsets.only(bottom: 12),
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Row(
                        mainAxisAlignment: MainAxisAlignment.spaceBetween,
                        children: [
                          Expanded(
                            child: Text(
                              item.carrier,
                              style: const TextStyle(
                                fontWeight: FontWeight.w600,
                                fontSize: 13,
                              ),
                              maxLines: 1,
                              overflow: TextOverflow.ellipsis,
                            ),
                          ),
                          Text(
                            '${item.count}',
                            style: TextStyle(
                              fontWeight: FontWeight.w700,
                              fontSize: 13,
                              color: colorScheme.primary,
                            ),
                          ),
                        ],
                      ),
                      const SizedBox(height: 6),
                      ClipRRect(
                        borderRadius: BorderRadius.circular(4),
                        child: LinearProgressIndicator(
                          value: ratio,
                          minHeight: 6,
                          backgroundColor:
                              colorScheme.primary.withValues(alpha: 0.08),
                          valueColor: AlwaysStoppedAnimation<Color>(
                            colorScheme.primary,
                          ),
                        ),
                      ),
                    ],
                  ),
                );
              }).toList(),
            ),
          ),
        ),
      ],
    );
  }
}

// ---------------------------------------------------------------------------
// 7) Orders by department (uses ordersByLocation state field)
// ---------------------------------------------------------------------------

class _OrdersByDepartmentSection extends StatelessWidget {
  final List<OrderCountByLocation> items;

  const _OrdersByDepartmentSection({required this.items});

  @override
  Widget build(BuildContext context) {
    final colorScheme = Theme.of(context).colorScheme;

    // Aggregate by state (department)
    final deptMap = <String, int>{};
    for (final item in items) {
      final dept = item.state.isNotEmpty ? item.state : 'Sin Departamento';
      deptMap[dept] = (deptMap[dept] ?? 0) + item.orderCount;
    }

    final departments = deptMap.entries.toList()
      ..sort((a, b) => b.value.compareTo(a.value));

    final maxCount =
        departments.fold<int>(0, (m, e) => math.max(m, e.value));
    final totalOrders = departments.fold<int>(0, (s, e) => s + e.value);

    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        const _SectionHeader(
          icon: Icons.map_rounded,
          title: 'Ordenes por Departamento',
        ),
        Card(
          elevation: 0,
          color: colorScheme.surfaceContainerLow,
          shape: RoundedRectangleBorder(
            borderRadius: BorderRadius.circular(16),
          ),
          child: Padding(
            padding: const EdgeInsets.all(16),
            child: Column(
              children: departments.take(10).map((entry) {
                final ratio = maxCount > 0 ? entry.value / maxCount : 0.0;
                final pct = totalOrders > 0
                    ? (entry.value / totalOrders * 100).toStringAsFixed(1)
                    : '0';
                return Padding(
                  padding: const EdgeInsets.only(bottom: 12),
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Row(
                        children: [
                          Icon(
                            Icons.location_on_rounded,
                            size: 16,
                            color: colorScheme.onSurfaceVariant,
                          ),
                          const SizedBox(width: 6),
                          Expanded(
                            child: Text(
                              entry.key,
                              style: const TextStyle(
                                fontWeight: FontWeight.w600,
                                fontSize: 13,
                              ),
                              maxLines: 1,
                              overflow: TextOverflow.ellipsis,
                            ),
                          ),
                          Text(
                            '${entry.value}',
                            style: TextStyle(
                              fontWeight: FontWeight.w700,
                              fontSize: 13,
                              color: colorScheme.primary,
                            ),
                          ),
                          const SizedBox(width: 6),
                          SizedBox(
                            width: 44,
                            child: Text(
                              '$pct%',
                              textAlign: TextAlign.right,
                              style: TextStyle(
                                fontSize: 11,
                                fontWeight: FontWeight.w500,
                                color: colorScheme.onSurfaceVariant,
                              ),
                            ),
                          ),
                        ],
                      ),
                      const SizedBox(height: 6),
                      ClipRRect(
                        borderRadius: BorderRadius.circular(4),
                        child: LinearProgressIndicator(
                          value: ratio,
                          minHeight: 6,
                          backgroundColor:
                              Colors.deepPurple.withValues(alpha: 0.08),
                          valueColor: const AlwaysStoppedAnimation<Color>(
                            Colors.deepPurple,
                          ),
                        ),
                      ),
                    ],
                  ),
                );
              }).toList(),
            ),
          ),
        ),
      ],
    );
  }
}

// ---------------------------------------------------------------------------
// Number formatting utility
// ---------------------------------------------------------------------------

String _formatNumber(int value) {
  if (value >= 1000000) {
    return '${(value / 1000000).toStringAsFixed(1)}M';
  }
  if (value >= 1000) {
    return '${(value / 1000).toStringAsFixed(1)}K';
  }
  return value.toString();
}
