import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../../../../../shared/theme/app_colors.dart';
import '../../../../../shared/theme/app_tokens.dart';
import '../../../../../shared/widgets/network_avatar.dart';
import '../../../../../shared/widgets/ui/ui.dart';
import '../../domain/entities.dart';
import '../providers/user_provider.dart';

class UserListScreen extends StatefulWidget {
  const UserListScreen({super.key});

  @override
  State<UserListScreen> createState() => _UserListScreenState();
}

class _UserListScreenState extends State<UserListScreen> {
  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addPostFrameCallback((_) => _refresh());
  }

  void _refresh() {
    context.read<UserProvider>().fetchUsers();
  }

  @override
  Widget build(BuildContext context) {
    return Consumer<UserProvider>(
      builder: (context, provider, _) {
        if (provider.isLoading && provider.users.isEmpty) {
          return const AppListSkeleton();
        }
        if (provider.error != null && provider.users.isEmpty) {
          return AppErrorState(message: provider.error!, onRetry: _refresh);
        }
        if (provider.users.isEmpty) {
          return AppEmptyState(
            icon: Icons.people_alt_outlined,
            title: 'Sin usuarios',
            message: 'Invita a tu equipo para que opere contigo.',
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
            itemCount: provider.users.length,
            separatorBuilder: (context, index) => const SizedBox(height: 10),
            itemBuilder: (context, index) => _UserCard(user: provider.users[index]),
          ),
        );
      },
    );
  }
}

class _UserCard extends StatelessWidget {
  const _UserCard({required this.user});

  final User user;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    return AppCard(
      padding: const EdgeInsets.all(13),
      child: Row(
        children: [
          NetworkAvatar(imageUrl: user.avatarUrl, fallbackText: user.name, radius: 22),
          const SizedBox(width: 12),
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  user.name,
                  style: theme.textTheme.titleSmall,
                  maxLines: 1,
                  overflow: TextOverflow.ellipsis,
                ),
                const SizedBox(height: 2),
                Text(
                  user.email,
                  style: theme.textTheme.labelSmall,
                  maxLines: 1,
                  overflow: TextOverflow.ellipsis,
                ),
                const SizedBox(height: 7),
                Wrap(
                  spacing: 6,
                  runSpacing: 6,
                  children: [
                    for (final assignment in user.businessRoleAssignments)
                      if ((assignment.roleName ?? '').isNotEmpty)
                        AppStatusChip(
                          dense: true,
                          label: assignment.roleName!,
                          tone: AppStatusTone.brand,
                        ),
                    AppStatusChip(
                      dense: true,
                      label: user.isActive ? 'Activo' : 'Inactivo',
                      tone: user.isActive ? AppStatusTone.success : AppStatusTone.neutral,
                    ),
                    if (user.isSuperUser)
                      const AppStatusChip(
                        dense: true,
                        label: 'Super admin',
                        tone: AppStatusTone.warning,
                        icon: Icons.shield_outlined,
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
