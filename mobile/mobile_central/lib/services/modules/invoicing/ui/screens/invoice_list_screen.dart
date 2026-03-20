import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../../domain/entities.dart';
import '../providers/invoicing_provider.dart';

class InvoiceListScreen extends StatefulWidget {
  final int? businessId;

  const InvoiceListScreen({super.key, this.businessId});

  @override
  State<InvoiceListScreen> createState() => _InvoiceListScreenState();
}

class _InvoiceListScreenState extends State<InvoiceListScreen> {
  int _currentPage = 1;
  String? _statusFilter;

  static const _statusOptions = [
    null,
    'pending',
    'completed',
    'failed',
    'cancelled',
  ];

  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addPostFrameCallback((_) {
      _loadInvoices();
    });
  }

  void _loadInvoices() {
    final provider = context.read<InvoicingProvider>();
    provider.setPage(_currentPage);
    provider.setFilters(status: _statusFilter);
    provider.fetchInvoices();
  }

  void _goToPage(int page) {
    setState(() => _currentPage = page);
    final provider = context.read<InvoicingProvider>();
    provider.setPage(page);
    provider.fetchInvoices();
  }

  void _onStatusChanged(String? status) {
    setState(() {
      _statusFilter = status;
      _currentPage = 1;
    });
    _loadInvoices();
  }

  Color _statusColor(String status) {
    switch (status.toLowerCase()) {
      case 'completed':
      case 'success':
        return Colors.green;
      case 'pending':
      case 'processing':
        return Colors.orange;
      case 'failed':
      case 'error':
        return Colors.red;
      case 'cancelled':
        return Colors.grey;
      default:
        return Colors.blue;
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('Facturas')),
      body: Consumer<InvoicingProvider>(
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
                      onPressed: _loadInvoices,
                      icon: const Icon(Icons.refresh),
                      label: const Text('Reintentar'),
                    ),
                  ],
                ),
              ),
            );
          }
          return RefreshIndicator(
            onRefresh: () async => _loadInvoices(),
            child: Column(
              children: [
                // Status filter chips
                SingleChildScrollView(
                  scrollDirection: Axis.horizontal,
                  padding:
                      const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
                  child: Row(
                    children: _statusOptions.map((status) {
                      final isSelected = _statusFilter == status;
                      return Padding(
                        padding: const EdgeInsets.only(right: 8),
                        child: FilterChip(
                          label: Text(status ?? 'Todos'),
                          selected: isSelected,
                          onSelected: (_) => _onStatusChanged(status),
                        ),
                      );
                    }).toList(),
                  ),
                ),
                Expanded(
                  child: provider.invoices.isEmpty
                      ? Center(
                          child: Column(
                            mainAxisAlignment: MainAxisAlignment.center,
                            children: [
                              Icon(Icons.receipt_long_outlined,
                                  size: 64, color: Colors.grey.shade400),
                              const SizedBox(height: 16),
                              Text('No hay facturas',
                                  style: TextStyle(
                                      fontSize: 16,
                                      color: Colors.grey.shade600)),
                            ],
                          ),
                        )
                      : ListView.builder(
                          padding: const EdgeInsets.symmetric(horizontal: 16),
                          itemCount: provider.invoices.length,
                          itemBuilder: (context, index) {
                            final invoice = provider.invoices[index];
                            return _InvoiceCard(
                              invoice: invoice,
                              statusColor: _statusColor(invoice.status),
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

class _InvoiceCard extends StatelessWidget {
  final Invoice invoice;
  final Color statusColor;

  const _InvoiceCard({required this.invoice, required this.statusColor});

  @override
  Widget build(BuildContext context) {
    return Card(
      margin: const EdgeInsets.only(bottom: 10),
      child: Padding(
        padding: const EdgeInsets.all(12),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Row(
              children: [
                const Icon(Icons.receipt_outlined, size: 18, color: Colors.blue),
                const SizedBox(width: 8),
                Expanded(
                  child: Text(
                    invoice.invoiceNumber.isNotEmpty
                        ? invoice.invoiceNumber
                        : 'Sin numero',
                    style: const TextStyle(
                        fontWeight: FontWeight.w600, fontSize: 15),
                  ),
                ),
                Container(
                  padding:
                      const EdgeInsets.symmetric(horizontal: 10, vertical: 4),
                  decoration: BoxDecoration(
                    color: statusColor.withValues(alpha: 0.12),
                    borderRadius: BorderRadius.circular(12),
                  ),
                  child: Text(
                    invoice.status,
                    style: TextStyle(
                      color: statusColor,
                      fontWeight: FontWeight.w500,
                      fontSize: 12,
                    ),
                  ),
                ),
              ],
            ),
            const SizedBox(height: 8),
            Row(
              children: [
                Icon(Icons.person_outline,
                    size: 14, color: Colors.grey.shade500),
                const SizedBox(width: 4),
                Expanded(
                  child: Text(invoice.customerName,
                      style:
                          TextStyle(fontSize: 13, color: Colors.grey.shade600),
                      overflow: TextOverflow.ellipsis),
                ),
              ],
            ),
            if (invoice.orderNumber != null) ...[
              const SizedBox(height: 4),
              Row(
                children: [
                  Icon(Icons.tag, size: 14, color: Colors.grey.shade500),
                  const SizedBox(width: 4),
                  Text('Orden: ${invoice.orderNumber}',
                      style:
                          TextStyle(fontSize: 13, color: Colors.grey.shade600)),
                ],
              ),
            ],
            const SizedBox(height: 8),
            Row(
              mainAxisAlignment: MainAxisAlignment.spaceBetween,
              children: [
                Text(
                  '\$${invoice.totalAmount.toStringAsFixed(0)} ${invoice.currency}',
                  style: const TextStyle(
                      fontWeight: FontWeight.bold, fontSize: 16),
                ),
                Text(
                  invoice.createdAt.length > 10
                      ? invoice.createdAt.substring(0, 10)
                      : invoice.createdAt,
                  style:
                      TextStyle(fontSize: 12, color: Colors.grey.shade500),
                ),
              ],
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
