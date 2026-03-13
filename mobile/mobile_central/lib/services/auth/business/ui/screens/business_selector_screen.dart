import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../../domain/entities.dart';
import '../providers/business_provider.dart';

class BusinessSelectorScreen extends StatefulWidget {
  final void Function(BusinessSimple business) onBusinessSelected;

  const BusinessSelectorScreen({
    super.key,
    required this.onBusinessSelected,
  });

  @override
  State<BusinessSelectorScreen> createState() =>
      _BusinessSelectorScreenState();
}

class _BusinessSelectorScreenState extends State<BusinessSelectorScreen> {
  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addPostFrameCallback((_) {
      context.read<BusinessProvider>().fetchBusinessesSimple();
    });
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Seleccionar negocio'),
      ),
      body: Consumer<BusinessProvider>(
        builder: (context, provider, child) {
          if (provider.isLoading) {
            return const Center(child: CircularProgressIndicator());
          }

          if (provider.businessesSimple.isEmpty) {
            return const Center(
              child: Text('No hay negocios disponibles'),
            );
          }

          return ListView.builder(
            padding: const EdgeInsets.all(16),
            itemCount: provider.businessesSimple.length,
            itemBuilder: (context, index) {
              final business = provider.businessesSimple[index];
              return Card(
                child: ListTile(
                  leading: business.logoUrl != null
                      ? CircleAvatar(
                          backgroundImage: NetworkImage(business.logoUrl!),
                        )
                      : CircleAvatar(
                          child: Text(business.name[0].toUpperCase()),
                        ),
                  title: Text(business.name),
                  subtitle: Text('ID: ${business.id}'),
                  trailing: const Icon(Icons.arrow_forward_ios, size: 16),
                  onTap: () => widget.onBusinessSelected(business),
                ),
              );
            },
          );
        },
      ),
    );
  }
}
