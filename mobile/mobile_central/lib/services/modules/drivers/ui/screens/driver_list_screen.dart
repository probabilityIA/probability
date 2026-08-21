import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../../../../../shared/theme/app_colors.dart';
import '../../../../../shared/theme/app_tokens.dart';
import '../../../../../shared/utils/formatters.dart';
import '../../../../../shared/widgets/network_avatar.dart';
import '../../../../../shared/widgets/ui/ui.dart';
import '../../domain/entities.dart';
import '../providers/drivers_provider.dart';

const Map<String, String> driverStatusLabels = {
  'available': 'Disponible',
  'on_route': 'En ruta',
  'inactive': 'Inactivo',
  'off_duty': 'Fuera de turno',
};

class DriverListScreen extends StatefulWidget {
  const DriverListScreen({super.key, this.businessId});

  final int? businessId;

  @override
  State<DriverListScreen> createState() => _DriverListScreenState();
}

class _DriverListScreenState extends State<DriverListScreen> {
  String _status = '';

  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addPostFrameCallback((_) => _refresh());
  }

  @override
  void didUpdateWidget(DriverListScreen oldWidget) {
    super.didUpdateWidget(oldWidget);
    if (oldWidget.businessId != widget.businessId) _refresh();
  }

  void _refresh() {
    context.read<DriverProvider>().fetchDrivers(businessId: widget.businessId);
  }

  @override
  Widget build(BuildContext context) {
    return Column(
      children: [
        const SizedBox(height: 12),
        AppFilterChips(
          options: const [
            (value: '', label: 'Todos'),
            (value: 'available', label: 'Disponibles'),
            (value: 'on_route', label: 'En ruta'),
            (value: 'inactive', label: 'Inactivos'),
          ],
          selected: _status,
          onSelected: (value) => setState(() => _status = value),
        ),
        const SizedBox(height: 4),
        Expanded(
          child: Consumer<DriverProvider>(
            builder: (context, provider, _) {
              if (provider.isLoading && provider.drivers.isEmpty) {
                return const AppListSkeleton();
              }
              if (provider.error != null && provider.drivers.isEmpty) {
                return AppErrorState(message: provider.error!, onRetry: _refresh);
              }

              final rows = _status.isEmpty
                  ? provider.drivers
                  : provider.drivers.where((d) => d.status == _status).toList();

              if (rows.isEmpty) {
                return AppEmptyState(
                  icon: Icons.badge_outlined,
                  title: 'Sin conductores',
                  message: 'Registra conductores para armar rutas con tu flota.',
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
                  itemCount: rows.length,
                  separatorBuilder: (context, index) => const SizedBox(height: 10),
                  itemBuilder: (context, index) => _DriverCard(driver: rows[index]),
                ),
              );
            },
          ),
        ),
      ],
    );
  }
}

class _DriverCard extends StatelessWidget {
  const _DriverCard({required this.driver});

  final DriverInfo driver;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final name = '${driver.firstName} ${driver.lastName}'.trim();

    return AppCard(
      padding: const EdgeInsets.all(13),
      child: Row(
        children: [
          NetworkAvatar(
            imageUrl: driver.photoUrl,
            fallbackText: name,
            radius: 22,
          ),
          const SizedBox(width: 12),
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  name,
                  style: theme.textTheme.titleSmall,
                  maxLines: 1,
                  overflow: TextOverflow.ellipsis,
                ),
                const SizedBox(height: 2),
                Text(
                  '${driver.phone}  ·  Licencia ${driver.licenseType}',
                  style: theme.textTheme.labelSmall,
                  maxLines: 1,
                  overflow: TextOverflow.ellipsis,
                ),
                const SizedBox(height: 7),
                Row(
                  children: [
                    AppStatusChip(
                      dense: true,
                      label: driverStatusLabels[driver.status] ?? driver.status,
                      tone: AppStatusChip.toneFromCode(driver.status),
                    ),
                    const SizedBox(width: 7),
                    Text(
                      'Vence ${AppFormat.date(AppFormat.parseDate(driver.licenseExpiry))}',
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
