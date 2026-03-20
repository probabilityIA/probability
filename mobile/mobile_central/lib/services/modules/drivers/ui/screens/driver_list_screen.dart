import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../providers/drivers_provider.dart';
import '../../domain/entities.dart';

class DriverListScreen extends StatefulWidget {
  final int? businessId;

  const DriverListScreen({super.key, this.businessId});

  @override
  State<DriverListScreen> createState() => _DriverListScreenState();
}

class _DriverListScreenState extends State<DriverListScreen> {
  final _searchController = TextEditingController();

  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addPostFrameCallback((_) {
      context.read<DriverProvider>().fetchDrivers(businessId: widget.businessId);
    });
  }

  @override
  void didUpdateWidget(DriverListScreen oldWidget) {
    super.didUpdateWidget(oldWidget);
    if (oldWidget.businessId != widget.businessId) {
      _searchController.clear();
      final provider = context.read<DriverProvider>();
      provider.resetFilters();
      provider.fetchDrivers(businessId: widget.businessId);
    }
  }

  @override
  void dispose() {
    _searchController.dispose();
    super.dispose();
  }

  void _onSearch(DriverProvider provider) {
    provider.setFilters(search: _searchController.text);
    provider.fetchDrivers(businessId: widget.businessId);
  }

  void _onClearSearch(DriverProvider provider) {
    _searchController.clear();
    provider.setFilters(search: '');
    provider.fetchDrivers(businessId: widget.businessId);
  }

  void _showDriverForm({DriverInfo? driver}) {
    final isEditing = driver != null;
    final firstNameCtrl = TextEditingController(text: driver?.firstName ?? '');
    final lastNameCtrl = TextEditingController(text: driver?.lastName ?? '');
    final identificationCtrl =
        TextEditingController(text: driver?.identification ?? '');
    final phoneCtrl = TextEditingController(text: driver?.phone ?? '');
    final emailCtrl = TextEditingController(text: driver?.email ?? '');
    final notesCtrl = TextEditingController(text: driver?.notes ?? '');
    String licenseType = driver?.licenseType ?? '';
    String status = driver?.status ?? 'active';
    final formKey = GlobalKey<FormState>();
    bool isSaving = false;

    const licenseTypes = ['', 'A1', 'A2', 'B1', 'B2', 'C1'];

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
                              isEditing
                                  ? 'Editar conductor'
                                  : 'Nuevo conductor',
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
                      TextFormField(
                        controller: firstNameCtrl,
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
                        controller: lastNameCtrl,
                        decoration: const InputDecoration(
                          labelText: 'Apellido *',
                          border: OutlineInputBorder(),
                        ),
                        validator: (v) => (v == null || v.trim().length < 2)
                            ? 'Minimo 2 caracteres'
                            : null,
                      ),
                      const SizedBox(height: 12),
                      TextFormField(
                        controller: identificationCtrl,
                        decoration: const InputDecoration(
                          labelText: 'Identificacion *',
                          border: OutlineInputBorder(),
                        ),
                        validator: (v) => (v == null || v.trim().isEmpty)
                            ? 'Requerido'
                            : null,
                      ),
                      const SizedBox(height: 12),
                      TextFormField(
                        controller: phoneCtrl,
                        decoration: const InputDecoration(
                          labelText: 'Telefono *',
                          border: OutlineInputBorder(),
                        ),
                        keyboardType: TextInputType.phone,
                        validator: (v) => (v == null || v.trim().isEmpty)
                            ? 'Requerido'
                            : null,
                      ),
                      const SizedBox(height: 12),
                      TextFormField(
                        controller: emailCtrl,
                        decoration: const InputDecoration(
                          labelText: 'Email',
                          border: OutlineInputBorder(),
                        ),
                        keyboardType: TextInputType.emailAddress,
                      ),
                      const SizedBox(height: 12),
                      DropdownButtonFormField<String>(
                        initialValue: licenseType,
                        decoration: const InputDecoration(
                          labelText: 'Tipo de licencia',
                          border: OutlineInputBorder(),
                        ),
                        items: licenseTypes
                            .map((t) => DropdownMenuItem(
                                  value: t,
                                  child: Text(
                                      t.isEmpty ? 'Sin especificar' : t),
                                ))
                            .toList(),
                        onChanged: (v) =>
                            setModalState(() => licenseType = v ?? ''),
                      ),
                      if (isEditing) ...[
                        const SizedBox(height: 12),
                        DropdownButtonFormField<String>(
                          initialValue: status,
                          decoration: const InputDecoration(
                            labelText: 'Estado',
                            border: OutlineInputBorder(),
                          ),
                          items: const [
                            DropdownMenuItem(
                                value: 'active', child: Text('Activo')),
                            DropdownMenuItem(
                                value: 'inactive', child: Text('Inactivo')),
                          ],
                          onChanged: (v) =>
                              setModalState(() => status = v ?? 'active'),
                        ),
                      ],
                      const SizedBox(height: 12),
                      TextFormField(
                        controller: notesCtrl,
                        decoration: const InputDecoration(
                          labelText: 'Notas',
                          border: OutlineInputBorder(),
                        ),
                        maxLines: 3,
                      ),
                      const SizedBox(height: 20),
                      FilledButton(
                        onPressed: isSaving
                            ? null
                            : () async {
                                if (!formKey.currentState!.validate()) return;
                                setModalState(() => isSaving = true);
                                final provider =
                                    context.read<DriverProvider>();
                                bool success;
                                if (isEditing) {
                                  final dto = UpdateDriverDTO(
                                    firstName: firstNameCtrl.text.trim(),
                                    lastName: lastNameCtrl.text.trim(),
                                    identification:
                                        identificationCtrl.text.trim(),
                                    phone: phoneCtrl.text.trim(),
                                    email: emailCtrl.text.trim().isNotEmpty
                                        ? emailCtrl.text.trim()
                                        : null,
                                    licenseType: licenseType.isNotEmpty
                                        ? licenseType
                                        : null,
                                    notes: notesCtrl.text.trim().isNotEmpty
                                        ? notesCtrl.text.trim()
                                        : null,
                                    status: status,
                                  );
                                  success = await provider.updateDriver(
                                    driver.id,
                                    dto,
                                    businessId: widget.businessId,
                                  );
                                } else {
                                  final dto = CreateDriverDTO(
                                    firstName: firstNameCtrl.text.trim(),
                                    lastName: lastNameCtrl.text.trim(),
                                    identification:
                                        identificationCtrl.text.trim(),
                                    phone: phoneCtrl.text.trim(),
                                    email: emailCtrl.text.trim().isNotEmpty
                                        ? emailCtrl.text.trim()
                                        : null,
                                    licenseType: licenseType.isNotEmpty
                                        ? licenseType
                                        : null,
                                    notes: notesCtrl.text.trim().isNotEmpty
                                        ? notesCtrl.text.trim()
                                        : null,
                                  );
                                  final result = await provider.createDriver(
                                    dto,
                                    businessId: widget.businessId,
                                  );
                                  success = result != null;
                                }
                                setModalState(() => isSaving = false);
                                if (success && ctx.mounted) {
                                  Navigator.pop(ctx);
                                  provider.fetchDrivers(
                                      businessId: widget.businessId);
                                  if (mounted) {
                                    ScaffoldMessenger.of(context).showSnackBar(
                                      SnackBar(
                                        content: Text(isEditing
                                            ? 'Conductor actualizado'
                                            : 'Conductor creado'),
                                      ),
                                    );
                                  }
                                } else if (mounted) {
                                  ScaffoldMessenger.of(context).showSnackBar(
                                    SnackBar(
                                      content: Text(
                                          provider.error ?? 'Error al guardar'),
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

  void _confirmDelete(DriverInfo driver) {
    showDialog(
      context: context,
      builder: (ctx) => AlertDialog(
        title: const Text('Eliminar conductor'),
        content: Text(
            'Eliminar a "${driver.firstName} ${driver.lastName}"? Esta accion no se puede deshacer.'),
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
              final provider = context.read<DriverProvider>();
              final ok = await provider.deleteDriver(driver.id,
                  businessId: widget.businessId);
              if (ok) {
                provider.fetchDrivers(businessId: widget.businessId);
                if (mounted) {
                  ScaffoldMessenger.of(context).showSnackBar(
                    const SnackBar(content: Text('Conductor eliminado')),
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

  Color _statusColor(String status) {
    switch (status) {
      case 'active':
        return Colors.green;
      case 'inactive':
        return Colors.grey;
      case 'on_route':
        return Colors.blue;
      default:
        return Colors.grey;
    }
  }

  String _statusLabel(String status) {
    switch (status) {
      case 'active':
        return 'Activo';
      case 'inactive':
        return 'Inactivo';
      case 'on_route':
        return 'En ruta';
      default:
        return status;
    }
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    return Scaffold(
      appBar: AppBar(title: const Text('Conductores')),
      floatingActionButton: FloatingActionButton(
        onPressed: () => _showDriverForm(),
        child: const Icon(Icons.add),
      ),
      body: Consumer<DriverProvider>(
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
                          hintText: 'Buscar por nombre, identificacion...',
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
              Expanded(
                child: _buildContent(provider, theme),
              ),

              // Pagination
              if (provider.pagination != null && !provider.isLoading)
                _buildPagination(provider),
            ],
          );
        },
      ),
    );
  }

  Widget _buildContent(DriverProvider provider, ThemeData theme) {
    if (provider.isLoading && provider.drivers.isEmpty) {
      return const Center(child: CircularProgressIndicator());
    }

    if (provider.error != null && provider.drivers.isEmpty) {
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
                  provider.fetchDrivers(businessId: widget.businessId),
              icon: const Icon(Icons.refresh),
              label: const Text('Reintentar'),
            ),
          ],
        ),
      );
    }

    if (provider.drivers.isEmpty) {
      return Center(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Icon(Icons.person_off, size: 48, color: theme.disabledColor),
            const SizedBox(height: 16),
            const Text('No hay conductores registrados'),
          ],
        ),
      );
    }

    return RefreshIndicator(
      onRefresh: () => provider.fetchDrivers(businessId: widget.businessId),
      child: ListView.builder(
        padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 4),
        itemCount: provider.drivers.length,
        itemBuilder: (context, index) {
          final driver = provider.drivers[index];
          return Card(
            margin: const EdgeInsets.only(bottom: 8),
            child: ListTile(
              leading: CircleAvatar(
                child: Text(
                  driver.firstName.isNotEmpty
                      ? driver.firstName[0].toUpperCase()
                      : '?',
                ),
              ),
              title: Text('${driver.firstName} ${driver.lastName}'),
              subtitle: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  if (driver.identification.isNotEmpty)
                    Text('ID: ${driver.identification}',
                        style: theme.textTheme.bodySmall),
                  if (driver.phone.isNotEmpty)
                    Text(driver.phone, style: theme.textTheme.bodySmall),
                  if (driver.licenseType.isNotEmpty)
                    Text('Licencia: ${driver.licenseType}',
                        style: theme.textTheme.bodySmall),
                ],
              ),
              trailing: Row(
                mainAxisSize: MainAxisSize.min,
                children: [
                  Container(
                    padding:
                        const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
                    decoration: BoxDecoration(
                      color: _statusColor(driver.status).withValues(alpha: 0.15),
                      borderRadius: BorderRadius.circular(12),
                    ),
                    child: Text(
                      _statusLabel(driver.status),
                      style: TextStyle(
                        fontSize: 11,
                        fontWeight: FontWeight.w600,
                        color: _statusColor(driver.status),
                      ),
                    ),
                  ),
                  PopupMenuButton<String>(
                    onSelected: (action) {
                      if (action == 'edit') {
                        _showDriverForm(driver: driver);
                      } else if (action == 'delete') {
                        _confirmDelete(driver);
                      }
                    },
                    itemBuilder: (_) => [
                      const PopupMenuItem(
                          value: 'edit', child: Text('Editar')),
                      const PopupMenuItem(
                          value: 'delete', child: Text('Eliminar')),
                    ],
                  ),
                ],
              ),
              isThreeLine: true,
            ),
          );
        },
      ),
    );
  }

  Widget _buildPagination(DriverProvider provider) {
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
                        provider.fetchDrivers(businessId: widget.businessId);
                      }
                    : null,
              ),
              IconButton(
                icon: const Icon(Icons.chevron_right),
                onPressed: pagination.hasNext
                    ? () {
                        provider.setPage(pagination.currentPage + 1);
                        provider.fetchDrivers(businessId: widget.businessId);
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
