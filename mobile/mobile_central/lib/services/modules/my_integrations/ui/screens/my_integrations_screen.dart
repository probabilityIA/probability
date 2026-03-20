import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../../domain/entities.dart';
import '../providers/my_integrations_provider.dart';

class MyIntegrationsScreen extends StatefulWidget {
  final int? businessId;

  const MyIntegrationsScreen({super.key, this.businessId});

  @override
  State<MyIntegrationsScreen> createState() => _MyIntegrationsScreenState();
}

class _MyIntegrationsScreenState extends State<MyIntegrationsScreen> {
  int _currentPage = 1;

  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addPostFrameCallback((_) {
      _loadIntegrations();
    });
  }

  void _loadIntegrations() {
    final provider = context.read<MyIntegrationsProvider>();
    provider.setPage(_currentPage);
    provider.fetchIntegrations();
  }

  void _goToPage(int page) {
    setState(() => _currentPage = page);
    final provider = context.read<MyIntegrationsProvider>();
    provider.setPage(page);
    provider.fetchIntegrations();
  }

  IconData _categoryIcon(String? category) {
    switch (category?.toLowerCase()) {
      case 'platform':
        return Icons.extension;
      case 'ecommerce':
        return Icons.shopping_cart_outlined;
      case 'invoicing':
        return Icons.receipt_long_outlined;
      case 'messaging':
        return Icons.chat_outlined;
      case 'payment':
        return Icons.credit_card;
      case 'shipping':
        return Icons.local_shipping_outlined;
      default:
        return Icons.integration_instructions;
    }
  }

  Color _categoryColor(String? category) {
    switch (category?.toLowerCase()) {
      case 'platform':
        return Colors.purple;
      case 'ecommerce':
        return Colors.blue;
      case 'invoicing':
        return Colors.teal;
      case 'messaging':
        return Colors.green;
      case 'payment':
        return Colors.orange;
      case 'shipping':
        return Colors.indigo;
      default:
        return Colors.grey;
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('Mis Integraciones')),
      body: Consumer<MyIntegrationsProvider>(
        builder: (context, provider, _) {
          if (provider.isLoading) {
            return const Center(child: CircularProgressIndicator());
          }
          if (provider.error != null) {
            return Center(
              child: Padding(
                padding: const EdgeInsets.all(24),
                child: Column(
                  mainAxisAlignment: MainAxisAlignment.center,
                  children: [
                    const Icon(Icons.error_outline,
                        size: 48, color: Colors.red),
                    const SizedBox(height: 16),
                    Text(provider.error!,
                        textAlign: TextAlign.center,
                        style: const TextStyle(color: Colors.red)),
                    const SizedBox(height: 16),
                    FilledButton.icon(
                      onPressed: _loadIntegrations,
                      icon: const Icon(Icons.refresh),
                      label: const Text('Reintentar'),
                    ),
                  ],
                ),
              ),
            );
          }
          if (provider.integrations.isEmpty) {
            return Center(
              child: Column(
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  Icon(Icons.integration_instructions_outlined,
                      size: 64, color: Colors.grey.shade400),
                  const SizedBox(height: 16),
                  Text('No hay integraciones conectadas',
                      style:
                          TextStyle(fontSize: 16, color: Colors.grey.shade600)),
                ],
              ),
            );
          }
          return RefreshIndicator(
            onRefresh: () async => _loadIntegrations(),
            child: Column(
              children: [
                Expanded(
                  child: ListView.builder(
                    padding: const EdgeInsets.all(16),
                    itemCount: provider.integrations.length,
                    itemBuilder: (context, index) {
                      final integration = provider.integrations[index];
                      return _IntegrationCard(
                        integration: integration,
                        icon: _categoryIcon(integration.categoryCode),
                        color: _categoryColor(integration.categoryCode),
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
      ),
    );
  }
}

class _IntegrationCard extends StatelessWidget {
  final MyIntegration integration;
  final IconData icon;
  final Color color;

  const _IntegrationCard({
    required this.integration,
    required this.icon,
    required this.color,
  });

  @override
  Widget build(BuildContext context) {
    return Card(
      margin: const EdgeInsets.only(bottom: 10),
      child: ListTile(
        leading: CircleAvatar(
          backgroundColor: color.withValues(alpha: 0.12),
          child: Icon(icon, color: color),
        ),
        title: Text(integration.name,
            style: const TextStyle(fontWeight: FontWeight.w600)),
        subtitle: Row(
          children: [
            if (integration.integrationTypeName != null)
              Flexible(
                child: Text(integration.integrationTypeName!,
                    style:
                        TextStyle(fontSize: 13, color: Colors.grey.shade600),
                    overflow: TextOverflow.ellipsis),
              ),
            if (integration.categoryCode != null) ...[
              const SizedBox(width: 8),
              Container(
                padding:
                    const EdgeInsets.symmetric(horizontal: 8, vertical: 2),
                decoration: BoxDecoration(
                  color: color.withValues(alpha: 0.1),
                  borderRadius: BorderRadius.circular(12),
                ),
                child: Text(
                  integration.categoryCode!,
                  style: TextStyle(
                      fontSize: 10,
                      color: color,
                      fontWeight: FontWeight.w500),
                ),
              ),
            ],
          ],
        ),
        trailing: Container(
          padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 4),
          decoration: BoxDecoration(
            color: integration.isActive
                ? Colors.green.withValues(alpha: 0.1)
                : Colors.grey.withValues(alpha: 0.1),
            borderRadius: BorderRadius.circular(12),
          ),
          child: Text(
            integration.isActive ? 'Activa' : 'Inactiva',
            style: TextStyle(
              fontSize: 12,
              color: integration.isActive ? Colors.green : Colors.grey,
              fontWeight: FontWeight.w500,
            ),
          ),
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
