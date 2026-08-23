import 'package:flutter/material.dart';

import '../../navigation/app_modules.dart';
import '../../theme/app_colors.dart';

class ModuleStageBadge extends StatelessWidget {
  const ModuleStageBadge({super.key, this.stage = ModuleStage.beta});

  final ModuleStage stage;

  @override
  Widget build(BuildContext context) {
    if (!stage.showsBadge) return const SizedBox.shrink();

    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 7, vertical: 2),
      decoration: BoxDecoration(
        color: AppColors.warningSoft,
        borderRadius: BorderRadius.circular(999),
      ),
      child: Text(
        stage.label.toUpperCase(),
        style: const TextStyle(
          fontFamily: 'Inter',
          fontSize: 9,
          fontWeight: FontWeight.w700,
          letterSpacing: 0.6,
          color: AppColors.warning,
          height: 1.2,
        ),
      ),
    );
  }
}
