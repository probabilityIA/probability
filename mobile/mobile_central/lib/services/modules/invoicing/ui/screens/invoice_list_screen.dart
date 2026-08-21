import 'dart:async';
import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../../../../../shared/theme/app_colors.dart';
import '../../../../../shared/theme/app_tokens.dart';
import '../../../../../shared/utils/formatters.dart';
import '../../../../../shared/widgets/ui/ui.dart';
import '../../domain/entities.dart';
import '../providers/invoicing_provider.dart';
import 'invoice_detail_screen.dart';

const Map<String, String> invoiceStatusLabels = {
  'issued': 'Emitida',
  'pending': 'Pendiente',
  'failed': 'Fallida',
  'cancelled': 'Cancelada',
  'draft': 'Borrador',
};

String invoiceStatusLabel(String code) => invoiceStatusLabels[code] ?? code;

class InvoiceListScreen extends StatefulWidget {
  const InvoiceListScreen({super.key, this.businessId});

  final int? businessId;

  @override
  State<InvoiceListScreen> createState() => _InvoiceListScreenState();
}

class _InvoiceListScreenState extends State<InvoiceListScreen> {
  final _searchController = TextEditingController();
  Timer? _debounce;
  String _status = '';

  static const List<({String value, String label})> _statusOptions = [
    (value: '', label: 'Todas'),
    (value: 'issued', label: 'Emitidas'),
    (value: 'pending', label: 'Pendientes'),
    (value: 'failed', label: 'Fallidas'),
    (value: 'cancelled', label: 'Canceladas'),
  ];

  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addPostFrameCallback((_) => _refresh());
  }

  @override
  void dispose() {
    _debounce?.cancel();
    _searchController.dispose();
    super.dispose();
  }

  void _refresh() {
    final provider = context.read<InvoicingProvider>();
    provider.setFilters(
      status: _status.isEmpty ? null : _status,
      invoiceNumber: _searchController.text.trim().isEmpty
          ? null
          : _searchController.text.trim(),
    );
    provider.fetchInvoices(businessId: widget.businessId);
  }

  void _onSearch(String value) {
    _debounce?.cancel();
    _debounce = Timer(const Duration(milliseconds: 400), () {
      if (mounted) _refresh();
    });
  }

  @override
  Widget build(BuildContext context) {
    return Column(
      children: [
        Padding(
          padding: const EdgeInsets.fromLTRB(16, 14, 16, 10),
          child: AppSearchField(
            controller: _searchController,
            hintText: 'Numero de factura, orden o cliente',
            onChanged: _onSearch,
          ),
        ),
        AppFilterChips(
          options: _statusOptions,
          selected: _status,
          onSelected: (value) {
            setState(() => _status = value);
            _refresh();
          },
        ),
        const SizedBox(height: 4),
        Expanded(
          child: Consumer<InvoicingProvider>(
            builder: (context, provider, _) {
              if (provider.isLoading && provider.invoices.isEmpty) {
                return const AppListSkeleton();
              }
              if (provider.error != null && provider.invoices.isEmpty) {
                return AppErrorState(message: provider.error!, onRetry: _refresh);
              }
              if (provider.invoices.isEmpty) {
                return AppEmptyState(
                  icon: Icons.description_outlined,
                  title: 'Sin facturas',
                  message: 'Cuando factures una orden vas a ver aqui el documento y su estado en la DIAN.',
                  actionLabel: 'Actualizar',
                  onAction: _refresh,
                );
              }

              return RefreshIndicator(
                onRefresh: () async => _refresh(),
                color: AppColors.primary,
                child: ListView.separated(
                  physics: const AlwaysScrollableScrollPhysics(),
                  padding: AppSpacing.page,
                  itemCount: provider.invoices.length,
                  separatorBuilder: (context, index) => const SizedBox(height: 10),
                  itemBuilder: (context, index) {
                    final invoice = provider.invoices[index];
                    return InvoiceCard(
                      invoice: invoice,
                      onTap: () => Navigator.of(context).push(
                        MaterialPageRoute(
                          builder: (_) => InvoiceDetailScreen(invoiceId: invoice.id),
                        ),
                      ),
                    );
                  },
                ),
              );
            },
          ),
        ),
      ],
    );
  }
}

class InvoiceCard extends StatelessWidget {
  const InvoiceCard({super.key, required this.invoice, this.onTap});

  final Invoice invoice;
  final VoidCallback? onTap;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    return AppCard(
      padding: const EdgeInsets.all(14),
      onTap: onTap,
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              BrandLogo(
                name: invoice.providerName ?? 'Facturador',
                imageUrl: invoice.providerLogoUrl,
                size: 40,
                radius: 11,
                padding: 6,
              ),
              const SizedBox(width: 12),
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(
                      invoice.invoiceNumber,
                      style: theme.textTheme.titleSmall,
                      maxLines: 1,
                      overflow: TextOverflow.ellipsis,
                    ),
                    const SizedBox(height: 2),
                    Text(
                      '${invoice.providerName ?? ''}  ·  ${invoice.orderNumber ?? ''}',
                      style: theme.textTheme.labelSmall,
                      maxLines: 1,
                      overflow: TextOverflow.ellipsis,
                    ),
                  ],
                ),
              ),
              const SizedBox(width: 8),
              Column(
                crossAxisAlignment: CrossAxisAlignment.end,
                children: [
                  Text(
                    AppFormat.money(invoice.totalAmount),
                    style: theme.textTheme.titleSmall?.copyWith(fontSize: 15),
                  ),
                  const SizedBox(height: 3),
                  Text(
                    AppFormat.relative(AppFormat.parseDate(invoice.createdAt)),
                    style: theme.textTheme.labelSmall,
                  ),
                ],
              ),
            ],
          ),
          const SizedBox(height: 11),
          Row(
            children: [
              AppStatusChip(
                dense: true,
                label: invoiceStatusLabel(invoice.status),
                tone: AppStatusChip.toneFromCode(invoice.status),
              ),
              const SizedBox(width: 8),
              Expanded(
                child: Text(
                  invoice.customerName,
                  style: theme.textTheme.labelSmall,
                  maxLines: 1,
                  overflow: TextOverflow.ellipsis,
                  textAlign: TextAlign.right,
                ),
              ),
            ],
          ),
          if ((invoice.errorMessage ?? '').isNotEmpty) ...[
            const SizedBox(height: 10),
            Container(
              padding: const EdgeInsets.symmetric(horizontal: 11, vertical: 9),
              decoration: BoxDecoration(
                color: AppColors.errorSoft,
                borderRadius: AppRadius.smAll,
              ),
              child: Row(
                children: [
                  const Icon(Icons.error_outline, size: 14, color: AppColors.error),
                  const SizedBox(width: 8),
                  Expanded(
                    child: Text(
                      invoice.errorMessage!,
                      style: theme.textTheme.labelSmall?.copyWith(
                        color: const Color(0xFFB91C1C),
                        fontWeight: FontWeight.w500,
                      ),
                      maxLines: 2,
                      overflow: TextOverflow.ellipsis,
                    ),
                  ),
                ],
              ),
            ),
          ],
        ],
      ),
    );
  }
}
