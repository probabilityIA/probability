import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:provider/provider.dart';

import '../../services/auth/login/ui/providers/login_provider.dart';
import '../navigation/app_modules.dart';
import '../theme/app_colors.dart';
import '../theme/app_tokens.dart';
import 'business_switcher_button.dart';

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
