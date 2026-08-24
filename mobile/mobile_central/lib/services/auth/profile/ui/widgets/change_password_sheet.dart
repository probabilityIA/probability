import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../../../../../shared/theme/app_colors.dart';
import '../../../../../shared/theme/app_tokens.dart';
import '../../../login/ui/providers/login_provider.dart';

Future<void> showChangePasswordSheet(BuildContext context) {
  return showModalBottomSheet<void>(
    context: context,
    isScrollControlled: true,
    backgroundColor: AppColors.surface,
    builder: (sheetContext) => Padding(
      padding: EdgeInsets.only(bottom: MediaQuery.of(sheetContext).viewInsets.bottom),
      child: const _ChangePasswordForm(),
    ),
  );
}

class _ChangePasswordForm extends StatefulWidget {
  const _ChangePasswordForm();

  @override
  State<_ChangePasswordForm> createState() => _ChangePasswordFormState();
}

class _ChangePasswordFormState extends State<_ChangePasswordForm> {
  final _current = TextEditingController();
  final _next = TextEditingController();
  final _confirm = TextEditingController();
  bool _busy = false;
  String? _error;

  @override
  void dispose() {
    _current.dispose();
    _next.dispose();
    _confirm.dispose();
    super.dispose();
  }

  Future<void> _submit() async {
    if (_next.text.length < 6) {
      setState(() => _error = 'La nueva contrase\u00F1a necesita al menos 6 caracteres');
      return;
    }
    if (_next.text != _confirm.text) {
      setState(() => _error = 'Las contrase\u00F1as no coinciden');
      return;
    }

    setState(() {
      _busy = true;
      _error = null;
    });

    final provider = context.read<LoginProvider>();
    final response = await provider.changePassword(_current.text, _next.text);

    if (!mounted) return;
    setState(() => _busy = false);

    if (response == null) {
      setState(() => _error = provider.error ?? 'No se pudo cambiar la contrase\u00F1a');
      return;
    }

    Navigator.pop(context);
    ScaffoldMessenger.of(context).showSnackBar(
      const SnackBar(content: Text('Contrase\u00F1a actualizada')),
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
          Text('Cambiar contrase\u00F1a', style: theme.textTheme.titleLarge),
          const SizedBox(height: 6),
          Text(
            'Necesitas tu contrase\u00F1a actual para confirmar el cambio.',
            style: theme.textTheme.bodySmall,
          ),
          const SizedBox(height: 20),
          TextField(
            controller: _current,
            obscureText: true,
            decoration: const InputDecoration(
              hintText: 'Contrase\u00F1a actual',
              prefixIcon: Icon(Icons.lock_outline, size: 20),
            ),
          ),
          const SizedBox(height: 12),
          TextField(
            controller: _next,
            obscureText: true,
            decoration: const InputDecoration(
              hintText: 'Nueva contrase\u00F1a',
              prefixIcon: Icon(Icons.lock_reset_outlined, size: 20),
            ),
          ),
          const SizedBox(height: 12),
          TextField(
            controller: _confirm,
            obscureText: true,
            decoration: const InputDecoration(
              hintText: 'Confirmar nueva contrase\u00F1a',
              prefixIcon: Icon(Icons.lock_reset_outlined, size: 20),
            ),
          ),
          if (_error != null) ...[
            const SizedBox(height: 14),
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
          const SizedBox(height: 20),
          FilledButton(
            onPressed: _busy ? null : _submit,
            style: FilledButton.styleFrom(minimumSize: const Size(0, 50)),
            child: _busy
                ? const SizedBox(
                    height: 19,
                    width: 19,
                    child: CircularProgressIndicator(strokeWidth: 2.2, color: Colors.white),
                  )
                : const Text('Guardar'),
          ),
        ],
      ),
    );
  }
}
