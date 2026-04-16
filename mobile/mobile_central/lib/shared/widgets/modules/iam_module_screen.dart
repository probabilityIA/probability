import 'package:flutter/material.dart';

import '../../../services/auth/actions/ui/screens/action_list_screen.dart';
import '../../../services/auth/permissions/ui/screens/permission_list_screen.dart';
import '../../../services/auth/resources/ui/screens/resource_list_screen.dart';
import '../../../services/auth/roles/ui/screens/role_list_screen.dart';
import '../../../services/auth/users/ui/screens/user_list_screen.dart';

/// Module wrapper that groups Users, Roles, Permissions, Resources and Actions
/// behind a TabBar, replicating the Next.js "subnavbar" pattern.
/// This is an auth/IAM module so it does NOT use BusinessSelectorWrapper.
class IamModuleScreen extends StatelessWidget {
  final int initialTab;

  const IamModuleScreen({super.key, this.initialTab = 0});

  @override
  Widget build(BuildContext context) {
    return DefaultTabController(
      length: 5,
      initialIndex: initialTab,
      child: Scaffold(
        appBar: AppBar(
          title: const Text('Administracion'),
          bottom: const TabBar(
            isScrollable: true,
            tabAlignment: TabAlignment.start,
            tabs: [
              Tab(icon: Icon(Icons.people), text: 'Usuarios'),
              Tab(icon: Icon(Icons.admin_panel_settings), text: 'Roles'),
              Tab(icon: Icon(Icons.lock), text: 'Permisos'),
              Tab(icon: Icon(Icons.category), text: 'Recursos'),
              Tab(icon: Icon(Icons.touch_app), text: 'Acciones'),
            ],
          ),
        ),
        body: const TabBarView(
          children: [
            UserListScreen(),
            RoleListScreen(),
            PermissionListScreen(),
            ResourceListScreen(),
            ActionListScreen(),
          ],
        ),
      ),
    );
  }
}
