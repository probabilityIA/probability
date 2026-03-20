import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../../domain/entities.dart';
import '../providers/permission_provider.dart';

class PermissionListScreen extends StatefulWidget {
  const PermissionListScreen({super.key});

  @override
  State<PermissionListScreen> createState() => _PermissionListScreenState();
}

class _PermissionListScreenState extends State<PermissionListScreen> {
  int _currentPage = 1;
  static const int _pageSize = 20;

  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addPostFrameCallback((_) {
      _loadPermissions();
    });
  }

  void _loadPermissions() {
    context.read<PermissionProvider>().fetchPermissions(
          params:
              GetPermissionsParams(page: _currentPage, pageSize: _pageSize),
        );
  }

  void _goToPage(int page) {
    setState(() => _currentPage = page);
    context.read<PermissionProvider>().fetchPermissions(
          params: GetPermissionsParams(page: page, pageSize: _pageSize),
        );
  }

  Future<void> _showFormDialog({Permission? permission}) async {
    final saved = await showDialog<bool>(
      context: context,
      builder: (ctx) => _PermissionFormDialog(permission: permission),
    );
    if (saved == true) _loadPermissions();
  }

  Future<void> _confirmDelete(Permission permission) async {
    final label = permission.resource != null && permission.action != null
        ? '${permission.resource}:${permission.action}'
        : 'ID ${permission.id}';

    final confirmed = await showDialog<bool>(
      context: context,
      builder: (ctx) => AlertDialog(
        title: const Text('Eliminar Permiso'),
        content: Text(
            '¿Estás seguro de que deseas eliminar el permiso "$label"? Esta acción no se puede deshacer.'),
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
      final success =
          await context.read<PermissionProvider>().deletePermission(permission.id);
      if (success && mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(content: Text('Permiso eliminado correctamente')),
        );
        _loadPermissions();
      } else if (mounted) {
        final error = context.read<PermissionProvider>().error;
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text(error ?? 'Error al eliminar el permiso')),
        );
      }
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Permisos'),
      ),
      floatingActionButton: FloatingActionButton(
        onPressed: () => _showFormDialog(),
        child: const Icon(Icons.add),
      ),
      body: Consumer<PermissionProvider>(
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
                      onPressed: _loadPermissions,
                      icon: const Icon(Icons.refresh),
                      label: const Text('Reintentar'),
                    ),
                  ],
                ),
              ),
            );
          }

          if (provider.permissions.isEmpty) {
            return Center(
              child: Column(
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  Icon(Icons.lock_outline,
                      size: 64, color: Colors.grey.shade400),
                  const SizedBox(height: 16),
                  Text('No hay permisos disponibles',
                      style: TextStyle(
                          fontSize: 16, color: Colors.grey.shade600)),
                ],
              ),
            );
          }

          return RefreshIndicator(
            onRefresh: () async => _loadPermissions(),
            child: Column(
              children: [
                Expanded(
                  child: ListView.builder(
                    padding: const EdgeInsets.all(16),
                    itemCount: provider.permissions.length,
                    itemBuilder: (context, index) {
                      final perm = provider.permissions[index];
                      return Card(
                        margin: const EdgeInsets.only(bottom: 8),
                        child: ListTile(
                          leading: CircleAvatar(
                            backgroundColor: Colors.purple.shade100,
                            child: Icon(Icons.vpn_key,
                                color: Colors.purple.shade700, size: 20),
                          ),
                          title: Text(
                            perm.resource != null && perm.action != null
                                ? '${perm.resource}:${perm.action}'
                                : 'Permiso ${perm.id}',
                            style:
                                const TextStyle(fontWeight: FontWeight.w600),
                          ),
                          subtitle: Column(
                            crossAxisAlignment: CrossAxisAlignment.start,
                            children: [
                              if (perm.description != null &&
                                  perm.description!.isNotEmpty)
                                Text(perm.description!,
                                    maxLines: 1,
                                    overflow: TextOverflow.ellipsis),
                              const SizedBox(height: 4),
                              Row(
                                children: [
                                  if (perm.resource != null)
                                    _InfoChip(
                                        label: perm.resource!,
                                        color: Colors.blue),
                                  if (perm.resource != null)
                                    const SizedBox(width: 4),
                                  if (perm.action != null)
                                    _InfoChip(
                                        label: perm.action!,
                                        color: Colors.teal),
                                  if (perm.businessTypeName != null) ...[
                                    const SizedBox(width: 4),
                                    _InfoChip(
                                        label: perm.businessTypeName!,
                                        color: Colors.orange),
                                  ],
                                ],
                              ),
                            ],
                          ),
                          isThreeLine: true,
                          trailing: PopupMenuButton<String>(
                            onSelected: (value) {
                              if (value == 'edit') {
                                _showFormDialog(permission: perm);
                              } else if (value == 'delete') {
                                _confirmDelete(perm);
                              }
                            },
                            itemBuilder: (ctx) => [
                              const PopupMenuItem(
                                value: 'edit',
                                child: ListTile(
                                  leading:
                                      Icon(Icons.edit, color: Colors.orange),
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
class _PermissionFormDialog extends StatefulWidget {
  final Permission? permission;
  const _PermissionFormDialog({this.permission});

  @override
  State<_PermissionFormDialog> createState() => _PermissionFormDialogState();
}

class _PermissionFormDialogState extends State<_PermissionFormDialog> {
  final _formKey = GlobalKey<FormState>();
  late TextEditingController _resourceCtrl;
  late TextEditingController _actionCtrl;
  late TextEditingController _descriptionCtrl;
  int? _scopeId;
  bool _saving = false;
  String? _error;

  bool get _isEditing => widget.permission != null;

  @override
  void initState() {
    super.initState();
    _resourceCtrl =
        TextEditingController(text: widget.permission?.resource ?? '');
    _actionCtrl =
        TextEditingController(text: widget.permission?.action ?? '');
    _descriptionCtrl =
        TextEditingController(text: widget.permission?.description ?? '');
    _scopeId = widget.permission?.scopeId;
  }

  @override
  void dispose() {
    _resourceCtrl.dispose();
    _actionCtrl.dispose();
    _descriptionCtrl.dispose();
    super.dispose();
  }

  Future<void> _save() async {
    if (!_formKey.currentState!.validate()) return;

    setState(() {
      _saving = true;
      _error = null;
    });

    final provider = context.read<PermissionProvider>();
    bool success;

    if (_isEditing) {
      success = await provider.updatePermission(
        widget.permission!.id,
        UpdatePermissionDTO(
          resource: _resourceCtrl.text.trim(),
          action: _actionCtrl.text.trim(),
          description: _descriptionCtrl.text.trim().isNotEmpty
              ? _descriptionCtrl.text.trim()
              : null,
          scopeId: _scopeId,
        ),
      );
    } else {
      final result = await provider.createPermission(
        CreatePermissionDTO(
          resource: _resourceCtrl.text.trim(),
          action: _actionCtrl.text.trim(),
          description: _descriptionCtrl.text.trim().isNotEmpty
              ? _descriptionCtrl.text.trim()
              : null,
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
        _error = provider.error ?? 'Error al guardar el permiso';
        _saving = false;
      });
    }
  }

  @override
  Widget build(BuildContext context) {
    return AlertDialog(
      title: Text(_isEditing ? 'Editar Permiso' : 'Crear Permiso'),
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
                controller: _resourceCtrl,
                decoration: const InputDecoration(
                  labelText: 'Recurso *',
                  border: OutlineInputBorder(),
                  hintText: 'ej: orders, users, products',
                ),
                validator: (v) =>
                    v == null || v.trim().isEmpty ? 'El recurso es requerido' : null,
              ),
              const SizedBox(height: 12),
              TextFormField(
                controller: _actionCtrl,
                decoration: const InputDecoration(
                  labelText: 'Acción *',
                  border: OutlineInputBorder(),
                  hintText: 'ej: read, create, update, delete',
                ),
                validator: (v) =>
                    v == null || v.trim().isEmpty ? 'La acción es requerida' : null,
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
              DropdownButtonFormField<int>(
                initialValue: _scopeId,
                decoration: const InputDecoration(
                  labelText: 'Scope',
                  border: OutlineInputBorder(),
                ),
                items: const [
                  DropdownMenuItem(value: 1, child: Text('Platform')),
                  DropdownMenuItem(value: 2, child: Text('Business')),
                ],
                onChanged: (v) => setState(() => _scopeId = v),
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
        border: Border(top: BorderSide(color: Colors.grey.shade300)),
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
                onPressed: currentPage > 1
                    ? () => onPageChanged(currentPage - 1)
                    : null,
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
