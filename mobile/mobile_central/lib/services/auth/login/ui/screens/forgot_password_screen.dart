import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:go_router/go_router.dart';
import 'package:provider/provider.dart';
import '../../../../../core/errors/error_parser.dart';
import '../../../../../shared/theme/app_colors.dart';
import '../../../../../shared/theme/app_tokens.dart';
import '../../../../../shared/widgets/ui/ui.dart';
import '../../domain/entities.dart';
import '../providers/login_provider.dart';

enum _Step { email, code, password, done }

class ForgotPasswordScreen extends StatefulWidget {
  const ForgotPasswordScreen({super.key});

  @override
  State<ForgotPasswordScreen> createState() => _ForgotPasswordScreenState();
}

class _ForgotPasswordScreenState extends State<ForgotPasswordScreen> {
  final _emailController = TextEditingController();
  final _codeController = TextEditingController();
  final _passwordController = TextEditingController();
  final _confirmController = TextEditingController();

  _Step _step = _Step.email;
  bool _busy = false;
  String? _error;
  String _channel = 'email';
  String? _resetToken;
  List<RecoveryChannel> _channels = const [];

  @override
  void dispose() {
    _emailController.dispose();
    _codeController.dispose();
    _passwordController.dispose();
    _confirmController.dispose();
    super.dispose();
  }

  Future<void> _run(Future<void> Function() action) async {
    setState(() {
      _busy = true;
      _error = null;
    });
    try {
      await action();
    } catch (e) {
      setState(() => _error = parseError(e));
    } finally {
      if (mounted) setState(() => _busy = false);
    }
  }

  Future<void> _loadChannels() async {
    final email = _emailController.text.trim();
    if (email.isEmpty || !email.contains('@')) {
      setState(() => _error = 'Ingresa un correo valido');
      return;
    }
    await _run(() async {
      final provider = context.read<LoginProvider>();
      try {
        _channels = await provider.recoveryChannels(email);
      } catch (_) {
        _channels = const [];
      }
      if (_channels.isNotEmpty) {
        _channel = _channels.firstWhere((c) => c.available, orElse: () => _channels.first).channel;
      }
      await provider.forgotPassword(email, _channel);
      if (mounted) setState(() => _step = _Step.code);
    });
  }

  Future<void> _verifyCode() async {
    final code = _codeController.text.trim();
    if (code.length != 6) {
      setState(() => _error = 'El codigo tiene 6 digitos');
      return;
    }
    await _run(() async {
      final response = await context.read<LoginProvider>().verifyOtp(
            _emailController.text.trim(),
            code,
          );
      _resetToken = response.token;
      if (mounted) setState(() => _step = _Step.password);
    });
  }

  Future<void> _savePassword() async {
    final password = _passwordController.text;
    if (password.length < 6) {
      setState(() => _error = 'Minimo 6 caracteres');
      return;
    }
    if (password != _confirmController.text) {
      setState(() => _error = 'Las contrase\u00F1as no coinciden');
      return;
    }
    await _run(() async {
      await context.read<LoginProvider>().resetPassword(_resetToken ?? '', password);
      if (mounted) setState(() => _step = _Step.done);
    });
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    return Scaffold(
      backgroundColor: AppColors.surface,
      appBar: AppBar(
        backgroundColor: AppColors.surface,
        leading: IconButton(
          icon: const Icon(Icons.arrow_back, size: 22),
          onPressed: () => context.go('/login'),
        ),
        title: const AppLogo(height: 26),
        shape: null,
      ),
      body: SafeArea(
        child: Center(
          child: ConstrainedBox(
            constraints: const BoxConstraints(maxWidth: 440),
            child: SingleChildScrollView(
              padding: const EdgeInsets.fromLTRB(24, 16, 24, 32),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.stretch,
                children: [
                  _StepIndicator(step: _step),
                  const SizedBox(height: 28),
                  Text(_titleFor(_step), style: theme.textTheme.displaySmall),
                  const SizedBox(height: 8),
                  Text(
                    _subtitleFor(_step),
                    style: theme.textTheme.bodyMedium?.copyWith(color: AppColors.textMuted),
                  ),
                  const SizedBox(height: 28),
                  ..._fieldsFor(_step),
                  if (_error != null) ...[
                    const SizedBox(height: 16),
                    Container(
                      padding: const EdgeInsets.symmetric(horizontal: 14, vertical: 12),
                      decoration: BoxDecoration(
                        color: AppColors.errorSoft,
                        borderRadius: AppRadius.mdAll,
                      ),
                      child: Row(
                        children: [
                          const Icon(Icons.error_outline, size: 18, color: AppColors.error),
                          const SizedBox(width: 10),
                          Expanded(
                            child: Text(
                              _error!,
                              style: theme.textTheme.bodySmall?.copyWith(
                                color: const Color(0xFFB91C1C),
                                fontWeight: FontWeight.w500,
                              ),
                            ),
                          ),
                        ],
                      ),
                    ),
                  ],
                  const SizedBox(height: 24),
                  FilledButton(
                    onPressed: _busy ? null : _primaryAction,
                    style: FilledButton.styleFrom(minimumSize: const Size(0, 52)),
                    child: _busy
                        ? const SizedBox(
                            height: 20,
                            width: 20,
                            child: CircularProgressIndicator(strokeWidth: 2.2, color: Colors.white),
                          )
                        : Text(_primaryLabelFor(_step)),
                  ),
                  if (_step == _Step.code) ...[
                    const SizedBox(height: 6),
                    TextButton(
                      onPressed: _busy ? null : _loadChannels,
                      child: const Text('Reenviar codigo'),
                    ),
                  ],
                ],
              ),
            ),
          ),
        ),
      ),
    );
  }

  VoidCallback get _primaryAction {
    switch (_step) {
      case _Step.email:
        return _loadChannels;
      case _Step.code:
        return _verifyCode;
      case _Step.password:
        return _savePassword;
      case _Step.done:
        return () => context.go('/login');
    }
  }

  String _titleFor(_Step step) {
    switch (step) {
      case _Step.email:
        return 'Recuperar acceso';
      case _Step.code:
        return 'Verifica el codigo';
      case _Step.password:
        return 'Nueva contrase\u00F1a';
      case _Step.done:
        return 'Listo';
    }
  }

  String _subtitleFor(_Step step) {
    switch (step) {
      case _Step.email:
        return 'Escribe el correo de tu cuenta y te enviamos un codigo de verificacion.';
      case _Step.code:
        return 'Enviamos un codigo de 6 digitos a ${_maskedTarget()}.';
      case _Step.password:
        return 'Elige una contrase\u00F1a de al menos 6 caracteres.';
      case _Step.done:
        return 'Tu contrase\u00F1a quedo actualizada. Ya puedes iniciar sesion.';
    }
  }

  String _primaryLabelFor(_Step step) {
    switch (step) {
      case _Step.email:
        return 'Enviar codigo';
      case _Step.code:
        return 'Verificar';
      case _Step.password:
        return 'Guardar contrase\u00F1a';
      case _Step.done:
        return 'Ir al inicio de sesion';
    }
  }

  String _maskedTarget() {
    final match = _channels.where((c) => c.channel == _channel).firstOrNull;
    return match?.masked ?? _emailController.text.trim();
  }

  List<Widget> _fieldsFor(_Step step) {
    switch (step) {
      case _Step.email:
        return [
          TextField(
            controller: _emailController,
            keyboardType: TextInputType.emailAddress,
            decoration: const InputDecoration(
              hintText: 'tucorreo@empresa.com',
              prefixIcon: Icon(Icons.mail_outline, size: 20),
            ),
          ),
          if (_channels.length > 1) ...[
            const SizedBox(height: 16),
            _ChannelPicker(
              channels: _channels,
              selected: _channel,
              onSelected: (value) => setState(() => _channel = value),
            ),
          ],
        ];
      case _Step.code:
        return [
          TextField(
            controller: _codeController,
            keyboardType: TextInputType.number,
            maxLength: 6,
            textAlign: TextAlign.center,
            inputFormatters: [FilteringTextInputFormatter.digitsOnly],
            style: const TextStyle(
              fontFamily: 'Inter',
              fontSize: 26,
              fontWeight: FontWeight.w700,
              letterSpacing: 10,
            ),
            decoration: const InputDecoration(counterText: '', hintText: '000000'),
          ),
        ];
      case _Step.password:
        return [
          TextField(
            controller: _passwordController,
            obscureText: true,
            decoration: const InputDecoration(
              hintText: 'Nueva contrase\u00F1a',
              prefixIcon: Icon(Icons.lock_outline, size: 20),
            ),
          ),
          const SizedBox(height: 14),
          TextField(
            controller: _confirmController,
            obscureText: true,
            decoration: const InputDecoration(
              hintText: 'Confirmar contrase\u00F1a',
              prefixIcon: Icon(Icons.lock_outline, size: 20),
            ),
          ),
        ];
      case _Step.done:
        return [
          Container(
            padding: const EdgeInsets.all(18),
            decoration: BoxDecoration(
              color: AppColors.successSoft,
              borderRadius: AppRadius.lgAll,
            ),
            child: const Row(
              children: [
                Icon(Icons.check_circle_outline, color: Color(0xFF047857)),
                SizedBox(width: 12),
                Expanded(
                  child: Text(
                    'Contrase\u00F1a actualizada correctamente',
                    style: TextStyle(
                      fontFamily: 'Inter',
                      fontWeight: FontWeight.w600,
                      color: Color(0xFF047857),
                    ),
                  ),
                ),
              ],
            ),
          ),
        ];
    }
  }
}

class _ChannelPicker extends StatelessWidget {
  const _ChannelPicker({
    required this.channels,
    required this.selected,
    required this.onSelected,
  });

  final List<RecoveryChannel> channels;
  final String selected;
  final ValueChanged<String> onSelected;

  @override
  Widget build(BuildContext context) {
    return Column(
      children: channels.map((channel) {
        final active = channel.channel == selected;
        return Padding(
          padding: const EdgeInsets.only(bottom: 8),
          child: AppCard(
            padding: const EdgeInsets.symmetric(horizontal: 14, vertical: 12),
            borderColor: active ? AppColors.primary : AppColors.border,
            onTap: channel.available ? () => onSelected(channel.channel) : null,
            child: Row(
              children: [
                Icon(
                  channel.channel == 'whatsapp' ? Icons.chat_outlined : Icons.mail_outline,
                  size: 20,
                  color: active ? AppColors.primary : AppColors.textMuted,
                ),
                const SizedBox(width: 12),
                Expanded(
                  child: Text(
                    channel.masked ?? channel.channel,
                    style: Theme.of(context).textTheme.titleSmall,
                  ),
                ),
                if (active)
                  const Icon(Icons.check_circle, size: 18, color: AppColors.primary),
              ],
            ),
          ),
        );
      }).toList(),
    );
  }
}

class _StepIndicator extends StatelessWidget {
  const _StepIndicator({required this.step});

  final _Step step;

  @override
  Widget build(BuildContext context) {
    final index = _Step.values.indexOf(step);
    return Row(
      children: List.generate(_Step.values.length, (i) {
        final done = i <= index;
        return Expanded(
          child: Container(
            height: 4,
            margin: EdgeInsets.only(right: i == _Step.values.length - 1 ? 0 : 6),
            decoration: BoxDecoration(
              color: done ? AppColors.primary : AppColors.border,
              borderRadius: BorderRadius.circular(2),
            ),
          ),
        );
      }),
    );
  }
}
