import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../providers/vehicle_provider.dart';
import '../../domain/entities.dart';

class VehicleListScreen extends StatefulWidget {
  final int? businessId;

  const VehicleListScreen({super.key, this.businessId});

  @override
  State<VehicleListScreen> createState() => _VehicleListScreenState();
}

class _VehicleListScreenState extends State<VehicleListScreen> {
  final _searchController = TextEditingController();

  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addPostFrameCallback((_) {
      context.read<VehicleProvider>().fetchVehicles(businessId: widget.businessId);
    });
  }

  @override
  void didUpdateWidget(VehicleListScreen oldWidget) {
    super.didUpdateWidget(oldWidget);
    if (oldWidget.businessId != widget.businessId) {
      _searchController.clear();
      context.read<VehicleProvider>().fetchVehicles(businessId: widget.businessId);
    }
  }

  @override
  void dispose() {
    _searchController.dispose();
    super.dispose();
  }

  void _onSearch(VehicleProvider provider) {
    provider.setPage(1);
    provider.fetchVehicles(businessId: widget.businessId);
  }

  void _onClearSearch(VehicleProvider provider) {
    _searchController.clear();
    provider.setPage(1);
    provider.fetchVehicles(businessId: widget.businessId);
  }

  static const _vehicleTypes = [
    {'value': 'motorcycle', 'label': 'Motocicleta'},
    {'value': 'car', 'label': 'Carro'},
    {'value': 'van', 'label': 'Van'},
    {'value': 'truck', 'label': 'Camion'},
  ];

  static const _vehicleStatuses = [
    {'value': 'active', 'label': 'Activo'},
    {'value': 'inactive', 'label': 'Inactivo'},
    {'value': 'in_maintenance', 'label': 'En mantenimiento'},
  ];

  IconData _typeIcon(String type) {
    switch (type) {
      case 'motorcycle':
        return Icons.two_wheeler;
      case 'car':
        return Icons.directions_car;
      case 'van':
        return Icons.airport_shuttle;
      case 'truck':
        return Icons.local_shipping;
      default:
        return Icons.directions_car;
    }
  }

  Color _statusColor(String status) {
    switch (status) {
      case 'active':
        return Colors.green;
      case 'inactive':
        return Colors.grey;
      case 'in_maintenance':
        return Colors.orange;
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
      case 'in_maintenance':
        return 'Mantenimiento';
      default:
        return status;
    }
  }

  String _typeLabel(String type) {
    for (final t in _vehicleTypes) {
      if (t['value'] == type) return t['label']!;
    }
    return type;
  }

  void _showVehicleForm({VehicleInfo? vehicle}) {
    final isEditing = vehicle != null;
    final licensePlateCtrl =
        TextEditingController(text: vehicle?.licensePlate ?? '');
    final brandCtrl = TextEditingController(text: vehicle?.brand ?? '');
    final modelCtrl = TextEditingController(text: vehicle?.model ?? '');
    final yearCtrl = TextEditingController(
        text: vehicle?.year != null ? vehicle!.year.toString() : '');
    final colorCtrl = TextEditingController(text: vehicle?.color ?? '');
    final weightCtrl = TextEditingController(
        text: vehicle?.weightCapacityKg != null
            ? vehicle!.weightCapacityKg.toString()
            : '');
    final volumeCtrl = TextEditingController(
        text: vehicle?.volumeCapacityM3 != null
            ? vehicle!.volumeCapacityM3.toString()
            : '');
    String selectedType = vehicle?.type ?? 'motorcycle';
    String selectedStatus = vehicle?.status ?? 'active';
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
                              isEditing ? 'Editar vehiculo' : 'Nuevo vehiculo',
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
                      DropdownButtonFormField<String>(
                        initialValue: selectedType,
                        decoration: const InputDecoration(
                          labelText: 'Tipo *',
                          border: OutlineInputBorder(),
                        ),
                        items: _vehicleTypes
                            .map((t) => DropdownMenuItem(
                                  value: t['value'],
                                  child: Text(t['label']!),
                                ))
                            .toList(),
                        onChanged: (v) =>
                            setModalState(() => selectedType = v ?? 'motorcycle'),
                      ),
                      const SizedBox(height: 12),
                      TextFormField(
                        controller: licensePlateCtrl,
                        decoration: const InputDecoration(
                          labelText: 'Placa *',
                          border: OutlineInputBorder(),
                          hintText: 'ABC123',
                        ),
                        textCapitalization: TextCapitalization.characters,
                        validator: (v) => (v == null || v.trim().isEmpty)
                            ? 'Requerido'
                            : null,
                      ),
                      const SizedBox(height: 12),
                      Row(
                        children: [
                          Expanded(
                            child: TextFormField(
                              controller: brandCtrl,
                              decoration: const InputDecoration(
                                labelText: 'Marca',
                                border: OutlineInputBorder(),
                              ),
                            ),
                          ),
                          const SizedBox(width: 12),
                          Expanded(
                            child: TextFormField(
                              controller: modelCtrl,
                              decoration: const InputDecoration(
                                labelText: 'Modelo',
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
                              controller: yearCtrl,
                              decoration: const InputDecoration(
                                labelText: 'Ano',
                                border: OutlineInputBorder(),
                              ),
                              keyboardType: TextInputType.number,
                            ),
                          ),
                          const SizedBox(width: 12),
                          Expanded(
                            child: TextFormField(
                              controller: colorCtrl,
                              decoration: const InputDecoration(
                                labelText: 'Color',
                                border: OutlineInputBorder(),
                              ),
                            ),
                          ),
                        ],
                      ),
                      if (isEditing) ...[
                        const SizedBox(height: 12),
                        DropdownButtonFormField<String>(
                          initialValue: selectedStatus,
                          decoration: const InputDecoration(
                            labelText: 'Estado',
                            border: OutlineInputBorder(),
                          ),
                          items: _vehicleStatuses
                              .map((s) => DropdownMenuItem(
                                    value: s['value'],
                                    child: Text(s['label']!),
                                  ))
                              .toList(),
                          onChanged: (v) =>
                              setModalState(() => selectedStatus = v ?? 'active'),
                        ),
                      ],
                      const SizedBox(height: 12),
                      Row(
                        children: [
                          Expanded(
                            child: TextFormField(
                              controller: weightCtrl,
                              decoration: const InputDecoration(
                                labelText: 'Peso max (kg)',
                                border: OutlineInputBorder(),
                              ),
                              keyboardType:
                                  const TextInputType.numberWithOptions(
                                      decimal: true),
                            ),
                          ),
                          const SizedBox(width: 12),
                          Expanded(
                            child: TextFormField(
                              controller: volumeCtrl,
                              decoration: const InputDecoration(
                                labelText: 'Volumen max (m3)',
                                border: OutlineInputBorder(),
                              ),
                              keyboardType:
                                  const TextInputType.numberWithOptions(
                                      decimal: true),
                            ),
                          ),
                        ],
                      ),
                      const SizedBox(height: 20),
                      FilledButton(
                        onPressed: isSaving
                            ? null
                            : () async {
                                if (!formKey.currentState!.validate()) return;
                                setModalState(() => isSaving = true);
                                final provider =
                                    context.read<VehicleProvider>();
                                bool success;
                                if (isEditing) {
                                  final dto = UpdateVehicleDTO(
                                    type: selectedType,
                                    licensePlate: licensePlateCtrl.text.trim(),
                                    brand: brandCtrl.text.trim().isNotEmpty
                                        ? brandCtrl.text.trim()
                                        : null,
                                    model: modelCtrl.text.trim().isNotEmpty
                                        ? modelCtrl.text.trim()
                                        : null,
                                    year: yearCtrl.text.trim().isNotEmpty
                                        ? int.tryParse(yearCtrl.text.trim())
                                        : null,
                                    color: colorCtrl.text.trim().isNotEmpty
                                        ? colorCtrl.text.trim()
                                        : null,
                                    status: selectedStatus,
                                    weightCapacityKg:
                                        weightCtrl.text.trim().isNotEmpty
                                            ? double.tryParse(
                                                weightCtrl.text.trim())
                                            : null,
                                    volumeCapacityM3:
                                        volumeCtrl.text.trim().isNotEmpty
                                            ? double.tryParse(
                                                volumeCtrl.text.trim())
                                            : null,
                                  );
                                  final result = await provider.updateVehicle(
                                    vehicle.id,
                                    dto,
                                    businessId: widget.businessId,
                                  );
                                  success = result != null;
                                } else {
                                  final dto = CreateVehicleDTO(
                                    type: selectedType,
                                    licensePlate: licensePlateCtrl.text.trim(),
                                    brand: brandCtrl.text.trim().isNotEmpty
                                        ? brandCtrl.text.trim()
                                        : null,
                                    model: modelCtrl.text.trim().isNotEmpty
                                        ? modelCtrl.text.trim()
                                        : null,
                                    year: yearCtrl.text.trim().isNotEmpty
                                        ? int.tryParse(yearCtrl.text.trim())
                                        : null,
                                    color: colorCtrl.text.trim().isNotEmpty
                                        ? colorCtrl.text.trim()
                                        : null,
                                    weightCapacityKg:
                                        weightCtrl.text.trim().isNotEmpty
                                            ? double.tryParse(
                                                weightCtrl.text.trim())
                                            : null,
                                    volumeCapacityM3:
                                        volumeCtrl.text.trim().isNotEmpty
                                            ? double.tryParse(
                                                volumeCtrl.text.trim())
                                            : null,
                                  );
                                  final result = await provider.createVehicle(
                                    dto,
                                    businessId: widget.businessId,
                                  );
                                  success = result != null;
                                }
                                setModalState(() => isSaving = false);
                                if (success && ctx.mounted) {
                                  Navigator.pop(ctx);
                                  provider.fetchVehicles(
                                      businessId: widget.businessId);
                                  if (mounted) {
                                    ScaffoldMessenger.of(context).showSnackBar(
                                      SnackBar(
                                        content: Text(isEditing
                                            ? 'Vehiculo actualizado'
                                            : 'Vehiculo creado'),
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

  void _confirmDelete(VehicleInfo vehicle) {
    showDialog(
      context: context,
      builder: (ctx) => AlertDialog(
        title: const Text('Eliminar vehiculo'),
        content: Text(
            'Eliminar el vehiculo "${vehicle.licensePlate}"? Esta accion no se puede deshacer.'),
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
              final provider = context.read<VehicleProvider>();
              final ok = await provider.deleteVehicle(vehicle.id,
                  businessId: widget.businessId);
              if (ok) {
                provider.fetchVehicles(businessId: widget.businessId);
                if (mounted) {
                  ScaffoldMessenger.of(context).showSnackBar(
                    const SnackBar(content: Text('Vehiculo eliminado')),
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
      appBar: AppBar(title: const Text('Vehiculos')),
      floatingActionButton: FloatingActionButton(
        onPressed: () => _showVehicleForm(),
        child: const Icon(Icons.add),
      ),
      body: Consumer<VehicleProvider>(
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
                          hintText: 'Buscar por placa, marca o modelo...',
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

  Widget _buildContent(VehicleProvider provider, ThemeData theme) {
    if (provider.isLoading && provider.vehicles.isEmpty) {
      return const Center(child: CircularProgressIndicator());
    }

    if (provider.error != null && provider.vehicles.isEmpty) {
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
                  provider.fetchVehicles(businessId: widget.businessId),
              icon: const Icon(Icons.refresh),
              label: const Text('Reintentar'),
            ),
          ],
        ),
      );
    }

    if (provider.vehicles.isEmpty) {
      return Center(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Icon(Icons.directions_car_outlined,
                size: 48, color: theme.disabledColor),
            const SizedBox(height: 16),
            const Text('No hay vehiculos registrados'),
          ],
        ),
      );
    }

    return RefreshIndicator(
      onRefresh: () =>
          provider.fetchVehicles(businessId: widget.businessId),
      child: ListView.builder(
        padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 4),
        itemCount: provider.vehicles.length,
        itemBuilder: (context, index) {
          final vehicle = provider.vehicles[index];
          final brandModel = [vehicle.brand, vehicle.model]
              .where((s) => s.isNotEmpty)
              .join(' ');

          return Card(
            margin: const EdgeInsets.only(bottom: 8),
            child: ListTile(
              leading: CircleAvatar(
                child: Icon(_typeIcon(vehicle.type)),
              ),
              title: Text(vehicle.licensePlate,
                  style: const TextStyle(fontWeight: FontWeight.w600)),
              subtitle: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(_typeLabel(vehicle.type),
                      style: theme.textTheme.bodySmall),
                  if (brandModel.isNotEmpty)
                    Text(brandModel, style: theme.textTheme.bodySmall),
                  if (vehicle.year != null)
                    Text('Ano: ${vehicle.year}',
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
                      color: _statusColor(vehicle.status).withValues(alpha: 0.15),
                      borderRadius: BorderRadius.circular(12),
                    ),
                    child: Text(
                      _statusLabel(vehicle.status),
                      style: TextStyle(
                        fontSize: 11,
                        fontWeight: FontWeight.w600,
                        color: _statusColor(vehicle.status),
                      ),
                    ),
                  ),
                  PopupMenuButton<String>(
                    onSelected: (action) {
                      if (action == 'edit') {
                        _showVehicleForm(vehicle: vehicle);
                      } else if (action == 'delete') {
                        _confirmDelete(vehicle);
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

  Widget _buildPagination(VehicleProvider provider) {
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
                        provider.fetchVehicles(businessId: widget.businessId);
                      }
                    : null,
              ),
              IconButton(
                icon: const Icon(Icons.chevron_right),
                onPressed: pagination.hasNext
                    ? () {
                        provider.setPage(pagination.currentPage + 1);
                        provider.fetchVehicles(businessId: widget.businessId);
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
