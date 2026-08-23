import 'package:flutter/material.dart';
import '../../../services/modules/notification_config/ui/screens/notification_config_screen.dart';
import 'module_tabs_scaffold.dart';

class NotificationsModuleScreen extends StatelessWidget {
  const NotificationsModuleScreen({super.key, this.initialTab = 0});

  final int initialTab;

  @override
  Widget build(BuildContext context) {
    return ModuleTabsScaffold(
      title: 'Notificaciones',
      subtitle: 'Eventos, canales y plantillas',
      initialTab: initialTab,
      tabs: const ['Configuraciones'],
      builder: (context, businessId) => [
        NotificationConfigScreen(key: ValueKey('notifications_$businessId'), businessId: businessId),
      ],
    );
  }
}
