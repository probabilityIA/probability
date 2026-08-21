import 'package:flutter/material.dart';
import '../../../../../shared/theme/app_colors.dart';
import '../../../../../shared/utils/formatters.dart';
import '../../../../../shared/widgets/ui/ui.dart';
import '../../../mobile/domain/entities.dart';

class OrderShipmentCard extends StatelessWidget {
  const OrderShipmentCard({super.key, required this.shipment});

  final MobileShipmentSummary shipment;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final margin = shipment.appliedMargin ??
        ((shipment.totalCost ?? 0) - (shipment.carrierCost ?? 0));

    return AppCard(
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(
            children: [
              BrandLogo(name: shipment.carrier ?? '', size: 40, radius: 11, padding: 6),
              const SizedBox(width: 12),
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(
                      shipment.trackingNumber ?? 'Sin guia',
                      style: theme.textTheme.titleSmall,
                      maxLines: 1,
                      overflow: TextOverflow.ellipsis,
                    ),
                    const SizedBox(height: 2),
                    Text(shipment.carrier ?? '', style: theme.textTheme.bodySmall),
                  ],
                ),
              ),
              AppStatusChip(
                dense: true,
                label: shipment.status,
                tone: AppStatusChip.toneFromCode(shipment.status),
              ),
            ],
          ),
          const Divider(height: 22),
          AppKeyValueRow(
            label: 'Costo transportadora',
            value: AppFormat.money(shipment.carrierCost),
            dense: true,
          ),
          AppKeyValueRow(label: 'Margen aplicado', value: AppFormat.money(margin), dense: true),
          AppKeyValueRow(
            label: 'Total de la guia',
            value: AppFormat.money(shipment.totalCost),
            dense: true,
            valueStyle: theme.textTheme.titleSmall?.copyWith(color: AppColors.primary),
          ),
          if ((shipment.codCarrierFee ?? 0) > 0)
            AppKeyValueRow(
              label: 'Comision contra entrega',
              value: AppFormat.money(shipment.codCarrierFee),
              dense: true,
            ),
        ],
      ),
    );
  }
}

class OrderInvoiceCard extends StatelessWidget {
  const OrderInvoiceCard({super.key, required this.invoice});

  final MobileInvoiceSummary invoice;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    return AppCard(
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(
            children: [
              Container(
                width: 40,
                height: 40,
                alignment: Alignment.center,
                decoration: BoxDecoration(
                  color: AppColors.primarySoft,
                  borderRadius: BorderRadius.circular(11),
                ),
                child: const Icon(Icons.description_outlined, size: 19, color: AppColors.primary),
              ),
              const SizedBox(width: 12),
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(invoice.invoiceNumber, style: theme.textTheme.titleSmall),
                    const SizedBox(height: 2),
                    Text(
                      AppFormat.date(AppFormat.parseDate(invoice.issuedAt)),
                      style: theme.textTheme.bodySmall,
                    ),
                  ],
                ),
              ),
              AppStatusChip(
                dense: true,
                label: invoice.status,
                tone: AppStatusChip.toneFromCode(invoice.status),
              ),
            ],
          ),
          const Divider(height: 22),
          AppKeyValueRow(
            label: 'Total facturado',
            value: AppFormat.money(invoice.totalAmount),
            dense: true,
          ),
          if ((invoice.cufe ?? '').isNotEmpty)
            AppKeyValueRow(label: 'CUFE', value: invoice.cufe!, dense: true),
        ],
      ),
    );
  }
}
