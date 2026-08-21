import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../../../../../shared/theme/app_colors.dart';
import '../../../../../shared/theme/app_tokens.dart';
import '../../../../../shared/utils/formatters.dart';
import '../../../../../shared/widgets/ui/ui.dart';
import '../../domain/entities.dart';
import '../providers/warehouse_provider.dart';
import 'warehouse_detail_screen.dart';

class WarehouseListScreen extends StatefulWidget {
  const WarehouseListScreen({super.key, this.businessId});

  final int? businessId;

  @override
  State<WarehouseListScreen> createState() => _WarehouseListScreenState();
}

class _WarehouseListScreenState extends State<WarehouseListScreen> {
  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addPostFrameCallback((_) => _refresh());
  }

  @override
  void didUpdateWidget(WarehouseListScreen oldWidget) {
    super.didUpdateWidget(oldWidget);
    if (oldWidget.businessId != widget.businessId) _refresh();
  }

  void _refresh() {
    context.read<WarehouseProvider>().fetchWarehouses(businessId: widget.businessId);
  }

  @override
  Widget build(BuildContext context) {
    return Consumer<WarehouseProvider>(
      builder: (context, provider, _) {
        if (provider.isLoading && provider.warehouses.isEmpty) {
          return const AppListSkeleton();
        }
        if (provider.error != null && provider.warehouses.isEmpty) {
          return AppErrorState(message: provider.error!, onRetry: _refresh);
        }
        if (provider.warehouses.isEmpty) {
          return AppEmptyState(
            icon: Icons.warehouse_outlined,
            title: 'Sin bodegas',
            message: 'Crea la primera bodega para poder generar guias y controlar stock.',
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
            itemCount: provider.warehouses.length,
            separatorBuilder: (context, index) => const SizedBox(height: 10),
            itemBuilder: (context, index) {
              final warehouse = provider.warehouses[index];
              return _WarehouseCard(
                warehouse: warehouse,
                onTap: () => Navigator.of(context).push(
                  MaterialPageRoute(
                    builder: (_) => WarehouseDetailScreen(
                      warehouseId: warehouse.id,
                      businessId: widget.businessId,
                    ),
                  ),
                ),
              );
            },
          ),
        );
      },
    );
  }
}

class _WarehouseCard extends StatelessWidget {
  const _WarehouseCard({required this.warehouse, required this.onTap});

  final Warehouse warehouse;
  final VoidCallback onTap;

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
              Container(
                width: 42,
                height: 42,
                alignment: Alignment.center,
                decoration: BoxDecoration(
                  color: warehouse.isActive ? AppColors.primarySoft : AppColors.surfaceMuted,
                  borderRadius: AppRadius.mdAll,
                ),
                child: Icon(
                  Icons.warehouse_outlined,
                  size: 20,
                  color: warehouse.isActive ? AppColors.primary : AppColors.textDisabled,
                ),
              ),
              const SizedBox(width: 12),
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(
                      warehouse.name,
                      style: theme.textTheme.titleSmall,
                      maxLines: 1,
                      overflow: TextOverflow.ellipsis,
                    ),
                    const SizedBox(height: 2),
                    Text(
                      '${warehouse.code}  ·  ${warehouse.city}',
                      style: theme.textTheme.labelSmall,
                      maxLines: 1,
                      overflow: TextOverflow.ellipsis,
                    ),
                  ],
                ),
              ),
              const Icon(Icons.chevron_right_rounded, color: AppColors.textDisabled),
            ],
          ),
          const SizedBox(height: 11),
          Wrap(
            spacing: 6,
            runSpacing: 6,
            children: [
              AppStatusChip(
                dense: true,
                label: warehouse.isActive ? 'Activa' : 'Inactiva',
                tone: warehouse.isActive ? AppStatusTone.success : AppStatusTone.neutral,
              ),
              if (warehouse.isDefault)
                const AppStatusChip(dense: true, label: 'Predeterminada', tone: AppStatusTone.brand),
              if (warehouse.isFulfillment)
                const AppStatusChip(dense: true, label: 'Despacha', tone: AppStatusTone.info),
            ],
          ),
          if (warehouse.address.isNotEmpty) ...[
            const SizedBox(height: 10),
            Row(
              children: [
                const Icon(Icons.place_outlined, size: 13, color: AppColors.textDisabled),
                const SizedBox(width: 5),
                Expanded(
                  child: Text(
                    warehouse.address,
                    style: theme.textTheme.labelSmall,
                    maxLines: 1,
                    overflow: TextOverflow.ellipsis,
                  ),
                ),
                Text(
                  AppFormat.relative(AppFormat.parseDate(warehouse.updatedAt)),
                  style: theme.textTheme.labelSmall,
                ),
              ],
            ),
          ],
        ],
      ),
    );
  }
}
