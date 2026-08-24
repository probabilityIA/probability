import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:provider/provider.dart';
import '../../../../../shared/theme/app_colors.dart';
import '../../../../../shared/theme/app_tokens.dart';
import '../../../../../shared/utils/formatters.dart';
import '../../../../../shared/widgets/ui/ui.dart';
import '../../domain/entities.dart';
import '../providers/invoicing_provider.dart';
import 'invoice_list_screen.dart';

class InvoiceDetailScreen extends StatefulWidget {
  const InvoiceDetailScreen({super.key, required this.invoiceId});

  final int invoiceId;

  @override
  State<InvoiceDetailScreen> createState() => _InvoiceDetailScreenState();
}

class _InvoiceDetailScreenState extends State<InvoiceDetailScreen> {
  Invoice? _invoice;
  bool _busy = false;

  @override
  void initState() {
    super.initState();
    _invoice = context.read<InvoicingProvider>().invoiceById(widget.invoiceId);
  }

  void _copyCufe() {
    final cufe = _invoice?.cufe;
    if (cufe == null || cufe.isEmpty) return;
    Clipboard.setData(ClipboardData(text: cufe));
    ScaffoldMessenger.of(context).showSnackBar(
      const SnackBar(content: Text('CUFE copiado')),
    );
  }

  Future<void> _retry() async {
    setState(() => _busy = true);
    final provider = context.read<InvoicingProvider>();
    final ok = await provider.retryInvoice(widget.invoiceId);
    if (!mounted) return;
    setState(() {
      _busy = false;
      if (ok) _invoice = provider.invoiceById(widget.invoiceId) ?? _invoice;
    });
    ScaffoldMessenger.of(context).showSnackBar(
      SnackBar(content: Text(ok ? 'Reintento encolado' : (provider.error ?? 'No se pudo reintentar'))),
    );
  }

  Future<void> _cancel() async {
    final confirmed = await showDialog<bool>(
      context: context,
      builder: (dialogContext) => AlertDialog(
        title: const Text('Cancelar factura'),
        content: const Text(
          'La factura se anula ante la DIAN. Esta accion no se puede deshacer.',
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(dialogContext, false),
            child: const Text('Volver'),
          ),
          FilledButton(
            style: FilledButton.styleFrom(backgroundColor: AppColors.error),
            onPressed: () => Navigator.pop(dialogContext, true),
            child: const Text('Cancelar factura'),
          ),
        ],
      ),
    );

    if (confirmed != true || !mounted) return;

    setState(() => _busy = true);
    final provider = context.read<InvoicingProvider>();
    final ok = await provider.cancelInvoice(widget.invoiceId);
    if (!mounted) return;
    setState(() {
      _busy = false;
      if (ok) _invoice = provider.invoiceById(widget.invoiceId) ?? _invoice;
    });
    ScaffoldMessenger.of(context).showSnackBar(
      SnackBar(content: Text(ok ? 'Factura cancelada' : (provider.error ?? 'No se pudo cancelar'))),
    );
  }

  @override
  Widget build(BuildContext context) {
    final invoice = _invoice;

    return AppScaffold(
      title: invoice?.invoiceNumber ?? 'Factura',
      subtitle: invoice?.providerName,
      onBack: () => Navigator.of(context).pop(),
      actions: [
        if ((invoice?.cufe ?? '').isNotEmpty)
          IconButton(
            icon: const Icon(Icons.copy_rounded, size: 20),
            tooltip: 'Copiar CUFE',
            onPressed: _copyCufe,
          ),
      ],
      body: invoice == null
          ? const AppEmptyState(
              icon: Icons.description_outlined,
              title: 'Factura no encontrada',
            )
          : ListView(
              padding: AppSpacing.page,
              children: [
                _Header(invoice: invoice),
                const SizedBox(height: 18),
                const AppSectionHeader(title: 'Totales'),
                _TotalsCard(invoice: invoice),
                const SizedBox(height: 18),
                const AppSectionHeader(title: 'Cliente'),
                _CustomerCard(invoice: invoice),
                if (invoice.items != null && invoice.items!.isNotEmpty) ...[
                  const SizedBox(height: 18),
                  AppSectionHeader(
                    title: 'Items facturados',
                    subtitle: '${invoice.items!.length} lineas',
                  ),
                  _ItemsCard(items: invoice.items!),
                ],
                const SizedBox(height: 18),
                const AppSectionHeader(title: 'DIAN'),
                _DianCard(invoice: invoice),
                const SizedBox(height: 20),
                if (invoice.status == 'failed')
                  FilledButton.icon(
                    onPressed: _busy ? null : _retry,
                    style: FilledButton.styleFrom(minimumSize: const Size(0, 50)),
                    icon: const Icon(Icons.refresh_rounded, size: 18),
                    label: const Text('Reintentar emision'),
                  ),
                if (invoice.status == 'issued') ...[
                  const SizedBox(height: 10),
                  OutlinedButton.icon(
                    onPressed: _busy ? null : _cancel,
                    style: OutlinedButton.styleFrom(
                      foregroundColor: AppColors.error,
                      side: const BorderSide(color: AppColors.error),
                    ),
                    icon: const Icon(Icons.cancel_outlined, size: 18),
                    label: const Text('Cancelar factura'),
                  ),
                ],
                const SizedBox(height: 28),
              ],
            ),
    );
  }
}

class _Header extends StatelessWidget {
  const _Header({required this.invoice});

  final Invoice invoice;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    return AppCard(
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(
            children: [
              BrandLogo(
                name: invoice.providerName ?? 'Facturador',
                imageUrl: invoice.providerLogoUrl,
                size: 46,
                radius: 12,
                padding: 7,
              ),
              const SizedBox(width: 13),
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(invoice.invoiceNumber, style: theme.textTheme.titleLarge),
                    const SizedBox(height: 2),
                    Text(
                      'Orden ${invoice.orderNumber ?? '-'}',
                      style: theme.textTheme.bodySmall,
                    ),
                  ],
                ),
              ),
            ],
          ),
          const SizedBox(height: 14),
          AppStatusChip(
            label: invoiceStatusLabel(invoice.status),
            tone: AppStatusChip.toneFromCode(invoice.status),
          ),
          if ((invoice.errorMessage ?? '').isNotEmpty) ...[
            const SizedBox(height: 12),
            Container(
              padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 11),
              decoration: BoxDecoration(
                color: AppColors.errorSoft,
                borderRadius: AppRadius.mdAll,
              ),
              child: Row(
                children: [
                  const Icon(Icons.error_outline, size: 17, color: AppColors.error),
                  const SizedBox(width: 9),
                  Expanded(
                    child: Text(
                      invoice.errorMessage!,
                      style: theme.textTheme.bodySmall?.copyWith(
                        color: const Color(0xFFB91C1C),
                        fontWeight: FontWeight.w500,
                      ),
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

class _TotalsCard extends StatelessWidget {
  const _TotalsCard({required this.invoice});

  final Invoice invoice;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    return AppCard(
      child: Column(
        children: [
          AppKeyValueRow(label: 'Subtotal', value: AppFormat.money(invoice.subtotal), dense: true),
          if (invoice.discount > 0)
            AppKeyValueRow(
              label: 'Descuento',
              value: '- ${AppFormat.money(invoice.discount)}',
              dense: true,
            ),
          AppKeyValueRow(label: 'IVA', value: AppFormat.money(invoice.tax), dense: true),
          const Divider(height: 18),
          Row(
            children: [
              Expanded(child: Text('Total facturado', style: theme.textTheme.titleMedium)),
              Text(
                AppFormat.money(invoice.totalAmount),
                style: theme.textTheme.titleLarge?.copyWith(color: AppColors.primary),
              ),
            ],
          ),
        ],
      ),
    );
  }
}

class _CustomerCard extends StatelessWidget {
  const _CustomerCard({required this.invoice});

  final Invoice invoice;

  @override
  Widget build(BuildContext context) {
    return AppCard(
      child: Column(
        children: [
          AppKeyValueRow(label: 'Nombre', value: invoice.customerName),
          const Divider(height: 18),
          AppKeyValueRow(label: 'Documento', value: invoice.customerDni ?? '-'),
          const Divider(height: 18),
          AppKeyValueRow(label: 'Correo', value: invoice.customerEmail ?? '-'),
        ],
      ),
    );
  }
}

class _ItemsCard extends StatelessWidget {
  const _ItemsCard({required this.items});

  final List<InvoiceItem> items;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    return AppCard(
      child: Column(
        children: [
          for (var i = 0; i < items.length; i++) ...[
            if (i > 0) const Divider(height: 20),
            Row(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Container(
                  width: 40,
                  height: 40,
                  alignment: Alignment.center,
                  decoration: BoxDecoration(
                    color: AppColors.surfaceMuted,
                    borderRadius: AppRadius.smAll,
                  ),
                  child: Text(
                    '${items[i].quantity}x',
                    style: theme.textTheme.labelMedium?.copyWith(
                      fontWeight: FontWeight.w700,
                      color: AppColors.textSecondary,
                    ),
                  ),
                ),
                const SizedBox(width: 12),
                Expanded(
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Text(
                        items[i].productName,
                        style: theme.textTheme.titleSmall,
                        maxLines: 2,
                        overflow: TextOverflow.ellipsis,
                      ),
                      const SizedBox(height: 3),
                      Text(
                        'SKU ${items[i].productSku}  ·  IVA ${(items[i].taxRate ?? 0).toStringAsFixed(0)}%',
                        style: theme.textTheme.labelSmall,
                      ),
                    ],
                  ),
                ),
                const SizedBox(width: 10),
                Text(
                  AppFormat.money(items[i].totalPrice),
                  style: theme.textTheme.titleSmall?.copyWith(fontSize: 14),
                ),
              ],
            ),
          ],
        ],
      ),
    );
  }
}

class _DianCard extends StatelessWidget {
  const _DianCard({required this.invoice});

  final Invoice invoice;

  @override
  Widget build(BuildContext context) {
    return AppCard(
      child: Column(
        children: [
          AppKeyValueRow(
            label: 'CUFE',
            value: (invoice.cufe ?? '').isEmpty ? 'Sin emitir' : invoice.cufe!,
          ),
          const Divider(height: 18),
          AppKeyValueRow(
            label: 'Emitida',
            value: AppFormat.dateTime(AppFormat.parseDate(invoice.issuedAt)),
          ),
          const Divider(height: 18),
          AppKeyValueRow(
            label: 'Proveedor',
            value: invoice.providerName ?? '-',
          ),
          const Divider(height: 18),
          AppKeyValueRow(
            label: 'PDF',
            value: (invoice.pdfUrl ?? '').isEmpty ? 'No disponible' : 'Disponible',
          ),
        ],
      ),
    );
  }
}
