import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../../../../../shared/theme/app_colors.dart';
import '../../../../../shared/theme/app_tokens.dart';
import '../../../../../shared/widgets/ui/ui.dart';
import '../../domain/entities.dart';
import '../providers/role_provider.dart';

class RoleListScreen extends StatefulWidget {
  const RoleListScreen({super.key});

  @override
  State<RoleListScreen> createState() => _RoleListScreenState();
}

class _RoleListScreenState extends State<RoleListScreen> {
  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addPostFrameCallback((_) => _refresh());
  }

  void _refresh() {
    context.read<RoleProvider>().fetchRoles();
  }

  @override
  Widget build(BuildContext context) {
    return Consumer<RoleProvider>(
      builder: (context, provider, _) {
        if (provider.isLoading && provider.roles.isEmpty) {
          return const AppListSkeleton();
        }
        if (provider.error != null && provider.roles.isEmpty) {
          return AppErrorState(message: provider.error!, onRetry: _refresh);
        }
        if (provider.roles.isEmpty) {
          return AppEmptyState(
            icon: Icons.admin_panel_settings_outlined,
            title: 'Sin roles',
            message: 'Crea roles para controlar que puede hacer cada persona.',
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
            itemCount: provider.roles.length,
            separatorBuilder: (context, index) => const SizedBox(height: 10),
            itemBuilder: (context, index) => _RoleCard(role: provider.roles[index]),
          ),
        );
      },
    );
  }
}

class _RoleCard extends StatelessWidget {
  const _RoleCard({required this.role});

  final Role role;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    return AppCard(
      padding: const EdgeInsets.all(14),
      child: Row(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Container(
            width: 40,
            height: 40,
            alignment: Alignment.center,
            decoration: BoxDecoration(
              color: AppColors.primarySoft,
              borderRadius: AppRadius.mdAll,
            ),
            child: Text(
              '${role.level ?? 0}',
              style: const TextStyle(
                fontFamily: 'Inter',
                fontSize: 15,
                fontWeight: FontWeight.w700,
                color: AppColors.primary,
              ),
            ),
          ),
          const SizedBox(width: 12),
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Row(
                  children: [
                    Expanded(
                      child: Text(
                        role.name,
                        style: theme.textTheme.titleSmall,
                        maxLines: 1,
                        overflow: TextOverflow.ellipsis,
                      ),
                    ),
                    if (role.isSystem)
                      const AppStatusChip(
                        dense: true,
                        label: 'Del sistema',
                        tone: AppStatusTone.neutral,
                      ),
                  ],
                ),
                if ((role.description ?? '').isNotEmpty) ...[
                  const SizedBox(height: 3),
                  Text(
                    role.description!,
                    style: theme.textTheme.labelSmall,
                    maxLines: 2,
                    overflow: TextOverflow.ellipsis,
                  ),
                ],
                const SizedBox(height: 8),
                Row(
                  children: [
                    const Icon(Icons.people_outline, size: 13, color: AppColors.textDisabled),
                    const SizedBox(width: 5),
                    Text(
                      role.code ?? '-',
                      style: theme.textTheme.labelSmall,
                    ),
                    const SizedBox(width: 12),
                    Text('Nivel ${role.level ?? 0}', style: theme.textTheme.labelSmall),
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
