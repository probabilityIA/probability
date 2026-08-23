import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:provider/provider.dart';
import '../../../../../shared/theme/app_colors.dart';
import '../../../../../shared/theme/app_tokens.dart';
import '../../../../../shared/utils/formatters.dart';
import '../../../../../shared/widgets/ui/ui.dart';
import '../../../orders/domain/entities.dart';
import '../../../orders/ui/providers/order_provider.dart';
import '../../../orders/ui/widgets/order_card.dart';
import '../../domain/entities.dart';
import '../providers/customer_provider.dart';

class CustomerDetailScreen extends StatefulWidget {
  const CustomerDetailScreen({
    super.key,
    required this.customerId,
    this.businessId,
  });

  final int customerId;
  final int? businessId;

  @override
  State<CustomerDetailScreen> createState() => _CustomerDetailScreenState();
}

class _CustomerDetailScreenState extends State<CustomerDetailScreen> {
  CustomerInfo? _customer;
  List<Order> _orders = const [];
  bool _loadingOrders = true;

  @override
  void initState() {
    super.initState();
    _customer = context.read<CustomerProvider>().customerById(widget.customerId);
    WidgetsBinding.instance.addPostFrameCallback((_) => _loadOrders());
  }

  Future<void> _loadOrders() async {
    final customer = _customer;
    if (customer == null || (customer.email ?? '').isEmpty) {
      setState(() => _loadingOrders = false);
      return;
    }

    final orders = await context.read<OrderProvider>().ordersByCustomerEmail(
          customer.email!,
          businessId: widget.businessId,
        );

    if (!mounted) return;
    setState(() {
      _orders = orders;
      _loadingOrders = false;
    });
  }

  void _copy(String label, String value) {
    Clipboard.setData(ClipboardData(text: value));
    ScaffoldMessenger.of(context).showSnackBar(
      SnackBar(content: Text('$label copiado')),
    );
  }

  @override
  Widget build(BuildContext context) {
    final customer = _customer;
    final detail = customer is CustomerDetail ? customer : null;

    return AppScaffold(
      title: customer?.name ?? 'Cliente',
      subtitle: customer?.email,
      onBack: () => Navigator.of(context).pop(),
      body: customer == null
          ? const AppEmptyState(
              icon: Icons.people_alt_outlined,
              title: 'Cliente no encontrado',
            )
          : ListView(
              padding: AppSpacing.page,
              children: [
                _HeaderCard(customer: customer, detail: detail),
                if (detail != null) ...[
                  const SizedBox(height: 14),
                  Row(
                    children: [
                      Expanded(
                        child: AppKpiTile(
                          label: 'Ordenes',
                          value: AppFormat.number(detail.orderCount),
                          icon: Icons.receipt_long_outlined,
                        ),
                      ),
                      const SizedBox(width: 11),
                      Expanded(
                        child: AppKpiTile(
                          label: 'Total comprado',
                          value: AppFormat.compact(detail.totalSpent),
                          icon: Icons.payments_outlined,
                          accent: AppColors.success,
                        ),
                      ),
                    ],
                  ),
                ],
                const SizedBox(height: 18),
                const AppSectionHeader(title: 'Contacto'),
                _ContactCard(customer: customer, onCopy: _copy),
                const SizedBox(height: 18),
                AppSectionHeader(
                  title: 'Ordenes recientes',
                  subtitle: _loadingOrders ? null : '${_orders.length} encontradas',
                ),
                if (_loadingOrders)
                  const AppCard(child: AppLoading())
                else if (_orders.isEmpty)
                  AppCard(
                    child: Text(
                      'Este cliente todavia no tiene ordenes registradas',
                      style: Theme.of(context).textTheme.bodySmall,
                    ),
                  )
                else
                  ..._orders.take(10).map(
                        (order) => Padding(
                          padding: const EdgeInsets.only(bottom: 10),
                          child: OrderCard(order: order),
                        ),
                      ),
                const SizedBox(height: 26),
              ],
            ),
    );
  }
}

class _HeaderCard extends StatelessWidget {
  const _HeaderCard({required this.customer, this.detail});

  final CustomerInfo customer;
  final CustomerDetail? detail;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    return AppCard(
      child: Row(
        children: [
          Container(
            width: 56,
            height: 56,
            alignment: Alignment.center,
            decoration: const BoxDecoration(
              color: AppColors.primarySoft,
              shape: BoxShape.circle,
            ),
            child: Text(
              AppFormat.initials(customer.name),
              style: const TextStyle(
                fontFamily: 'Inter',
                fontSize: 17,
                fontWeight: FontWeight.w700,
                color: AppColors.primary,
              ),
            ),
          ),
          const SizedBox(width: 14),
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(customer.name, style: theme.textTheme.titleMedium),
                const SizedBox(height: 3),
                Text(
                  'Cliente desde ${AppFormat.date(AppFormat.parseDate(customer.createdAt))}',
                  style: theme.textTheme.bodySmall,
                ),
                if (detail?.lastOrderAt != null) ...[
                  const SizedBox(height: 6),
                  AppStatusChip(
                    dense: true,
                    label: 'Ultima compra ${AppFormat.relative(AppFormat.parseDate(detail!.lastOrderAt))}',
                    tone: AppStatusTone.brand,
                  ),
                ],
              ],
            ),
          ),
        ],
      ),
    );
  }
}

class _ContactCard extends StatelessWidget {
  const _ContactCard({required this.customer, required this.onCopy});

  final CustomerInfo customer;
  final void Function(String label, String value) onCopy;

  @override
  Widget build(BuildContext context) {
    return AppCard(
      padding: EdgeInsets.zero,
      child: Column(
        children: [
          _Row(
            icon: Icons.mail_outline,
            label: 'Correo',
            value: customer.email ?? '-',
            onCopy: customer.email == null ? null : () => onCopy('Correo', customer.email!),
          ),
          const Divider(height: 1),
          _Row(
            icon: Icons.phone_outlined,
            label: 'Telefono',
            value: customer.phone.isEmpty ? '-' : customer.phone,
            onCopy: customer.phone.isEmpty ? null : () => onCopy('Telefono', customer.phone),
          ),
          const Divider(height: 1),
          _Row(
            icon: Icons.badge_outlined,
            label: 'Documento',
            value: customer.dni ?? '-',
            onCopy: customer.dni == null ? null : () => onCopy('Documento', customer.dni!),
          ),
        ],
      ),
    );
  }
}

class _Row extends StatelessWidget {
  const _Row({
    required this.icon,
    required this.label,
    required this.value,
    this.onCopy,
  });

  final IconData icon;
  final String label;
  final String value;
  final VoidCallback? onCopy;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    return Padding(
      padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 13),
      child: Row(
        children: [
          Icon(icon, size: 19, color: AppColors.textMuted),
          const SizedBox(width: 13),
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(label, style: theme.textTheme.labelSmall),
                const SizedBox(height: 2),
                Text(
                  value,
                  style: theme.textTheme.titleSmall,
                  maxLines: 1,
                  overflow: TextOverflow.ellipsis,
                ),
              ],
            ),
          ),
          if (onCopy != null)
            IconButton(
              icon: const Icon(Icons.copy_rounded, size: 17),
              onPressed: onCopy,
              tooltip: 'Copiar',
            ),
        ],
      ),
    );
  }
}
