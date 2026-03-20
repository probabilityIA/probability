import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../../../../../shared/widgets/network_avatar.dart';
import '../providers/user_provider.dart';

class UserListScreen extends StatefulWidget {
  final int? businessId;

  const UserListScreen({super.key, this.businessId});

  @override
  State<UserListScreen> createState() => _UserListScreenState();
}

class _UserListScreenState extends State<UserListScreen> {
  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addPostFrameCallback((_) {
      context.read<UserProvider>().fetchUsers(businessId: widget.businessId);
    });
  }

  @override
  void didUpdateWidget(UserListScreen oldWidget) {
    super.didUpdateWidget(oldWidget);
    if (oldWidget.businessId != widget.businessId) {
      context.read<UserProvider>().resetFilters();
      context.read<UserProvider>().fetchUsers(businessId: widget.businessId);
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Usuarios'),
      ),
      body: Consumer<UserProvider>(
        builder: (context, provider, child) {
          if (provider.isLoading) {
            return const Center(child: CircularProgressIndicator());
          }

          if (provider.error != null) {
            return Center(
              child: Column(
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  Text('Error: ${provider.error}'),
                  const SizedBox(height: 16),
                  FilledButton(
                    onPressed: () =>
                        provider.fetchUsers(businessId: widget.businessId),
                    child: const Text('Reintentar'),
                  ),
                ],
              ),
            );
          }

          if (provider.users.isEmpty) {
            return const Center(child: Text('No hay usuarios'));
          }

          return RefreshIndicator(
            onRefresh: () =>
                provider.fetchUsers(businessId: widget.businessId),
            child: ListView.builder(
              padding: const EdgeInsets.all(16),
              itemCount: provider.users.length,
              itemBuilder: (context, index) {
                final user = provider.users[index];
                return Card(
                  child: ListTile(
                    leading: NetworkAvatar(
                      imageUrl: user.avatarUrl,
                      fallbackText: user.name,
                      radius: 20,
                    ),
                    title: Text(user.name),
                    subtitle: Text(user.email),
                    trailing: Row(
                      mainAxisSize: MainAxisSize.min,
                      children: [
                        if (user.isSuperUser)
                          const Chip(
                            label: Text('Super',
                                style: TextStyle(fontSize: 10)),
                            padding: EdgeInsets.zero,
                          ),
                        const SizedBox(width: 4),
                        Icon(
                          Icons.circle,
                          size: 12,
                          color:
                              user.isActive ? Colors.green : Colors.grey,
                        ),
                      ],
                    ),
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
