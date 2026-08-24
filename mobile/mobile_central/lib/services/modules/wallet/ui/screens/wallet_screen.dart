import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../../../../../shared/theme/app_colors.dart';
import '../../../../../shared/theme/app_tokens.dart';
import '../../../../../shared/utils/formatters.dart';
import '../../../../../shared/widgets/ui/ui.dart';
import '../../domain/entities.dart';
import '../providers/wallet_provider.dart';
import 'wallet_recharge_sheet.dart';

const Map<String, String> walletConceptLabels = {
  'GUIDE': 'Cobro de guia',
  'RECHARGE': 'Recarga',
  'SUBSCRIPTION': 'Suscripcion',
  'REFUND': 'Reembolso',
  'ADJUSTMENT': 'Ajuste manual',
  'OTHER': 'Otro cobro',
};

const Map<String, IconData> walletConceptIcons = {
  'GUIDE': Icons.local_shipping_outlined,
  'RECHARGE': Icons.add_card_outlined,
  'SUBSCRIPTION': Icons.workspace_premium_outlined,
  'REFUND': Icons.undo_rounded,
  'ADJUSTMENT': Icons.tune_rounded,
  'OTHER': Icons.receipt_long_outlined,
};

class WalletScreen extends StatefulWidget {
  const WalletScreen({super.key, this.businessId});

  final int? businessId;

  @override
  State<WalletScreen> createState() => _WalletScreenState();
}

class _WalletScreenState extends State<WalletScreen> {
  String _filter = '';

  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addPostFrameCallback((_) => _refresh());
  }

  void _refresh() {
    final provider = context.read<WalletProvider>();
    provider.fetchBalance(businessId: widget.businessId);
    provider.fetchHistory(businessId: widget.businessId);
  }

  Future<void> _openRecharge() async {
    final done = await showWalletRechargeSheet(context, businessId: widget.businessId);
    if (done == true) _refresh();
  }

  @override
  Widget build(BuildContext context) {
    return AppScaffold(
      title: 'Billetera',
      subtitle: 'Saldo y movimientos',
      actions: [
        IconButton(
          icon: const Icon(Icons.refresh_rounded, size: 20),
          tooltip: 'Actualizar',
          onPressed: _refresh,
        ),
      ],
      body: Consumer<WalletProvider>(
        builder: (context, provider, _) {
          if (provider.isLoading && provider.wallet == null) {
            return const AppListSkeleton();
          }
          if (provider.error != null && provider.wallet == null) {
            return AppErrorState(message: provider.error!, onRetry: _refresh);
          }

          final movements = provider.movements;
          final filtered = _filter.isEmpty
              ? movements
              : movements.where((m) => m.concept == _filter).toList();

          final spentOnGuides = movements
              .where((m) => m.concept == 'GUIDE' && !m.isCredit)
              .fold<double>(0, (acc, m) => acc + m.amount);
          final recharged = movements
              .where((m) => m.isCredit)
              .fold<double>(0, (acc, m) => acc + m.amount);

          return RefreshIndicator(
            onRefresh: () async => _refresh(),
            color: AppColors.primary,
            child: ListView(
              physics: const AlwaysScrollableScrollPhysics(),
              padding: AppSpacing.page,
              children: [
                _BalanceCard(
                  balance: provider.wallet?.balance ?? 0,
                  onRecharge: _openRecharge,
                ),
                const SizedBox(height: 14),
                Row(
                  children: [
                    Expanded(
                      child: AppKpiTile(
                        label: 'Recargado',
                        value: AppFormat.compact(recharged),
                        icon: Icons.add_card_outlined,
                        accent: AppColors.success,
                      ),
                    ),
                    const SizedBox(width: 11),
                    Expanded(
                      child: AppKpiTile(
                        label: 'Gastado en guias',
                        value: AppFormat.compact(spentOnGuides),
                        icon: Icons.local_shipping_outlined,
                        accent: const Color(0xFFF97316),
                      ),
                    ),
                  ],
                ),
                const SizedBox(height: 20),
                const AppSectionHeader(title: 'Movimientos'),
                _ConceptFilter(
                  selected: _filter,
                  onSelected: (value) => setState(() => _filter = value),
                ),
                const SizedBox(height: 12),
                if (filtered.isEmpty)
                  AppCard(
                    child: Text(
                      'Sin movimientos para este filtro',
                      style: Theme.of(context).textTheme.bodySmall,
                    ),
                  )
                else
                  AppCard(
                    child: Column(
                      children: [
                        for (var i = 0; i < filtered.length; i++) ...[
                          if (i > 0) const Divider(height: 1),
                          _MovementTile(movement: filtered[i]),
                        ],
                      ],
                    ),
                  ),
                const SizedBox(height: 26),
              ],
            ),
          );
        },
      ),
    );
  }
}

class _BalanceCard extends StatelessWidget {
  const _BalanceCard({required this.balance, required this.onRecharge});

  final double balance;
  final VoidCallback onRecharge;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final low = balance < 100000;

    return Container(
      padding: const EdgeInsets.all(18),
      decoration: BoxDecoration(
        color: AppColors.primary,
        borderRadius: AppRadius.lgAll,
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(
            children: [
              const Icon(Icons.account_balance_wallet_outlined, size: 18, color: Colors.white70),
              const SizedBox(width: 8),
              Text(
                'Saldo disponible',
                style: theme.textTheme.labelMedium?.copyWith(color: Colors.white70),
              ),
            ],
          ),
          const SizedBox(height: 12),
          Text(
            AppFormat.money(balance),
            style: theme.textTheme.displaySmall?.copyWith(
              color: Colors.white,
              fontSize: 30,
            ),
          ),
          if (low) ...[
            const SizedBox(height: 10),
            Container(
              padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 7),
              decoration: BoxDecoration(
                color: Colors.white.withValues(alpha: 0.16),
                borderRadius: AppRadius.smAll,
              ),
              child: Row(
                mainAxisSize: MainAxisSize.min,
                children: [
                  const Icon(Icons.warning_amber_rounded, size: 14, color: Colors.white),
                  const SizedBox(width: 7),
                  Text(
                    'Saldo bajo: las guias pueden fallar',
                    style: theme.textTheme.labelSmall?.copyWith(
                      color: Colors.white,
                      fontWeight: FontWeight.w600,
                    ),
                  ),
                ],
              ),
            ),
          ],
          const SizedBox(height: 16),
          SizedBox(
            width: double.infinity,
            child: FilledButton.icon(
              onPressed: onRecharge,
              style: FilledButton.styleFrom(
                backgroundColor: Colors.white,
                foregroundColor: AppColors.primary,
                minimumSize: const Size(0, 46),
              ),
              icon: const Icon(Icons.add_rounded, size: 19),
              label: const Text('Recargar saldo'),
            ),
          ),
        ],
      ),
    );
  }
}

class _ConceptFilter extends StatelessWidget {
  const _ConceptFilter({required this.selected, required this.onSelected});

  final String selected;
  final ValueChanged<String> onSelected;

  @override
  Widget build(BuildContext context) {
    return SizedBox(
      height: 38,
      child: ListView(
        scrollDirection: Axis.horizontal,
        children: [
          for (final entry in [
            (value: '', label: 'Todos'),
            (value: 'GUIDE', label: 'Guias'),
            (value: 'RECHARGE', label: 'Recargas'),
            (value: 'SUBSCRIPTION', label: 'Suscripcion'),
            (value: 'REFUND', label: 'Reembolsos'),
            (value: 'ADJUSTMENT', label: 'Ajustes'),
          ])
            Padding(
              padding: const EdgeInsets.only(right: 8),
              child: ChoiceChip(
                label: Text(entry.label),
                selected: entry.value == selected,
                onSelected: (_) => onSelected(entry.value),
                labelStyle: TextStyle(
                  fontFamily: 'Inter',
                  fontSize: 12.5,
                  fontWeight: FontWeight.w600,
                  color: entry.value == selected
                      ? AppColors.primaryDark
                      : AppColors.textSecondary,
                ),
              ),
            ),
        ],
      ),
    );
  }
}

class _MovementTile extends StatelessWidget {
  const _MovementTile({required this.movement});

  final WalletMovement movement;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final credit = movement.isCredit;
    final color = credit ? AppColors.success : AppColors.textPrimary;

    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 11),
      child: Row(
        children: [
          Container(
            width: 36,
            height: 36,
            decoration: BoxDecoration(
              color: credit
                  ? AppColors.successSoft
                  : AppColors.surfaceMuted,
              borderRadius: AppRadius.smAll,
            ),
            child: Icon(
              walletConceptIcons[movement.concept] ?? Icons.receipt_long_outlined,
              size: 17,
              color: credit ? const Color(0xFF047857) : AppColors.textMuted,
            ),
          ),
          const SizedBox(width: 12),
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  walletConceptLabels[movement.concept] ?? movement.concept,
                  style: theme.textTheme.titleSmall,
                ),
                const SizedBox(height: 2),
                Text(
                  '${movement.reference ?? ''}  ·  ${AppFormat.relative(AppFormat.parseDate(movement.createdAt))}',
                  style: theme.textTheme.labelSmall,
                  maxLines: 1,
                  overflow: TextOverflow.ellipsis,
                ),
              ],
            ),
          ),
          const SizedBox(width: 10),
          Text(
            '${credit ? '+' : '-'} ${AppFormat.money(movement.amount)}',
            style: theme.textTheme.titleSmall?.copyWith(color: color, fontSize: 14),
          ),
        ],
      ),
    );
  }
}
