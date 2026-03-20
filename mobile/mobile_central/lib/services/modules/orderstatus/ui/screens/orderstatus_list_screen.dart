import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../../domain/entities.dart';
import '../providers/orderstatus_provider.dart';

class OrderStatusListScreen extends StatefulWidget {
  const OrderStatusListScreen({super.key});

  @override
  State<OrderStatusListScreen> createState() => _OrderStatusListScreenState();
}

class _OrderStatusListScreenState extends State<OrderStatusListScreen>
    with SingleTickerProviderStateMixin {
  late TabController _tabController;
  int _currentPage = 1;

  @override
  void initState() {
    super.initState();
    _tabController = TabController(length: 2, vsync: this);
    WidgetsBinding.instance.addPostFrameCallback((_) {
      _loadData();
    });
  }

  @override
  void dispose() {
    _tabController.dispose();
    super.dispose();
  }

  void _loadData() {
    final provider = context.read<OrderStatusProvider>();
    provider.setPage(_currentPage);
    provider.fetchMappings();
    provider.fetchStatuses();
  }

  void _goToPage(int page) {
    setState(() => _currentPage = page);
    final provider = context.read<OrderStatusProvider>();
    provider.setPage(page);
    provider.fetchMappings();
  }

  Color _parseColor(String? colorStr) {
    if (colorStr == null || colorStr.isEmpty) return Colors.grey;
    try {
      final hex = colorStr.replaceFirst('#', '');
      if (hex.length == 6) {
        return Color(int.parse('FF$hex', radix: 16));
      }
    } catch (_) {}
    return Colors.grey;
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Estados de Orden'),
        bottom: TabBar(
          controller: _tabController,
          tabs: const [
            Tab(text: 'Mapeos'),
            Tab(text: 'Estados'),
          ],
        ),
      ),
      body: TabBarView(
        controller: _tabController,
        children: [
          _buildMappingsTab(),
          _buildStatusesTab(),
        ],
      ),
    );
  }

  Widget _buildMappingsTab() {
    return Consumer<OrderStatusProvider>(
      builder: (context, provider, _) {
        if (provider.isLoading) {
          return const Center(child: CircularProgressIndicator());
        }
        if (provider.error != null) {
          return _ErrorView(
              message: provider.error!, onRetry: _loadData);
        }
        if (provider.mappings.isEmpty) {
          return _EmptyView(
            icon: Icons.swap_horiz,
            message: 'No hay mapeos de estados',
          );
        }
        return RefreshIndicator(
          onRefresh: () async => _loadData(),
          child: Column(
            children: [
              Expanded(
                child: ListView.builder(
                  padding: const EdgeInsets.all(16),
                  itemCount: provider.mappings.length,
                  itemBuilder: (context, index) {
                    final mapping = provider.mappings[index];
                    return _MappingCard(
                      mapping: mapping,
                      parseColor: _parseColor,
                    );
                  },
                ),
              ),
              if (provider.pagination != null &&
                  provider.pagination!.lastPage > 1)
                _PaginationBar(
                  currentPage: provider.pagination!.currentPage,
                  totalPages: provider.pagination!.lastPage,
                  total: provider.pagination!.total,
                  onPageChanged: _goToPage,
                ),
            ],
          ),
        );
      },
    );
  }

  Widget _buildStatusesTab() {
    return Consumer<OrderStatusProvider>(
      builder: (context, provider, _) {
        if (provider.isLoading && provider.statuses.isEmpty) {
          return const Center(child: CircularProgressIndicator());
        }
        if (provider.statuses.isEmpty) {
          return _EmptyView(
            icon: Icons.flag_outlined,
            message: 'No hay estados definidos',
          );
        }
        return RefreshIndicator(
          onRefresh: () async {
            await provider.fetchStatuses();
          },
          child: ListView.builder(
            padding: const EdgeInsets.all(16),
            itemCount: provider.statuses.length,
            itemBuilder: (context, index) {
              final status = provider.statuses[index];
              final color = _parseColor(status.color);
              return Card(
                margin: const EdgeInsets.only(bottom: 8),
                child: ListTile(
                  leading: CircleAvatar(
                    backgroundColor: color.withValues(alpha: 0.15),
                    child: Icon(Icons.circle, color: color, size: 16),
                  ),
                  title: Text(status.name,
                      style: const TextStyle(fontWeight: FontWeight.w600)),
                  subtitle: Text(status.code),
                  trailing: Row(
                    mainAxisSize: MainAxisSize.min,
                    children: [
                      if (status.category != null)
                        _InfoChip(
                            label: status.category!, color: Colors.indigo),
                      const SizedBox(width: 4),
                      _InfoChip(
                        label: status.isActive == true ? 'Activo' : 'Inactivo',
                        color: status.isActive == true
                            ? Colors.green
                            : Colors.grey,
                      ),
                    ],
                  ),
                ),
              );
            },
          ),
        );
      },
    );
  }
}

class _MappingCard extends StatelessWidget {
  final OrderStatusMapping mapping;
  final Color Function(String?) parseColor;

  const _MappingCard({required this.mapping, required this.parseColor});

  @override
  Widget build(BuildContext context) {
    final statusColor = parseColor(mapping.orderStatus?.color);
    return Card(
      margin: const EdgeInsets.only(bottom: 8),
      child: Padding(
        padding: const EdgeInsets.all(12),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Row(
              children: [
                Expanded(
                  child: Text(
                    mapping.originalStatus,
                    style: const TextStyle(
                        fontWeight: FontWeight.w600, fontSize: 15),
                  ),
                ),
                Icon(Icons.arrow_forward, size: 16, color: Colors.grey.shade400),
                const SizedBox(width: 8),
                Container(
                  padding:
                      const EdgeInsets.symmetric(horizontal: 10, vertical: 4),
                  decoration: BoxDecoration(
                    color: statusColor.withValues(alpha: 0.12),
                    borderRadius: BorderRadius.circular(12),
                  ),
                  child: Text(
                    mapping.orderStatus?.name ?? 'N/A',
                    style: TextStyle(
                      color: statusColor,
                      fontWeight: FontWeight.w500,
                      fontSize: 13,
                    ),
                  ),
                ),
              ],
            ),
            const SizedBox(height: 8),
            Row(
              children: [
                if (mapping.integrationType != null)
                  _InfoChip(
                    label: mapping.integrationType!.name,
                    color: Colors.blue,
                  ),
                const SizedBox(width: 6),
                _InfoChip(
                  label: mapping.isActive ? 'Activo' : 'Inactivo',
                  color: mapping.isActive ? Colors.green : Colors.grey,
                ),
              ],
            ),
            if (mapping.description.isNotEmpty) ...[
              const SizedBox(height: 6),
              Text(
                mapping.description,
                style: TextStyle(fontSize: 13, color: Colors.grey.shade600),
                maxLines: 2,
                overflow: TextOverflow.ellipsis,
              ),
            ],
          ],
        ),
      ),
    );
  }
}

class _InfoChip extends StatelessWidget {
  final String label;
  final Color color;
  const _InfoChip({required this.label, required this.color});

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 2),
      decoration: BoxDecoration(
        color: color.withValues(alpha: 0.1),
        borderRadius: BorderRadius.circular(12),
      ),
      child: Text(label,
          style: TextStyle(
              fontSize: 11, color: color, fontWeight: FontWeight.w500)),
    );
  }
}

class _EmptyView extends StatelessWidget {
  final IconData icon;
  final String message;
  const _EmptyView({required this.icon, required this.message});

  @override
  Widget build(BuildContext context) {
    return Center(
      child: Column(
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          Icon(icon, size: 64, color: Colors.grey.shade400),
          const SizedBox(height: 16),
          Text(message,
              style: TextStyle(fontSize: 16, color: Colors.grey.shade600)),
        ],
      ),
    );
  }
}

class _ErrorView extends StatelessWidget {
  final String message;
  final VoidCallback onRetry;
  const _ErrorView({required this.message, required this.onRetry});

  @override
  Widget build(BuildContext context) {
    return Center(
      child: Padding(
        padding: const EdgeInsets.all(24),
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            const Icon(Icons.error_outline, size: 48, color: Colors.red),
            const SizedBox(height: 16),
            Text(message,
                textAlign: TextAlign.center,
                style: const TextStyle(color: Colors.red)),
            const SizedBox(height: 16),
            FilledButton.icon(
              onPressed: onRetry,
              icon: const Icon(Icons.refresh),
              label: const Text('Reintentar'),
            ),
          ],
        ),
      ),
    );
  }
}

class _PaginationBar extends StatelessWidget {
  final int currentPage;
  final int totalPages;
  final int total;
  final ValueChanged<int> onPageChanged;

  const _PaginationBar({
    required this.currentPage,
    required this.totalPages,
    required this.total,
    required this.onPageChanged,
  });

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 12),
      decoration: BoxDecoration(
        color: Theme.of(context).colorScheme.surface,
        border: Border(top: BorderSide(color: Colors.grey.shade300)),
      ),
      child: Row(
        mainAxisAlignment: MainAxisAlignment.spaceBetween,
        children: [
          Text('$total resultados',
              style: TextStyle(fontSize: 13, color: Colors.grey.shade600)),
          Row(
            children: [
              IconButton(
                icon: const Icon(Icons.chevron_left),
                onPressed: currentPage > 1
                    ? () => onPageChanged(currentPage - 1)
                    : null,
                iconSize: 20,
                visualDensity: VisualDensity.compact,
              ),
              Text('$currentPage / $totalPages',
                  style: const TextStyle(fontSize: 13)),
              IconButton(
                icon: const Icon(Icons.chevron_right),
                onPressed: currentPage < totalPages
                    ? () => onPageChanged(currentPage + 1)
                    : null,
                iconSize: 20,
                visualDensity: VisualDensity.compact,
              ),
            ],
          ),
        ],
      ),
    );
  }
}
