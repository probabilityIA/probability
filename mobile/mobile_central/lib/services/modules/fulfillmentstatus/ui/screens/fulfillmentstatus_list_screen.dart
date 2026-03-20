import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../providers/fulfillmentstatus_provider.dart';

class FulfillmentStatusListScreen extends StatefulWidget {
  const FulfillmentStatusListScreen({super.key});

  @override
  State<FulfillmentStatusListScreen> createState() =>
      _FulfillmentStatusListScreenState();
}

class _FulfillmentStatusListScreenState
    extends State<FulfillmentStatusListScreen> {
  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addPostFrameCallback((_) {
      _loadStatuses();
    });
  }

  void _loadStatuses() {
    context.read<FulfillmentStatusProvider>().fetchStatuses();
  }

  Color _parseColor(String? colorStr) {
    if (colorStr == null || colorStr.isEmpty) return Colors.grey;
    try {
      final hex = colorStr.replaceFirst('#', '');
      if (hex.length == 6) {
        return Color(int.parse('FF$hex', radix: 16));
      }
    } catch (_) {}
    return Colors.grey;
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('Estados de Fulfillment')),
      body: Consumer<FulfillmentStatusProvider>(
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
                      onPressed: _loadStatuses,
                      icon: const Icon(Icons.refresh),
                      label: const Text('Reintentar'),
                    ),
                  ],
                ),
              ),
            );
          }
          if (provider.statuses.isEmpty) {
            return Center(
              child: Column(
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  Icon(Icons.local_shipping_outlined,
                      size: 64, color: Colors.grey.shade400),
                  const SizedBox(height: 16),
                  Text('No hay estados de fulfillment',
                      style:
                          TextStyle(fontSize: 16, color: Colors.grey.shade600)),
                ],
              ),
            );
          }
          return RefreshIndicator(
            onRefresh: () async => _loadStatuses(),
            child: ListView.builder(
              padding: const EdgeInsets.all(16),
              itemCount: provider.statuses.length,
              itemBuilder: (context, index) {
                final status = provider.statuses[index];
                final color = _parseColor(status.color);
                return Card(
                  margin: const EdgeInsets.only(bottom: 8),
                  child: ListTile(
                    leading: CircleAvatar(
                      backgroundColor: color.withValues(alpha: 0.15),
                      child: Icon(Icons.circle, color: color, size: 16),
                    ),
                    title: Text(status.name,
                        style: const TextStyle(fontWeight: FontWeight.w600)),
                    subtitle: Row(
                      children: [
                        Text(status.code),
                        if (status.category != null) ...[
                          const SizedBox(width: 8),
                          Container(
                            padding: const EdgeInsets.symmetric(
                                horizontal: 8, vertical: 2),
                            decoration: BoxDecoration(
                              color: Colors.indigo.withValues(alpha: 0.1),
                              borderRadius: BorderRadius.circular(12),
                            ),
                            child: Text(status.category!,
                                style: const TextStyle(
                                    fontSize: 11,
                                    color: Colors.indigo,
                                    fontWeight: FontWeight.w500)),
                          ),
                        ],
                      ],
                    ),
                    trailing: status.description != null
                        ? Tooltip(
                            message: status.description!,
                            child: const Icon(Icons.info_outline, size: 20),
                          )
                        : null,
                  ),
                );
              },
            ),
          );
        },
      ),
    );
  }
}
