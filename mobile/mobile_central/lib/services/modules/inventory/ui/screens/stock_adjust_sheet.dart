import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:provider/provider.dart';
import '../../../../../shared/theme/app_colors.dart';
import '../../../../../shared/theme/app_tokens.dart';
import '../../../../../shared/utils/formatters.dart';
import '../../domain/entities.dart';
import '../providers/inventory_provider.dart';

Future<bool?> showStockAdjustSheet(
  BuildContext context, {
  required InventoryLevel level,
  int? businessId,
}) {
  return showModalBottomSheet<bool>(
    context: context,
    isScrollControlled: true,
    backgroundColor: AppColors.surface,
    builder: (sheetContext) => Padding(
      padding: EdgeInsets.only(bottom: MediaQuery.of(sheetContext).viewInsets.bottom),
      child: _StockAdjustForm(level: level, businessId: businessId),
    ),
  );
}

class _StockAdjustForm extends StatefulWidget {
  const _StockAdjustForm({required this.level, this.businessId});

  final InventoryLevel level;
  final int? businessId;

  @override
  State<_StockAdjustForm> createState() => _StockAdjustFormState();
}

class _StockAdjustFormState extends State<_StockAdjustForm> {
  final _quantityController = TextEditingController(text: '1');
  final _notesController = TextEditingController();
  bool _isEntry = true;
  bool _busy = false;
  String? _error;
  String _reason = 'Ajuste por conteo';

  static const List<String> _reasons = [
    'Ajuste por conteo',
    'Compra a proveedor',
    'Devolucion de cliente',
    'Producto averiado',
    'Perdida o robo',
  ];

  int get _delta {
    final value = int.tryParse(_quantityController.text) ?? 0;
    return _isEntry ? value : -value;
  }

  @override
  void dispose() {
    _quantityController.dispose();
    _notesController.dispose();
    super.dispose();
  }

  Future<void> _submit() async {
    final value = int.tryParse(_quantityController.text) ?? 0;
    if (value <= 0) {
      setState(() => _error = 'La cantidad debe ser mayor a cero');
      return;
    }
    if (!_isEntry && value > widget.level.quantity) {
      setState(() => _error = 'No puedes descontar mas de las ${widget.level.quantity} unidades en bodega');
      return;
    }

    setState(() {
      _busy = true;
      _error = null;
    });

    final provider = context.read<InventoryProvider>();
    final result = await provider.adjustStock(
      AdjustStockDTO(
        productId: widget.level.productId,
        warehouseId: widget.level.warehouseId,
        quantity: _delta,
        reason: _reason,
        notes: _notesController.text.trim().isEmpty ? null : _notesController.text.trim(),
      ),
      businessId: widget.businessId,
    );

    if (!mounted) return;
    setState(() => _busy = false);

    if (result == null) {
      setState(() => _error = provider.error ?? 'No se pudo registrar el ajuste');
      return;
    }

    Navigator.pop(context, true);
    ScaffoldMessenger.of(context).showSnackBar(
      const SnackBar(content: Text('Ajuste registrado')),
    );
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final level = widget.level;
    final resulting = level.quantity + _delta;

    return Padding(
      padding: const EdgeInsets.fromLTRB(20, 4, 20, 24),
      child: Column(
        mainAxisSize: MainAxisSize.min,
        crossAxisAlignment: CrossAxisAlignment.stretch,
        children: [
          Text('Ajustar stock', style: theme.textTheme.titleLarge),
          const SizedBox(height: 4),
          Text(
            '${level.productName ?? ''}  ·  ${level.warehouseName ?? ''}',
            style: theme.textTheme.bodySmall,
          ),
          const SizedBox(height: 18),
          Row(
            children: [
              Expanded(
                child: _ModeButton(
                  label: 'Entrada',
                  icon: Icons.south_west_rounded,
                  active: _isEntry,
                  color: AppColors.success,
                  onTap: () => setState(() => _isEntry = true),
                ),
              ),
              const SizedBox(width: 10),
              Expanded(
                child: _ModeButton(
                  label: 'Salida',
                  icon: Icons.north_east_rounded,
                  active: !_isEntry,
                  color: AppColors.error,
                  onTap: () => setState(() => _isEntry = false),
                ),
              ),
            ],
          ),
          const SizedBox(height: 16),
          TextField(
            controller: _quantityController,
            keyboardType: TextInputType.number,
            inputFormatters: [FilteringTextInputFormatter.digitsOnly],
            onChanged: (_) => setState(() {}),
            decoration: const InputDecoration(
              labelText: 'Cantidad',
              prefixIcon: Icon(Icons.numbers_rounded, size: 20),
            ),
          ),
          const SizedBox(height: 12),
          DropdownButtonFormField<String>(
            initialValue: _reason,
            decoration: const InputDecoration(labelText: 'Motivo'),
            items: _reasons
                .map((reason) => DropdownMenuItem(value: reason, child: Text(reason)))
                .toList(),
            onChanged: (value) => setState(() => _reason = value ?? _reason),
          ),
          const SizedBox(height: 12),
          TextField(
            controller: _notesController,
            maxLines: 2,
            decoration: const InputDecoration(labelText: 'Notas (opcional)'),
          ),
          const SizedBox(height: 16),
          Container(
            padding: const EdgeInsets.symmetric(horizontal: 14, vertical: 12),
            decoration: BoxDecoration(
              color: AppColors.surfaceMuted,
              borderRadius: AppRadius.mdAll,
            ),
            child: Row(
              children: [
                Expanded(
                  child: Text('Stock resultante', style: theme.textTheme.titleSmall),
                ),
                Text(
                  '${AppFormat.number(level.quantity)} a ${AppFormat.number(resulting)}',
                  style: theme.textTheme.titleSmall?.copyWith(
                    color: resulting < 0 ? AppColors.error : AppColors.primary,
                  ),
                ),
              ],
            ),
          ),
          if (_error != null) ...[
            const SizedBox(height: 12),
            Container(
              padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 10),
              decoration: BoxDecoration(
                color: AppColors.errorSoft,
                borderRadius: AppRadius.mdAll,
              ),
              child: Text(
                _error!,
                style: theme.textTheme.bodySmall?.copyWith(
                  color: const Color(0xFFB91C1C),
                  fontWeight: FontWeight.w500,
                ),
              ),
            ),
          ],
          const SizedBox(height: 18),
          FilledButton(
            onPressed: _busy ? null : _submit,
            style: FilledButton.styleFrom(minimumSize: const Size(0, 50)),
            child: _busy
                ? const SizedBox(
                    height: 19,
                    width: 19,
                    child: CircularProgressIndicator(strokeWidth: 2.2, color: Colors.white),
                  )
                : const Text('Registrar ajuste'),
          ),
        ],
      ),
    );
  }
}

class _ModeButton extends StatelessWidget {
  const _ModeButton({
    required this.label,
    required this.icon,
    required this.active,
    required this.color,
    required this.onTap,
  });

  final String label;
  final IconData icon;
  final bool active;
  final Color color;
  final VoidCallback onTap;

  @override
  Widget build(BuildContext context) {
    return Material(
      color: active ? color.withValues(alpha: 0.10) : AppColors.surface,
      borderRadius: AppRadius.mdAll,
      child: InkWell(
        borderRadius: AppRadius.mdAll,
        onTap: onTap,
        child: Container(
          padding: const EdgeInsets.symmetric(vertical: 13),
          decoration: BoxDecoration(
            borderRadius: AppRadius.mdAll,
            border: Border.all(color: active ? color : AppColors.border),
          ),
          child: Row(
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              Icon(icon, size: 17, color: active ? color : AppColors.textMuted),
              const SizedBox(width: 8),
              Text(
                label,
                style: TextStyle(
                  fontFamily: 'Inter',
                  fontSize: 13.5,
                  fontWeight: FontWeight.w600,
                  color: active ? color : AppColors.textSecondary,
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }
}
