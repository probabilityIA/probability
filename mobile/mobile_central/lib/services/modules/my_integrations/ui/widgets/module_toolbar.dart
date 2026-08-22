import 'package:flutter/material.dart';

import '../../../../../shared/theme/app_colors.dart';
import '../../../../../shared/theme/app_tokens.dart';
import '../../domain/sync_entities.dart';

class EnvironmentSpec {
  const EnvironmentSpec({
    required this.environment,
    required this.label,
    required this.icon,
    required this.hint,
    this.enabled = true,
    this.runLabel,
  });

  final SyncEnvironment environment;
  final String label;
  final IconData icon;
  final String hint;
  final bool enabled;
  final String? runLabel;
}

const List<EnvironmentSpec> coreEnvironments = <EnvironmentSpec>[
  EnvironmentSpec(
    environment: SyncEnvironment.overview,
    label: 'Vista general',
    icon: Icons.grid_view_rounded,
    hint: 'Cuantas ordenes entraron por cada origen y en que estado van.',
  ),
  EnvironmentSpec(
    environment: SyncEnvironment.products,
    label: 'Comparar productos',
    icon: Icons.swap_horiz_rounded,
    hint: 'Que producto esta en cada canal y cual falta por publicar.',
    runLabel: 'Comparar productos',
  ),
  EnvironmentSpec(
    environment: SyncEnvironment.data,
    label: 'Actualizar productos',
    icon: Icons.download_rounded,
    hint: 'Que dato del canal (nombre, imagen, categoria) puede entrar a Probability.',
  ),
  EnvironmentSpec(
    environment: SyncEnvironment.inventory,
    label: 'Sincronizar inventario',
    icon: Icons.refresh_rounded,
    hint: 'Cuanto stock tiene cada canal y cual quedaria distinto al de Probability.',
    runLabel: 'Sincronizar inventario',
  ),
  EnvironmentSpec(
    environment: SyncEnvironment.ordersCompare,
    label: 'Comparar ordenes',
    icon: Icons.fact_check_outlined,
    hint: 'Que orden existe en el canal y no en Probability, y crearla aca.',
  ),
  EnvironmentSpec(
    environment: SyncEnvironment.invoicing,
    label: 'Facturar',
    icon: Icons.receipt_long_outlined,
    hint: 'Facturacion desde el hub: proximamente.',
    enabled: false,
  ),
];

EnvironmentSpec environmentSpec(SyncEnvironment environment) =>
    coreEnvironments.firstWhere((spec) => spec.environment == environment);

class ModuleToolbar extends StatelessWidget implements PreferredSizeWidget {
  const ModuleToolbar({
    super.key,
    required this.environment,
    required this.onEnvironment,
    this.running = false,
  });

  final SyncEnvironment environment;
  final ValueChanged<SyncEnvironment> onEnvironment;
  final bool running;

  @override
  Size get preferredSize => const Size.fromHeight(48);

  @override
  Widget build(BuildContext context) {
    return Column(
      mainAxisSize: MainAxisSize.min,
      children: [
        SizedBox(
          height: 44,
          child: ListView.separated(
            scrollDirection: Axis.horizontal,
            padding: const EdgeInsets.fromLTRB(16, 0, 16, 0),
            itemCount: coreEnvironments.length,
            separatorBuilder: (context, index) => const SizedBox(width: 7),
            itemBuilder: (context, index) {
              final spec = coreEnvironments[index];
              return _EnvironmentChip(
                spec: spec,
                selected: spec.environment == environment,
                disabled: !spec.enabled || running,
                onTap: () => onEnvironment(spec.environment),
              );
            },
          ),
        ),
      ],
    );
  }
}

class _EnvironmentChip extends StatelessWidget {
  const _EnvironmentChip({
    required this.spec,
    required this.selected,
    required this.disabled,
    required this.onTap,
  });

  final EnvironmentSpec spec;
  final bool selected;
  final bool disabled;
  final VoidCallback onTap;

  @override
  Widget build(BuildContext context) {
    final scheme = Theme.of(context).colorScheme;
    final foreground = disabled
        ? AppColors.textDisabled
        : selected
            ? Colors.white
            : AppColors.textSecondary;

    return GestureDetector(
      onTap: disabled ? null : onTap,
      behavior: HitTestBehavior.opaque,
      child: Container(
        alignment: Alignment.center,
        padding: const EdgeInsets.symmetric(horizontal: 13),
        decoration: BoxDecoration(
          color: selected ? scheme.primary : AppColors.surface,
          borderRadius: AppRadius.pillAll,
          border: Border.all(
            color: selected ? scheme.primary : AppColors.border,
          ),
        ),
        child: Row(
          mainAxisSize: MainAxisSize.min,
          children: [
            Icon(spec.icon, size: 14, color: foreground),
            const SizedBox(width: 7),
            Text(
              spec.label,
              style: TextStyle(
                fontFamily: 'Inter',
                fontSize: 12.5,
                fontWeight: FontWeight.w600,
                color: foreground,
              ),
            ),
          ],
        ),
      ),
    );
  }
}
