import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../../services/auth/business/ui/providers/business_provider.dart';
import '../../services/auth/login/ui/providers/login_provider.dart';

/// Widget que muestra un selector de negocio para super admins.
/// Si el usuario es super admin y no ha seleccionado negocio, muestra el selector.
/// Si el usuario es normal (tiene business_id en JWT), renderiza el child directamente.
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

    // Usuario normal → no necesita selector
    if (!login.isSuperAdmin) {
      return widget.builder(context, 0);
    }

    // Super admin → necesita seleccionar negocio
    return Consumer<BusinessProvider>(
      builder: (context, bizProvider, _) {
        final selectedId = bizProvider.selectedBusinessId;
        final businesses = bizProvider.businessesSimple;
        final hasError = bizProvider.error != null;

        return Column(
          children: [
            // Selector de negocio (barra superior)
            Container(
              width: double.infinity,
              padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 8),
              decoration: BoxDecoration(
                color: Theme.of(context).colorScheme.primaryContainer.withAlpha(60),
                border: Border(
                  bottom: BorderSide(
                    color: Theme.of(context).colorScheme.outlineVariant,
                  ),
                ),
              ),
              child: Row(
                children: [
                  Icon(Icons.business,
                      size: 20, color: Theme.of(context).colorScheme.primary),
                  const SizedBox(width: 8),
                  Expanded(
                    child: _buildSelector(bizProvider, businesses, selectedId, hasError),
                  ),
                ],
              ),
            ),

            // Contenido
            Expanded(
              child: selectedId == null
                  ? _buildPlaceholder(hasError)
                  : widget.builder(context, selectedId),
            ),
          ],
        );
      },
    );
  }

  Widget _buildSelector(
    BusinessProvider provider,
    List businesses,
    int? selectedId,
    bool hasError,
  ) {
    if (provider.isLoading) {
      return const Row(
        children: [
          SizedBox(
            width: 16, height: 16,
            child: CircularProgressIndicator(strokeWidth: 2),
          ),
          SizedBox(width: 8),
          Text('Cargando negocios...', style: TextStyle(fontSize: 13)),
        ],
      );
    }

    if (hasError || businesses.isEmpty) {
      return InkWell(
        onTap: () => provider.fetchBusinessesSimple(),
        child: Row(
          children: [
            Icon(Icons.refresh, size: 16, color: Colors.red.shade400),
            const SizedBox(width: 4),
            Expanded(
              child: Text(
                hasError
                    ? 'Error al cargar negocios. Toca para reintentar.'
                    : 'Sin negocios. Toca para reintentar.',
                style: TextStyle(fontSize: 12, color: Colors.red.shade400),
              ),
            ),
          ],
        ),
      );
    }

    return DropdownButtonHideUnderline(
      child: DropdownButton<int>(
        value: selectedId,
        isExpanded: true,
        isDense: true,
        hint: const Text('Selecciona un negocio',
            style: TextStyle(fontSize: 13)),
        items: businesses.map<DropdownMenuItem<int>>((b) {
          return DropdownMenuItem<int>(
            value: b.id,
            child: Text('${b.name} (ID: ${b.id})',
                style: const TextStyle(fontSize: 13)),
          );
        }).toList(),
        onChanged: (id) {
          if (id != null) {
            provider.setSelectedBusinessId(id);
          }
        },
      ),
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
          if (hasError) ...[
            const SizedBox(height: 12),
            FilledButton.icon(
              onPressed: _loadBusinesses,
              icon: const Icon(Icons.refresh),
              label: const Text('Reintentar'),
            ),
          ],
        ],
      ),
    );
  }
}
