import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../../../../../shared/theme/app_colors.dart';
import '../../../../../shared/theme/app_tokens.dart';
import '../../../../../shared/widgets/ui/ui.dart';
import '../../domain/entities.dart';
import '../providers/permission_provider.dart';

const Map<String, String> permissionActionLabels = {
  'read': 'Ver',
  'create': 'Crear',
  'update': 'Editar',
  'delete': 'Eliminar',
  'export': 'Exportar',
};

class PermissionListScreen extends StatefulWidget {
  const PermissionListScreen({super.key});

  @override
  State<PermissionListScreen> createState() => _PermissionListScreenState();
}

class _PermissionListScreenState extends State<PermissionListScreen> {
  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addPostFrameCallback((_) => _refresh());
  }

  void _refresh() {
    context.read<PermissionProvider>().fetchPermissions();
  }

  @override
  Widget build(BuildContext context) {
    return Consumer<PermissionProvider>(
      builder: (context, provider, _) {
        if (provider.isLoading && provider.permissions.isEmpty) {
          return const AppListSkeleton();
        }
        if (provider.error != null && provider.permissions.isEmpty) {
          return AppErrorState(message: provider.error!, onRetry: _refresh);
        }
        if (provider.permissions.isEmpty) {
          return AppEmptyState(
            icon: Icons.lock_outline_rounded,
            title: 'Sin permisos',
            message: 'No hay permisos definidos para este tipo de negocio.',
            actionLabel: 'Actualizar',
            onAction: _refresh,
          );
        }

        final grouped = <String, List<Permission>>{};
        for (final permission in provider.permissions) {
          grouped.putIfAbsent(permission.resource ?? 'otros', () => []).add(permission);
        }
        final resources = grouped.keys.toList()..sort();

        return RefreshIndicator(
          onRefresh: () async => _refresh(),
          color: AppColors.primary,
          child: ListView.builder(
            physics: const AlwaysScrollableScrollPhysics(),
            padding: AppSpacing.page,
            cacheExtent: 600,
            addAutomaticKeepAlives: false,
            itemCount: resources.length,
            itemBuilder: (context, index) {
              final resource = resources[index];
              return Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  AppSectionHeader(title: resource),
                  AppCard(
                    child: Wrap(
                      spacing: 7,
                      runSpacing: 7,
                      children: grouped[resource]!
                          .map((permission) => AppStatusChip(
                                label: permissionActionLabels[permission.action] ??
                                    permission.action ??
                                    '-',
                                tone: AppStatusTone.brand,
                              ))
                          .toList(),
                    ),
                  ),
                  SizedBox(height: index == resources.length - 1 ? 28 : 16),
                ],
              );
            },
          ),
        );
      },
    );
  }
}
