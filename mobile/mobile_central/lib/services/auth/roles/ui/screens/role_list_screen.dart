import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../../domain/entities.dart';
import '../providers/role_provider.dart';

class RoleListScreen extends StatefulWidget {
  const RoleListScreen({super.key});

  @override
  State<RoleListScreen> createState() => _RoleListScreenState();
}

class _RoleListScreenState extends State<RoleListScreen> {
  int _currentPage = 1;
  static const int _pageSize = 20;

  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addPostFrameCallback((_) {
      _loadRoles();
    });
  }

  void _loadRoles() {
    context.read<RoleProvider>().fetchRoles(
          params: GetRolesParams(page: _currentPage, pageSize: _pageSize),
        );
  }

  void _goToPage(int page) {
    setState(() => _currentPage = page);
    context.read<RoleProvider>().fetchRoles(
          params: GetRolesParams(page: page, pageSize: _pageSize),
        );
  }

  Future<void> _showFormDialog({Role? role}) async {
    final saved = await showDialog<bool>(
      context: context,
      builder: (ctx) => _RoleFormDialog(role: role),
    );
    if (saved == true) _loadRoles();
  }

  Future<void> _confirmDelete(Role role) async {
    final confirmed = await showDialog<bool>(
      context: context,
      builder: (ctx) => AlertDialog(
        title: const Text('Eliminar Rol'),
        content: Text(
            '¿Estás seguro de que deseas eliminar el rol "${role.name}"? Esta acción no se puede deshacer.'),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(ctx, false),
            child: const Text('Cancelar'),
          ),
          FilledButton(
            onPressed: () => Navigator.pop(ctx, true),
            style: FilledButton.styleFrom(backgroundColor: Colors.red),
            child: const Text('Eliminar'),
          ),
        ],
      ),
    );

    if (confirmed == true && mounted) {
      final success = await context.read<RoleProvider>().deleteRole(role.id);
      if (success && mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(content: Text('Rol eliminado correctamente')),
        );
        _loadRoles();
      } else if (mounted) {
        final error = context.read<RoleProvider>().error;
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text(error ?? 'Error al eliminar el rol')),
        );
      }
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Roles'),
      ),
      floatingActionButton: FloatingActionButton(
        onPressed: () => _showFormDialog(),
        child: const Icon(Icons.add),
      ),
      body: Consumer<RoleProvider>(
        builder: (context, provider, child) {
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
                    Text(
                      provider.error!,
                      textAlign: TextAlign.center,
                      style: const TextStyle(color: Colors.red),
                    ),
                    const SizedBox(height: 16),
                    FilledButton.icon(
                      onPressed: _loadRoles,
                      icon: const Icon(Icons.refresh),
                      label: const Text('Reintentar'),
                    ),
                  ],
                ),
              ),
            );
          }

          if (provider.roles.isEmpty) {
            return Center(
              child: Column(
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  Icon(Icons.admin_panel_settings_outlined,
                      size: 64, color: Colors.grey.shade400),
                  const SizedBox(height: 16),
                  Text('No hay roles disponibles',
                      style: TextStyle(
                          fontSize: 16, color: Colors.grey.shade600)),
                ],
              ),
            );
          }

          return RefreshIndicator(
            onRefresh: () async => _loadRoles(),
            child: Column(
              children: [
                Expanded(
                  child: ListView.builder(
                    padding: const EdgeInsets.all(16),
                    itemCount: provider.roles.length,
                    itemBuilder: (context, index) {
                      final role = provider.roles[index];
                      return Card(
                        margin: const EdgeInsets.only(bottom: 8),
                        child: ListTile(
                          leading: CircleAvatar(
                            backgroundColor: role.isSystem
                                ? Colors.blue.shade100
                                : Colors.grey.shade200,
                            child: Icon(
                              role.isSystem
                                  ? Icons.shield
                                  : Icons.person_outline,
                              color: role.isSystem
                                  ? Colors.blue.shade700
                                  : Colors.grey.shade600,
                            ),
                          ),
                          title: Text(role.name,
                              style:
                                  const TextStyle(fontWeight: FontWeight.w600)),
                          subtitle: Column(
                            crossAxisAlignment: CrossAxisAlignment.start,
                            children: [
                              if (role.description != null &&
                                  role.description!.isNotEmpty)
                                Text(role.description!,
                                    maxLines: 1,
                                    overflow: TextOverflow.ellipsis),
                              const SizedBox(height: 4),
                              Row(
                                children: [
                                  if (role.level != null)
                                    _InfoChip(
                                        label: 'Nivel ${role.level}',
                                        color: Colors.indigo),
                                  if (role.level != null)
                                    const SizedBox(width: 4),
                                  _InfoChip(
                                    label:
                                        role.isSystem ? 'Sistema' : 'Custom',
                                    color: role.isSystem
                                        ? Colors.blue
                                        : Colors.grey,
                                  ),
                                ],
                              ),
                            ],
                          ),
                          isThreeLine: true,
                          trailing: PopupMenuButton<String>(
                            onSelected: (value) {
                              if (value == 'edit') {
                                _showFormDialog(role: role);
                              } else if (value == 'delete') {
                                _confirmDelete(role);
                              }
                            },
                            itemBuilder: (ctx) => [
                              const PopupMenuItem(
                                value: 'edit',
                                child: ListTile(
                                  leading: Icon(Icons.edit,
                                      color: Colors.orange),
                                  title: Text('Editar'),
                                  contentPadding: EdgeInsets.zero,
                                  dense: true,
                                ),
                              ),
                              const PopupMenuItem(
                                value: 'delete',
                                child: ListTile(
                                  leading:
                                      Icon(Icons.delete, color: Colors.red),
                                  title: Text('Eliminar'),
                                  contentPadding: EdgeInsets.zero,
                                  dense: true,
                                ),
                              ),
                            ],
                          ),
                        ),
                      );
                    },
                  ),
                ),
                // Pagination
                if (provider.pagination != null &&
                    provider.pagination!.lastPage > 1)
                  _PaginationBar(
                    currentPage: provider.pagination!.currentPage,
                    totalPages: provider.pagination!.lastPage,
                    total: provider.pagination!.total,
                    onPageChanged: _goToPage,
                  ),
              ],
            ),
          );
        },
      ),
    );
  }
}

// --------------------------------------------------
// Form Dialog
// --------------------------------------------------
class _RoleFormDialog extends StatefulWidget {
  final Role? role;
  const _RoleFormDialog({this.role});

  @override
  State<_RoleFormDialog> createState() => _RoleFormDialogState();
}

class _RoleFormDialogState extends State<_RoleFormDialog> {
  final _formKey = GlobalKey<FormState>();
  late TextEditingController _nameCtrl;
  late TextEditingController _descriptionCtrl;
  late TextEditingController _levelCtrl;
  int? _scopeId;
  bool _saving = false;
  String? _error;

  bool get _isEditing => widget.role != null;

  @override
  void initState() {
    super.initState();
    _nameCtrl = TextEditingController(text: widget.role?.name ?? '');
    _descriptionCtrl =
        TextEditingController(text: widget.role?.description ?? '');
    _levelCtrl =
        TextEditingController(text: widget.role?.level?.toString() ?? '');
    _scopeId = widget.role?.scopeId;
  }

  @override
  void dispose() {
    _nameCtrl.dispose();
    _descriptionCtrl.dispose();
    _levelCtrl.dispose();
    super.dispose();
  }

  Future<void> _save() async {
    if (!_formKey.currentState!.validate()) return;

    setState(() {
      _saving = true;
      _error = null;
    });

    final provider = context.read<RoleProvider>();
    bool success;

    final level =
        _levelCtrl.text.isNotEmpty ? int.tryParse(_levelCtrl.text) : null;

    if (_isEditing) {
      success = await provider.updateRole(
        widget.role!.id,
        UpdateRoleDTO(
          name: _nameCtrl.text.trim(),
          description: _descriptionCtrl.text.trim().isNotEmpty
              ? _descriptionCtrl.text.trim()
              : null,
          level: level,
          scopeId: _scopeId,
        ),
      );
    } else {
      final result = await provider.createRole(
        CreateRoleDTO(
          name: _nameCtrl.text.trim(),
          description: _descriptionCtrl.text.trim().isNotEmpty
              ? _descriptionCtrl.text.trim()
              : null,
          level: level,
          scopeId: _scopeId,
        ),
      );
      success = result != null;
    }

    if (!mounted) return;

    if (success) {
      Navigator.pop(context, true);
    } else {
      setState(() {
        _error = provider.error ?? 'Error al guardar el rol';
        _saving = false;
      });
    }
  }

  @override
  Widget build(BuildContext context) {
    return AlertDialog(
      title: Text(_isEditing ? 'Editar Rol' : 'Crear Rol'),
      content: SingleChildScrollView(
        child: Form(
          key: _formKey,
          child: Column(
            mainAxisSize: MainAxisSize.min,
            children: [
              if (_error != null) ...[
                Container(
                  padding: const EdgeInsets.all(8),
                  decoration: BoxDecoration(
                    color: Colors.red.shade50,
                    borderRadius: BorderRadius.circular(8),
                  ),
                  child: Row(
                    children: [
                      Icon(Icons.error, color: Colors.red.shade700, size: 20),
                      const SizedBox(width: 8),
                      Expanded(
                        child: Text(_error!,
                            style: TextStyle(
                                color: Colors.red.shade700, fontSize: 13)),
                      ),
                    ],
                  ),
                ),
                const SizedBox(height: 12),
              ],
              TextFormField(
                controller: _nameCtrl,
                decoration: const InputDecoration(
                  labelText: 'Nombre *',
                  border: OutlineInputBorder(),
                ),
                validator: (v) =>
                    v == null || v.trim().isEmpty ? 'El nombre es requerido' : null,
              ),
              const SizedBox(height: 12),
              TextFormField(
                controller: _descriptionCtrl,
                decoration: const InputDecoration(
                  labelText: 'Descripción',
                  border: OutlineInputBorder(),
                ),
                maxLines: 2,
              ),
              const SizedBox(height: 12),
              TextFormField(
                controller: _levelCtrl,
                decoration: const InputDecoration(
                  labelText: 'Nivel (1-10)',
                  border: OutlineInputBorder(),
                ),
                keyboardType: TextInputType.number,
                validator: (v) {
                  if (v == null || v.isEmpty) return null;
                  final n = int.tryParse(v);
                  if (n == null || n < 1 || n > 10) {
                    return 'Ingresa un número entre 1 y 10';
                  }
                  return null;
                },
              ),
              const SizedBox(height: 12),
              DropdownButtonFormField<int>(
                initialValue: _scopeId,
                decoration: const InputDecoration(
                  labelText: 'Scope *',
                  border: OutlineInputBorder(),
                ),
                items: const [
                  DropdownMenuItem(value: 1, child: Text('Platform')),
                  DropdownMenuItem(value: 2, child: Text('Business')),
                ],
                onChanged: (v) => setState(() => _scopeId = v),
                validator: (v) => v == null ? 'Selecciona un scope' : null,
              ),
            ],
          ),
        ),
      ),
      actions: [
        TextButton(
          onPressed: _saving ? null : () => Navigator.pop(context),
          child: const Text('Cancelar'),
        ),
        FilledButton(
          onPressed: _saving ? null : _save,
          child: _saving
              ? const SizedBox(
                  width: 20,
                  height: 20,
                  child: CircularProgressIndicator(strokeWidth: 2),
                )
              : Text(_isEditing ? 'Actualizar' : 'Crear'),
        ),
      ],
    );
  }
}

// --------------------------------------------------
// Shared Widgets
// --------------------------------------------------
class _InfoChip extends StatelessWidget {
  final String label;
  final Color color;
  const _InfoChip({required this.label, required this.color});

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 2),
      decoration: BoxDecoration(
        color: color.withValues(alpha: 0.1),
        borderRadius: BorderRadius.circular(12),
      ),
      child: Text(
        label,
        style: TextStyle(fontSize: 11, color: color, fontWeight: FontWeight.w500),
      ),
    );
  }
}

class _PaginationBar extends StatelessWidget {
  final int currentPage;
  final int totalPages;
  final int total;
  final ValueChanged<int> onPageChanged;

  const _PaginationBar({
    required this.currentPage,
    required this.totalPages,
    required this.total,
    required this.onPageChanged,
  });

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 12),
      decoration: BoxDecoration(
        color: Theme.of(context).colorScheme.surface,
        border: Border(
            top: BorderSide(color: Colors.grey.shade300)),
      ),
      child: Row(
        mainAxisAlignment: MainAxisAlignment.spaceBetween,
        children: [
          Text('$total resultados',
              style: TextStyle(fontSize: 13, color: Colors.grey.shade600)),
          Row(
            children: [
              IconButton(
                icon: const Icon(Icons.chevron_left),
                onPressed:
                    currentPage > 1 ? () => onPageChanged(currentPage - 1) : null,
                iconSize: 20,
                visualDensity: VisualDensity.compact,
              ),
              Text('$currentPage / $totalPages',
                  style: const TextStyle(fontSize: 13)),
              IconButton(
                icon: const Icon(Icons.chevron_right),
                onPressed: currentPage < totalPages
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
