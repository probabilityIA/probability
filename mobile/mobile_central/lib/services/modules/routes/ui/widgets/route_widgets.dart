import 'package:flutter/material.dart';
import '../../../../../shared/theme/app_colors.dart';
import '../../../../../shared/theme/app_tokens.dart';
import '../../../../../shared/utils/formatters.dart';
import '../../../../../shared/widgets/ui/ui.dart';
import '../../domain/entities.dart';

const Map<String, String> routeStatusLabels = {
  'draft': 'Borrador',
  'planned': 'Planeada',
  'in_progress': 'En curso',
  'completed': 'Completada',
  'cancelled': 'Cancelada',
};

const Map<String, String> stopStatusLabels = {
  'pending': 'Pendiente',
  'in_transit': 'En camino',
  'delivered': 'Entregada',
  'failed': 'Fallida',
};

String routeStatusLabel(String code) => routeStatusLabels[code] ?? code;

class RouteCard extends StatelessWidget {
  const RouteCard({super.key, required this.route, this.onTap});

  final RouteInfo route;
  final VoidCallback? onTap;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final progress = route.totalStops == 0
        ? 0.0
        : (route.completedStops / route.totalStops).clamp(0.0, 1.0);

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
                  color: AppColors.primarySoft,
                  borderRadius: AppRadius.mdAll,
                ),
                child: const Icon(Icons.alt_route_outlined, size: 20, color: AppColors.primary),
              ),
              const SizedBox(width: 12),
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(
                      'Ruta ${route.id}',
                      style: theme.textTheme.titleSmall,
                      maxLines: 1,
                      overflow: TextOverflow.ellipsis,
                    ),
                    const SizedBox(height: 2),
                    Text(
                      '${AppFormat.date(AppFormat.parseDate(route.date))}  ·  ${route.originAddress ?? ''}',
                      style: theme.textTheme.labelSmall,
                      maxLines: 1,
                      overflow: TextOverflow.ellipsis,
                    ),
                  ],
                ),
              ),
              AppStatusChip(
                dense: true,
                label: routeStatusLabel(route.status),
                tone: AppStatusChip.toneFromCode(route.status),
              ),
            ],
          ),
          const SizedBox(height: 12),
          ClipRRect(
            borderRadius: BorderRadius.circular(3),
            child: LinearProgressIndicator(
              value: progress,
              minHeight: 5,
              backgroundColor: AppColors.surfaceMuted,
              valueColor: const AlwaysStoppedAnimation(AppColors.primary),
            ),
          ),
          const SizedBox(height: 9),
          Row(
            children: [
              Text(
                '${route.completedStops} de ${route.totalStops} paradas',
                style: theme.textTheme.labelSmall,
              ),
              if (route.failedStops > 0) ...[
                const SizedBox(width: 8),
                Text(
                  '${route.failedStops} fallidas',
                  style: theme.textTheme.labelSmall?.copyWith(color: AppColors.error),
                ),
              ],
              const Spacer(),
              if (route.totalDistanceKm != null)
                Text(
                  '${AppFormat.number(route.totalDistanceKm)} km',
                  style: theme.textTheme.labelSmall,
                ),
            ],
          ),
          const SizedBox(height: 10),
          Row(
            children: [
              const Icon(Icons.person_outline, size: 13, color: AppColors.textDisabled),
              const SizedBox(width: 5),
              Expanded(
                child: Text(
                  route.driverName ?? 'Sin conductor',
                  style: theme.textTheme.labelSmall,
                  maxLines: 1,
                  overflow: TextOverflow.ellipsis,
                ),
              ),
              const Icon(Icons.local_shipping_outlined, size: 13, color: AppColors.textDisabled),
              const SizedBox(width: 5),
              Text(route.vehiclePlate ?? '-', style: theme.textTheme.labelSmall),
            ],
          ),
        ],
      ),
    );
  }
}

class RouteStopTile extends StatelessWidget {
  const RouteStopTile({super.key, required this.stop, required this.isLast});

  final RouteStopInfo stop;
  final bool isLast;

  Color get _color {
    switch (stop.status) {
      case 'delivered':
        return AppColors.success;
      case 'failed':
        return AppColors.error;
      case 'in_transit':
        return AppColors.info;
      default:
        return AppColors.borderStrong;
    }
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    return IntrinsicHeight(
      child: Row(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Column(
            children: [
              Container(
                width: 26,
                height: 26,
                alignment: Alignment.center,
                decoration: BoxDecoration(
                  color: _color,
                  shape: BoxShape.circle,
                ),
                child: Text(
                  '${stop.sequence}',
                  style: const TextStyle(
                    fontFamily: 'Inter',
                    fontSize: 11,
                    fontWeight: FontWeight.w700,
                    color: Colors.white,
                  ),
                ),
              ),
              if (!isLast)
                Expanded(
                  child: Container(width: 2, color: AppColors.border),
                ),
            ],
          ),
          const SizedBox(width: 12),
          Expanded(
            child: Padding(
              padding: EdgeInsets.only(bottom: isLast ? 0 : 18),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Row(
                    children: [
                      Expanded(
                        child: Text(
                          stop.customerName,
                          style: theme.textTheme.titleSmall,
                          maxLines: 1,
                          overflow: TextOverflow.ellipsis,
                        ),
                      ),
                      AppStatusChip(
                        dense: true,
                        label: stopStatusLabels[stop.status] ?? stop.status,
                        tone: AppStatusChip.toneFromCode(stop.status),
                      ),
                    ],
                  ),
                  const SizedBox(height: 3),
                  Text(
                    stop.address,
                    style: theme.textTheme.labelSmall,
                    maxLines: 2,
                    overflow: TextOverflow.ellipsis,
                  ),
                  const SizedBox(height: 4),
                  Text(
                    [stop.city, stop.customerPhone]
                        .where((e) => (e ?? '').isNotEmpty)
                        .join('  \u00B7  '),
                    style: theme.textTheme.labelSmall,
                  ),
                  if ((stop.failureReason ?? '').isNotEmpty) ...[
                    const SizedBox(height: 7),
                    Container(
                      padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 8),
                      decoration: BoxDecoration(
                        color: AppColors.errorSoft,
                        borderRadius: AppRadius.smAll,
                      ),
                      child: Text(
                        stop.failureReason!,
                        style: theme.textTheme.labelSmall?.copyWith(
                          color: const Color(0xFFB91C1C),
                          fontWeight: FontWeight.w500,
                        ),
                      ),
                    ),
                  ],
                ],
              ),
            ),
          ),
        ],
      ),
    );
  }
}
