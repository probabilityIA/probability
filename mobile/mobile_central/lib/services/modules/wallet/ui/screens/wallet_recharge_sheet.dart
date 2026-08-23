import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:provider/provider.dart';
import '../../../../../shared/theme/app_colors.dart';
import '../../../../../shared/theme/app_tokens.dart';
import '../../../../../shared/utils/formatters.dart';
import '../providers/wallet_provider.dart';

Future<bool?> showWalletRechargeSheet(BuildContext context, {int? businessId}) {
  return showModalBottomSheet<bool>(
    context: context,
    isScrollControlled: true,
    backgroundColor: AppColors.surface,
    builder: (sheetContext) => Padding(
      padding: EdgeInsets.only(bottom: MediaQuery.of(sheetContext).viewInsets.bottom),
      child: _RechargeForm(businessId: businessId),
    ),
  );
}

class _RechargeForm extends StatefulWidget {
  const _RechargeForm({this.businessId});

  final int? businessId;

  @override
  State<_RechargeForm> createState() => _RechargeFormState();
}

class _RechargeFormState extends State<_RechargeForm> {
  final _controller = TextEditingController(text: '200000');
  bool _busy = false;
  String? _error;

  static const List<int> _presets = [100000, 200000, 500000, 1000000];

  @override
  void dispose() {
    _controller.dispose();
    super.dispose();
  }

  Future<void> _submit() async {
    final amount = double.tryParse(_controller.text) ?? 0;
    if (amount < 10000) {
      setState(() => _error = 'El monto minimo de recarga es 10.000');
      return;
    }

    setState(() {
      _busy = true;
      _error = null;
    });

    final provider = context.read<WalletProvider>();
    final ok = await provider.rechargeWallet(amount: amount, businessId: widget.businessId);

    if (!mounted) return;
    setState(() => _busy = false);

    if (!ok) {
      setState(() => _error = provider.error ?? 'No se pudo iniciar la recarga');
      return;
    }

    Navigator.pop(context, true);
    ScaffoldMessenger.of(context).showSnackBar(
      const SnackBar(content: Text('Recarga iniciada, completa el pago en la pasarela')),
    );
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    return Padding(
      padding: const EdgeInsets.fromLTRB(20, 4, 20, 24),
      child: Column(
        mainAxisSize: MainAxisSize.min,
        crossAxisAlignment: CrossAxisAlignment.stretch,
        children: [
          Text('Recargar billetera', style: theme.textTheme.titleLarge),
          const SizedBox(height: 4),
          Text(
            'El saldo se usa para pagar las guias que generes.',
            style: theme.textTheme.bodySmall,
          ),
          const SizedBox(height: 18),
          Wrap(
            spacing: 8,
            runSpacing: 8,
            children: _presets
                .map((preset) => ChoiceChip(
                      label: Text(AppFormat.money(preset)),
                      selected: _controller.text == preset.toString(),
                      onSelected: (_) => setState(() => _controller.text = preset.toString()),
                      labelStyle: TextStyle(
                        fontFamily: 'Inter',
                        fontSize: 12.5,
                        fontWeight: FontWeight.w600,
                        color: _controller.text == preset.toString()
                            ? AppColors.primaryDark
                            : AppColors.textSecondary,
                      ),
                    ))
                .toList(),
          ),
          const SizedBox(height: 16),
          TextField(
            controller: _controller,
            keyboardType: TextInputType.number,
            inputFormatters: [FilteringTextInputFormatter.digitsOnly],
            onChanged: (_) => setState(() {}),
            decoration: const InputDecoration(
              labelText: 'Monto a recargar',
              prefixIcon: Icon(Icons.attach_money_rounded, size: 20),
            ),
          ),
          const SizedBox(height: 14),
          Container(
            padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 11),
            decoration: BoxDecoration(
              color: AppColors.warningSoft,
              borderRadius: AppRadius.mdAll,
            ),
            child: Row(
              children: [
                const Icon(Icons.info_outline, size: 17, color: Color(0xFFB45309)),
                const SizedBox(width: 9),
                Expanded(
                  child: Text(
                    'La recarga se cobra de verdad en la pasarela de pagos, no es una simulacion.',
                    style: theme.textTheme.labelSmall?.copyWith(
                      color: const Color(0xFFB45309),
                      fontWeight: FontWeight.w600,
                    ),
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
                : const Text('Continuar al pago'),
          ),
        ],
      ),
    );
  }
}
