import 'package:flutter/material.dart';
import '../../../services/auth/actions/ui/screens/action_list_screen.dart';
import '../../../services/auth/permissions/ui/screens/permission_list_screen.dart';
import '../../../services/auth/resources/ui/screens/resource_list_screen.dart';
import '../../../services/auth/roles/ui/screens/role_list_screen.dart';
import '../../../services/auth/users/ui/screens/user_list_screen.dart';
import 'module_tabs_scaffold.dart';

class IamModuleScreen extends StatelessWidget {
  const IamModuleScreen({super.key, this.initialTab = 0});

  final int initialTab;

  @override
  Widget build(BuildContext context) {
    return ModuleTabsScaffold(
      title: 'Administracion',
      subtitle: 'Usuarios, roles y permisos',
      initialTab: initialTab,
      tabs: const ['Usuarios', 'Roles', 'Permisos', 'Recursos', 'Acciones'],
      builder: (context, businessId) => const [
        UserListScreen(),
        RoleListScreen(),
        PermissionListScreen(),
        ResourceListScreen(),
        ActionListScreen(),
      ],
    );
  }
}
