import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../../../../../shared/widgets/ui/ui.dart';
import '../../domain/entities.dart';
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
    context.read<RouteProvider>().fetchRoutes(
          businessId: widget.businessId,
          status: _status.isEmpty ? null : _status,
        );
  }

  void _onStatus(String value) {
    setState(() => _status = value);
    _refresh();
  }

  @override
  Widget build(BuildContext context) {
    return Column(
      children: [
        const SizedBox(height: 12),
        AppFilterChips(
          options: _statusOptions,
          selected: _status,
          onSelected: _onStatus,
        ),
        const SizedBox(height: 4),
        Expanded(
          child: Consumer<RouteProvider>(
            builder: (context, provider, _) {
              return PaginatedListView<RouteInfo>(
                controller: provider.list,
                unitLabel: 'rutas',
                placeholderHeight: 128,
                emptyIcon: Icons.alt_route_outlined,
                emptyTitle: 'Sin rutas',
                emptyMessage: _status.isEmpty
                    ? 'Arma una ruta para repartir con tu propia flota.'
                    : 'Ninguna ruta coincide con el filtro aplicado.',
                itemBuilder: (context, route, index) => RouteCard(
                  route: route,
                  onTap: () => Navigator.of(context).push(
                    MaterialPageRoute(
                      builder: (_) => RouteDetailScreen(
                        routeId: route.id,
                        businessId: widget.businessId,
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
