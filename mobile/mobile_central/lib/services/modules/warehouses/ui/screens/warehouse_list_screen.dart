import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../providers/warehouse_provider.dart';
import '../../domain/entities.dart';

class WarehouseListScreen extends StatefulWidget {
  final int? businessId;

  const WarehouseListScreen({super.key, this.businessId});

  @override
  State<WarehouseListScreen> createState() => _WarehouseListScreenState();
}

class _WarehouseListScreenState extends State<WarehouseListScreen> {
  final _searchController = TextEditingController();

  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addPostFrameCallback((_) {
      context.read<WarehouseProvider>().fetchWarehouses(businessId: widget.businessId);
    });
  }

  @override
  void didUpdateWidget(WarehouseListScreen oldWidget) {
    super.didUpdateWidget(oldWidget);
    if (oldWidget.businessId != widget.businessId) {
      _searchController.clear();
      final provider = context.read<WarehouseProvider>();
      provider.setPage(1);
      provider.fetchWarehouses(businessId: widget.businessId);
    }
  }

  @override
  void dispose() {
    _searchController.dispose();
    super.dispose();
  }

  void _onSearch(WarehouseProvider provider) {
    provider.setPage(1);
    provider.fetchWarehouses(businessId: widget.businessId);
  }

  void _onClearSearch(WarehouseProvider provider) {
    _searchController.clear();
    provider.setPage(1);
    provider.fetchWarehouses(businessId: widget.businessId);
  }

  void _showWarehouseForm({Warehouse? warehouse}) {
    final isEditing = warehouse != null;
    final nameCtrl = TextEditingController(text: warehouse?.name ?? '');
    final codeCtrl = TextEditingController(text: warehouse?.code ?? '');
    final addressCtrl = TextEditingController(text: warehouse?.address ?? '');
    final cityCtrl = TextEditingController(text: warehouse?.city ?? '');
    final stateCtrl = TextEditingController(text: warehouse?.state ?? '');
    final countryCtrl = TextEditingController(text: warehouse?.country ?? '');
    final zipCodeCtrl = TextEditingController(text: warehouse?.zipCode ?? '');
    final phoneCtrl = TextEditingController(text: warehouse?.phone ?? '');
    final contactNameCtrl =
        TextEditingController(text: warehouse?.contactName ?? '');
    final contactEmailCtrl =
        TextEditingController(text: warehouse?.contactEmail ?? '');
    bool isDefault = warehouse?.isDefault ?? false;
    bool isFulfillment = warehouse?.isFulfillment ?? false;
    bool isActive = warehouse?.isActive ?? true;
    final formKey = GlobalKey<FormState>();
    bool isSaving = false;

    showModalBottomSheet(
      context: context,
      isScrollControlled: true,
      useSafeArea: true,
      shape: const RoundedRectangleBorder(
        borderRadius: BorderRadius.vertical(top: Radius.circular(16)),
      ),
      builder: (ctx) {
        return StatefulBuilder(
          builder: (ctx, setModalState) {
            return Padding(
              padding: EdgeInsets.only(
                left: 16,
                right: 16,
                top: 16,
                bottom: MediaQuery.of(ctx).viewInsets.bottom + 16,
              ),
              child: Form(
                key: formKey,
                child: SingleChildScrollView(
                  child: Column(
                    mainAxisSize: MainAxisSize.min,
                    crossAxisAlignment: CrossAxisAlignment.stretch,
                    children: [
                      Row(
                        children: [
                          Expanded(
                            child: Text(
                              isEditing ? 'Editar bodega' : 'Nueva bodega',
                              style: Theme.of(ctx).textTheme.titleLarge,
                            ),
                          ),
                          IconButton(
                            icon: const Icon(Icons.close),
                            onPressed: () => Navigator.pop(ctx),
                          ),
                        ],
                      ),
                      const SizedBox(height: 16),

                      // Basic info
                      Text('Informacion basica',
                          style: Theme.of(ctx).textTheme.titleSmall),
                      const SizedBox(height: 8),
                      TextFormField(
                        controller: nameCtrl,
                        decoration: const InputDecoration(
                          labelText: 'Nombre *',
                          border: OutlineInputBorder(),
                        ),
                        validator: (v) => (v == null || v.trim().length < 2)
                            ? 'Minimo 2 caracteres'
                            : null,
                      ),
                      const SizedBox(height: 12),
                      TextFormField(
                        controller: codeCtrl,
                        decoration: const InputDecoration(
                          labelText: 'Codigo *',
                          border: OutlineInputBorder(),
                          hintText: 'BOD-001',
                        ),
                        textCapitalization: TextCapitalization.characters,
                        validator: (v) => (v == null || v.trim().isEmpty)
                            ? 'Requerido'
                            : null,
                      ),
                      const SizedBox(height: 16),

                      // Address
                      Text('Direccion',
                          style: Theme.of(ctx).textTheme.titleSmall),
                      const SizedBox(height: 8),
                      TextFormField(
                        controller: addressCtrl,
                        decoration: const InputDecoration(
                          labelText: 'Direccion',
                          border: OutlineInputBorder(),
                        ),
                      ),
                      const SizedBox(height: 12),
                      Row(
                        children: [
                          Expanded(
                            child: TextFormField(
                              controller: cityCtrl,
                              decoration: const InputDecoration(
                                labelText: 'Ciudad',
                                border: OutlineInputBorder(),
                              ),
                            ),
                          ),
                          const SizedBox(width: 12),
                          Expanded(
                            child: TextFormField(
                              controller: stateCtrl,
                              decoration: const InputDecoration(
                                labelText: 'Departamento',
                                border: OutlineInputBorder(),
                              ),
                            ),
                          ),
                        ],
                      ),
                      const SizedBox(height: 12),
                      Row(
                        children: [
                          Expanded(
                            child: TextFormField(
                              controller: countryCtrl,
                              decoration: const InputDecoration(
                                labelText: 'Pais',
                                border: OutlineInputBorder(),
                              ),
                            ),
                          ),
                          const SizedBox(width: 12),
                          Expanded(
                            child: TextFormField(
                              controller: zipCodeCtrl,
                              decoration: const InputDecoration(
                                labelText: 'Codigo postal',
                                border: OutlineInputBorder(),
                              ),
                            ),
                          ),
                        ],
                      ),
                      const SizedBox(height: 16),

                      // Contact
                      Text('Contacto',
                          style: Theme.of(ctx).textTheme.titleSmall),
                      const SizedBox(height: 8),
                      TextFormField(
                        controller: phoneCtrl,
                        decoration: const InputDecoration(
                          labelText: 'Telefono',
                          border: OutlineInputBorder(),
                        ),
                        keyboardType: TextInputType.phone,
                      ),
                      const SizedBox(height: 12),
                      TextFormField(
                        controller: contactNameCtrl,
                        decoration: const InputDecoration(
                          labelText: 'Nombre contacto',
                          border: OutlineInputBorder(),
                        ),
                      ),
                      const SizedBox(height: 12),
                      TextFormField(
                        controller: contactEmailCtrl,
                        decoration: const InputDecoration(
                          labelText: 'Email contacto',
                          border: OutlineInputBorder(),
                        ),
                        keyboardType: TextInputType.emailAddress,
                      ),
                      const SizedBox(height: 16),

                      // Toggles
                      CheckboxListTile(
                        title: const Text('Bodega principal'),
                        subtitle: const Text(
                            'Por defecto para nuevas ordenes',
                            style: TextStyle(fontSize: 12)),
                        value: isDefault,
                        onChanged: (v) =>
                            setModalState(() => isDefault = v ?? false),
                        controlAffinity: ListTileControlAffinity.leading,
                        contentPadding: EdgeInsets.zero,
                      ),
                      CheckboxListTile(
                        title: const Text('Fulfillment'),
                        subtitle: const Text('Despacho de pedidos',
                            style: TextStyle(fontSize: 12)),
                        value: isFulfillment,
                        onChanged: (v) =>
                            setModalState(() => isFulfillment = v ?? false),
                        controlAffinity: ListTileControlAffinity.leading,
                        contentPadding: EdgeInsets.zero,
                      ),
                      if (isEditing)
                        CheckboxListTile(
                          title: const Text('Activa'),
                          value: isActive,
                          onChanged: (v) =>
                              setModalState(() => isActive = v ?? true),
                          controlAffinity: ListTileControlAffinity.leading,
                          contentPadding: EdgeInsets.zero,
                        ),

                      const SizedBox(height: 20),
                      FilledButton(
                        onPressed: isSaving
                            ? null
                            : () async {
                                if (!formKey.currentState!.validate()) return;
                                setModalState(() => isSaving = true);
                                final provider =
                                    context.read<WarehouseProvider>();
                                bool success;
                                if (isEditing) {
                                  final dto = UpdateWarehouseDTO(
                                    name: nameCtrl.text.trim(),
                                    code: codeCtrl.text.trim(),
                                    address:
                                        addressCtrl.text.trim().isNotEmpty
                                            ? addressCtrl.text.trim()
                                            : null,
                                    city: cityCtrl.text.trim().isNotEmpty
                                        ? cityCtrl.text.trim()
                                        : null,
                                    state: stateCtrl.text.trim().isNotEmpty
                                        ? stateCtrl.text.trim()
                                        : null,
                                    country:
                                        countryCtrl.text.trim().isNotEmpty
                                            ? countryCtrl.text.trim()
                                            : null,
                                    zipCode:
                                        zipCodeCtrl.text.trim().isNotEmpty
                                            ? zipCodeCtrl.text.trim()
                                            : null,
                                    phone: phoneCtrl.text.trim().isNotEmpty
                                        ? phoneCtrl.text.trim()
                                        : null,
                                    contactName: contactNameCtrl.text
                                            .trim()
                                            .isNotEmpty
                                        ? contactNameCtrl.text.trim()
                                        : null,
                                    contactEmail: contactEmailCtrl.text
                                            .trim()
                                            .isNotEmpty
                                        ? contactEmailCtrl.text.trim()
                                        : null,
                                    isDefault: isDefault,
                                    isFulfillment: isFulfillment,
                                    isActive: isActive,
                                  );
                                  final result =
                                      await provider.updateWarehouse(
                                    warehouse.id,
                                    dto,
                                    businessId: widget.businessId,
                                  );
                                  success = result != null;
                                } else {
                                  final dto = CreateWarehouseDTO(
                                    name: nameCtrl.text.trim(),
                                    code: codeCtrl.text.trim(),
                                    address:
                                        addressCtrl.text.trim().isNotEmpty
                                            ? addressCtrl.text.trim()
                                            : null,
                                    city: cityCtrl.text.trim().isNotEmpty
                                        ? cityCtrl.text.trim()
                                        : null,
                                    state: stateCtrl.text.trim().isNotEmpty
                                        ? stateCtrl.text.trim()
                                        : null,
                                    country:
                                        countryCtrl.text.trim().isNotEmpty
                                            ? countryCtrl.text.trim()
                                            : null,
                                    zipCode:
                                        zipCodeCtrl.text.trim().isNotEmpty
                                            ? zipCodeCtrl.text.trim()
                                            : null,
                                    phone: phoneCtrl.text.trim().isNotEmpty
                                        ? phoneCtrl.text.trim()
                                        : null,
                                    contactName: contactNameCtrl.text
                                            .trim()
                                            .isNotEmpty
                                        ? contactNameCtrl.text.trim()
                                        : null,
                                    contactEmail: contactEmailCtrl.text
                                            .trim()
                                            .isNotEmpty
                                        ? contactEmailCtrl.text.trim()
                                        : null,
                                    isDefault: isDefault,
                                    isFulfillment: isFulfillment,
                                  );
                                  final result =
                                      await provider.createWarehouse(
                                    dto,
                                    businessId: widget.businessId,
                                  );
                                  success = result != null;
                                }
                                setModalState(() => isSaving = false);
                                if (success && ctx.mounted) {
                                  Navigator.pop(ctx);
                                  provider.fetchWarehouses(
                                      businessId: widget.businessId);
                                  if (mounted) {
                                    ScaffoldMessenger.of(context).showSnackBar(
                                      SnackBar(
                                        content: Text(isEditing
                                            ? 'Bodega actualizada'
                                            : 'Bodega creada'),
                                      ),
                                    );
                                  }
                                } else if (mounted) {
                                  ScaffoldMessenger.of(context).showSnackBar(
                                    SnackBar(
                                      content: Text(provider.error ??
                                          'Error al guardar'),
                                      backgroundColor:
                                          Theme.of(context).colorScheme.error,
                                    ),
                                  );
                                }
                              },
                        child: isSaving
                            ? const SizedBox(
                                height: 20,
                                width: 20,
                                child: CircularProgressIndicator(
                                    strokeWidth: 2, color: Colors.white),
                              )
                            : Text(isEditing ? 'Actualizar' : 'Crear'),
                      ),
                    ],
                  ),
                ),
              ),
            );
          },
        );
      },
    );
  }

  void _confirmDelete(Warehouse warehouse) {
    showDialog(
      context: context,
      builder: (ctx) => AlertDialog(
        title: const Text('Eliminar bodega'),
        content: Text(
            'Eliminar la bodega "${warehouse.name}"? Esta accion no se puede deshacer.'),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(ctx),
            child: const Text('Cancelar'),
          ),
          FilledButton(
            style: FilledButton.styleFrom(
              backgroundColor: Theme.of(ctx).colorScheme.error,
            ),
            onPressed: () async {
              Navigator.pop(ctx);
              final provider = context.read<WarehouseProvider>();
              final ok = await provider.deleteWarehouse(warehouse.id,
                  businessId: widget.businessId);
              if (ok) {
                provider.fetchWarehouses(businessId: widget.businessId);
                if (mounted) {
                  ScaffoldMessenger.of(context).showSnackBar(
                    const SnackBar(content: Text('Bodega eliminada')),
                  );
                }
              } else if (mounted) {
                ScaffoldMessenger.of(context).showSnackBar(
                  SnackBar(
                    content: Text(provider.error ?? 'Error al eliminar'),
                    backgroundColor: Theme.of(context).colorScheme.error,
                  ),
                );
              }
            },
            child: const Text('Eliminar'),
          ),
        ],
      ),
    );
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    return Scaffold(
      appBar: AppBar(title: const Text('Bodegas')),
      floatingActionButton: FloatingActionButton(
        onPressed: () => _showWarehouseForm(),
        child: const Icon(Icons.add),
      ),
      body: Consumer<WarehouseProvider>(
        builder: (context, provider, _) {
          return Column(
            children: [
              // Search bar
              Padding(
                padding: const EdgeInsets.fromLTRB(16, 12, 16, 8),
                child: Row(
                  children: [
                    Expanded(
                      child: TextField(
                        controller: _searchController,
                        decoration: InputDecoration(
                          hintText: 'Buscar por nombre o codigo...',
                          prefixIcon: const Icon(Icons.search),
                          suffixIcon: _searchController.text.isNotEmpty
                              ? IconButton(
                                  icon: const Icon(Icons.clear),
                                  onPressed: () => _onClearSearch(provider),
                                )
                              : null,
                          border: const OutlineInputBorder(),
                          contentPadding:
                              const EdgeInsets.symmetric(horizontal: 12),
                        ),
                        onSubmitted: (_) => _onSearch(provider),
                      ),
                    ),
                    const SizedBox(width: 8),
                    IconButton.filled(
                      onPressed: () => _onSearch(provider),
                      icon: const Icon(Icons.search),
                    ),
                  ],
                ),
              ),

              // Content
              Expanded(child: _buildContent(provider, theme)),

              // Pagination
              if (provider.pagination != null && !provider.isLoading)
                _buildPagination(provider),
            ],
          );
        },
      ),
    );
  }

  Widget _buildContent(WarehouseProvider provider, ThemeData theme) {
    if (provider.isLoading && provider.warehouses.isEmpty) {
      return const Center(child: CircularProgressIndicator());
    }

    if (provider.error != null && provider.warehouses.isEmpty) {
      return Center(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Icon(Icons.error_outline, size: 48, color: theme.colorScheme.error),
            const SizedBox(height: 16),
            Text(provider.error!, textAlign: TextAlign.center),
            const SizedBox(height: 16),
            FilledButton.icon(
              onPressed: () =>
                  provider.fetchWarehouses(businessId: widget.businessId),
              icon: const Icon(Icons.refresh),
              label: const Text('Reintentar'),
            ),
          ],
        ),
      );
    }

    if (provider.warehouses.isEmpty) {
      return Center(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Icon(Icons.warehouse_outlined, size: 48, color: theme.disabledColor),
            const SizedBox(height: 16),
            const Text('No hay bodegas registradas'),
          ],
        ),
      );
    }

    return RefreshIndicator(
      onRefresh: () =>
          provider.fetchWarehouses(businessId: widget.businessId),
      child: ListView.builder(
        padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 4),
        itemCount: provider.warehouses.length,
        itemBuilder: (context, index) {
          final wh = provider.warehouses[index];
          final location = [wh.city, wh.state]
              .where((s) => s.isNotEmpty)
              .join(', ');

          return Card(
            margin: const EdgeInsets.only(bottom: 8),
            child: ListTile(
              leading: CircleAvatar(
                backgroundColor: wh.isActive
                    ? theme.colorScheme.primaryContainer
                    : theme.disabledColor.withValues(alpha: 0.2),
                child: Icon(
                  Icons.warehouse,
                  color: wh.isActive
                      ? theme.colorScheme.onPrimaryContainer
                      : theme.disabledColor,
                ),
              ),
              title: Row(
                children: [
                  Expanded(
                    child: Text(wh.name,
                        style: const TextStyle(fontWeight: FontWeight.w600)),
                  ),
                  Text(wh.code,
                      style: theme.textTheme.bodySmall?.copyWith(
                          fontFamily: 'monospace',
                          color: theme.colorScheme.secondary)),
                ],
              ),
              subtitle: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  if (wh.address.isNotEmpty)
                    Text(wh.address,
                        style: theme.textTheme.bodySmall,
                        maxLines: 1,
                        overflow: TextOverflow.ellipsis),
                  if (location.isNotEmpty)
                    Text(location, style: theme.textTheme.bodySmall),
                  const SizedBox(height: 4),
                  Row(
                    children: [
                      if (wh.isDefault)
                        _buildBadge('Principal', Colors.blue),
                      if (wh.isFulfillment)
                        _buildBadge('Fulfillment', Colors.purple),
                      _buildBadge(
                        wh.isActive ? 'Activa' : 'Inactiva',
                        wh.isActive ? Colors.green : Colors.grey,
                      ),
                    ],
                  ),
                ],
              ),
              trailing: PopupMenuButton<String>(
                onSelected: (action) {
                  if (action == 'edit') {
                    _showWarehouseForm(warehouse: wh);
                  } else if (action == 'delete') {
                    _confirmDelete(wh);
                  }
                },
                itemBuilder: (_) => [
                  const PopupMenuItem(value: 'edit', child: Text('Editar')),
                  const PopupMenuItem(
                      value: 'delete', child: Text('Eliminar')),
                ],
              ),
              isThreeLine: true,
            ),
          );
        },
      ),
    );
  }

  Widget _buildBadge(String label, Color color) {
    return Container(
      margin: const EdgeInsets.only(right: 4),
      padding: const EdgeInsets.symmetric(horizontal: 6, vertical: 2),
      decoration: BoxDecoration(
        color: color.withValues(alpha: 0.15),
        borderRadius: BorderRadius.circular(8),
      ),
      child: Text(
        label,
        style: TextStyle(
          fontSize: 10,
          fontWeight: FontWeight.w600,
          color: color,
        ),
      ),
    );
  }

  Widget _buildPagination(WarehouseProvider provider) {
    final pagination = provider.pagination!;
    return Padding(
      padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
      child: Row(
        mainAxisAlignment: MainAxisAlignment.spaceBetween,
        children: [
          Text(
            'Pag. ${pagination.currentPage} de ${pagination.lastPage}  (${pagination.total} total)',
            style: Theme.of(context).textTheme.bodySmall,
          ),
          Row(
            children: [
              IconButton(
                icon: const Icon(Icons.chevron_left),
                onPressed: pagination.hasPrev
                    ? () {
                        provider.setPage(pagination.currentPage - 1);
                        provider.fetchWarehouses(
                            businessId: widget.businessId);
                      }
                    : null,
              ),
              IconButton(
                icon: const Icon(Icons.chevron_right),
                onPressed: pagination.hasNext
                    ? () {
                        provider.setPage(pagination.currentPage + 1);
                        provider.fetchWarehouses(
                            businessId: widget.businessId);
                      }
                    : null,
              ),
            ],
          ),
        ],
      ),
    );
  }
}
