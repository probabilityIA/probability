import 'package:flutter/material.dart';
import '../../theme/app_colors.dart';
import '../../theme/app_tokens.dart';
import 'app_logo.dart';

class AppScaffold extends StatelessWidget {
  const AppScaffold({
    super.key,
    required this.body,
    this.title,
    this.subtitle,
    this.actions = const [],
    this.bottom,
    this.floatingActionButton,
    this.showLogo = false,
    this.showMenu = true,
    this.onBack,
    this.backgroundColor,
    this.padBottomForNav = true,
  });

  final Widget body;
  final String? title;
  final String? subtitle;
  final List<Widget> actions;
  final PreferredSizeWidget? bottom;
  final Widget? floatingActionButton;
  final bool showLogo;
  final bool showMenu;
  final VoidCallback? onBack;
  final Color? backgroundColor;
  final bool padBottomForNav;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final isStacked = onBack != null;

    if (isStacked) {
      return Scaffold(
        backgroundColor: backgroundColor ?? AppColors.background,
        appBar: AppBar(
          toolbarHeight: subtitle == null ? 58 : 66,
          leadingWidth: 56,
          leading: IconButton(
            icon: const Icon(Icons.arrow_back, size: 22),
            onPressed: onBack,
            tooltip: 'Volver',
          ),
          automaticallyImplyLeading: false,
          title: showLogo
              ? const AppLogo(height: 24)
              : Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  mainAxisAlignment: MainAxisAlignment.center,
                  children: [
                    Text(
                      title ?? '',
                      style: theme.textTheme.titleLarge,
                      maxLines: 1,
                      overflow: TextOverflow.ellipsis,
                    ),
                    if (subtitle != null)
                      Text(
                        subtitle!,
                        style: theme.textTheme.labelMedium,
                        maxLines: 1,
                        overflow: TextOverflow.ellipsis,
                      ),
                  ],
                ),
          actions: [...actions, const SizedBox(width: 4)],
          bottom: bottom,
        ),
        floatingActionButton: floatingActionButton,
        body: SafeArea(top: false, bottom: padBottomForNav, child: body),
      );
    }

    return Scaffold(
      backgroundColor: backgroundColor ?? AppColors.background,
      floatingActionButton: floatingActionButton,
      body: SafeArea(
        bottom: padBottomForNav,
        child: Column(
          children: [
            if (actions.isNotEmpty)
              Padding(
                padding: const EdgeInsets.fromLTRB(8, 4, 8, 0),
                child: Row(
                  mainAxisAlignment: MainAxisAlignment.end,
                  children: actions,
                ),
              ),
            if (bottom != null)
              Container(
                decoration: const BoxDecoration(
                  color: AppColors.surface,
                  border: Border(bottom: AppBorders.hairline),
                ),
                child: bottom,
              ),
            Expanded(child: body),
          ],
        ),
      ),
    );
  }
}
