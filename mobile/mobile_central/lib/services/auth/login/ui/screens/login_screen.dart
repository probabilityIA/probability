import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../../../../../core/config/environment.dart';
import '../../../../../shared/theme/app_colors.dart';
import '../../../../../shared/theme/app_tokens.dart';
import '../../../../../shared/widgets/ui/ui.dart';
import '../providers/login_provider.dart';

class LoginScreen extends StatefulWidget {
  const LoginScreen({super.key});

  @override
  State<LoginScreen> createState() => _LoginScreenState();
}

class _LoginScreenState extends State<LoginScreen> {
  final _formKey = GlobalKey<FormState>();
  late final _emailController = TextEditingController(text: Environment.devEmail);
  late final _passwordController = TextEditingController(text: Environment.devPassword);
  bool _obscurePassword = true;

  @override
  void dispose() {
    _emailController.dispose();
    _passwordController.dispose();
    super.dispose();
  }

  Future<void> _handleLogin() async {
    if (!_formKey.currentState!.validate()) return;
    FocusScope.of(context).unfocus();
    await context.read<LoginProvider>().login(
          _emailController.text.trim(),
          _passwordController.text,
        );
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    return Scaffold(
      backgroundColor: AppColors.surface,
      body: SafeArea(
        child: Center(
          child: ConstrainedBox(
            constraints: const BoxConstraints(maxWidth: 440),
            child: SingleChildScrollView(
              padding: const EdgeInsets.fromLTRB(24, 32, 24, 32),
              child: Form(
                key: _formKey,
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.stretch,
                  children: [
                    const SizedBox(height: 12),
                    const Align(
                      alignment: Alignment.centerLeft,
                      child: AppLogo(height: 34),
                    ),
                    const SizedBox(height: 44),
                    Text(
                      'Hola de nuevo',
                      style: theme.textTheme.displaySmall,
                    ),
                    const SizedBox(height: 8),
                    Text(
                      'Ingresa a tu operacion: ordenes, envios, inventario y pagos en un solo lugar.',
                      style: theme.textTheme.bodyMedium?.copyWith(color: AppColors.textMuted),
                    ),
                    const SizedBox(height: 32),
                    _FieldLabel('Correo electr\u00F3nico'),
                    const SizedBox(height: 7),
                    TextFormField(
                      controller: _emailController,
                      keyboardType: TextInputType.emailAddress,
                      autofillHints: const [AutofillHints.email],
                      decoration: const InputDecoration(
                        hintText: 'tucorreo@empresa.com',
                        prefixIcon: Icon(Icons.mail_outline, size: 20),
                      ),
                      validator: (value) {
                        if (value == null || value.trim().isEmpty) return 'Ingresa tu correo';
                        if (!value.contains('@')) return 'Correo inv\u00E1lido';
                        return null;
                      },
                    ),
                    const SizedBox(height: 18),
                    _FieldLabel('Contrase\u00F1a'),
                    const SizedBox(height: 7),
                    TextFormField(
                      controller: _passwordController,
                      obscureText: _obscurePassword,
                      autofillHints: const [AutofillHints.password],
                      decoration: InputDecoration(
                        hintText: '\u2022' * 8,
                        prefixIcon: const Icon(Icons.lock_outline, size: 20),
                        suffixIcon: IconButton(
                          icon: Icon(
                            _obscurePassword ? Icons.visibility_off_outlined : Icons.visibility_outlined,
                            size: 20,
                          ),
                          onPressed: () => setState(() => _obscurePassword = !_obscurePassword),
                        ),
                      ),
                      validator: (value) {
                        if (value == null || value.isEmpty) return 'Ingresa tu contrase\u00F1a';
                        return null;
                      },
                      onFieldSubmitted: (_) => _handleLogin(),
                    ),
                    const SizedBox(height: 10),
                    Align(
                      alignment: Alignment.centerRight,
                      child: TextButton(
                        onPressed: () {},
                        child: const Text('\u00BFOlvidaste tu contrase\u00F1a?'),
                      ),
                    ),
                    const SizedBox(height: 10),
                    Consumer<LoginProvider>(
                      builder: (context, provider, _) {
                        if (provider.error == null) return const SizedBox.shrink();
                        return Container(
                          margin: const EdgeInsets.only(bottom: 16),
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
                                  provider.error!,
                                  style: theme.textTheme.bodySmall?.copyWith(
                                    color: const Color(0xFFB91C1C),
                                    fontWeight: FontWeight.w500,
                                  ),
                                ),
                              ),
                            ],
                          ),
                        );
                      },
                    ),
                    Consumer<LoginProvider>(
                      builder: (context, provider, _) {
                        return FilledButton(
                          onPressed: provider.isLoading ? null : _handleLogin,
                          style: FilledButton.styleFrom(minimumSize: const Size(0, 52)),
                          child: provider.isLoading
                              ? const SizedBox(
                                  height: 20,
                                  width: 20,
                                  child: CircularProgressIndicator(strokeWidth: 2.2, color: Colors.white),
                                )
                              : const Text('Iniciar sesi\u00F3n'),
                        );
                      },
                    ),
                    const SizedBox(height: 36),
                    Row(
                      mainAxisAlignment: MainAxisAlignment.center,
                      children: [
                        const Icon(Icons.shield_outlined, size: 14, color: AppColors.textDisabled),
                        const SizedBox(width: 6),
                        Text(
                          'Conexi\u00F3n segura \u00B7 ProbabilityIA',
                          style: theme.textTheme.labelSmall,
                        ),
                      ],
                    ),
                  ],
                ),
              ),
            ),
          ),
        ),
      ),
    );
  }
}

class _FieldLabel extends StatelessWidget {
  const _FieldLabel(this.text);

  final String text;

  @override
  Widget build(BuildContext context) {
    return Text(
      text,
      style: Theme.of(context).textTheme.titleSmall?.copyWith(
            color: AppColors.textSecondary,
            fontSize: 13.5,
          ),
    );
  }
}
