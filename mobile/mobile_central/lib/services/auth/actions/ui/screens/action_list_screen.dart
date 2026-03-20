import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../../domain/entities.dart';
import '../providers/action_provider.dart';

class ActionListScreen extends StatefulWidget {
  const ActionListScreen({super.key});

  @override
  State<ActionListScreen> createState() => _ActionListScreenState();
}

class _ActionListScreenState extends State<ActionListScreen> {
  int _currentPage = 1;
  static const int _pageSize = 20;

  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addPostFrameCallback((_) {
      _loadActions();
    });
  }

  void _loadActions() {
    context.read<ActionProvider>().fetchActions(
          params: GetActionsParams(page: _currentPage, pageSize: _pageSize),
        );
  }

  void _goToPage(int page) {
    setState(() => _currentPage = page);
    context.read<ActionProvider>().fetchActions(
          params: GetActionsParams(page: page, pageSize: _pageSize),
        );
  }

  Future<void> _showFormDialog({ActionEntity? action}) async {
    final saved = await showDialog<bool>(
      context: context,
      builder: (ctx) => _ActionFormDialog(action: action),
    );
    if (saved == true) _loadActions();
  }

  Future<void> _confirmDelete(ActionEntity action) async {
    final confirmed = await showDialog<bool>(
      context: context,
      builder: (ctx) => AlertDialog(
        title: const Text('Eliminar Acción'),
        content: Text(
            '¿Estás seguro de que deseas eliminar la acción "${action.name}"? Esta acción no se puede deshacer.'),
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
          await context.read<ActionProvider>().deleteAction(action.id);
      if (success && mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(content: Text('Acción eliminada correctamente')),
        );
        _loadActions();
      } else if (mounted) {
        final error = context.read<ActionProvider>().error;
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text(error ?? 'Error al eliminar la acción')),
        );
      }
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Acciones'),
      ),
      floatingActionButton: FloatingActionButton(
        onPressed: () => _showFormDialog(),
        child: const Icon(Icons.add),
      ),
      body: Consumer<ActionProvider>(
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
                      onPressed: _loadActions,
                      icon: const Icon(Icons.refresh),
                      label: const Text('Reintentar'),
                    ),
                  ],
                ),
              ),
            );
          }

          if (provider.actions.isEmpty) {
            return Center(
              child: Column(
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  Icon(Icons.touch_app_outlined,
                      size: 64, color: Colors.grey.shade400),
                  const SizedBox(height: 16),
                  Text('No hay acciones disponibles',
                      style: TextStyle(
                          fontSize: 16, color: Colors.grey.shade600)),
                ],
              ),
            );
          }

          return RefreshIndicator(
            onRefresh: () async => _loadActions(),
            child: Column(
              children: [
                Expanded(
                  child: ListView.builder(
                    padding: const EdgeInsets.all(16),
                    itemCount: provider.actions.length,
                    itemBuilder: (context, index) {
                      final action = provider.actions[index];
                      return Card(
                        margin: const EdgeInsets.only(bottom: 8),
                        child: ListTile(
                          leading: CircleAvatar(
                            backgroundColor: Colors.deepPurple.shade100,
                            child: Icon(Icons.touch_app,
                                color: Colors.deepPurple.shade700, size: 20),
                          ),
                          title: Text(action.name,
                              style: const TextStyle(
                                  fontWeight: FontWeight.w600)),
                          subtitle: action.description != null &&
                                  action.description!.isNotEmpty
                              ? Text(action.description!,
                                  maxLines: 2,
                                  overflow: TextOverflow.ellipsis)
                              : null,
                          trailing: PopupMenuButton<String>(
                            onSelected: (value) {
                              if (value == 'edit') {
                                _showFormDialog(action: action);
                              } else if (value == 'delete') {
                                _confirmDelete(action);
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
class _ActionFormDialog extends StatefulWidget {
  final ActionEntity? action;
  const _ActionFormDialog({this.action});

  @override
  State<_ActionFormDialog> createState() => _ActionFormDialogState();
}

class _ActionFormDialogState extends State<_ActionFormDialog> {
  final _formKey = GlobalKey<FormState>();
  late TextEditingController _nameCtrl;
  late TextEditingController _descriptionCtrl;
  bool _saving = false;
  String? _error;

  bool get _isEditing => widget.action != null;

  @override
  void initState() {
    super.initState();
    _nameCtrl = TextEditingController(text: widget.action?.name ?? '');
    _descriptionCtrl =
        TextEditingController(text: widget.action?.description ?? '');
  }

  @override
  void dispose() {
    _nameCtrl.dispose();
    _descriptionCtrl.dispose();
    super.dispose();
  }

  Future<void> _save() async {
    if (!_formKey.currentState!.validate()) return;

    setState(() {
      _saving = true;
      _error = null;
    });

    final provider = context.read<ActionProvider>();
    bool success;

    if (_isEditing) {
      success = await provider.updateAction(
        widget.action!.id,
        UpdateActionDTO(
          name: _nameCtrl.text.trim(),
          description: _descriptionCtrl.text.trim().isNotEmpty
              ? _descriptionCtrl.text.trim()
              : null,
        ),
      );
    } else {
      final result = await provider.createAction(
        CreateActionDTO(
          name: _nameCtrl.text.trim(),
          description: _descriptionCtrl.text.trim().isNotEmpty
              ? _descriptionCtrl.text.trim()
              : null,
        ),
      );
      success = result != null;
    }

    if (!mounted) return;

    if (success) {
      Navigator.pop(context, true);
    } else {
      setState(() {
        _error = provider.error ?? 'Error al guardar la acción';
        _saving = false;
      });
    }
  }

  @override
  Widget build(BuildContext context) {
    return AlertDialog(
      title: Text(_isEditing ? 'Editar Acción' : 'Crear Acción'),
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
                  hintText: 'ej: read, create, update, delete',
                ),
                validator: (v) => v == null || v.trim().isEmpty
                    ? 'El nombre es requerido'
                    : null,
              ),
              const SizedBox(height: 12),
              TextFormField(
                controller: _descriptionCtrl,
                decoration: const InputDecoration(
                  labelText: 'Descripción',
                  border: OutlineInputBorder(),
                  hintText: 'Descripción de la acción',
                ),
                maxLines: 3,
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
