import 'package:flutter/material.dart';
import 'app_colors.dart';

class AppSpacing {
  const AppSpacing._();

  static const double xs = 4;
  static const double sm = 8;
  static const double md = 12;
  static const double lg = 16;
  static const double xl = 24;
  static const double xxl = 32;

  static const EdgeInsets page = EdgeInsets.fromLTRB(16, 16, 16, 24);
  static const EdgeInsets card = EdgeInsets.all(16);
  static const EdgeInsets listItem = EdgeInsets.symmetric(horizontal: 16, vertical: 12);
}

class AppRadius {
  const AppRadius._();

  static const double sm = 8;
  static const double md = 12;
  static const double lg = 16;
  static const double xl = 24;
  static const double pill = 999;

  static final BorderRadius smAll = BorderRadius.circular(sm);
  static final BorderRadius mdAll = BorderRadius.circular(md);
  static final BorderRadius lgAll = BorderRadius.circular(lg);
  static final BorderRadius pillAll = BorderRadius.circular(pill);
}

class AppShadows {
  const AppShadows._();

  static const List<BoxShadow> soft = [
    BoxShadow(color: Color(0x0D000000), blurRadius: 2, offset: Offset(0, 1)),
  ];

  static const List<BoxShadow> card = [
    BoxShadow(color: Color(0x0F000000), blurRadius: 10, offset: Offset(0, 4)),
    BoxShadow(color: Color(0x08000000), blurRadius: 2, offset: Offset(0, 1)),
  ];

  static const List<BoxShadow> raised = [
    BoxShadow(color: Color(0x1A000000), blurRadius: 20, offset: Offset(0, 10)),
  ];
}

class AppBorders {
  const AppBorders._();

  static const BorderSide hairline = BorderSide(color: AppColors.border, width: 1);
  static Border all = Border.all(color: AppColors.border, width: 1);
}
