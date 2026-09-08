import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:provider/provider.dart';
import '../../../../../core/security/biometric_auth.dart';
import '../../../../../shared/theme/app_colors.dart';
import '../../../../../shared/theme/app_tokens.dart';
import '../../../../../shared/utils/formatters.dart';
import '../../../../../shared/widgets/network_avatar.dart';
import '../../../../../shared/widgets/ui/ui.dart';
import '../../../login/ui/providers/login_provider.dart';
import '../widgets/change_password_sheet.dart';

class ProfileScreen extends StatelessWidget {
  const ProfileScreen({super.key});

  @override
  Widget build(BuildContext context) {
    final login = context.watch<LoginProvider>();
    final user = login.user;
    final permissions = login.rolesPermissions;

    return AppScaffold(
      title: 'Mi perfil',
      onBack: () => context.go('/more'),
      body: ListView(
        padding: AppSpacing.page,
        children: [
          AppCard(
            child: Column(
              children: [
                NetworkAvatar(
                  imageUrl: user?.avatarUrl,
                  fallbackText: user?.name ?? '?',
                  radius: 34,
                  backgroundColor: AppColors.primarySoft,
                  foregroundColor: AppColors.primary,
                ),
                const SizedBox(height: 14),
                Text(
                  user?.name ?? 'Usuario',
                  style: Theme.of(context).textTheme.titleLarge,
                  textAlign: TextAlign.center,
                ),
                const SizedBox(height: 4),
                Text(
                  user?.email ?? '',
                  style: Theme.of(context).textTheme.bodySmall,
                  textAlign: TextAlign.center,
                ),
                const SizedBox(height: 12),
                Wrap(
                  spacing: 8,
                  runSpacing: 8,
                  alignment: WrapAlignment.center,
                  children: [
                    if (login.roleName != null)
                      AppStatusChip(
                        label: login.roleName!,
                        tone: AppStatusTone.brand,
                        icon: Icons.badge_outlined,
                      ),
                    if (login.isSuperAdmin)
                      const AppStatusChip(
                        label: 'Super admin',
                        tone: AppStatusTone.warning,
                        icon: Icons.shield_outlined,
                      ),
                  ],
                ),
              ],
            ),
          ),
          const SizedBox(height: 18),
          const AppSectionHeader(title: 'Negocio'),
          AppCard(
            child: Column(
              children: [
                AppKeyValueRow(label: 'Nombre', value: login.businessName ?? '-'),
                const Divider(height: 18),
                AppKeyValueRow(
                  label: 'Tipo',
                  value: permissions?.businessTypeName ?? '-',
                ),
                const Divider(height: 18),
                AppKeyValueRow(
                  label: 'Suscripcion',
                  value: permissions?.subscriptionStatus ?? '-',
                ),
                const Divider(height: 18),
                AppKeyValueRow(
                  label: 'Modulos con acceso',
                  value: '${permissions?.resources.length ?? 0}',
                ),
              ],
            ),
          ),
          const SizedBox(height: 18),
          const AppSectionHeader(title: 'Datos de contacto'),
          AppCard(
            child: Column(
              children: [
                AppKeyValueRow(label: 'Correo', value: user?.email ?? '-'),
                const Divider(height: 18),
                AppKeyValueRow(label: 'Tel\u00e9fono', value: user?.phone ?? '-'),
                const Divider(height: 18),
                AppKeyValueRow(
                  label: 'Ultimo ingreso',
                  value: AppFormat.dateTime(AppFormat.parseDate(user?.lastLoginAt)),
                ),
              ],
            ),
          ),
          const SizedBox(height: 18),
          const AppSectionHeader(title: 'Seguridad'),
          AppCard(
            padding: EdgeInsets.zero,
            child: Column(
              children: [
                _ActionRow(
                  icon: Icons.lock_reset_outlined,
                  label: 'Cambiar contrase\u00F1a',
                  onTap: () => showChangePasswordSheet(context),
                ),
                const Divider(height: 1),
                const _BiometricRow(),
                const Divider(height: 1),
                _ActionRow(
                  icon: Icons.logout_rounded,
                  label: 'Cerrar sesi\u00f3n',
                  tone: AppColors.error,
                  onTap: () {
                    context.read<LoginProvider>().logout();
                    context.go('/login');
                  },
                ),
              ],
            ),
          ),
          const SizedBox(height: 26),
          Center(
            child: Text(
              'ProbabilityIA \u00B7 Central Mobile',
              style: Theme.of(context).textTheme.labelSmall,
            ),
          ),
          const SizedBox(height: 24),
        ],
      ),
    );
  }
}

class _ActionRow extends StatelessWidget {
  const _ActionRow({
    required this.icon,
    required this.label,
    required this.onTap,
    this.tone = AppColors.textSecondary,
  });

  final IconData icon;
  final String label;
  final VoidCallback onTap;
  final Color tone;

  @override
  Widget build(BuildContext context) {
    return InkWell(
      onTap: onTap,
      borderRadius: AppRadius.lgAll,
      child: Padding(
        padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 15),
        child: Row(
          children: [
            Icon(icon, size: 20, color: tone),
            const SizedBox(width: 13),
            Expanded(
              child: Text(
                label,
                style: Theme.of(context).textTheme.titleSmall?.copyWith(color: tone),
              ),
            ),
            const Icon(Icons.chevron_right_rounded, size: 20, color: AppColors.textDisabled),
          ],
        ),
      ),
    );
  }
}

class _BiometricRow extends StatefulWidget {
  const _BiometricRow();

  @override
  State<_BiometricRow> createState() => _BiometricRowState();
}

class _BiometricRowState extends State<_BiometricRow> {
  BiometricStatus? _status;
  bool _busy = false;

  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addPostFrameCallback((_) => _check());
  }

  Future<void> _check() async {
    final s = await context.read<LoginProvider>().biometricStatus();
    if (!mounted) return;
    setState(() => _status = s);
  }

  String get _subtitle {
    switch (_status) {
      case BiometricStatus.noHardware:
        return 'Este dispositivo no tiene lector de huella';
      case BiometricStatus.notEnrolled:
        return 'Registra una huella en los ajustes del tel\u00e9fono';
      case BiometricStatus.unavailable:
        return 'No disponible en este dispositivo';
      default:
        return 'Pide tu huella al abrir la app';
    }
  }

  Future<void> _toggle(bool value) async {
    final login = context.read<LoginProvider>();
    setState(() => _busy = true);
    var message = '';
    if (value) {
      final ok = await login.enableBiometric();
      message = ok
          ? 'Ingreso con huella activado'
          : 'No se pudo activar el ingreso con huella';
    } else {
      await login.disableBiometric();
      message = 'Ingreso con huella desactivado';
    }
    if (!mounted) return;
    setState(() => _busy = false);
    ScaffoldMessenger.of(context)
        .showSnackBar(SnackBar(content: Text(message)));
  }

  @override
  Widget build(BuildContext context) {
    final enabled = context.watch<LoginProvider>().biometricEnabled;
    final usable = _status == BiometricStatus.available;

    return ListTile(
      leading: const Icon(Icons.fingerprint_rounded, size: 20),
      title: Text(
        'Ingresar con huella',
        style: Theme.of(context).textTheme.titleSmall,
      ),
      subtitle: Text(
        _subtitle,
        style: Theme.of(context).textTheme.labelSmall,
      ),
      trailing: _busy
          ? const SizedBox(
              width: 18,
              height: 18,
              child: CircularProgressIndicator(strokeWidth: 2),
            )
          : Switch(
              value: enabled,
              onChanged: usable ? _toggle : null,
            ),
    );
  }
}
