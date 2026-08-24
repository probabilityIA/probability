import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../../../../../shared/theme/app_colors.dart';
import '../../../../../shared/theme/app_tokens.dart';
import '../../../../../shared/utils/formatters.dart';
import '../../../../../shared/widgets/ui/ui.dart';
import '../../domain/entities.dart';
import '../providers/warehouse_provider.dart';

class WarehouseDetailScreen extends StatefulWidget {
  const WarehouseDetailScreen({
    super.key,
    required this.warehouseId,
    this.businessId,
  });

  final int warehouseId;
  final int? businessId;

  @override
  State<WarehouseDetailScreen> createState() => _WarehouseDetailScreenState();
}

class _WarehouseDetailScreenState extends State<WarehouseDetailScreen> {
  Warehouse? _warehouse;

  @override
  void initState() {
    super.initState();
    _warehouse = context.read<WarehouseProvider>().warehouseById(widget.warehouseId);
    WidgetsBinding.instance.addPostFrameCallback((_) {
      context
          .read<WarehouseProvider>()
          .fetchLocations(widget.warehouseId, businessId: widget.businessId);
    });
  }

  @override
  Widget build(BuildContext context) {
    final warehouse = _warehouse;

    return AppScaffold(
      title: warehouse?.name ?? 'Bodega',
      subtitle: warehouse?.code,
      onBack: () => Navigator.of(context).pop(),
      body: warehouse == null
          ? const AppEmptyState(
              icon: Icons.warehouse_outlined,
              title: 'Bodega no encontrada',
            )
          : ListView(
              padding: AppSpacing.page,
              children: [
                _StatusCard(warehouse: warehouse),
                const SizedBox(height: 18),
                const AppSectionHeader(title: 'Direccion'),
                _AddressCard(warehouse: warehouse),
                const SizedBox(height: 18),
                const AppSectionHeader(
                  title: 'Contacto',
                  subtitle: 'Datos que viajan al generar una guia',
                ),
                _ContactCard(warehouse: warehouse),
                const SizedBox(height: 18),
                const AppSectionHeader(
                  title: 'Ubicaciones',
                  subtitle: 'Pasillos y zonas de la bodega',
                ),
                const _LocationsCard(),
                const SizedBox(height: 26),
              ],
            ),
    );
  }
}

class _StatusCard extends StatelessWidget {
  const _StatusCard({required this.warehouse});

  final Warehouse warehouse;

  @override
  Widget build(BuildContext context) {
    return AppCard(
      child: Row(
        children: [
          Container(
            width: 48,
            height: 48,
            alignment: Alignment.center,
            decoration: BoxDecoration(
              color: warehouse.isActive ? AppColors.primarySoft : AppColors.surfaceMuted,
              borderRadius: AppRadius.mdAll,
            ),
            child: Icon(
              Icons.warehouse_outlined,
              size: 23,
              color: warehouse.isActive ? AppColors.primary : AppColors.textDisabled,
            ),
          ),
          const SizedBox(width: 14),
          Expanded(
            child: Wrap(
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
          ),
        ],
      ),
    );
  }
}

class _AddressCard extends StatelessWidget {
  const _AddressCard({required this.warehouse});

  final Warehouse warehouse;

  @override
  Widget build(BuildContext context) {
    return AppCard(
      child: Column(
        children: [
          AppKeyValueRow(label: 'Direccion', value: warehouse.address),
          const Divider(height: 18),
          AppKeyValueRow(label: 'Barrio', value: warehouse.suburb.isEmpty ? '-' : warehouse.suburb),
          const Divider(height: 18),
          AppKeyValueRow(label: 'Ciudad', value: '${warehouse.city}, ${warehouse.state}'),
          const Divider(height: 18),
          AppKeyValueRow(
            label: 'Codigo DANE',
            value: warehouse.cityDaneCode.isEmpty ? '-' : warehouse.cityDaneCode,
          ),
          const Divider(height: 18),
          AppKeyValueRow(
            label: 'Codigo postal',
            value: warehouse.postalCode.isEmpty ? warehouse.zipCode : warehouse.postalCode,
          ),
        ],
      ),
    );
  }
}

class _ContactCard extends StatelessWidget {
  const _ContactCard({required this.warehouse});

  final Warehouse warehouse;

  bool get _missingContact =>
      warehouse.phone.trim().isEmpty || warehouse.contactEmail.trim().isEmpty;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    return AppCard(
      child: Column(
        children: [
          AppKeyValueRow(
            label: 'Responsable',
            value: warehouse.contactName.isEmpty ? '-' : warehouse.contactName,
          ),
          const Divider(height: 18),
          AppKeyValueRow(
            label: 'Telefono',
            value: warehouse.phone.isEmpty ? '-' : warehouse.phone,
          ),
          const Divider(height: 18),
          AppKeyValueRow(
            label: 'Correo',
            value: warehouse.contactEmail.isEmpty ? '-' : warehouse.contactEmail,
          ),
          if (_missingContact) ...[
            const SizedBox(height: 12),
            Container(
              padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 10),
              decoration: BoxDecoration(
                color: AppColors.warningSoft,
                borderRadius: AppRadius.mdAll,
              ),
              child: Row(
                children: [
                  const Icon(Icons.warning_amber_rounded, size: 17, color: Color(0xFFB45309)),
                  const SizedBox(width: 9),
                  Expanded(
                    child: Text(
                      'Sin telefono o correo la transportadora rechaza la guia',
                      style: theme.textTheme.labelSmall?.copyWith(
                        color: const Color(0xFFB45309),
                        fontWeight: FontWeight.w600,
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

class _LocationsCard extends StatelessWidget {
  const _LocationsCard();

  static const Map<String, String> _typeLabels = {
    'storage': 'Almacenamiento',
    'picking': 'Alistamiento',
    'receiving': 'Recepcion',
    'shipping': 'Despacho',
  };

  @override
  Widget build(BuildContext context) {
    return Consumer<WarehouseProvider>(
      builder: (context, provider, _) {
        if (provider.isLoadingLocations) {
          return const AppCard(child: AppLoading());
        }
        if (provider.locations.isEmpty) {
          return AppCard(
            child: Text(
              'Esta bodega no tiene ubicaciones configuradas',
              style: Theme.of(context).textTheme.bodySmall,
            ),
          );
        }

        return AppCard(
          child: Column(
            children: [
              for (var i = 0; i < provider.locations.length; i++) ...[
                if (i > 0) const Divider(height: 20),
                Row(
                  children: [
                    Container(
                      width: 34,
                      height: 34,
                      alignment: Alignment.center,
                      decoration: BoxDecoration(
                        color: AppColors.surfaceMuted,
                        borderRadius: AppRadius.smAll,
                      ),
                      child: const Icon(Icons.grid_view_rounded, size: 16, color: AppColors.textMuted),
                    ),
                    const SizedBox(width: 12),
                    Expanded(
                      child: Column(
                        crossAxisAlignment: CrossAxisAlignment.start,
                        children: [
                          Text(
                            provider.locations[i].name,
                            style: Theme.of(context).textTheme.titleSmall,
                          ),
                          const SizedBox(height: 2),
                          Text(
                            '${provider.locations[i].code}  ·  ${_typeLabels[provider.locations[i].type] ?? provider.locations[i].type}',
                            style: Theme.of(context).textTheme.labelSmall,
                          ),
                        ],
                      ),
                    ),
                    if (provider.locations[i].capacity != null)
                      Text(
                        '${AppFormat.number(provider.locations[i].capacity)} uds',
                        style: Theme.of(context).textTheme.labelMedium,
                      ),
                  ],
                ),
              ],
            ],
          ),
        );
      },
    );
  }
}
