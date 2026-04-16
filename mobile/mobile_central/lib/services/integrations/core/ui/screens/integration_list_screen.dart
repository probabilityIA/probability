import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../../domain/entities.dart';
import '../providers/integration_provider.dart';

class IntegrationListScreen extends StatefulWidget {
  const IntegrationListScreen({super.key});

  @override
  State<IntegrationListScreen> createState() => _IntegrationListScreenState();
}

class _IntegrationListScreenState extends State<IntegrationListScreen> {
  int _currentPage = 1;
  String? _selectedCategory;

  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addPostFrameCallback((_) {
      _loadData();
    });
  }

  void _loadData() {
    final provider = context.read<IntegrationProvider>();
    provider.setPage(_currentPage);
    provider.fetchIntegrations();
    provider.fetchIntegrationCategories();
  }

  void _goToPage(int page) {
    setState(() => _currentPage = page);
    final provider = context.read<IntegrationProvider>();
    provider.setPage(page);
    provider.fetchIntegrations();
  }

  void _onCategoryChanged(String? category) {
    setState(() {
      _selectedCategory = category;
      _currentPage = 1;
    });
    final provider = context.read<IntegrationProvider>();
    provider.setFilters(category: category);
    provider.setPage(1);
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
      appBar: AppBar(title: const Text('Catalogo de Integraciones')),
      body: Consumer<IntegrationProvider>(
        builder: (context, provider, _) {
          if (provider.isLoading && provider.integrations.isEmpty) {
            return const Center(child: CircularProgressIndicator());
          }
          if (provider.error != null && provider.integrations.isEmpty) {
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
                      onPressed: _loadData,
                      icon: const Icon(Icons.refresh),
                      label: const Text('Reintentar'),
                    ),
                  ],
                ),
              ),
            );
          }
          return RefreshIndicator(
            onRefresh: () async => _loadData(),
            child: Column(
              children: [
                // Category filter chips
                if (provider.integrationCategories.isNotEmpty)
                  SingleChildScrollView(
                    scrollDirection: Axis.horizontal,
                    padding:
                        const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
                    child: Row(
                      children: [
                        Padding(
                          padding: const EdgeInsets.only(right: 8),
                          child: FilterChip(
                            label: const Text('Todos'),
                            selected: _selectedCategory == null,
                            onSelected: (_) => _onCategoryChanged(null),
                          ),
                        ),
                        ...provider.integrationCategories.map((cat) {
                          return Padding(
                            padding: const EdgeInsets.only(right: 8),
                            child: FilterChip(
                              avatar: Icon(
                                _categoryIcon(cat.code),
                                size: 16,
                                color: _selectedCategory == cat.code
                                    ? null
                                    : _categoryColor(cat.code),
                              ),
                              label: Text(cat.name),
                              selected: _selectedCategory == cat.code,
                              onSelected: (_) =>
                                  _onCategoryChanged(cat.code),
                            ),
                          );
                        }),
                      ],
                    ),
                  ),
                Expanded(
                  child: provider.integrations.isEmpty
                      ? Center(
                          child: Column(
                            mainAxisAlignment: MainAxisAlignment.center,
                            children: [
                              Icon(Icons.integration_instructions_outlined,
                                  size: 64, color: Colors.grey.shade400),
                              const SizedBox(height: 16),
                              Text('No hay integraciones disponibles',
                                  style: TextStyle(
                                      fontSize: 16,
                                      color: Colors.grey.shade600)),
                            ],
                          ),
                        )
                      : ListView.builder(
                          padding: const EdgeInsets.symmetric(horizontal: 16),
                          itemCount: provider.integrations.length,
                          itemBuilder: (context, index) {
                            final integration = provider.integrations[index];
                            return _IntegrationCatalogCard(
                              integration: integration,
                              icon: _categoryIcon(integration.category),
                              color: _categoryColor(integration.category),
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

class _IntegrationCatalogCard extends StatelessWidget {
  final Integration integration;
  final IconData icon;
  final Color color;

  const _IntegrationCatalogCard({
    required this.integration,
    required this.icon,
    required this.color,
  });

  @override
  Widget build(BuildContext context) {
    return Card(
      margin: const EdgeInsets.only(bottom: 10),
      child: Padding(
        padding: const EdgeInsets.all(12),
        child: Row(
          children: [
            if (integration.integrationType?.imageUrl != null &&
                integration.integrationType!.imageUrl!.isNotEmpty)
              ClipRRect(
                borderRadius: BorderRadius.circular(8),
                child: Image.network(
                  integration.integrationType!.imageUrl!,
                  width: 44,
                  height: 44,
                  fit: BoxFit.contain,
                  errorBuilder: (context, error, stackTrace) => CircleAvatar(
                    backgroundColor: color.withValues(alpha: 0.12),
                    child: Icon(icon, color: color),
                  ),
                ),
              )
            else
              CircleAvatar(
                backgroundColor: color.withValues(alpha: 0.12),
                child: Icon(icon, color: color),
              ),
            const SizedBox(width: 12),
            Expanded(
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(integration.name,
                      style: const TextStyle(
                          fontWeight: FontWeight.w600, fontSize: 14)),
                  const SizedBox(height: 2),
                  Row(
                    children: [
                      Container(
                        padding: const EdgeInsets.symmetric(
                            horizontal: 8, vertical: 2),
                        decoration: BoxDecoration(
                          color: color.withValues(alpha: 0.1),
                          borderRadius: BorderRadius.circular(12),
                        ),
                        child: Text(
                          integration.categoryName ?? integration.category,
                          style: TextStyle(
                              fontSize: 10,
                              color: color,
                              fontWeight: FontWeight.w500),
                        ),
                      ),
                      const SizedBox(width: 6),
                      Container(
                        padding: const EdgeInsets.symmetric(
                            horizontal: 8, vertical: 2),
                        decoration: BoxDecoration(
                          color: integration.isActive
                              ? Colors.green.withValues(alpha: 0.1)
                              : Colors.grey.withValues(alpha: 0.1),
                          borderRadius: BorderRadius.circular(12),
                        ),
                        child: Text(
                          integration.isActive ? 'Activa' : 'Inactiva',
                          style: TextStyle(
                            fontSize: 10,
                            color: integration.isActive
                                ? Colors.green
                                : Colors.grey,
                            fontWeight: FontWeight.w500,
                          ),
                        ),
                      ),
                    ],
                  ),
                  if (integration.description != null &&
                      integration.description!.isNotEmpty) ...[
                    const SizedBox(height: 4),
                    Text(integration.description!,
                        style: TextStyle(
                            fontSize: 12, color: Colors.grey.shade600),
                        maxLines: 2,
                        overflow: TextOverflow.ellipsis),
                  ],
                ],
              ),
            ),
            const Icon(Icons.chevron_right, color: Colors.grey),
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
