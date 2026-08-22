import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:provider/provider.dart';
import '../../services/auth/login/ui/providers/login_provider.dart';
import '../navigation/app_modules.dart';
import '../theme/app_colors.dart';
import '../theme/app_tokens.dart';
import 'ui/ui.dart';

class ModulesScreen extends StatelessWidget {
  const ModulesScreen({super.key});

  @override
  Widget build(BuildContext context) {
    final login = context.watch<LoginProvider>();

    return AppScaffold(
      title: 'Modulos',
      subtitle: 'Todo lo que puedes gestionar',
      actions: [
        IconButton(
          icon: const Icon(Icons.person_outline, size: 22),
          tooltip: 'Perfil',
          onPressed: () => context.go('/profile'),
        ),
      ],
      body: ListView(
        padding: AppSpacing.page,
        children: [
          for (final group in AppModules.visibleGroups) ...[
            AppSectionHeader(title: group.title),
            GridView.builder(
              shrinkWrap: true,
              physics: const NeverScrollableScrollPhysics(),
              padding: EdgeInsets.zero,
              gridDelegate: const SliverGridDelegateWithFixedCrossAxisCount(
                crossAxisCount: 2,
                mainAxisSpacing: 12,
                crossAxisSpacing: 12,
                mainAxisExtent: 116,
              ),
              itemCount: group.modules.length,
              itemBuilder: (context, index) => _ModuleTile(module: group.modules[index]),
            ),
            const SizedBox(height: 20),
          ],
          const SizedBox(height: 4),
          AppCard(
            onTap: () {
              context.read<LoginProvider>().logout();
              context.go('/login');
            },
            child: Row(
              children: [
                const Icon(Icons.logout_rounded, size: 20, color: AppColors.error),
                const SizedBox(width: 12),
                Expanded(
                  child: Text(
                    'Cerrar sesion',
                    style: Theme.of(context).textTheme.titleSmall?.copyWith(color: AppColors.error),
                  ),
                ),
              ],
            ),
          ),
          const SizedBox(height: 14),
          Center(
            child: Text(
              login.user?.email ?? '',
              style: Theme.of(context).textTheme.labelSmall,
            ),
          ),
          const SizedBox(height: 24),
        ],
      ),
    );
  }
}

class _ModuleTile extends StatelessWidget {
  const _ModuleTile({required this.module});

  final AppModule module;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    return AppCard(
      padding: const EdgeInsets.all(14),
      onTap: () => context.go(module.route),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        mainAxisAlignment: MainAxisAlignment.spaceBetween,
        children: [
          Container(
            width: 36,
            height: 36,
            decoration: BoxDecoration(
              color: AppColors.primarySoft,
              borderRadius: AppRadius.smAll,
            ),
            child: Icon(module.icon, size: 19, color: AppColors.primary),
          ),
          const Spacer(),
          if (module.stage.showsBadge) ...[
            const ModuleStageBadge(),
            const SizedBox(height: 6),
          ],
          Text(
            module.label,
            style: theme.textTheme.titleSmall,
            maxLines: 1,
            overflow: TextOverflow.ellipsis,
          ),
          if (module.description != null) ...[
            const SizedBox(height: 3),
            Text(
              module.description!,
              style: theme.textTheme.labelSmall,
              maxLines: 2,
              overflow: TextOverflow.ellipsis,
            ),
          ],
        ],
      ),
    );
  }
}
