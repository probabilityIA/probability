import 'package:flutter/material.dart';
import '../../../../../shared/theme/app_colors.dart';
import '../../../../../shared/utils/formatters.dart';
import '../../../../../shared/widgets/ui/ui.dart';
import '../../domain/entities.dart';

const Map<String, String> shipmentStatusLabels = {
  'created': 'Generada',
  'pending': 'Pendiente',
  'in_transit': 'En transito',
  'delivered': 'Entregada',
  'returned': 'Devuelta',
  'cancelled': 'Cancelada',
  'failed': 'Fallida',
};

String shipmentStatusLabel(String code) =>
    shipmentStatusLabels[code] ?? (code.isEmpty ? 'Sin estado' : code);

class ShipmentCard extends StatelessWidget {
  const ShipmentCard({super.key, required this.shipment, this.onTap});

  final Shipment shipment;
  final VoidCallback? onTap;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final carrier = shipment.carrier ?? 'Transportadora';

    return AppCard(
      padding: const EdgeInsets.all(14),
      onTap: onTap,
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              BrandLogo(name: carrier, size: 40, radius: 11, padding: 6),
              const SizedBox(width: 12),
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(
                      shipment.trackingNumber?.isNotEmpty == true
                          ? shipment.trackingNumber!
                          : 'Sin guia',
                      style: theme.textTheme.titleSmall,
                      maxLines: 1,
                      overflow: TextOverflow.ellipsis,
                    ),
                    const SizedBox(height: 2),
                    Text(
                      '$carrier  ·  ${shipment.orderNumber ?? ''}',
                      style: theme.textTheme.bodySmall,
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
                    AppFormat.money(shipment.listCost),
                    style: theme.textTheme.titleSmall?.copyWith(fontSize: 15),
                  ),
                  const SizedBox(height: 3),
                  Text(
                    AppFormat.relative(AppFormat.parseDate(shipment.createdAt)),
                    style: theme.textTheme.labelSmall,
                  ),
                ],
              ),
            ],
          ),
          const SizedBox(height: 12),
          Wrap(
            spacing: 6,
            runSpacing: 6,
            children: [
              AppStatusChip(
                dense: true,
                label: shipmentStatusLabel(shipment.status),
                tone: AppStatusChip.toneFromCode(shipment.status),
              ),
              if (shipment.isCod)
                const AppStatusChip(
                  dense: true,
                  label: 'Contra entrega',
                  tone: AppStatusTone.brand,
                  icon: Icons.payments_outlined,
                ),
              if (shipment.isTest)
                const AppStatusChip(dense: true, label: 'Prueba', tone: AppStatusTone.warning),
            ],
          ),
          if ((shipment.destinationCity ?? '').isNotEmpty) ...[
            const SizedBox(height: 11),
            Row(
              children: [
                const Icon(Icons.place_outlined, size: 13, color: AppColors.textDisabled),
                const SizedBox(width: 5),
                Expanded(
                  child: Text(
                    shipment.destinationCity!,
                    style: theme.textTheme.labelSmall,
                    maxLines: 1,
                    overflow: TextOverflow.ellipsis,
                  ),
                ),
                Text(
                  shipment.clientName ?? shipment.customerName ?? '',
                  style: theme.textTheme.labelSmall,
                  maxLines: 1,
                  overflow: TextOverflow.ellipsis,
                ),
              ],
            ),
          ],
        ],
      ),
    );
  }
}
