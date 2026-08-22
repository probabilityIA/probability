import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:provider/provider.dart';
import '../../services/auth/login/ui/providers/login_provider.dart';
import '../navigation/app_modules.dart';
import '../theme/app_colors.dart';
import '../theme/app_tokens.dart';
import '../utils/formatters.dart';
import 'app_shell_scope.dart';
import 'network_avatar.dart';
import 'ui/ui.dart';

class AppShell extends StatefulWidget {
  const AppShell({super.key, required this.child});

  final Widget child;

  @override
  State<AppShell> createState() => _AppShellState();
}

class _AppShellState extends State<AppShell> {
  final GlobalKey<ScaffoldState> _scaffoldKey = GlobalKey<ScaffoldState>();

  void _openDrawer() => _scaffoldKey.currentState?.openDrawer();

  @override
  Widget build(BuildContext context) {
    final location = GoRouterState.of(context).matchedLocation;

    return Scaffold(
      key: _scaffoldKey,
      backgroundColor: AppColors.background,
      drawer: AppNavigationDrawer(location: location),
      body: AppShellScope(
        openDrawer: _openDrawer,
        location: location,
        child: widget.child,
      ),
      bottomNavigationBar: _BottomBar(location: location),
    );
  }
}

class _BottomBar extends StatelessWidget {
  const _BottomBar({required this.location});

  final String location;

  int get _index {
    var best = appBottomTabs.length - 1;
    var bestScore = 0;
    for (var i = 0; i < appBottomTabs.length; i++) {
      final score = appBottomTabs[i].matchScore(location);
      if (score > bestScore) {
        bestScore = score;
        best = i;
      }
    }
    return best;
  }

  @override
  Widget build(BuildContext context) {
    return Container(
      decoration: const BoxDecoration(
        color: AppColors.surface,
        border: Border(top: AppBorders.hairline),
      ),
      child: SafeArea(
        top: false,
        child: NavigationBar(
          selectedIndex: _index,
          onDestinationSelected: (index) {
            final tab = appBottomTabs[index];
            if (!tab.matches(location)) context.go(tab.route);
          },
          destinations: appBottomTabs
              .map((tab) => NavigationDestination(
                    icon: Icon(tab.icon),
                    selectedIcon: Icon(tab.activeIcon),
                    label: tab.label,
                  ))
              .toList(),
        ),
      ),
    );
  }
}

class AppNavigationDrawer extends StatelessWidget {
  const AppNavigationDrawer({super.key, required this.location});

  final String location;

  @override
  Widget build(BuildContext context) {
    final login = context.watch<LoginProvider>();

    return Drawer(
      child: Column(
        children: [
          _DrawerHeader(login: login),
          Expanded(
            child: ListView(
              padding: const EdgeInsets.fromLTRB(10, 8, 10, 16),
              children: [
                _DrawerItem(module: AppModules.dashboard, location: location),
                for (final group in AppModules.visibleGroups) ...[
                  Padding(
                    padding: const EdgeInsets.fromLTRB(12, 18, 12, 8),
                    child: Text(
                      group.title.toUpperCase(),
                      style: Theme.of(context).textTheme.labelSmall?.copyWith(
                            letterSpacing: 0.8,
                            fontWeight: FontWeight.w700,
                            color: AppColors.textDisabled,
                          ),
                    ),
                  ),
                  for (final module in group.modules)
                    _DrawerItem(module: module, location: location),
                ],
              ],
            ),
          ),
          const Divider(height: 1),
          Padding(
            padding: const EdgeInsets.fromLTRB(10, 8, 10, 12),
            child: ListTile(
              shape: RoundedRectangleBorder(borderRadius: AppRadius.mdAll),
              leading: const Icon(Icons.logout_rounded, size: 20, color: AppColors.error),
              title: Text(
                'Cerrar sesion',
                style: Theme.of(context).textTheme.titleSmall?.copyWith(color: AppColors.error),
              ),
              onTap: () {
                Navigator.pop(context);
                context.read<LoginProvider>().logout();
                context.go('/login');
              },
            ),
          ),
        ],
      ),
    );
  }
}

class _DrawerHeader extends StatelessWidget {
  const _DrawerHeader({required this.login});

  final LoginProvider login;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final user = login.user;
    final business = login.businesses.isNotEmpty ? login.businesses.first : null;

    return Container(
      width: double.infinity,
      padding: const EdgeInsets.fromLTRB(20, 0, 16, 18),
      decoration: const BoxDecoration(
        color: AppColors.surface,
        border: Border(bottom: AppBorders.hairline),
      ),
      child: SafeArea(
        bottom: false,
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            const SizedBox(height: 16),
            const AppLogo(height: 26),
            const SizedBox(height: 22),
            Row(
              children: [
                NetworkAvatar(
                  imageUrl: user?.avatarUrl,
                  fallbackText: user?.name ?? '?',
                  radius: 21,
                  backgroundColor: AppColors.primarySoft,
                  foregroundColor: AppColors.primary,
                ),
                const SizedBox(width: 12),
                Expanded(
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Text(
                        user?.name ?? 'Usuario',
                        style: theme.textTheme.titleSmall,
                        maxLines: 1,
                        overflow: TextOverflow.ellipsis,
                      ),
                      const SizedBox(height: 1),
                      Text(
                        user?.email ?? '',
                        style: theme.textTheme.labelSmall,
                        maxLines: 1,
                        overflow: TextOverflow.ellipsis,
                      ),
                    ],
                  ),
                ),
              ],
            ),
            if (business != null || login.isSuperAdmin) ...[
              const SizedBox(height: 14),
              _BusinessPill(
                name: login.isSuperAdmin ? 'Super admin' : (business?.name ?? ''),
                isSuperAdmin: login.isSuperAdmin,
              ),
            ],
          ],
        ),
      ),
    );
  }
}

class _BusinessPill extends StatelessWidget {
  const _BusinessPill({required this.name, required this.isSuperAdmin});

  final String name;
  final bool isSuperAdmin;

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 8),
      decoration: BoxDecoration(
        color: AppColors.surfaceMuted,
        borderRadius: AppRadius.mdAll,
      ),
      child: Row(
        children: [
          Container(
            width: 26,
            height: 26,
            alignment: Alignment.center,
            decoration: BoxDecoration(
              color: AppColors.primary,
              borderRadius: BorderRadius.circular(7),
            ),
            child: Text(
              AppFormat.initials(name),
              style: const TextStyle(
                fontFamily: 'Inter',
                fontSize: 11,
                fontWeight: FontWeight.w700,
                color: Colors.white,
              ),
            ),
          ),
          const SizedBox(width: 9),
          Expanded(
            child: Text(
              name,
              style: Theme.of(context).textTheme.titleSmall?.copyWith(fontSize: 13),
              maxLines: 1,
              overflow: TextOverflow.ellipsis,
            ),
          ),
          if (isSuperAdmin)
            const Icon(Icons.unfold_more_rounded, size: 18, color: AppColors.textMuted),
        ],
      ),
    );
  }
}

class _DrawerItem extends StatelessWidget {
  const _DrawerItem({required this.module, required this.location});

  final AppModule module;
  final String location;

  @override
  Widget build(BuildContext context) {
    final active = module.isActive(location);
    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 1),
      child: Material(
        color: active ? AppColors.primarySoft : Colors.transparent,
        borderRadius: AppRadius.mdAll,
        child: InkWell(
          borderRadius: AppRadius.mdAll,
          onTap: () {
            Navigator.pop(context);
            if (!active) context.go(module.route);
          },
          child: Padding(
            padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 11),
            child: Row(
              children: [
                Icon(
                  module.icon,
                  size: 20,
                  color: active ? AppColors.primary : AppColors.textMuted,
                ),
                const SizedBox(width: 12),
                Expanded(
                  child: Text(
                    module.label,
                    style: Theme.of(context).textTheme.titleSmall?.copyWith(
                          color: active ? AppColors.primaryDark : AppColors.textSecondary,
                          fontWeight: active ? FontWeight.w600 : FontWeight.w500,
                        ),
                  ),
                ),
                if (module.stage.showsBadge) const ModuleStageBadge(),
              ],
            ),
          ),
        ),
      ),
    );
  }
}
