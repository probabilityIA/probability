import 'dart:async';
import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../../../../../shared/types/paginated_response.dart';
import '../../domain/entities.dart';
import '../providers/customer_provider.dart';

class CustomerListScreen extends StatefulWidget {
  final int? businessId;

  const CustomerListScreen({super.key, this.businessId});

  @override
  State<CustomerListScreen> createState() => _CustomerListScreenState();
}

class _CustomerListScreenState extends State<CustomerListScreen> {
  final _searchController = TextEditingController();
  Timer? _debounce;

  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addPostFrameCallback((_) {
      _fetchData();
    });
  }

  @override
  void didUpdateWidget(CustomerListScreen oldWidget) {
    super.didUpdateWidget(oldWidget);
    if (oldWidget.businessId != widget.businessId) {
      final provider = context.read<CustomerProvider>();
      provider.resetFilters();
      _searchController.clear();
      provider.fetchCustomers(businessId: widget.businessId);
    }
  }

  @override
  void dispose() {
    _searchController.dispose();
    _debounce?.cancel();
    super.dispose();
  }

  void _fetchData() {
    context.read<CustomerProvider>().fetchCustomers(
          businessId: widget.businessId,
        );
  }

  Future<void> _refresh() async {
    await context.read<CustomerProvider>().fetchCustomers(
          businessId: widget.businessId,
        );
  }

  void _onSearchChanged(String value) {
    _debounce?.cancel();
    _debounce = Timer(const Duration(milliseconds: 500), () {
      final provider = context.read<CustomerProvider>();
      provider.setSearch(value.trim());
      provider.fetchCustomers(businessId: widget.businessId);
    });
  }

  void _clearSearch() {
    _searchController.clear();
    final provider = context.read<CustomerProvider>();
    provider.setSearch('');
    provider.fetchCustomers(businessId: widget.businessId);
  }

  void _goToPage(int page) {
    final provider = context.read<CustomerProvider>();
    provider.setPage(page);
    provider.fetchCustomers(businessId: widget.businessId);
  }

  void _showCreateForm() {
    _showCustomerForm(null);
  }

  void _showEditForm(CustomerInfo customer) {
    _showCustomerForm(customer);
  }

  void _showCustomerForm(CustomerInfo? customer) {
    final isEditing = customer != null;
    final nameController = TextEditingController(text: customer?.name ?? '');
    final emailController = TextEditingController(text: customer?.email ?? '');
    final phoneController = TextEditingController(text: customer?.phone ?? '');
    final dniController = TextEditingController(text: customer?.dni ?? '');
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
                      isEditing ? 'Editar Cliente' : 'Crear Cliente',
                      style: Theme.of(builderContext).textTheme.titleLarge,
                    ),
                    const SizedBox(height: 24),

                    // Name
                    TextFormField(
                      controller: nameController,
                      decoration: const InputDecoration(
                        labelText: 'Nombre *',
                        prefixIcon: Icon(Icons.person),
                      ),
                      textCapitalization: TextCapitalization.words,
                    ),
                    const SizedBox(height: 16),

                    // Email
                    TextFormField(
                      controller: emailController,
                      decoration: const InputDecoration(
                        labelText: 'Email',
                        prefixIcon: Icon(Icons.email_outlined),
                        hintText: 'correo@ejemplo.com',
                      ),
                      keyboardType: TextInputType.emailAddress,
                    ),
                    const SizedBox(height: 16),

                    // Phone
                    TextFormField(
                      controller: phoneController,
                      decoration: const InputDecoration(
                        labelText: 'Telefono',
                        prefixIcon: Icon(Icons.phone_outlined),
                        hintText: '3001234567',
                      ),
                      keyboardType: TextInputType.phone,
                    ),
                    const SizedBox(height: 16),

                    // DNI
                    TextFormField(
                      controller: dniController,
                      decoration: const InputDecoration(
                        labelText: 'Documento de identidad',
                        prefixIcon: Icon(Icons.badge_outlined),
                        hintText: 'CC, NIT, pasaporte...',
                      ),
                    ),
                    const SizedBox(height: 24),

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
                                    final name = nameController.text.trim();
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

                                    setSheetState(() => isSaving = true);

                                    final prov =
                                        context.read<CustomerProvider>();
                                    bool success;

                                    final email =
                                        emailController.text.trim();
                                    final phone =
                                        phoneController.text.trim();
                                    final dni = dniController.text.trim();

                                    if (isEditing) {
                                      success = await prov.updateCustomer(
                                        customer.id,
                                        UpdateCustomerDTO(
                                          name: name,
                                          email: email.isNotEmpty
                                              ? email
                                              : null,
                                          phone: phone.isNotEmpty
                                              ? phone
                                              : null,
                                          dni: dni.isNotEmpty
                                              ? dni
                                              : null,
                                        ),
                                      );
                                    } else {
                                      final created =
                                          await prov.createCustomer(
                                        CreateCustomerDTO(
                                          name: name,
                                          email: email.isNotEmpty
                                              ? email
                                              : null,
                                          phone: phone.isNotEmpty
                                              ? phone
                                              : null,
                                          dni: dni.isNotEmpty
                                              ? dni
                                              : null,
                                        ),
                                      );
                                      success = created != null;
                                    }

                                    setSheetState(() => isSaving = false);

                                    if (!builderContext.mounted) return;

                                    if (success) {
                                      Navigator.pop(builderContext);
                                      _refresh();
                                      if (mounted) {
                                        ScaffoldMessenger.of(context)
                                            .showSnackBar(
                                          SnackBar(
                                            content: Text(isEditing
                                                ? 'Cliente actualizado'
                                                : 'Cliente creado'),
                                          ),
                                        );
                                      }
                                    } else {
                                      ScaffoldMessenger.of(builderContext)
                                          .showSnackBar(
                                        SnackBar(
                                          content: Text(prov.error ??
                                              'Error al guardar'),
                                          backgroundColor: Colors.red,
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
                                : Text(
                                    isEditing ? 'Actualizar' : 'Crear'),
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

  void _showCustomerDetail(CustomerInfo customer) {
    showModalBottomSheet(
      context: context,
      isScrollControlled: true,
      useSafeArea: true,
      shape: const RoundedRectangleBorder(
        borderRadius: BorderRadius.vertical(top: Radius.circular(16)),
      ),
      builder: (sheetContext) {
        return Padding(
          padding: const EdgeInsets.all(16),
          child: Column(
            mainAxisSize: MainAxisSize.min,
            crossAxisAlignment: CrossAxisAlignment.start,
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

              // Header
              Row(
                children: [
                  CircleAvatar(
                    radius: 28,
                    backgroundColor:
                        Theme.of(context).colorScheme.primaryContainer,
                    child: Text(
                      customer.name.isNotEmpty
                          ? customer.name[0].toUpperCase()
                          : '?',
                      style: TextStyle(
                        fontSize: 22,
                        fontWeight: FontWeight.bold,
                        color: Theme.of(context).colorScheme.primary,
                      ),
                    ),
                  ),
                  const SizedBox(width: 16),
                  Expanded(
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        Text(
                          customer.name,
                          style: Theme.of(context)
                              .textTheme
                              .titleLarge
                              ?.copyWith(fontWeight: FontWeight.w600),
                        ),
                        if (customer.email != null &&
                            customer.email!.isNotEmpty)
                          Text(
                            customer.email!,
                            style:
                                Theme.of(context).textTheme.bodyMedium?.copyWith(
                                      color: Colors.grey[600],
                                    ),
                          ),
                      ],
                    ),
                  ),
                ],
              ),
              const SizedBox(height: 24),

              // Detail fields
              _DetailField(
                icon: Icons.phone_outlined,
                label: 'Telefono',
                value: customer.phone.isNotEmpty ? customer.phone : null,
              ),
              _DetailField(
                icon: Icons.badge_outlined,
                label: 'Documento',
                value: customer.dni,
              ),
              _DetailField(
                icon: Icons.calendar_today_outlined,
                label: 'Creado',
                value: _formatDate(customer.createdAt),
              ),
              _DetailField(
                icon: Icons.update,
                label: 'Actualizado',
                value: _formatDate(customer.updatedAt),
              ),
              const SizedBox(height: 16),

              // Actions
              Row(
                children: [
                  Expanded(
                    child: OutlinedButton.icon(
                      onPressed: () {
                        Navigator.pop(sheetContext);
                        _showEditForm(customer);
                      },
                      icon: const Icon(Icons.edit_outlined),
                      label: const Text('Editar'),
                    ),
                  ),
                  const SizedBox(width: 12),
                  Expanded(
                    child: FilledButton.icon(
                      onPressed: () {
                        Navigator.pop(sheetContext);
                      },
                      icon: const Icon(Icons.close),
                      label: const Text('Cerrar'),
                    ),
                  ),
                ],
              ),
              const SizedBox(height: 8),
            ],
          ),
        );
      },
    );
  }

  String? _formatDate(String? dateStr) {
    if (dateStr == null || dateStr.isEmpty) return null;
    try {
      final date = DateTime.parse(dateStr);
      return '${date.day.toString().padLeft(2, '0')}/'
          '${date.month.toString().padLeft(2, '0')}/'
          '${date.year}';
    } catch (_) {
      return dateStr;
    }
  }

  void _confirmDelete(CustomerInfo customer) {
    showDialog(
      context: context,
      builder: (dialogContext) {
        return AlertDialog(
          title: const Text('Eliminar Cliente'),
          content: Text(
              '¿Eliminar a "${customer.name}"? Esta accion no se puede deshacer.'),
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
                final provider = context.read<CustomerProvider>();
                final success = await provider.deleteCustomer(customer.id);
                if (success && mounted) {
                  ScaffoldMessenger.of(context).showSnackBar(
                    const SnackBar(content: Text('Cliente eliminado')),
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

  @override
  Widget build(BuildContext context) {
    final colorScheme = Theme.of(context).colorScheme;

    return Scaffold(
      appBar: AppBar(
        title: const Text('Clientes'),
      ),
      floatingActionButton: FloatingActionButton(
        onPressed: _showCreateForm,
        child: const Icon(Icons.person_add),
      ),
      body: Consumer<CustomerProvider>(
        builder: (context, provider, child) {
          return Column(
            children: [
              // Search bar
              Padding(
                padding: const EdgeInsets.fromLTRB(16, 8, 16, 0),
                child: TextField(
                  controller: _searchController,
                  onChanged: _onSearchChanged,
                  decoration: InputDecoration(
                    hintText: 'Buscar por nombre, email o telefono...',
                    prefixIcon: const Icon(Icons.search),
                    suffixIcon: _searchController.text.isNotEmpty
                        ? IconButton(
                            icon: const Icon(Icons.clear),
                            onPressed: _clearSearch,
                          )
                        : null,
                    isDense: true,
                    contentPadding: const EdgeInsets.symmetric(
                        horizontal: 16, vertical: 10),
                  ),
                ),
              ),
              const SizedBox(height: 8),

              // Content
              Expanded(
                child: _buildContent(provider, colorScheme),
              ),

              // Pagination
              if (provider.pagination != null &&
                  provider.pagination!.lastPage > 1)
                _PaginationBar(
                  pagination: provider.pagination!,
                  currentPage: provider.page,
                  onPageChanged: _goToPage,
                ),
            ],
          );
        },
      ),
    );
  }

  Widget _buildContent(CustomerProvider provider, ColorScheme colorScheme) {
    if (provider.isLoading && provider.customers.isEmpty) {
      return const Center(child: CircularProgressIndicator());
    }

    if (provider.error != null && provider.customers.isEmpty) {
      return Center(
        child: Padding(
          padding: const EdgeInsets.all(24),
          child: Column(
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              Icon(Icons.error_outline, size: 48, color: colorScheme.error),
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

    if (provider.customers.isEmpty) {
      return Center(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Icon(Icons.people_outline, size: 64, color: Colors.grey[400]),
            const SizedBox(height: 16),
            Text(
              'No hay clientes',
              style: Theme.of(context).textTheme.titleMedium?.copyWith(
                    color: Colors.grey[600],
                  ),
            ),
            const SizedBox(height: 8),
            Text(
              _searchController.text.isNotEmpty
                  ? 'No se encontraron resultados'
                  : 'Agrega tu primer cliente',
              style: Theme.of(context).textTheme.bodyMedium?.copyWith(
                    color: Colors.grey[500],
                  ),
            ),
          ],
        ),
      );
    }

    return RefreshIndicator(
      onRefresh: _refresh,
      child: ListView.builder(
        padding: const EdgeInsets.fromLTRB(16, 0, 16, 80),
        itemCount: provider.customers.length,
        itemBuilder: (context, index) {
          final customer = provider.customers[index];
          return _CustomerCard(
            customer: customer,
            onTap: () => _showCustomerDetail(customer),
            onEdit: () => _showEditForm(customer),
            onDelete: () => _confirmDelete(customer),
          );
        },
      ),
    );
  }
}

class _CustomerCard extends StatelessWidget {
  final CustomerInfo customer;
  final VoidCallback onTap;
  final VoidCallback onEdit;
  final VoidCallback onDelete;

  const _CustomerCard({
    required this.customer,
    required this.onTap,
    required this.onEdit,
    required this.onDelete,
  });

  @override
  Widget build(BuildContext context) {
    return Card(
      margin: const EdgeInsets.only(bottom: 8),
      child: InkWell(
        borderRadius: BorderRadius.circular(12),
        onTap: onTap,
        onLongPress: onDelete,
        child: Padding(
          padding: const EdgeInsets.all(12),
          child: Row(
            children: [
              // Avatar
              CircleAvatar(
                radius: 22,
                backgroundColor:
                    Theme.of(context).colorScheme.primaryContainer,
                child: Text(
                  customer.name.isNotEmpty
                      ? customer.name[0].toUpperCase()
                      : '?',
                  style: TextStyle(
                    fontWeight: FontWeight.bold,
                    color: Theme.of(context).colorScheme.primary,
                  ),
                ),
              ),
              const SizedBox(width: 12),

              // Info
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(
                      customer.name,
                      style: Theme.of(context)
                          .textTheme
                          .titleSmall
                          ?.copyWith(fontWeight: FontWeight.w600),
                      maxLines: 1,
                      overflow: TextOverflow.ellipsis,
                    ),
                    const SizedBox(height: 4),
                    Row(
                      children: [
                        if (customer.email != null &&
                            customer.email!.isNotEmpty) ...[
                          Icon(Icons.email_outlined,
                              size: 13, color: Colors.grey[500]),
                          const SizedBox(width: 4),
                          Flexible(
                            child: Text(
                              customer.email!,
                              style: Theme.of(context)
                                  .textTheme
                                  .bodySmall
                                  ?.copyWith(color: Colors.grey[600]),
                              maxLines: 1,
                              overflow: TextOverflow.ellipsis,
                            ),
                          ),
                        ],
                      ],
                    ),
                    const SizedBox(height: 2),
                    Row(
                      children: [
                        if (customer.phone.isNotEmpty) ...[
                          Icon(Icons.phone_outlined,
                              size: 13, color: Colors.grey[500]),
                          const SizedBox(width: 4),
                          Text(
                            customer.phone,
                            style: Theme.of(context)
                                .textTheme
                                .bodySmall
                                ?.copyWith(color: Colors.grey[600]),
                          ),
                        ],
                        if (customer.phone.isNotEmpty &&
                            customer.dni != null &&
                            customer.dni!.isNotEmpty)
                          const SizedBox(width: 12),
                        if (customer.dni != null &&
                            customer.dni!.isNotEmpty) ...[
                          Icon(Icons.badge_outlined,
                              size: 13, color: Colors.grey[500]),
                          const SizedBox(width: 4),
                          Text(
                            customer.dni!,
                            style: Theme.of(context)
                                .textTheme
                                .bodySmall
                                ?.copyWith(color: Colors.grey[600]),
                          ),
                        ],
                      ],
                    ),
                  ],
                ),
              ),

              // Action buttons
              Column(
                children: [
                  InkWell(
                    borderRadius: BorderRadius.circular(8),
                    onTap: onEdit,
                    child: Padding(
                      padding: const EdgeInsets.all(6),
                      child: Icon(Icons.edit_outlined,
                          size: 20, color: Colors.amber[700]),
                    ),
                  ),
                  const SizedBox(height: 4),
                  InkWell(
                    borderRadius: BorderRadius.circular(8),
                    onTap: onDelete,
                    child: Padding(
                      padding: const EdgeInsets.all(6),
                      child: Icon(Icons.delete_outline,
                          size: 20, color: Colors.red[400]),
                    ),
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

class _DetailField extends StatelessWidget {
  final IconData icon;
  final String label;
  final String? value;

  const _DetailField({
    required this.icon,
    required this.label,
    this.value,
  });

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.only(bottom: 12),
      child: Row(
        children: [
          Icon(icon, size: 20, color: Colors.grey[500]),
          const SizedBox(width: 12),
          Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Text(
                label,
                style: Theme.of(context)
                    .textTheme
                    .bodySmall
                    ?.copyWith(color: Colors.grey[500]),
              ),
              const SizedBox(height: 2),
              Text(
                value ?? '--',
                style: Theme.of(context).textTheme.bodyMedium?.copyWith(
                      color: value != null ? null : Colors.grey[400],
                    ),
              ),
            ],
          ),
        ],
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
            '${pagination.total} clientes',
            style: Theme.of(context).textTheme.bodySmall?.copyWith(
                  color: Colors.grey[600],
                ),
          ),
          Row(
            children: [
              IconButton(
                icon: const Icon(Icons.chevron_left),
                onPressed: pagination.hasPrev
                    ? () => onPageChanged(currentPage - 1)
                    : null,
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
                onPressed: pagination.hasNext
                    ? () => onPageChanged(currentPage + 1)
                    : null,
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
