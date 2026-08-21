import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../../../../../shared/theme/app_colors.dart';
import '../../../../../shared/theme/app_tokens.dart';
import '../../../../../shared/widgets/ui/ui.dart';
import '../providers/route_provider.dart';
import '../widgets/route_widgets.dart';
import 'route_detail_screen.dart';

class RouteListScreen extends StatefulWidget {
  const RouteListScreen({super.key, this.businessId});

  final int? businessId;

  @override
  State<RouteListScreen> createState() => _RouteListScreenState();
}

class _RouteListScreenState extends State<RouteListScreen> {
  String _status = '';

  static const List<({String value, String label})> _statusOptions = [
    (value: '', label: 'Todas'),
    (value: 'draft', label: 'Borrador'),
    (value: 'in_progress', label: 'En curso'),
    (value: 'completed', label: 'Completadas'),
  ];

  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addPostFrameCallback((_) => _refresh());
  }

  @override
  void didUpdateWidget(RouteListScreen oldWidget) {
    super.didUpdateWidget(oldWidget);
    if (oldWidget.businessId != widget.businessId) _refresh();
  }

  void _refresh() {
    context.read<RouteProvider>().fetchRoutes(businessId: widget.businessId);
  }

  @override
  Widget build(BuildContext context) {
    return Column(
      children: [
        const SizedBox(height: 12),
        AppFilterChips(
          options: _statusOptions,
          selected: _status,
          onSelected: (value) => setState(() => _status = value),
        ),
        const SizedBox(height: 4),
        Expanded(
          child: Consumer<RouteProvider>(
            builder: (context, provider, _) {
              if (provider.isLoading && provider.routes.isEmpty) {
                return const AppListSkeleton();
              }
              if (provider.error != null && provider.routes.isEmpty) {
                return AppErrorState(message: provider.error!, onRetry: _refresh);
              }

              final rows = _status.isEmpty
                  ? provider.routes
                  : provider.routes.where((r) => r.status == _status).toList();

              if (rows.isEmpty) {
                return AppEmptyState(
                  icon: Icons.alt_route_outlined,
                  title: 'Sin rutas',
                  message: 'Arma una ruta para repartir con tu propia flota.',
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
                  itemBuilder: (context, index) => RouteCard(
                    route: rows[index],
                    onTap: () => Navigator.of(context).push(
                      MaterialPageRoute(
                        builder: (_) => RouteDetailScreen(
                          routeId: rows[index].id,
                          businessId: widget.businessId,
                        ),
                      ),
                    ),
                  ),
                ),
              );
            },
          ),
        ),
      ],
    );
  }
}
