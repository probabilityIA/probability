import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../../../../../shared/theme/app_colors.dart';
import '../../../../../shared/theme/app_tokens.dart';
import '../../../../../shared/utils/formatters.dart';
import '../../../../../shared/widgets/ui/ui.dart';
import '../../domain/entities.dart';
import '../providers/vehicle_provider.dart';

const Map<String, String> vehicleTypeLabels = {
  'van': 'Furgoneta',
  'truck': 'Camion',
  'motorcycle': 'Moto',
  'car': 'Automovil',
};

const Map<String, IconData> vehicleTypeIcons = {
  'van': Icons.airport_shuttle_outlined,
  'truck': Icons.local_shipping_outlined,
  'motorcycle': Icons.two_wheeler_outlined,
  'car': Icons.directions_car_outlined,
};

const Map<String, String> vehicleStatusLabels = {
  'available': 'Disponible',
  'in_use': 'En uso',
  'maintenance': 'En mantenimiento',
  'inactive': 'Inactivo',
};

class VehicleListScreen extends StatefulWidget {
  const VehicleListScreen({super.key, this.businessId});

  final int? businessId;

  @override
  State<VehicleListScreen> createState() => _VehicleListScreenState();
}

class _VehicleListScreenState extends State<VehicleListScreen> {
  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addPostFrameCallback((_) => _refresh());
  }

  @override
  void didUpdateWidget(VehicleListScreen oldWidget) {
    super.didUpdateWidget(oldWidget);
    if (oldWidget.businessId != widget.businessId) _refresh();
  }

  void _refresh() {
    context.read<VehicleProvider>().fetchVehicles(businessId: widget.businessId);
  }

  @override
  Widget build(BuildContext context) {
    return Consumer<VehicleProvider>(
      builder: (context, provider, _) {
        return PaginatedListView<VehicleInfo>(
          controller: provider.list,
          unitLabel: 'vehiculos',
          placeholderHeight: 86,
          emptyIcon: Icons.local_shipping_outlined,
          emptyTitle: 'Sin vehiculos',
          emptyMessage:
              'Registra los vehiculos de tu flota para asignarlos a rutas.',
          itemBuilder: (context, vehicle, index) => _VehicleCard(vehicle: vehicle),
        );
      },
    );
  }
}

class _VehicleCard extends StatelessWidget {
  const _VehicleCard({required this.vehicle});

  final VehicleInfo vehicle;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    return AppCard(
      padding: const EdgeInsets.all(13),
      child: Row(
        children: [
          Container(
            width: 44,
            height: 44,
            alignment: Alignment.center,
            decoration: BoxDecoration(
              color: AppColors.primarySoft,
              borderRadius: AppRadius.mdAll,
            ),
            child: Icon(
              vehicleTypeIcons[vehicle.type] ?? Icons.local_shipping_outlined,
              size: 21,
              color: AppColors.primary,
            ),
          ),
          const SizedBox(width: 12),
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(vehicle.licensePlate, style: theme.textTheme.titleSmall),
                const SizedBox(height: 2),
                Text(
                  '${vehicle.brand} ${vehicle.model}  ·  ${vehicleTypeLabels[vehicle.type] ?? vehicle.type}',
                  style: theme.textTheme.labelSmall,
                  maxLines: 1,
                  overflow: TextOverflow.ellipsis,
                ),
                const SizedBox(height: 7),
                Row(
                  children: [
                    AppStatusChip(
                      dense: true,
                      label: vehicleStatusLabels[vehicle.status] ?? vehicle.status,
                      tone: vehicle.status == 'available'
                          ? AppStatusTone.success
                          : vehicle.status == 'maintenance'
                              ? AppStatusTone.warning
                              : AppStatusTone.neutral,
                    ),
                    const SizedBox(width: 7),
                    if (vehicle.weightCapacityKg != null)
                      Text(
                        '${AppFormat.number(vehicle.weightCapacityKg)} kg',
                        style: theme.textTheme.labelSmall,
                      ),
                  ],
                ),
              ],
            ),
          ),
        ],
      ),
    );
  }
}
