import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../../../../../shared/theme/app_colors.dart';
import '../../../../../shared/theme/app_tokens.dart';
import '../../../../../shared/utils/formatters.dart';
import '../../../../../shared/widgets/ui/ui.dart';
import '../../domain/entities.dart';
import '../providers/route_provider.dart';
import '../widgets/route_widgets.dart';

class RouteDetailScreen extends StatefulWidget {
  const RouteDetailScreen({super.key, required this.routeId, this.businessId});

  final int routeId;
  final int? businessId;

  @override
  State<RouteDetailScreen> createState() => _RouteDetailScreenState();
}

class _RouteDetailScreenState extends State<RouteDetailScreen> {
  RouteDetail? _route;
  bool _loading = true;
  bool _busy = false;

  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addPostFrameCallback((_) => _load());
  }

  Future<void> _load() async {
    setState(() => _loading = true);
    final provider = context.read<RouteProvider>();
    await provider.fetchRouteDetail(widget.routeId, businessId: widget.businessId);
    if (!mounted) return;
    setState(() {
      _route = provider.selectedRoute;
      _loading = false;
    });
  }

  Future<void> _start() async {
    setState(() => _busy = true);
    final provider = context.read<RouteProvider>();
    final ok = await provider.startRoute(widget.routeId, businessId: widget.businessId) != null;
    if (!mounted) return;
    setState(() => _busy = false);
    ScaffoldMessenger.of(context).showSnackBar(
      SnackBar(content: Text(ok ? 'Ruta iniciada' : (provider.error ?? 'No se pudo iniciar'))),
    );
    if (ok) _load();
  }

  Future<void> _complete() async {
    setState(() => _busy = true);
    final provider = context.read<RouteProvider>();
    final ok = await provider.completeRoute(widget.routeId, businessId: widget.businessId) != null;
    if (!mounted) return;
    setState(() => _busy = false);
    ScaffoldMessenger.of(context).showSnackBar(
      SnackBar(content: Text(ok ? 'Ruta completada' : (provider.error ?? 'No se pudo completar'))),
    );
    if (ok) _load();
  }

  @override
  Widget build(BuildContext context) {
    final route = _route;

    return AppScaffold(
      title: route == null ? 'Ruta' : 'Ruta ${route.id}',
      subtitle: route == null ? null : routeStatusLabel(route.status),
      onBack: () => Navigator.of(context).pop(),
      body: _loading
          ? const AppLoading(label: 'Cargando ruta')
          : route == null
              ? AppErrorState(
                  message: 'No se pudo cargar la ruta',
                  onRetry: _load,
                )
              : ListView(
                  padding: AppSpacing.page,
                  children: [
                    _SummaryCard(route: route),
                    const SizedBox(height: 18),
                    const AppSectionHeader(title: 'Asignacion'),
                    _AssignmentCard(route: route),
                    const SizedBox(height: 18),
                    AppSectionHeader(
                      title: 'Paradas',
                      subtitle: '${route.completedStops} de ${route.totalStops} completadas',
                    ),
                    if (route.stops.isEmpty)
                      AppCard(
                        child: Text(
                          'Esta ruta no tiene paradas asignadas',
                          style: Theme.of(context).textTheme.bodySmall,
                        ),
                      )
                    else
                      AppCard(
                        child: Column(
                          children: [
                            for (var i = 0; i < route.stops.length; i++)
                              RouteStopTile(
                                stop: route.stops[i],
                                isLast: i == route.stops.length - 1,
                              ),
                          ],
                        ),
                      ),
                    const SizedBox(height: 20),
                    if (route.status == 'draft' || route.status == 'planned')
                      FilledButton.icon(
                        onPressed: _busy ? null : _start,
                        style: FilledButton.styleFrom(minimumSize: const Size(0, 50)),
                        icon: const Icon(Icons.play_arrow_rounded, size: 19),
                        label: const Text('Iniciar ruta'),
                      ),
                    if (route.status == 'in_progress')
                      FilledButton.icon(
                        onPressed: _busy ? null : _complete,
                        style: FilledButton.styleFrom(minimumSize: const Size(0, 50)),
                        icon: const Icon(Icons.flag_rounded, size: 19),
                        label: const Text('Completar ruta'),
                      ),
                    const SizedBox(height: 28),
                  ],
                ),
    );
  }
}

class _SummaryCard extends StatelessWidget {
  const _SummaryCard({required this.route});

  final RouteDetail route;

  @override
  Widget build(BuildContext context) {
    final progress = route.totalStops == 0
        ? 0.0
        : (route.completedStops / route.totalStops).clamp(0.0, 1.0);

    return AppCard(
      child: Column(
        children: [
          Row(
            children: [
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(
                      AppFormat.date(AppFormat.parseDate(route.date)),
                      style: Theme.of(context).textTheme.titleMedium,
                    ),
                    const SizedBox(height: 3),
                    Text(
                      route.originAddress ?? 'Sin origen',
                      style: Theme.of(context).textTheme.bodySmall,
                      maxLines: 1,
                      overflow: TextOverflow.ellipsis,
                    ),
                  ],
                ),
              ),
              AppStatusChip(
                label: routeStatusLabel(route.status),
                tone: AppStatusChip.toneFromCode(route.status),
              ),
            ],
          ),
          const SizedBox(height: 14),
          ClipRRect(
            borderRadius: BorderRadius.circular(3),
            child: LinearProgressIndicator(
              value: progress,
              minHeight: 5,
              backgroundColor: AppColors.surfaceMuted,
              valueColor: const AlwaysStoppedAnimation(AppColors.primary),
            ),
          ),
          const SizedBox(height: 14),
          Row(
            children: [
              Expanded(
                child: _Metric(
                  label: 'Paradas',
                  value: '${route.completedStops}/${route.totalStops}',
                ),
              ),
              Expanded(
                child: _Metric(
                  label: 'Distancia',
                  value: route.totalDistanceKm == null
                      ? '-'
                      : '${AppFormat.number(route.totalDistanceKm)} km',
                ),
              ),
              Expanded(
                child: _Metric(
                  label: 'Duracion',
                  value: route.totalDurationMin == null
                      ? '-'
                      : '${AppFormat.number(route.totalDurationMin)} min',
                ),
              ),
            ],
          ),
        ],
      ),
    );
  }
}

class _Metric extends StatelessWidget {
  const _Metric({required this.label, required this.value});

  final String label;
  final String value;

  @override
  Widget build(BuildContext context) {
    return Column(
      children: [
        Text(value, style: Theme.of(context).textTheme.titleSmall),
        const SizedBox(height: 2),
        Text(label, style: Theme.of(context).textTheme.labelSmall),
      ],
    );
  }
}

class _AssignmentCard extends StatelessWidget {
  const _AssignmentCard({required this.route});

  final RouteDetail route;

  @override
  Widget build(BuildContext context) {
    return AppCard(
      child: Column(
        children: [
          AppKeyValueRow(label: 'Conductor', value: route.driverName ?? 'Sin asignar'),
          const Divider(height: 18),
          AppKeyValueRow(label: 'Vehiculo', value: route.vehiclePlate ?? 'Sin asignar'),
          const Divider(height: 18),
          AppKeyValueRow(
            label: 'Inicio real',
            value: AppFormat.dateTime(AppFormat.parseDate(route.actualStartTime)),
          ),
          const Divider(height: 18),
          AppKeyValueRow(
            label: 'Fin real',
            value: AppFormat.dateTime(AppFormat.parseDate(route.actualEndTime)),
          ),
        ],
      ),
    );
  }
}
