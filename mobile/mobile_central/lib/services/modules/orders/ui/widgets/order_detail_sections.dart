import 'package:flutter/material.dart';
import '../../../../../shared/theme/app_colors.dart';
import '../../../../../shared/theme/app_tokens.dart';
import '../../../../../shared/utils/formatters.dart';
import '../../../../../shared/widgets/ui/ui.dart';
import '../../domain/entities.dart';

class OrderHeaderCard extends StatelessWidget {
  const OrderHeaderCard({super.key, required this.order});

  final Order order;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final channel = order.integrationName?.isNotEmpty == true
        ? order.integrationName!
        : order.integrationType;

    return AppCard(
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(
            children: [
              BrandLogo(
                name: channel,
                imageUrl: order.integrationLogoUrl,
                size: 46,
                radius: 12,
                padding: 7,
              ),
              const SizedBox(width: 13),
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(order.orderNumber, style: theme.textTheme.titleLarge),
                    const SizedBox(height: 2),
                    Text(
                      '$channel  ·  ${AppFormat.dateTime(AppFormat.parseDate(order.createdAt))}',
                      style: theme.textTheme.bodySmall,
                      maxLines: 1,
                      overflow: TextOverflow.ellipsis,
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
                label: order.orderStatus?.name ?? order.status,
                tone: AppStatusChip.toneFromCode(order.status),
              ),
              AppStatusChip(
                label: order.paymentStatus?.name ?? (order.isPaid ? 'Pagada' : 'Sin pago'),
                tone: order.isPaid ? AppStatusTone.success : AppStatusTone.warning,
              ),
              if (order.fulfillmentStatus != null)
                AppStatusChip(
                  label: order.fulfillmentStatus!.name,
                  tone: AppStatusChip.toneFromCode(order.fulfillmentStatus!.code),
                ),
              if (order.codTotal != null && order.codTotal! > 0)
                const AppStatusChip(
                  label: 'Contra entrega',
                  tone: AppStatusTone.brand,
                  icon: Icons.payments_outlined,
                ),
            ],
          ),
        ],
      ),
    );
  }
}

class OrderItemsCard extends StatelessWidget {
  const OrderItemsCard({super.key, required this.items});

  final List<OrderLineItem> items;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    if (items.isEmpty) {
      return AppCard(
        child: Text('Esta orden no tiene items', style: theme.textTheme.bodySmall),
      );
    }

    return AppCard(
      child: Column(
        children: [
          for (var i = 0; i < items.length; i++) ...[
            if (i > 0) const Divider(height: 20),
            Row(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Container(
                  width: 42,
                  height: 42,
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
                        items[i].name,
                        style: theme.textTheme.titleSmall,
                        maxLines: 2,
                        overflow: TextOverflow.ellipsis,
                      ),
                      const SizedBox(height: 3),
                      Text(
                        'SKU ${items[i].sku}  ·  ${AppFormat.money(items[i].unitPrice)} c/u',
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

class OrderTotalsCard extends StatelessWidget {
  const OrderTotalsCard({super.key, required this.order});

  final Order order;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    return AppCard(
      child: Column(
        children: [
          AppKeyValueRow(label: 'Subtotal', value: AppFormat.money(order.subtotal), dense: true),
          if (order.discount > 0)
            AppKeyValueRow(
              label: 'Descuento',
              value: '- ${AppFormat.money(order.discount)}',
              dense: true,
            ),
          AppKeyValueRow(label: 'Envio', value: AppFormat.money(order.shippingCost), dense: true),
          AppKeyValueRow(label: 'Impuestos', value: AppFormat.money(order.tax), dense: true),
          const Divider(height: 18),
          Row(
            children: [
              Expanded(child: Text('Total', style: theme.textTheme.titleMedium)),
              Text(
                AppFormat.money(order.totalAmount),
                style: theme.textTheme.titleLarge?.copyWith(color: AppColors.primary),
              ),
            ],
          ),
          if (order.codTotal != null && order.codTotal! > 0) ...[
            const SizedBox(height: 10),
            Container(
              padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 10),
              decoration: BoxDecoration(
                color: AppColors.primarySoft,
                borderRadius: AppRadius.mdAll,
              ),
              child: Row(
                children: [
                  const Icon(Icons.payments_outlined, size: 17, color: AppColors.primaryDark),
                  const SizedBox(width: 9),
                  Expanded(
                    child: Text(
                      'A recaudar contra entrega',
                      style: theme.textTheme.labelMedium?.copyWith(
                        color: AppColors.primaryDark,
                        fontWeight: FontWeight.w600,
                      ),
                    ),
                  ),
                  Text(
                    AppFormat.money(order.codTotal),
                    style: theme.textTheme.titleSmall?.copyWith(color: AppColors.primaryDark),
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

class OrderCustomerCard extends StatelessWidget {
  const OrderCustomerCard({super.key, required this.order});

  final Order order;

  @override
  Widget build(BuildContext context) {
    final address = [
      order.shippingStreet,
      order.shippingCity,
      order.shippingState,
    ].where((e) => e.isNotEmpty).join(', ');

    return AppCard(
      child: Column(
        children: [
          AppKeyValueRow(
            label: 'Cliente',
            value: order.customerName.isEmpty ? '-' : order.customerName,
          ),
          const Divider(height: 18),
          AppKeyValueRow(
            label: 'Documento',
            value: order.customerDni.isEmpty ? '-' : order.customerDni,
          ),
          const Divider(height: 18),
          AppKeyValueRow(
            label: 'Telefono',
            value: order.customerPhone.isEmpty ? '-' : order.customerPhone,
          ),
          const Divider(height: 18),
          AppKeyValueRow(
            label: 'Correo',
            value: order.customerEmail.isEmpty ? '-' : order.customerEmail,
          ),
          const Divider(height: 18),
          AppKeyValueRow(label: 'Direccion', value: address.isEmpty ? '-' : address),
        ],
      ),
    );
  }
}

class OrderShippingCard extends StatelessWidget {
  const OrderShippingCard({super.key, required this.order, this.onOpenGuide});

  final Order order;
  final VoidCallback? onOpenGuide;

  @override
  Widget build(BuildContext context) {
    return AppCard(
      child: Column(
        children: [
          AppKeyValueRow(
            label: 'Guia',
            value: (order.trackingNumber ?? '').isEmpty ? 'Sin generar' : order.trackingNumber!,
          ),
          const Divider(height: 18),
          AppKeyValueRow(
            label: 'Bodega',
            value: order.warehouseName.isEmpty ? '-' : order.warehouseName,
          ),
          const Divider(height: 18),
          AppKeyValueRow(
            label: 'Peso',
            value: order.weight == null ? '-' : '${AppFormat.number(order.weight)} g',
          ),
          if ((order.guideLink ?? '').isNotEmpty) ...[
            const SizedBox(height: 12),
            SizedBox(
              width: double.infinity,
              child: OutlinedButton.icon(
                onPressed: onOpenGuide,
                icon: const Icon(Icons.description_outlined, size: 18),
                label: const Text('Ver guia'),
              ),
            ),
          ],
        ],
      ),
    );
  }
}
