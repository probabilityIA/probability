import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../../../../../shared/theme/app_colors.dart';
import '../providers/login_provider.dart';

class LockScreen extends StatefulWidget {
  const LockScreen({super.key});

  @override
  State<LockScreen> createState() => _LockScreenState();
}

class _LockScreenState extends State<LockScreen> {
  bool _running = false;
  bool _failed = false;

  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addPostFrameCallback((_) => _unlock());
  }

  Future<void> _unlock() async {
    if (_running) return;
    setState(() {
      _running = true;
      _failed = false;
    });
    final ok = await context.read<LoginProvider>().unlockWithBiometric();
    if (!mounted) return;
    setState(() {
      _running = false;
      _failed = !ok;
    });
  }

  Future<void> _useAnotherAccount() async {
    await context.read<LoginProvider>().logout();
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final name = context.watch<LoginProvider>().user?.name;

    return Scaffold(
      backgroundColor: const Color(0xFF120326),
      body: SafeArea(
        child: Padding(
          padding: const EdgeInsets.symmetric(horizontal: 28),
          child: Column(
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              Container(
                width: 96,
                height: 96,
                decoration: const BoxDecoration(
                  shape: BoxShape.circle,
                  color: Color(0xFF2A0F5C),
                ),
                child: const Icon(
                  Icons.fingerprint_rounded,
                  size: 52,
                  color: Colors.white,
                ),
              ),
              const SizedBox(height: 28),
              Text(
                name == null || name.isEmpty
                    ? 'Sesi\u00f3n bloqueada'
                    : 'Hola, $name',
                textAlign: TextAlign.center,
                style: theme.textTheme.headlineSmall
                    ?.copyWith(color: Colors.white, fontWeight: FontWeight.w700),
              ),
              const SizedBox(height: 10),
              Text(
                _failed
                    ? 'No se pudo verificar tu identidad. Intenta de nuevo.'
                    : 'Usa tu huella para volver a entrar.',
                textAlign: TextAlign.center,
                style: theme.textTheme.bodyMedium
                    ?.copyWith(color: Colors.white70),
              ),
              const SizedBox(height: 32),
              SizedBox(
                width: double.infinity,
                child: FilledButton.icon(
                  onPressed: _running ? null : _unlock,
                  icon: _running
                      ? const SizedBox(
                          width: 18,
                          height: 18,
                          child: CircularProgressIndicator(
                            strokeWidth: 2,
                            color: Colors.white,
                          ),
                        )
                      : const Icon(Icons.fingerprint_rounded),
                  label: Text(_running ? 'Verificando' : 'Desbloquear'),
                  style: FilledButton.styleFrom(
                    backgroundColor: AppColors.primary,
                    padding: const EdgeInsets.symmetric(vertical: 16),
                  ),
                ),
              ),
              const SizedBox(height: 12),
              TextButton(
                onPressed: _running ? null : _useAnotherAccount,
                child: const Text(
                  'Entrar con otra cuenta',
                  style: TextStyle(color: Colors.white70),
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }
}
