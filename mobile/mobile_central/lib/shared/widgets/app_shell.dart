import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:provider/provider.dart';

import '../../services/auth/login/ui/providers/login_provider.dart';
import '../navigation/app_modules.dart';
import '../theme/app_colors.dart';
import '../theme/app_tokens.dart';
import 'business_switcher_button.dart';
import 'core_button.dart';

class AppShell extends StatelessWidget {
  const AppShell({super.key, required this.child});

  final Widget child;

  @override
  Widget build(BuildContext context) {
    final location = GoRouterState.of(context).matchedLocation;
    final isSuperAdmin = context.watch<LoginProvider>().isSuperAdmin;

    return Scaffold(
      backgroundColor: AppColors.background,
      body: Stack(
        children: [
          child,
          if (isSuperAdmin) const DraggableBusinessButton(),
        ],
      ),
      bottomNavigationBar: _BottomBar(
        location: location,
        onCore: () => context.push('/core'),
      ),
    );
  }
}

class _BottomBar extends StatelessWidget {
  const _BottomBar({required this.location, required this.onCore});

  final String location;
  final VoidCallback onCore;

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
    final left = appBottomTabs.take(2).toList();
    final right = appBottomTabs.skip(2).toList();

    return Container(
      decoration: const BoxDecoration(
        color: AppColors.surface,
        border: Border(top: AppBorders.hairline),
      ),
      child: SafeArea(
        top: false,
        child: SizedBox(
          height: 62,
          child: Stack(
            clipBehavior: Clip.none,
            alignment: Alignment.topCenter,
            children: [
              Row(
                children: [
                  for (var i = 0; i < left.length; i++)
                    Expanded(
                      child: _Tab(
                        tab: left[i],
                        active: _index == i,
                        onTap: () => _go(context, left[i]),
                      ),
                    ),
                  const SizedBox(width: 76),
                  for (var i = 0; i < right.length; i++)
                    Expanded(
                      child: _Tab(
                        tab: right[i],
                        active: _index == left.length + i,
                        onTap: () => _go(context, right[i]),
                      ),
                    ),
                ],
              ),
              Positioned(
                top: -30,
                child: CoreButton(onTap: onCore),
              ),
            ],
          ),
        ),
      ),
    );
  }

  void _go(BuildContext context, AppBottomTab tab) {
    if (!tab.matches(location)) context.go(tab.route);
  }
}

class _Tab extends StatelessWidget {
  const _Tab({required this.tab, required this.active, required this.onTap});

  final AppBottomTab tab;
  final bool active;
  final VoidCallback onTap;

  @override
  Widget build(BuildContext context) {
    final scheme = Theme.of(context).colorScheme;
    final color = active ? scheme.primary : AppColors.textDisabled;

    return InkResponse(
      onTap: onTap,
      radius: 42,
      child: Column(
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          Icon(active ? tab.activeIcon : tab.icon, size: 22, color: color),
          const SizedBox(height: 3),
          Text(
            tab.label,
            style: TextStyle(
              fontFamily: 'Inter',
              fontSize: 11,
              fontWeight: active ? FontWeight.w600 : FontWeight.w500,
              color: color,
            ),
          ),
        ],
      ),
    );
  }
}
