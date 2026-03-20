import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../../domain/entities.dart';
import '../providers/notification_config_provider.dart';

class NotificationConfigScreen extends StatefulWidget {
  final int? businessId;

  const NotificationConfigScreen({super.key, this.businessId});

  @override
  State<NotificationConfigScreen> createState() =>
      _NotificationConfigScreenState();
}

class _NotificationConfigScreenState extends State<NotificationConfigScreen> {
  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addPostFrameCallback((_) {
      _loadConfigs();
    });
  }

  void _loadConfigs() {
    context.read<NotificationConfigProvider>().fetchConfigs();
  }

  IconData _channelIcon(String? channel) {
    switch (channel?.toLowerCase()) {
      case 'whatsapp':
        return Icons.chat;
      case 'email':
        return Icons.email_outlined;
      case 'sms':
        return Icons.sms_outlined;
      case 'push':
        return Icons.notifications_outlined;
      default:
        return Icons.notifications_outlined;
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('Configuracion de Notificaciones')),
      body: Consumer<NotificationConfigProvider>(
        builder: (context, provider, _) {
          if (provider.isLoading) {
            return const Center(child: CircularProgressIndicator());
          }
          if (provider.error != null) {
            return Center(
              child: Padding(
                padding: const EdgeInsets.all(24),
                child: Column(
                  mainAxisAlignment: MainAxisAlignment.center,
                  children: [
                    const Icon(Icons.error_outline,
                        size: 48, color: Colors.red),
                    const SizedBox(height: 16),
                    Text(provider.error!,
                        textAlign: TextAlign.center,
                        style: const TextStyle(color: Colors.red)),
                    const SizedBox(height: 16),
                    FilledButton.icon(
                      onPressed: _loadConfigs,
                      icon: const Icon(Icons.refresh),
                      label: const Text('Reintentar'),
                    ),
                  ],
                ),
              ),
            );
          }
          if (provider.configs.isEmpty) {
            return Center(
              child: Column(
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  Icon(Icons.notifications_off_outlined,
                      size: 64, color: Colors.grey.shade400),
                  const SizedBox(height: 16),
                  Text('No hay configuraciones',
                      style:
                          TextStyle(fontSize: 16, color: Colors.grey.shade600)),
                ],
              ),
            );
          }
          return RefreshIndicator(
            onRefresh: () async => _loadConfigs(),
            child: ListView.builder(
              padding: const EdgeInsets.all(16),
              itemCount: provider.configs.length,
              itemBuilder: (context, index) {
                final config = provider.configs[index];
                return _NotificationConfigCard(
                  config: config,
                  channelIcon: _channelIcon(
                      config.notificationType?.code ?? config.eventType),
                );
              },
            ),
          );
        },
      ),
    );
  }
}

class _NotificationConfigCard extends StatelessWidget {
  final NotificationConfig config;
  final IconData channelIcon;

  const _NotificationConfigCard({
    required this.config,
    required this.channelIcon,
  });

  @override
  Widget build(BuildContext context) {
    final typeName = config.notificationTypeName ??
        config.notificationType?.name ??
        'Tipo ${config.notificationTypeId}';
    final eventName = config.notificationEventName ??
        config.notificationEventType?.eventName ??
        'Evento ${config.notificationEventTypeId}';

    return Card(
      margin: const EdgeInsets.only(bottom: 10),
      child: Padding(
        padding: const EdgeInsets.all(12),
        child: Row(
          children: [
            CircleAvatar(
              backgroundColor: config.enabled
                  ? Colors.blue.shade100
                  : Colors.grey.shade200,
              child: Icon(channelIcon,
                  color: config.enabled
                      ? Colors.blue.shade700
                      : Colors.grey.shade500),
            ),
            const SizedBox(width: 12),
            Expanded(
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(typeName,
                      style: const TextStyle(
                          fontWeight: FontWeight.w600, fontSize: 14)),
                  const SizedBox(height: 2),
                  Text(eventName,
                      style: TextStyle(
                          fontSize: 13, color: Colors.grey.shade600)),
                  if (config.description != null &&
                      config.description!.isNotEmpty) ...[
                    const SizedBox(height: 2),
                    Text(config.description!,
                        style: TextStyle(
                            fontSize: 12, color: Colors.grey.shade500),
                        maxLines: 1,
                        overflow: TextOverflow.ellipsis),
                  ],
                  if (config.channels != null &&
                      config.channels!.isNotEmpty) ...[
                    const SizedBox(height: 4),
                    Wrap(
                      spacing: 4,
                      children: config.channels!
                          .map((ch) => Container(
                                padding: const EdgeInsets.symmetric(
                                    horizontal: 6, vertical: 1),
                                decoration: BoxDecoration(
                                  color: Colors.indigo.withValues(alpha: 0.1),
                                  borderRadius: BorderRadius.circular(8),
                                ),
                                child: Text(ch,
                                    style: const TextStyle(
                                        fontSize: 10, color: Colors.indigo)),
                              ))
                          .toList(),
                    ),
                  ],
                ],
              ),
            ),
            Switch(
              value: config.enabled,
              onChanged: null, // Read-only display
            ),
          ],
        ),
      ),
    );
  }
}
