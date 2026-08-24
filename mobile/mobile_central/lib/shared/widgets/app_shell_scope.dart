import 'package:flutter/material.dart';

class AppShellScope extends InheritedWidget {
  const AppShellScope({
    super.key,
    required this.openDrawer,
    required this.location,
    required super.child,
  });

  final VoidCallback openDrawer;
  final String location;

  static AppShellScope? maybeOf(BuildContext context) =>
      context.dependOnInheritedWidgetOfExactType<AppShellScope>();

  @override
  bool updateShouldNotify(AppShellScope oldWidget) =>
      oldWidget.location != location;
}
