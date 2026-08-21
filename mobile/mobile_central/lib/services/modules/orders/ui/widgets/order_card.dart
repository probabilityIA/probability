import 'package:flutter/material.dart';
import '../../../../../shared/theme/app_colors.dart';
import '../../../../../shared/utils/formatters.dart';
import '../../../../../shared/widgets/ui/ui.dart';
import '../../domain/entities.dart';

class OrderCard extends StatelessWidget {
  const OrderCard({super.key, required this.order, this.onTap});

  final Order order;
  final VoidCallback? onTap;

  int get _itemCount {
    final raw = order.orderItems;
    if (raw is List) return raw.length;
    return 0;
  }

  String get _channelName =>
      order.integrationName?.isNotEmpty == true ? order.integrationName! : order.integrationType;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final statusName = order.orderStatus?.name ?? order.status;

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
                name: _channelName,
                imageUrl: order.integrationLogoUrl,
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
                      order.orderNumber,
                      style: theme.textTheme.titleSmall,
                      maxLines: 1,
                      overflow: TextOverflow.ellipsis,
                    ),
                    const SizedBox(height: 2),
                    Text(
                      order.customerName.isEmpty ? 'Sin cliente' : order.customerName,
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
                    AppFormat.money(order.totalAmount),
                    style: theme.textTheme.titleSmall?.copyWith(fontSize: 15),
                  ),
                  const SizedBox(height: 3),
                  Text(
                    AppFormat.relative(AppFormat.parseDate(order.createdAt)),
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
                label: statusName,
                tone: AppStatusChip.toneFromCode(order.status),
              ),
              AppStatusChip(
                dense: true,
                label: order.isPaid ? 'Pagada' : 'Sin pago',
                tone: order.isPaid ? AppStatusTone.success : AppStatusTone.warning,
              ),
              if (order.codTotal != null && order.codTotal! > 0)
                const AppStatusChip(
                  dense: true,
                  label: 'Contra entrega',
                  tone: AppStatusTone.brand,
                  icon: Icons.payments_outlined,
                ),
              if ((order.trackingNumber ?? '').isNotEmpty)
                AppStatusChip(
                  dense: true,
                  label: order.trackingNumber!,
                  tone: AppStatusTone.info,
                  icon: Icons.local_shipping_outlined,
                ),
            ],
          ),
          if (order.shippingCity.isNotEmpty) ...[
            const SizedBox(height: 11),
            Row(
              children: [
                const Icon(Icons.place_outlined, size: 13, color: AppColors.textDisabled),
                const SizedBox(width: 5),
                Expanded(
                  child: Text(
                    order.shippingCity,
                    style: theme.textTheme.labelSmall,
                    maxLines: 1,
                    overflow: TextOverflow.ellipsis,
                  ),
                ),
                if (_itemCount > 0)
                  Text('$_itemCount items', style: theme.textTheme.labelSmall),
              ],
            ),
          ],
        ],
      ),
    );
  }
}
