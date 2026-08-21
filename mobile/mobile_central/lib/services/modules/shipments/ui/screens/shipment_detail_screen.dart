import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:provider/provider.dart';
import '../../../../../shared/theme/app_colors.dart';
import '../../../../../shared/theme/app_tokens.dart';
import '../../../../../shared/utils/formatters.dart';
import '../../../../../shared/widgets/ui/ui.dart';
import '../../domain/entities.dart';
import '../providers/shipment_provider.dart';
import '../widgets/shipment_card.dart';

class ShipmentDetailScreen extends StatefulWidget {
  const ShipmentDetailScreen({super.key, required this.shipmentId});

  final int shipmentId;

  @override
  State<ShipmentDetailScreen> createState() => _ShipmentDetailScreenState();
}

class _ShipmentDetailScreenState extends State<ShipmentDetailScreen> {
  Shipment? _shipment;
  bool _cancelling = false;

  @override
  void initState() {
    super.initState();
    _shipment = context.read<ShipmentProvider>().shipmentById(widget.shipmentId);
  }

  void _copyTracking() {
    final tracking = _shipment?.trackingNumber;
    if (tracking == null || tracking.isEmpty) return;
    Clipboard.setData(ClipboardData(text: tracking));
    ScaffoldMessenger.of(context).showSnackBar(
      SnackBar(content: Text('Copiado: $tracking')),
    );
  }

  Future<void> _cancel() async {
    final shipment = _shipment;
    if (shipment == null) return;

    final confirmed = await showDialog<bool>(
      context: context,
      builder: (dialogContext) => AlertDialog(
        title: const Text('Cancelar guia'),
        content: const Text(
          'Cancelar la guia en la transportadora no devuelve el saldo debitado de la billetera. Confirmas?',
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(dialogContext, false),
            child: const Text('Volver'),
          ),
          FilledButton(
            style: FilledButton.styleFrom(backgroundColor: AppColors.error),
            onPressed: () => Navigator.pop(dialogContext, true),
            child: const Text('Cancelar guia'),
          ),
        ],
      ),
    );

    if (confirmed != true || !mounted) return;

    setState(() => _cancelling = true);
    final provider = context.read<ShipmentProvider>();
    final ok = await provider.cancelShipment(shipment.id);

    if (!mounted) return;
    setState(() {
      _cancelling = false;
      if (ok) _shipment = provider.shipmentById(widget.shipmentId) ?? _shipment;
    });
    ScaffoldMessenger.of(context).showSnackBar(
      SnackBar(content: Text(ok ? 'Guia cancelada' : (provider.error ?? 'No se pudo cancelar'))),
    );
  }

  @override
  Widget build(BuildContext context) {
    final shipment = _shipment;

    return AppScaffold(
      title: shipment?.trackingNumber ?? 'Guia',
      subtitle: shipment?.carrier,
      onBack: () => Navigator.of(context).pop(),
      actions: [
        if (shipment != null)
          IconButton(
            icon: const Icon(Icons.copy_rounded, size: 20),
            tooltip: 'Copiar guia',
            onPressed: _copyTracking,
          ),
      ],
      body: shipment == null
          ? const AppEmptyState(
              icon: Icons.local_shipping_outlined,
              title: 'Guia no encontrada',
              message: 'Vuelve al listado y abrela de nuevo.',
            )
          : ListView(
              padding: AppSpacing.page,
              children: [
                _Header(shipment: shipment),
                const SizedBox(height: 18),
                const AppSectionHeader(
                  title: 'Costo de la guia',
                  subtitle: 'Costo transportadora + margen aplicado',
                ),
                _CostCard(shipment: shipment),
                if (shipment.isCod) ...[
                  const SizedBox(height: 18),
                  const AppSectionHeader(
                    title: 'Contra entrega',
                    subtitle: 'Comision que descuenta la transportadora del recaudo',
                  ),
                  _CodCard(shipment: shipment),
                ],
                const SizedBox(height: 18),
                const AppSectionHeader(title: 'Destino'),
                _DestinationCard(shipment: shipment),
                const SizedBox(height: 18),
                const AppSectionHeader(title: 'Paquete'),
                _PackageCard(shipment: shipment),
                const SizedBox(height: 20),
                if (shipment.status != 'cancelled')
                  OutlinedButton.icon(
                    onPressed: _cancelling ? null : _cancel,
                    style: OutlinedButton.styleFrom(
                      foregroundColor: AppColors.error,
                      side: const BorderSide(color: AppColors.error),
                    ),
                    icon: _cancelling
                        ? const SizedBox(
                            height: 16,
                            width: 16,
                            child: CircularProgressIndicator(strokeWidth: 2),
                          )
                        : const Icon(Icons.cancel_outlined, size: 18),
                    label: const Text('Cancelar guia'),
                  ),
                const SizedBox(height: 28),
              ],
            ),
    );
  }
}

class _Header extends StatelessWidget {
  const _Header({required this.shipment});

  final Shipment shipment;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    return AppCard(
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(
            children: [
              BrandLogo(name: shipment.carrier ?? '', size: 46, radius: 12, padding: 7),
              const SizedBox(width: 13),
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(shipment.carrier ?? 'Transportadora', style: theme.textTheme.titleMedium),
                    const SizedBox(height: 2),
                    Text(
                      shipment.orderNumber ?? '',
                      style: theme.textTheme.bodySmall,
                    ),
                  ],
                ),
              ),
            ],
          ),
          const SizedBox(height: 14),
          Wrap(
            spacing: 6,
            runSpacing: 6,
            children: [
              AppStatusChip(
                label: shipmentStatusLabel(shipment.status),
                tone: AppStatusChip.toneFromCode(shipment.status),
              ),
              if (shipment.isCod)
                const AppStatusChip(
                  label: 'Contra entrega',
                  tone: AppStatusTone.brand,
                  icon: Icons.payments_outlined,
                ),
            ],
          ),
          if ((shipment.carrierStatusDetail ?? '').isNotEmpty) ...[
            const SizedBox(height: 12),
            Text(shipment.carrierStatusDetail!, style: theme.textTheme.bodySmall),
          ],
        ],
      ),
    );
  }
}

class _CostCard extends StatelessWidget {
  const _CostCard({required this.shipment});

  final Shipment shipment;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    return AppCard(
      child: Column(
        children: [
          AppKeyValueRow(
            label: 'Costo transportadora',
            value: AppFormat.money(shipment.carrierCost),
            dense: true,
          ),
          if ((shipment.insuranceCost ?? 0) > 0)
            AppKeyValueRow(
              label: 'Seguro',
              value: AppFormat.money(shipment.insuranceCost),
              dense: true,
            ),
          AppKeyValueRow(
            label: 'Margen aplicado',
            value: AppFormat.money(shipment.marginAmount),
            dense: true,
          ),
          const Divider(height: 18),
          Row(
            children: [
              Expanded(child: Text('Total de la guia', style: theme.textTheme.titleMedium)),
              Text(
                AppFormat.money(shipment.totalCost),
                style: theme.textTheme.titleLarge?.copyWith(color: AppColors.primary),
              ),
            ],
          ),
        ],
      ),
    );
  }
}

class _CodCard extends StatelessWidget {
  const _CodCard({required this.shipment});

  final Shipment shipment;

  @override
  Widget build(BuildContext context) {
    return AppCard(
      child: Column(
        children: [
          AppKeyValueRow(
            label: 'Comision transportadora',
            value: AppFormat.money(shipment.codCarrierFee),
            dense: true,
          ),
          AppKeyValueRow(
            label: 'Margen Probability',
            value: AppFormat.money(shipment.codProbabilityMargin),
            dense: true,
          ),
        ],
      ),
    );
  }
}

class _DestinationCard extends StatelessWidget {
  const _DestinationCard({required this.shipment});

  final Shipment shipment;

  @override
  Widget build(BuildContext context) {
    return AppCard(
      child: Column(
        children: [
          AppKeyValueRow(
            label: 'Destinatario',
            value: shipment.clientName ?? shipment.customerName ?? '-',
          ),
          const Divider(height: 18),
          AppKeyValueRow(label: 'Telefono', value: shipment.customerPhone ?? '-'),
          const Divider(height: 18),
          AppKeyValueRow(label: 'Direccion', value: shipment.destinationAddress ?? '-'),
          const Divider(height: 18),
          AppKeyValueRow(
            label: 'Ciudad',
            value: [shipment.destinationCity, shipment.destinationState]
                .where((e) => (e ?? '').isNotEmpty)
                .join(', '),
          ),
          const Divider(height: 18),
          AppKeyValueRow(
            label: 'Entrega estimada',
            value: AppFormat.date(AppFormat.parseDate(shipment.estimatedDelivery)),
          ),
        ],
      ),
    );
  }
}

class _PackageCard extends StatelessWidget {
  const _PackageCard({required this.shipment});

  final Shipment shipment;

  @override
  Widget build(BuildContext context) {
    final dimensions = [shipment.length, shipment.width, shipment.height]
        .map((e) => e == null ? '-' : AppFormat.number(e))
        .join(' x ');

    return AppCard(
      child: Column(
        children: [
          AppKeyValueRow(
            label: 'Peso',
            value: shipment.weight == null ? '-' : '${AppFormat.number(shipment.weight)} g',
          ),
          const Divider(height: 18),
          AppKeyValueRow(label: 'Dimensiones', value: '$dimensions cm'),
          const Divider(height: 18),
          AppKeyValueRow(label: 'Bodega', value: shipment.warehouseName ?? '-'),
          const Divider(height: 18),
          AppKeyValueRow(
            label: 'Generada',
            value: AppFormat.dateTime(AppFormat.parseDate(shipment.createdAt)),
          ),
        ],
      ),
    );
  }
}
