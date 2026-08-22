import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../../services/auth/business/ui/providers/business_provider.dart';
import '../../services/auth/login/ui/providers/login_provider.dart';

class BusinessSelectorWrapper extends StatefulWidget {
  final Widget Function(BuildContext context, int businessId) builder;

  const BusinessSelectorWrapper({super.key, required this.builder});

  @override
  State<BusinessSelectorWrapper> createState() =>
      _BusinessSelectorWrapperState();
}

class _BusinessSelectorWrapperState extends State<BusinessSelectorWrapper> {
  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addPostFrameCallback((_) {
      _loadBusinesses();
    });
  }

  void _loadBusinesses() {
    final login = context.read<LoginProvider>();
    if (login.isSuperAdmin) {
      final biz = context.read<BusinessProvider>();
      if (biz.businessesSimple.isEmpty) {
        biz.fetchBusinessesSimple();
      }
    }
  }

  @override
  Widget build(BuildContext context) {
    final login = context.watch<LoginProvider>();

    if (!login.isSuperAdmin) {
      return widget.builder(context, 0);
    }

    return Consumer<BusinessProvider>(
      builder: (context, bizProvider, _) {
        final selectedId = bizProvider.selectedBusinessId;
        if (selectedId == null) {
          return _buildPlaceholder(bizProvider.error != null);
        }
        return widget.builder(context, selectedId);
      },
    );
  }

  Widget _buildPlaceholder(bool hasError) {
    return Center(
      child: Column(
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          Icon(Icons.store, size: 64, color: Colors.grey.shade400),
          const SizedBox(height: 16),
          Text(
            hasError
                ? 'Error al cargar negocios'
                : 'Selecciona un negocio para continuar',
            style: TextStyle(fontSize: 16, color: Colors.grey.shade600),
          ),
          const SizedBox(height: 8),
          if (hasError)
            FilledButton.icon(
              onPressed: _loadBusinesses,
              icon: const Icon(Icons.refresh),
              label: const Text('Reintentar'),
            )
          else
            Text(
              'Toca la pestania morada del borde para elegirlo.',
              style: TextStyle(fontSize: 13, color: Colors.grey.shade500),
            ),
        ],
      ),
    );
  }
}
