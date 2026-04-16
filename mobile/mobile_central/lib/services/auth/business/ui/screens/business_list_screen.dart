import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../../../../../shared/types/paginated_response.dart';
import '../../../../../shared/widgets/network_avatar.dart';
import '../../domain/entities.dart';
import '../providers/business_provider.dart';

class BusinessListScreen extends StatefulWidget {
  const BusinessListScreen({super.key});

  @override
  State<BusinessListScreen> createState() => _BusinessListScreenState();
}

class _BusinessListScreenState extends State<BusinessListScreen> {
  int _currentPage = 1;
  static const int _pageSize = 20;

  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addPostFrameCallback((_) {
      _fetchData();
    });
  }

  void _fetchData() {
    final provider = context.read<BusinessProvider>();
    provider.fetchBusinesses(
      params: GetBusinessesParams(page: _currentPage, pageSize: _pageSize),
    );
    provider.fetchBusinessTypes();
  }

  Future<void> _refresh() async {
    await context.read<BusinessProvider>().fetchBusinesses(
          params:
              GetBusinessesParams(page: _currentPage, pageSize: _pageSize),
        );
  }

  void _goToPage(int page) {
    setState(() => _currentPage = page);
    context.read<BusinessProvider>().fetchBusinesses(
          params: GetBusinessesParams(page: page, pageSize: _pageSize),
        );
  }

  void _showCreateForm() {
    _showBusinessForm(null);
  }

  void _showEditForm(Business business) {
    _showBusinessForm(business);
  }

  void _showBusinessForm(Business? business) {
    final isEditing = business != null;
    final nameController = TextEditingController(text: business?.name ?? '');
    final domainController =
        TextEditingController(text: business?.domain ?? '');
    final primaryColorController =
        TextEditingController(text: business?.primaryColor ?? '');
    final secondaryColorController =
        TextEditingController(text: business?.secondaryColor ?? '');
    final accentColorController =
        TextEditingController(text: business?.accentColor ?? '');
    final navbarColorController =
        TextEditingController(text: business?.navbarColor ?? '');

    bool hasDelivery = business?.hasDelivery ?? false;
    bool hasPickup = business?.hasPickup ?? false;
    int? selectedTypeId = business?.businessTypeId;
    bool isSaving = false;

    showModalBottomSheet(
      context: context,
      isScrollControlled: true,
      useSafeArea: true,
      shape: const RoundedRectangleBorder(
        borderRadius: BorderRadius.vertical(top: Radius.circular(16)),
      ),
      builder: (sheetContext) {
        return StatefulBuilder(
          builder: (builderContext, setSheetState) {
            final provider = context.read<BusinessProvider>();
            final businessTypes = provider.businessTypes;

            return Padding(
              padding: EdgeInsets.only(
                bottom: MediaQuery.of(builderContext).viewInsets.bottom,
                left: 16,
                right: 16,
                top: 16,
              ),
              child: SingleChildScrollView(
                child: Column(
                  mainAxisSize: MainAxisSize.min,
                  crossAxisAlignment: CrossAxisAlignment.stretch,
                  children: [
                    // Handle bar
                    Center(
                      child: Container(
                        width: 40,
                        height: 4,
                        decoration: BoxDecoration(
                          color: Colors.grey[300],
                          borderRadius: BorderRadius.circular(2),
                        ),
                      ),
                    ),
                    const SizedBox(height: 16),
                    Text(
                      isEditing ? 'Editar Negocio' : 'Crear Negocio',
                      style:
                          Theme.of(builderContext).textTheme.titleLarge,
                    ),
                    const SizedBox(height: 24),

                    // Name
                    TextFormField(
                      controller: nameController,
                      decoration: const InputDecoration(
                        labelText: 'Nombre *',
                        prefixIcon: Icon(Icons.business),
                      ),
                    ),
                    const SizedBox(height: 16),

                    // Domain
                    TextFormField(
                      controller: domainController,
                      decoration: const InputDecoration(
                        labelText: 'Dominio',
                        prefixIcon: Icon(Icons.language),
                        hintText: 'ejemplo.com',
                      ),
                    ),
                    const SizedBox(height: 16),

                    // Business Type dropdown
                    DropdownButtonFormField<int?>(
                      initialValue: selectedTypeId,
                      decoration: const InputDecoration(
                        labelText: 'Tipo de negocio',
                        prefixIcon: Icon(Icons.category),
                      ),
                      items: [
                        const DropdownMenuItem<int?>(
                          value: null,
                          child: Text('Sin tipo'),
                        ),
                        ...businessTypes.map(
                          (type) => DropdownMenuItem<int?>(
                            value: type.id,
                            child: Text(type.name),
                          ),
                        ),
                      ],
                      onChanged: (value) {
                        setSheetState(() => selectedTypeId = value);
                      },
                    ),
                    const SizedBox(height: 16),

                    // Colors section
                    Text(
                      'Colores',
                      style: Theme.of(builderContext)
                          .textTheme
                          .titleSmall
                          ?.copyWith(
                            color: Colors.grey[700],
                          ),
                    ),
                    const SizedBox(height: 8),
                    Row(
                      children: [
                        Expanded(
                          child: TextFormField(
                            controller: primaryColorController,
                            decoration: const InputDecoration(
                              labelText: 'Primario',
                              hintText: '#000000',
                              isDense: true,
                            ),
                          ),
                        ),
                        const SizedBox(width: 8),
                        Expanded(
                          child: TextFormField(
                            controller: secondaryColorController,
                            decoration: const InputDecoration(
                              labelText: 'Secundario',
                              hintText: '#FFFFFF',
                              isDense: true,
                            ),
                          ),
                        ),
                      ],
                    ),
                    const SizedBox(height: 8),
                    Row(
                      children: [
                        Expanded(
                          child: TextFormField(
                            controller: accentColorController,
                            decoration: const InputDecoration(
                              labelText: 'Acento',
                              hintText: '#3B82F6',
                              isDense: true,
                            ),
                          ),
                        ),
                        const SizedBox(width: 8),
                        Expanded(
                          child: TextFormField(
                            controller: navbarColorController,
                            decoration: const InputDecoration(
                              labelText: 'Navbar',
                              hintText: '#1E3A5F',
                              isDense: true,
                            ),
                          ),
                        ),
                      ],
                    ),
                    const SizedBox(height: 16),

                    // Toggles
                    SwitchListTile(
                      title: const Text('Delivery'),
                      subtitle: const Text('Habilitar servicio de entrega'),
                      value: hasDelivery,
                      contentPadding: EdgeInsets.zero,
                      onChanged: (value) {
                        setSheetState(() => hasDelivery = value);
                      },
                    ),
                    SwitchListTile(
                      title: const Text('Pickup'),
                      subtitle:
                          const Text('Habilitar recogida en tienda'),
                      value: hasPickup,
                      contentPadding: EdgeInsets.zero,
                      onChanged: (value) {
                        setSheetState(() => hasPickup = value);
                      },
                    ),
                    const SizedBox(height: 16),

                    // Actions
                    Row(
                      children: [
                        Expanded(
                          child: OutlinedButton(
                            onPressed: isSaving
                                ? null
                                : () => Navigator.pop(builderContext),
                            child: const Text('Cancelar'),
                          ),
                        ),
                        const SizedBox(width: 12),
                        Expanded(
                          child: FilledButton(
                            onPressed: isSaving
                                ? null
                                : () async {
                                    final name =
                                        nameController.text.trim();
                                    if (name.isEmpty) {
                                      ScaffoldMessenger.of(builderContext)
                                          .showSnackBar(
                                        const SnackBar(
                                          content: Text(
                                              'El nombre es obligatorio'),
                                        ),
                                      );
                                      return;
                                    }

                                    setSheetState(
                                        () => isSaving = true);

                                    final prov =
                                        context.read<BusinessProvider>();
                                    bool success;

                                    if (isEditing) {
                                      success =
                                          await prov.updateBusiness(
                                        business.id,
                                        UpdateBusinessDTO(
                                          name: name,
                                          primaryColor:
                                              primaryColorController
                                                      .text
                                                      .trim()
                                                      .isNotEmpty
                                                  ? primaryColorController
                                                      .text
                                                      .trim()
                                                  : null,
                                          secondaryColor:
                                              secondaryColorController
                                                      .text
                                                      .trim()
                                                      .isNotEmpty
                                                  ? secondaryColorController
                                                      .text
                                                      .trim()
                                                  : null,
                                          accentColor:
                                              accentColorController
                                                      .text
                                                      .trim()
                                                      .isNotEmpty
                                                  ? accentColorController
                                                      .text
                                                      .trim()
                                                  : null,
                                          navbarColor:
                                              navbarColorController
                                                      .text
                                                      .trim()
                                                      .isNotEmpty
                                                  ? navbarColorController
                                                      .text
                                                      .trim()
                                                  : null,
                                          domain: domainController
                                                  .text
                                                  .trim()
                                                  .isNotEmpty
                                              ? domainController.text
                                                  .trim()
                                              : null,
                                          hasDelivery: hasDelivery,
                                          hasPickup: hasPickup,
                                          businessTypeId:
                                              selectedTypeId,
                                        ),
                                      );
                                    } else {
                                      final created =
                                          await prov.createBusiness(
                                        CreateBusinessDTO(
                                          name: name,
                                          primaryColor:
                                              primaryColorController
                                                      .text
                                                      .trim()
                                                      .isNotEmpty
                                                  ? primaryColorController
                                                      .text
                                                      .trim()
                                                  : null,
                                          secondaryColor:
                                              secondaryColorController
                                                      .text
                                                      .trim()
                                                      .isNotEmpty
                                                  ? secondaryColorController
                                                      .text
                                                      .trim()
                                                  : null,
                                          accentColor:
                                              accentColorController
                                                      .text
                                                      .trim()
                                                      .isNotEmpty
                                                  ? accentColorController
                                                      .text
                                                      .trim()
                                                  : null,
                                          navbarColor:
                                              navbarColorController
                                                      .text
                                                      .trim()
                                                      .isNotEmpty
                                                  ? navbarColorController
                                                      .text
                                                      .trim()
                                                  : null,
                                          domain: domainController
                                                  .text
                                                  .trim()
                                                  .isNotEmpty
                                              ? domainController.text
                                                  .trim()
                                              : null,
                                          hasDelivery: hasDelivery,
                                          hasPickup: hasPickup,
                                          businessTypeId:
                                              selectedTypeId,
                                        ),
                                      );
                                      success = created != null;
                                    }

                                    setSheetState(
                                        () => isSaving = false);

                                    if (!builderContext.mounted) return;

                                    if (success) {
                                      Navigator.pop(builderContext);
                                      _refresh();
                                      if (mounted) {
                                        ScaffoldMessenger.of(context)
                                            .showSnackBar(
                                          SnackBar(
                                            content: Text(isEditing
                                                ? 'Negocio actualizado'
                                                : 'Negocio creado'),
                                          ),
                                        );
                                      }
                                    } else {
                                      ScaffoldMessenger.of(builderContext)
                                          .showSnackBar(
                                        SnackBar(
                                          content: Text(
                                              prov.error ??
                                                  'Error al guardar'),
                                          backgroundColor:
                                              Colors.red,
                                        ),
                                      );
                                    }
                                  },
                            child: isSaving
                                ? const SizedBox(
                                    height: 20,
                                    width: 20,
                                    child: CircularProgressIndicator(
                                      strokeWidth: 2,
                                      color: Colors.white,
                                    ),
                                  )
                                : Text(isEditing
                                    ? 'Actualizar'
                                    : 'Crear'),
                          ),
                        ),
                      ],
                    ),
                    const SizedBox(height: 16),
                  ],
                ),
              ),
            );
          },
        );
      },
    );
  }

  void _confirmDelete(Business business) {
    showDialog(
      context: context,
      builder: (dialogContext) {
        return AlertDialog(
          title: const Text('Eliminar Negocio'),
          content: Text(
              '¿Eliminar "${business.name}"? Esta accion no se puede deshacer.'),
          actions: [
            TextButton(
              onPressed: () => Navigator.pop(dialogContext),
              child: const Text('Cancelar'),
            ),
            FilledButton(
              style: FilledButton.styleFrom(
                backgroundColor: Colors.red,
              ),
              onPressed: () async {
                Navigator.pop(dialogContext);
                final provider = context.read<BusinessProvider>();
                final success = await provider.deleteBusiness(business.id);
                if (success && mounted) {
                  ScaffoldMessenger.of(context).showSnackBar(
                    const SnackBar(content: Text('Negocio eliminado')),
                  );
                  _refresh();
                } else if (mounted) {
                  ScaffoldMessenger.of(context).showSnackBar(
                    SnackBar(
                      content: Text(provider.error ?? 'Error al eliminar'),
                      backgroundColor: Colors.red,
                    ),
                  );
                }
              },
              child: const Text('Eliminar'),
            ),
          ],
        );
      },
    );
  }

  Future<void> _toggleActive(Business business) async {
    final provider = context.read<BusinessProvider>();
    bool success;
    if (business.isActive) {
      success = await provider.deactivateBusiness(business.id);
    } else {
      success = await provider.activateBusiness(business.id);
    }
    if (success && mounted) {
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(
          content: Text(business.isActive
              ? 'Negocio desactivado'
              : 'Negocio activado'),
        ),
      );
      _refresh();
    } else if (mounted) {
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(
          content: Text(provider.error ?? 'Error al cambiar estado'),
          backgroundColor: Colors.red,
        ),
      );
    }
  }

  @override
  Widget build(BuildContext context) {
    final colorScheme = Theme.of(context).colorScheme;

    return Scaffold(
      appBar: AppBar(
        title: const Text('Negocios'),
      ),
      floatingActionButton: FloatingActionButton(
        onPressed: _showCreateForm,
        child: const Icon(Icons.add),
      ),
      body: Consumer<BusinessProvider>(
        builder: (context, provider, child) {
          if (provider.isLoading && provider.businesses.isEmpty) {
            return const Center(child: CircularProgressIndicator());
          }

          if (provider.error != null && provider.businesses.isEmpty) {
            return Center(
              child: Padding(
                padding: const EdgeInsets.all(24),
                child: Column(
                  mainAxisAlignment: MainAxisAlignment.center,
                  children: [
                    Icon(Icons.error_outline,
                        size: 48, color: colorScheme.error),
                    const SizedBox(height: 16),
                    Text(
                      provider.error!,
                      textAlign: TextAlign.center,
                      style: TextStyle(color: colorScheme.error),
                    ),
                    const SizedBox(height: 16),
                    FilledButton.icon(
                      onPressed: _fetchData,
                      icon: const Icon(Icons.refresh),
                      label: const Text('Reintentar'),
                    ),
                  ],
                ),
              ),
            );
          }

          if (provider.businesses.isEmpty) {
            return Center(
              child: Column(
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  Icon(Icons.store_outlined,
                      size: 64, color: Colors.grey[400]),
                  const SizedBox(height: 16),
                  Text(
                    'No hay negocios',
                    style: Theme.of(context).textTheme.titleMedium?.copyWith(
                          color: Colors.grey[600],
                        ),
                  ),
                  const SizedBox(height: 8),
                  Text(
                    'Crea tu primer negocio',
                    style: Theme.of(context).textTheme.bodyMedium?.copyWith(
                          color: Colors.grey[500],
                        ),
                  ),
                ],
              ),
            );
          }

          return Column(
            children: [
              // List
              Expanded(
                child: RefreshIndicator(
                  onRefresh: _refresh,
                  child: ListView.builder(
                    padding: const EdgeInsets.fromLTRB(16, 8, 16, 80),
                    itemCount: provider.businesses.length,
                    itemBuilder: (context, index) {
                      final business = provider.businesses[index];
                      return _BusinessCard(
                        business: business,
                        onEdit: () => _showEditForm(business),
                        onDelete: () => _confirmDelete(business),
                        onToggleActive: () => _toggleActive(business),
                      );
                    },
                  ),
                ),
              ),

              // Pagination
              if (provider.pagination != null &&
                  provider.pagination!.lastPage > 1)
                _PaginationBar(
                  pagination: provider.pagination!,
                  currentPage: _currentPage,
                  onPageChanged: _goToPage,
                ),
            ],
          );
        },
      ),
    );
  }
}

class _BusinessCard extends StatelessWidget {
  final Business business;
  final VoidCallback onEdit;
  final VoidCallback onDelete;
  final VoidCallback onToggleActive;

  const _BusinessCard({
    required this.business,
    required this.onEdit,
    required this.onDelete,
    required this.onToggleActive,
  });

  Color? _parseColor(String? hex) {
    if (hex == null || hex.isEmpty) return null;
    final cleaned = hex.replaceFirst('#', '');
    if (cleaned.length != 6) return null;
    final value = int.tryParse(cleaned, radix: 16);
    if (value == null) return null;
    return Color(0xFF000000 | value);
  }

  @override
  Widget build(BuildContext context) {
    final primaryColor = _parseColor(business.primaryColor);
    final secondaryColor = _parseColor(business.secondaryColor);

    return Card(
      margin: const EdgeInsets.only(bottom: 8),
      child: InkWell(
        borderRadius: BorderRadius.circular(12),
        onTap: onEdit,
        child: Padding(
          padding: const EdgeInsets.all(12),
          child: Row(
            children: [
              // Logo / Avatar
              NetworkAvatar(
                imageUrl: business.logoUrl,
                fallbackText: business.name,
                fallbackIcon: Icons.business,
                radius: 24,
                backgroundColor: primaryColor ?? Colors.deepPurple,
                foregroundColor: secondaryColor ?? Colors.white,
              ),
              const SizedBox(width: 12),

              // Info
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(
                      business.name,
                      style: Theme.of(context).textTheme.titleSmall?.copyWith(
                            fontWeight: FontWeight.w600,
                          ),
                      maxLines: 1,
                      overflow: TextOverflow.ellipsis,
                    ),
                    const SizedBox(height: 4),
                    Row(
                      children: [
                        if (business.businessTypeName != null) ...[
                          Icon(Icons.category,
                              size: 14, color: Colors.grey[500]),
                          const SizedBox(width: 4),
                          Text(
                            business.businessTypeName!,
                            style: Theme.of(context)
                                .textTheme
                                .bodySmall
                                ?.copyWith(color: Colors.grey[600]),
                          ),
                          const SizedBox(width: 12),
                        ],
                        Text(
                          'ID: ${business.id}',
                          style: Theme.of(context)
                              .textTheme
                              .bodySmall
                              ?.copyWith(color: Colors.grey[500]),
                        ),
                      ],
                    ),
                    if (business.primaryColor != null ||
                        business.secondaryColor != null) ...[
                      const SizedBox(height: 6),
                      Row(
                        children: [
                          if (primaryColor != null)
                            Container(
                              width: 16,
                              height: 16,
                              margin: const EdgeInsets.only(right: 4),
                              decoration: BoxDecoration(
                                color: primaryColor,
                                borderRadius: BorderRadius.circular(4),
                                border: Border.all(
                                    color: Colors.grey[300]!, width: 0.5),
                              ),
                            ),
                          if (secondaryColor != null)
                            Container(
                              width: 16,
                              height: 16,
                              margin: const EdgeInsets.only(right: 4),
                              decoration: BoxDecoration(
                                color: secondaryColor,
                                borderRadius: BorderRadius.circular(4),
                                border: Border.all(
                                    color: Colors.grey[300]!, width: 0.5),
                              ),
                            ),
                          if (_parseColor(business.accentColor) != null)
                            Container(
                              width: 16,
                              height: 16,
                              decoration: BoxDecoration(
                                color: _parseColor(business.accentColor),
                                borderRadius: BorderRadius.circular(4),
                                border: Border.all(
                                    color: Colors.grey[300]!, width: 0.5),
                              ),
                            ),
                        ],
                      ),
                    ],
                  ],
                ),
              ),

              // Active status toggle + actions
              Column(
                crossAxisAlignment: CrossAxisAlignment.end,
                children: [
                  InkWell(
                    borderRadius: BorderRadius.circular(12),
                    onTap: onToggleActive,
                    child: Container(
                      padding: const EdgeInsets.symmetric(
                          horizontal: 10, vertical: 4),
                      decoration: BoxDecoration(
                        color: business.isActive
                            ? Colors.green[50]
                            : Colors.red[50],
                        borderRadius: BorderRadius.circular(12),
                      ),
                      child: Text(
                        business.isActive ? 'Activo' : 'Inactivo',
                        style: TextStyle(
                          fontSize: 11,
                          fontWeight: FontWeight.w600,
                          color: business.isActive
                              ? Colors.green[700]
                              : Colors.red[700],
                        ),
                      ),
                    ),
                  ),
                  const SizedBox(height: 8),
                  Row(
                    mainAxisSize: MainAxisSize.min,
                    children: [
                      InkWell(
                        borderRadius: BorderRadius.circular(8),
                        onTap: onEdit,
                        child: Padding(
                          padding: const EdgeInsets.all(4),
                          child: Icon(Icons.edit_outlined,
                              size: 20, color: Colors.amber[700]),
                        ),
                      ),
                      const SizedBox(width: 4),
                      InkWell(
                        borderRadius: BorderRadius.circular(8),
                        onTap: onDelete,
                        child: Padding(
                          padding: const EdgeInsets.all(4),
                          child:
                              Icon(Icons.delete_outline,
                                  size: 20, color: Colors.red[400]),
                        ),
                      ),
                    ],
                  ),
                ],
              ),
            ],
          ),
        ),
      ),
    );
  }
}

class _PaginationBar extends StatelessWidget {
  final Pagination pagination;
  final int currentPage;
  final ValueChanged<int> onPageChanged;

  const _PaginationBar({
    required this.pagination,
    required this.currentPage,
    required this.onPageChanged,
  });

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
      decoration: BoxDecoration(
        color: Theme.of(context).colorScheme.surface,
        border: Border(
          top: BorderSide(color: Colors.grey[200]!),
        ),
      ),
      child: Row(
        mainAxisAlignment: MainAxisAlignment.spaceBetween,
        children: [
          Text(
            '${pagination.total} negocios',
            style: Theme.of(context).textTheme.bodySmall?.copyWith(
                  color: Colors.grey[600],
                ),
          ),
          Row(
            children: [
              IconButton(
                icon: const Icon(Icons.chevron_left),
                onPressed:
                    pagination.hasPrev ? () => onPageChanged(currentPage - 1) : null,
                iconSize: 20,
                visualDensity: VisualDensity.compact,
              ),
              Text(
                '$currentPage / ${pagination.lastPage}',
                style: Theme.of(context).textTheme.bodySmall?.copyWith(
                      fontWeight: FontWeight.w600,
                    ),
              ),
              IconButton(
                icon: const Icon(Icons.chevron_right),
                onPressed:
                    pagination.hasNext ? () => onPageChanged(currentPage + 1) : null,
                iconSize: 20,
                visualDensity: VisualDensity.compact,
              ),
            ],
          ),
        ],
      ),
    );
  }
}
